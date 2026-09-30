package baseline

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestLoopClauseNamesTheDeliveryQueueAndItsRecoveryActs(t *testing.T) {
	for source, clauses := range loopClauseSurfaces(t) {
		t.Run(source, func(t *testing.T) {
			for _, phrase := range []string{
				"apply the configured pre-PR review policy, archive and commit the candidate",
				"Delivery Queue", "roundfix deliver", "roundfix reopen --spec <slug>",
				"Task Carry-Forward", "roundfix reconcile <run-id> --carry-forward",
				"roundfix deliver retry <slug>",
			} {
				if !strings.Contains(clauses[0], phrase) {
					t.Errorf("loop clause lacks %q", phrase)
				}
			}
		})
	}
}

func TestLoopClauseCitesNoSpecOrDecisionNumber(t *testing.T) {
	citation := regexp.MustCompile(`ADR-[0-9]+|Spec [0-9]{4}`)
	for source, clauses := range loopClauseSurfaces(t) {
		t.Run(source, func(t *testing.T) {
			for _, clause := range clauses {
				if match := citation.FindString(clause); match != "" {
					t.Errorf("clause cites repository record %q", match)
				}
			}
		})
	}
}

// Inspect both embedded source clauses and the paragraphs an adopter reads.
func loopClauseSurfaces(t *testing.T) map[string][]string {
	t.Helper()
	ids := []string{"clause.autonomous.loop-01-qa-once", "clause.autonomous.hook-strictness"}
	var source []string
	for _, id := range ids {
		found := false
		for _, clause := range embeddedBaselineClauses(t) {
			if clause.ID == id {
				if clause.Enforcement != "mandatory" {
					t.Fatalf("%s enforcement = %s, want mandatory", id, clause.Enforcement)
				}
				source = append(source, clause.Guidance)
				found = true
			}
		}
		if !found {
			t.Fatalf("missing clause %s", id)
		}
	}
	surfaces := map[string][]string{"module": source}
	for _, path := range []string{
		"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md",
		"../../docs/agents/autonomous-work.md",
	} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, marker := range []string{"Follow one order per Spec:", "A commit hook must not be stricter"} {
			found := false
			for _, paragraph := range strings.Split(string(content), "\n\n") {
				if strings.Contains(paragraph, marker) {
					surfaces[path] = append(surfaces[path], strings.Join(strings.Fields(paragraph), " "))
					found = true
				}
			}
			if !found {
				t.Fatalf("%s lacks clause marker %q", path, marker)
			}
		}
	}
	return surfaces
}
