package baseline

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func repositoryRecordCitations(clauses []characterizedBaselineClause) [][2]string {
	pattern := regexp.MustCompile(`ADR-[0-9]+|Spec [0-9]{4}`)
	var citations [][2]string
	for _, clause := range clauses {
		for _, match := range pattern.FindAllString(clause.Guidance, -1) {
			citations = append(citations, [2]string{clause.ID, match})
		}
	}
	return citations
}

func TestShippedGuidanceCitesNoRepositoryRecord(t *testing.T) {
	if citations := repositoryRecordCitations(embeddedBaselineClauses(t)); len(citations) != 0 {
		t.Fatalf("shipped guidance cites repository records: %v", citations)
	}
}

func TestARepositoryRecordCitationIsReported(t *testing.T) {
	clauses := []characterizedBaselineClause{{ID: "clause.test.citation", Guidance: "Follow ADR-0186."}}
	want := [][2]string{{"clause.test.citation", "ADR-0186"}}
	if got := repositoryRecordCitations(clauses); !reflect.DeepEqual(got, want) {
		t.Fatalf("citations = %v, want %v", got, want)
	}
}

func TestASpecCitationInRuleGuidanceIsReported(t *testing.T) {
	clauses := []characterizedBaselineClause{{ID: "rule.test.citation", RuleLevel: true, Guidance: "Follow Spec 0193."}}
	want := [][2]string{{"rule.test.citation", "Spec 0193"}}
	if got := repositoryRecordCitations(clauses); !reflect.DeepEqual(got, want) {
		t.Fatalf("citations = %v, want %v", got, want)
	}
}

func TestLifecycleClausesCarryTheScopedWording(t *testing.T) {
	want := map[string][]string{
		"clause.context.adr-02-active-status": {"For an ADR that carries lifecycle frontmatter, only `accepted` is active."},
		"clause.context.backlog-01-operational-contract": {
			"# open | promoted | declined | deferred | done",
			"Set `status: declined` or `status: deferred` only with a non-null `reason`.",
			"(`declined`, `deferred`, `done`,",
		},
		"clause.context.inbox-01-triage":            {"Preserve the boundary between evidence and intent:"},
		"clause.context.docs-one-job-per-directory": {"A finished orphan Review Artifact retires"},
	}
	for _, clause := range embeddedBaselineClauses(t) {
		phrases, ok := want[clause.ID]
		if !ok {
			continue
		}
		for _, phrase := range phrases {
			if !strings.Contains(clause.Guidance, phrase) {
				t.Errorf("%s lacks %q", clause.ID, phrase)
			}
		}
		delete(want, clause.ID)
	}
	if len(want) != 0 {
		t.Fatalf("missing lifecycle clauses: %v", want)
	}
	for _, path := range []string{
		"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md",
		"../../docs/agents/docs-layout.md",
	} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := strings.Join(strings.Fields(string(content)), " ")
		for _, phrase := range []string{"For an ADR that carries lifecycle frontmatter, only `accepted` is active.", "status: deferred", "Preserve the boundary between evidence and intent:"} {
			if !strings.Contains(text, phrase) {
				t.Errorf("%s lacks %q", path, phrase)
			}
		}
	}
}
