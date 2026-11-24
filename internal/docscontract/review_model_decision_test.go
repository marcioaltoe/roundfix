//go:build docscontract

//verify:always

package docscontract

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
)

func TestTheReviewModelDecisionAmendsTheLineageRule(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	recordPath, record := reviewModelDecisionRecord(t)
	header := spec.ReadADRHeader(record)
	if !header.Frontmatter || header.Status != "accepted" || header.Knowledge == nil || !*header.Knowledge {
		t.Fatalf("review model decision header = %+v, want accepted with true Knowledge Mark", header)
	}
	number, _, ok := strings.Cut(filepath.Base(recordPath), "-")
	if !ok {
		t.Fatalf("decision record has no numeric identifier: %s", recordPath)
	}
	identifier := "ADR-" + number
	var metadata struct {
		CreatedAt string `yaml:"created_at"`
	}
	parts := strings.SplitN(string(record), "---", 3)
	if len(parts) != 3 {
		t.Fatal("decision record lacks lifecycle frontmatter")
	}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		t.Fatal(err)
	}
	created, err := time.Parse(time.RFC3339, metadata.CreatedAt)
	if err != nil {
		t.Fatalf("decision created_at: %v", err)
	}
	prior, err := os.ReadFile(filepath.Join(root, "docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md"))
	if err != nil {
		t.Fatal(err)
	}
	amendments := spec.ReadADRHeader(prior).AmendedBy
	if len(amendments) == 0 || amendments[len(amendments)-1] != identifier {
		t.Errorf("ADR-0197 amended_by = %v, missing %s", amendments, identifier)
	}
	// Preserve the original decision and all three earlier amendment notices.
	const historicalBody = "\n# A pre-PR reviewer lineage spans at most two rounds\n\nEach `roundfix review` opened a fresh reviewer session and reread the whole\ncandidate diff, even when the head had only gained the fix for the last\nround's findings. The reviewer did not know what it had raised before, nor how\neach finding was disposed. The maintainer's rule of 2026-09-09 allows two\nreview rounds per candidate, but it lived only in operating notes.\n\nNow reviews of one candidate form a Reviewer Lineage: the same checkout,\nprovider and merge base, each head descending from the one reviewed before.\n\n- Round 1 reviews the full candidate diff.\n- Round 2 reviews only the diff from the round-1 head to the current head. Its\n  prompt always carries the round-1 findings with their validation and\n  disposition, so it never depends on the reviewer's memory.\n- When round 1 leaves standing findings, the round-1 Agent Session is left open\n  and round 2 continues it through acpx, on the same selection. Whether the\n  runtime resumed it is measured, not assumed: the round records the ACP\n  session id of each round, and `continued` is true only when both ids match.\n  A runtime that cannot resume gets a new session, and round 2 still has the\n  recorded conversation.\n- A third review in the lineage never calls the reviewer. It closes the\n  lineage as `ceiling-closed` when every standing round-2 finding has an\n  operator disposition: fixed by a commit the head contains, or dismissed with\n  evidence. Otherwise it prints a blocked result that names the findings to\n  dispose. It keeps the round-2 record in place, so `roundfix review dispose`\n  still reads it.\n\n## Consequences\n\nA rebase moves the merge base, so it starts a new lineage at round 1. A head\nthat does not descend from the recorded one also starts a new lineage, and\nRoundfix closes the old lineage's open session first. The session is closed\nafter round 2, after any round that leaves nothing standing, and when a reuse\nat the same head resolves to `findings-dismissed`.\n\nA pre-PR review is not a Run. ADR-0018 and ADR-0051 still forbid carrying a\nsession across Runs, and this lineage does not touch Run sessions.\n\n`ceiling-closed` is not a reviewer pass. It records that the change after\nround 2 is the final correction, with the operator dispositions that close it.\nIt names the head the reviewer last saw. The Delivery Queue advances on it as\nit advances on `findings-dismissed`.\n\nThis decision refines\n[ADR-0153](0153-pre-pr-review-is-an-explicit-provider-policy.md) and\n[ADR-0051](0051-tasks-and-qa-own-agent-sessions.md), which consolidated\nADR-0018. It supersedes none of them.\n\n**Amended by [ADR-0286](0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md) (2026-10-18).** When round 2 ends with findings and every standing round-2 finding is disposed as fixed by a commit the current head contains, the next review is one Correction Review of the delta from the round-2 head instead of the ceiling; it never repeats, and the ceiling's purpose, no unbounded review loop, stays.\n\n**Amended by [ADR-0366](0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md) (2026-11-05).** A conflict-free merge of the default branch into a descending head keeps the Reviewer Lineage although it moves the merge base; round 2 reviews the delta from the reviewed head carried onto the new merge base. A rebase, or a merge that brings in another branch or conflicts, still starts round 1.\n\n**Amended by [ADR-0391](0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md) (2026-11-11).** A merge of the default branch that conflicts keeps the Reviewer Lineage too; its delta starts from the reviewed head's Conflicted Carry. A rebase or a merge that brings in another branch still starts round 1, and the two-round ceiling is unchanged.\n"
	noticeLead := "**Amended by [" + identifier + "](" + filepath.Base(recordPath) + ") (" + created.Format("2006-01-02") + ").**"
	_, body, ok := strings.Cut(string(prior), "---\n\n")
	if !ok || !strings.HasPrefix(body, strings.TrimPrefix(historicalBody, "\n")+"\n"+noticeLead) {
		t.Error("ADR-0197 must preserve its decision and three earlier notices before the new dated notice")
	}
	if !strings.Contains(string(prior), "updated_at: "+metadata.CreatedAt+"\n") {
		t.Error("ADR-0197 updated_at must equal the review model decision's created_at")
	}
	export, err := os.ReadFile(filepath.Join(root, ".secondbrain-export"))
	if err != nil {
		t.Fatal(err)
	}
	const begin = "# roundfix:standard-export:begin"
	const end = "# roundfix:standard-export:end"
	lines := strings.Split(string(export), "\n")
	if strings.Count(string(export), begin) != 1 || strings.Count(string(export), end) != 1 {
		t.Fatal("Standard Export List must have one pair of markers")
	}
	start, stop := slices.Index(lines, begin), slices.Index(lines, end)
	if start < 0 || stop <= start {
		t.Fatal("Standard Export List markers are out of order")
	}
	var adrPaths []string
	count := 0
	for i, line := range lines {
		if line == recordPath {
			count++
			if i <= start || i >= stop {
				t.Error("review model decision path is outside the Standard Export List")
			}
		}
		if i > start && i < stop && strings.HasPrefix(line, "docs/adr/") {
			adrPaths = append(adrPaths, line)
		}
	}
	if count != 1 {
		t.Errorf("review model decision path occurs %d times in the export, want once", count)
	}
	if !slices.IsSorted(adrPaths) {
		t.Errorf("Standard Export List ADR paths are out of order: %v", adrPaths)
	}
}

