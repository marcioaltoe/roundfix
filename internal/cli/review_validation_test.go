// Suite: candidate anchor validation.
// Invariant: only a finding naming a candidate diff line can park review.
// Boundary IN: untrusted final answers and Git unified diffs through the public CLI.
// Boundary OUT: real reviewers, convention judgment, and reviewer lineage.
package cli

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/agent"
)

func TestReviewFindingAnchorParsesEachForm(t *testing.T) {
	tests := []struct {
		text       string
		path       string
		start, end int
		ok         bool
	}{
		{"`path.go:2` Failure: broken", "path.go", 2, 2, true},
		{"path.go:2: broken", "path.go", 2, 2, true},
		{"`path.go:2-4` broken", "path.go", 2, 4, true},
		{"path.go:2-4,8-9 broken", "path.go", 2, 4, true},
		{"path.go:2)", "path.go", 2, 2, true},
		{"path.go:2", "path.go", 2, 2, true},
		{"path.go:2\u00a0broken", "path.go", 2, 2, true},
		{"bad\u00a0path.go:2 broken", "", 0, 0, false},
		{"path.go: broken", "", 0, 0, false},
		{"path.go:2a broken", "", 0, 0, false},
		{"path.go:2-4a broken", "", 0, 0, false},
		{"finding path.go:2 broken", "", 0, 0, false},
		{" path.go:2 broken", "", 0, 0, false},
		{"path.go:4-2 broken", "", 0, 0, false},
		{"path.go:9999999999999999999999999999 broken", "", 0, 0, false},
	}
	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			got, ok := parseReviewFindingAnchor(test.text)
			if ok != test.ok || (ok && got != (reviewFindingAnchor{test.path, test.start, test.end})) {
				t.Fatalf("anchor=%+v readable=%v", got, ok)
			}
		})
	}
}

func TestReviewDiffIndexHoldsHunkLinesAndWholeFiles(t *testing.T) {
	diff := `diff --git a/hunk.go b/hunk.go
--- a/hunk.go
+++ b/hunk.go
@@ -10,2 +10,3 @@
 context
-old
+new
+added
@@ -30 +31 @@
-old
+new
 diff content
 diff content
 diff content
diff --git a/delete.go b/delete.go
--- a/delete.go
+++ b/delete.go
@@ -5,2 +5,0 @@
-old
-old
diff --git a/gone.go b/gone.go
--- a/gone.go
+++ /dev/null
@@ -1 +0,0 @@
-old
diff --git a/binary b/binary
Binary files a/binary and b/binary differ
diff --git a/old.go b/renamed.go
similarity index 100%
rename from old.go
rename to renamed.go
diff --git a/mode.go b/mode.go
old mode 100644
new mode 100755
`
	index := indexReviewDiff(diff)
	tests := []struct {
		path       string
		start, end int
		want       bool
	}{
		{"hunk.go", 10, 10, true}, {"hunk.go", 11, 11, true}, {"hunk.go", 12, 12, true},
		{"hunk.go", 9, 9, false}, {"hunk.go", 13, 30, false}, {"hunk.go", 31, 31, true},
		{"hunk.go", 9, 10, true}, {"delete.go", 5, 5, true}, {"delete.go", 6, 6, false},
		{"gone.go", 80, 90, true}, {"binary", 80, 90, true}, {"renamed.go", 80, 90, true},
		{"mode.go", 80, 90, true}, {"missing", 1, 1, false},
	}
	for _, test := range tests {
		t.Run(test.path+":"+strconv.Itoa(test.start), func(t *testing.T) {
			if got := index.holds(reviewFindingAnchor{test.path, test.start, test.end}); got != test.want {
				t.Fatalf("holds(%+v)=%v", test, got)
			}
		})
	}
}