func reviewModelDecisionRecord(t *testing.T) (string, []byte) {
	t.Helper()
	root := characterizationRepositoryRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "docs/adr"))
	if err != nil {
		t.Fatal(err)
	}
	const title = "# Round 2 of a pre-PR review opens a new session when the configured model changed"
	var recordPath string
	var record []byte
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join("docs/adr", entry.Name())
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(content), "\n") {
			if line != title {
				continue
			}
			if recordPath != "" {
				t.Fatalf("decision title occurs in both %s and %s", recordPath, path)
			}
			recordPath, record = filepath.ToSlash(path), content
		}
	}
	if recordPath == "" {
		t.Fatal("no decision record has the review model title")
	}
	return recordPath, record
}

func TestTheReviewerLineageEntryStatesTheModelRule(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	path, _ := reviewModelDecisionRecord(t)
	content, err := os.ReadFile(filepath.Join(root, "CONTEXT.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, entry, ok := strings.Cut(string(content), "**Reviewer Lineage**:")
	if !ok {
		t.Fatal("Reviewer Lineage entry missing")
	}
	entry, _, _ = strings.Cut(entry, "\n_Avoid_")
	text := strings.Join(strings.Fields(entry), " ")
	const clause = "continuing the same Agent Session when the runtime can resume it and the review profile still names the Agent Model that session was opened with, and otherwise opening a new one with the same prompt."
	if !strings.Contains(text, clause) {
		t.Errorf("Reviewer Lineage lacks %q", clause)
	}
	number, _, _ := strings.Cut(filepath.Base(path), "-")
	if !strings.Contains(entry, "[ADR-"+number+"]("+path+")") {
		t.Error("Reviewer Lineage lacks the model decision link")
	}
}

func TestTheReviewGuideAndReferenceDescribeTheModelRule(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	for _, doc := range []struct{ name, path string }{
		{"guide", "docs/user-guide/commands/review.md"},
		{"canonical_reference", ".agents/skills/roundfix/references/review.md"},
		{"mirror_reference", "skills/roundfix/references/review.md"},
	} {
		t.Run(doc.name, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(root, doc.path))
			if err != nil {
				t.Fatal(err)
			}
			text := strings.Join(strings.Fields(string(content)), " ")
			for _, phrase := range []string{
				"the model of that selection as `lineage.model`",
				"only when the selection at that index still names the recorded model",
				"or when the record holds no model because an earlier Roundfix wrote it",
				"and the model of its session in `lineage.model`",
			} {
				if !strings.Contains(text, phrase) {
					t.Errorf("%s lacks %q", doc.path, phrase)
				}
			}
			if strings.Contains(text, "Round 2 prepares that name on that selection through acpx. If preparation") {
				t.Error("replaced unconditional continuation sentence remains")
			}
		})
	}
}