func TestReviewPromptCarriesTheDeliveryConventionsAndTheFindingGrammar(t *testing.T) {
	prompt := buildReviewPrompt("base", "head", "diff")
	for _, text := range []string{
		"Start each finding with its anchor, `path:line` or `path:start-end`, naming a line of the candidate diff, and state what breaks in a clause that starts with `Failure:`.",
		"If there are findings, respond with Findings: on the first line, followed by each finding as a list item that starts with `- ` and names its file and line. If there are no findings, respond exactly: No findings. Include no other prose.\n\n",
		"Delivery Conventions (" + deliveryConventionsVersion + "). A Roundfix delivery writes the following by design; do not report them as defects:",
	} {
		if !strings.Contains(prompt, text) {
			t.Fatalf("prompt missing %q", text)
		}
	}
	for _, convention := range deliveryConventions() {
		if !strings.Contains(prompt, convention.ID+". "+convention.Prompt+"\n") {
			t.Fatalf("missing convention %s", convention.ID)
		}
	}
	if strings.Index(prompt, "Delivery Conventions") < strings.Index(prompt, "Head commit: head") {
		t.Fatal("conventions precede head")
	}
}

func validationReviewFixture(t *testing.T, answer string) (reviewCommandFixture, *reviewCommandRunner) {
	t.Helper()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: answer, StopReason: "end_turn"}}}}
	return newReviewCommandFixture(t, "codex", runner), runner
}

func TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne(t *testing.T) {
	const findings = "- internal/untouched.go:10: outside. Failure: broken.\n- `review.txt:2` inside. Failure: broken."
	fixture, _ := validationReviewFixture(t, "Findings:\n"+findings)
	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || record.Findings != findings {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	want := "roundfix: review finding F1 dismissed by validation (unanchored): anchor names no line of the candidate diff\n"
	if stderr != want {
		t.Fatalf("stderr=%q", stderr)
	}
	if len(record.FindingItems) != 2 {
		t.Fatalf("items=%+v", record.FindingItems)
	}
	first, second := record.FindingItems[0], record.FindingItems[1]
	if first.Anchor == nil || *first.Anchor != (reviewFindingAnchor{"internal/untouched.go", 10, 10}) || !reviewFindingDismissedByValidation(first) || first.Validation.Rule != "unanchored" {
		t.Fatalf("F1=%+v", first)
	}
	if second.Anchor == nil || *second.Anchor != (reviewFindingAnchor{"review.txt", 2, 2}) || second.Validation.Status != "stands" || second.Validation.Reason != "anchored in the candidate diff" {
		t.Fatalf("F2=%+v", second)
	}
	if record.Validation == nil || record.Validation.Conventions != deliveryConventionsVersion || record.Validation.Validator != "not-needed" {
		t.Fatalf("validation=%+v", record.Validation)
	}
	disk, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil || disk.FindingItems[0].Validation.Rule != "unanchored" {
		t.Fatalf("persisted record=%+v err=%v", disk, err)
	}
}

func TestReviewFindingsAllDismissedByValidationExitZero(t *testing.T) {
	fixture, _ := validationReviewFixture(t, "Findings:\n- review.txt:80: outside\n- missing.go:2: outside")
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || len(record.Dispositions) != 0 {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	if strings.Count(stderr, "dismissed by validation (unanchored)") != 2 {
		t.Fatalf("stderr=%q", stderr)
	}
	if _, err := os.Stat(filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName)); !os.IsNotExist(err) {
		t.Fatalf("ledger err=%v", err)
	}
	if err := validateReviewRecord(record); err != nil {
		t.Fatal(err)
	}
	record.Outcome = reviewOutcomeFindings
	if err := validateReviewRecord(record); err == nil {
		t.Fatal("findings accepted with no standing finding")
	}
}

func TestReviewBlocksFindingsThatNameNoFileAndLine(t *testing.T) {
	fixture, _ := validationReviewFixture(t, "Findings:\n- missing anchor\n- review.txt:2a: malformed anchor")
	code, record, stderr := fixture.run(t)
	if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || record.Reason != "findings name no file and line" || record.Findings != "" || len(record.FindingItems) != 0 || record.AnswerPath == "" {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	if stderr != "roundfix: review blocked: findings name no file and line\n" {
		t.Fatalf("stderr=%q", stderr)
	}
	if record.Validation == nil || record.Validation.Validator != "not-needed" {
		t.Fatalf("validation=%+v", record.Validation)
	}
}

func TestReviewReuseCountsValidationDismissalsAsDismissed(t *testing.T) {
	fixture, runner := validationReviewFixture(t, "Findings:\n- missing.go:2: outside\n- review.txt:2: standing")
	code, _, stderr := fixture.run(t)
	if code != exitRunFailed {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	code, _, stderr = runReviewDispose(t, "F2", "--dismiss", "--evidence", "checked supported path")
	if code != exitOK {
		t.Fatalf("dispose=%d stderr=%q", code, stderr)
	}
	resetReviewCommandRunner(runner, agent.ExecuteResult{})
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || !record.Reused || len(record.Dispositions) != 1 {
		t.Fatalf("reuse exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	assertNoReviewAgentActivity(t, runner)
	if strings.Contains(stderr, "no disposition") {
		t.Fatalf("stderr=%q", stderr)
	}
}

func TestReviewDisposeRefusesAFindingDismissedByValidation(t *testing.T) {
	fixture, _ := validationReviewFixture(t, "Findings:\n- missing.go:2: outside\n- review.txt:2: standing")
	if code, _, stderr := fixture.run(t); code != exitRunFailed {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if code, _, stderr := runReviewDispose(t, "F2", "--dismiss", "--evidence", "checked"); code != exitOK {
		t.Fatalf("dispose=%d stderr=%q", code, stderr)
	}
	path := filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "ignored")
	if code != exitPreflight || stdout != "" || stderr != "roundfix: review dispose refused: finding \"F1\" was dismissed by validation (unanchored)\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(before) {
		t.Fatalf("ledger changed err=%v", err)
	}
}

func TestReviewFindingFailureClauseRequiresContent(t *testing.T) {
	for _, test := range []struct {
		text string
		want bool
	}{
		{"Failure: breaks", true}, {"fAiLuRe: breaks", true}, {"Failure: \n\t", false}, {"no clause", false},
	} {
		t.Run(test.text, func(t *testing.T) {
			if got := reviewFindingHasFailureClause(test.text); got != test.want {
				t.Fatalf("got=%v", got)
			}
		})
	}
}

func TestReviewMissingAnchorIsRecordedWhenAnotherAnchorIsReadable(t *testing.T) {
	fixture, _ := validationReviewFixture(t, "Findings:\n- no anchor\n- review.txt:2: standing")
	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || record.FindingItems[0].Anchor != nil || record.FindingItems[0].Validation.Reason != "finding has no path:line anchor" {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	if stderr != "roundfix: review finding F1 dismissed by validation (unanchored): finding has no path:line anchor\n" {
		t.Fatalf("stderr=%q", stderr)
	}
}

func TestReviewDiffIndexDoesNotReadHunkContentAsFileHeaders(t *testing.T) {
	index := indexReviewDiff("diff --git a/review.txt b/review.txt\n--- a/review.txt\n+++ b/review.txt\n@@ -1,2 +1,2 @@\n--- missing.go\n+++ missing.go\n context\n")
	if !index.holds(reviewFindingAnchor{"review.txt", 2, 2}) || index.holds(reviewFindingAnchor{"missing.go", 2, 2}) {
		t.Fatalf("index=%+v", index)
	}
}

func TestReviewDiffIndexReadsGitQuotedPaths(t *testing.T) {
	index := indexReviewDiff(`diff --git "a/caf\303\251.go" "b/caf\303\251.go"
Binary files "a/caf\303\251.go" and "b/caf\303\251.go" differ
`)
	if !index.holds(reviewFindingAnchor{"café.go", 1, 1}) {
		t.Fatalf("index=%+v", index)
	}
}

func TestReviewReuseListsOnlyStandingFindingsAsMissing(t *testing.T) {
	fixture, runner := validationReviewFixture(t, "Findings:\n- missing.go:2: outside\n- review.txt:2: standing")
	if code, _, stderr := fixture.run(t); code != exitRunFailed {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	resetReviewCommandRunner(runner, agent.ExecuteResult{})
	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || !record.Reused || record.Outcome != reviewOutcomeFindings {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	if strings.Contains(stderr, "F1 has no disposition") || !strings.Contains(stderr, "F2 has no disposition") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoReviewAgentActivity(t, runner)
}

func TestReviewDiffIndexReadsAQuotedPathEndingInBackslash(t *testing.T) {
	index := indexReviewDiff(`diff --git "a/name\\" "b/name\\"
Binary files "a/name\\" and "b/name\\" differ
`)
	if !index.holds(reviewFindingAnchor{"name\\", 1, 1}) {
		t.Fatalf("index=%+v", index)
	}
}
