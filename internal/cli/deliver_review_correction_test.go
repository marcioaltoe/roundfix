// Boundary: real Git history and persisted review artifacts, all in disposable directories.
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

type reviewCorrectionFixture struct {
	workflow        *commandDeliveryWorkflow
	review          reviewCommandFixture
	record          reviewRecord
	candidate, head string
	dispositions    []reviewFindingDisposition
}

const correctionArchivePath = "docs/history/specs/0225-example/_prd.md"

func commitReviewCorrectionFile(t *testing.T, review reviewCommandFixture, name, content string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(review.repository, name)), 0755); err != nil {
		t.Fatal(err)
	}
	return commitLineageFile(t, review, name, content)
}

func newReviewCorrectionFixture(t *testing.T) *reviewCorrectionFixture {
	t.Helper()
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.PersistIdentity(t, repo)
	review := reviewCommandFixture{repository: repo, artifactDir: filepath.Join(root, "artifacts"), homeDir: filepath.Join(root, "home")}
	base := commitReviewCorrectionFile(t, review, "docs/specs/.keep", "Spec root\n")
	candidate := commitReviewCorrectionFile(t, review, correctionArchivePath, "archived record\n")
	head := commitReviewCorrectionFile(t, review, correctionArchivePath, "corrected record\n")
	record := newReviewRecord(repo, base, candidate, reviewPolicyForLineage(), reviewOutcomeFindings)
	record.Findings = "- Fix archive record\n- Check archive evidence"
	record.FindingItems = splitReviewFindings(record.Findings)
	record.Lineage = &reviewLineage{Round: 1}
	record.ArchivedSpecs = []string{"0225-example"}
	config := roundconfig.Builtin()
	config.Defaults.ArtifactDir = review.artifactDir
	fixture := &reviewCorrectionFixture{workflow: &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{Config: config, HomeDir: review.homeDir, GitRoot: repo}}, review: review, record: record, candidate: candidate, head: head}
	fixture.dispositions = []reviewFindingDisposition{
		{Repository: repo, HeadCommit: candidate, Finding: record.FindingItems[0].ID, Text: record.FindingItems[0].Text, Disposition: "fixed", FixedBy: head},
		{Repository: repo, HeadCommit: candidate, Finding: record.FindingItems[1].ID, Text: record.FindingItems[1].Text, Disposition: "dismissed", Evidence: "archived QA already records the evidence"},
	}
	fixture.persist(t)
	return fixture
}
func (f *reviewCorrectionFixture) persist(t *testing.T) {
	t.Helper()
	var encoded bytes.Buffer
	if err := writeReviewRecord(&encoded, f.record); err != nil {
		t.Fatal(err)
	}
	checkout := reviewCheckoutDir(f.review.artifactDir, f.review.repository)
	if err := os.MkdirAll(checkout, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkout, reviewRecordFileName), encoded.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	ledger := filepath.Join(f.review.artifactDir, reviewDispositionLedgerFileName)
	if err := os.Remove(ledger); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	for _, d := range f.dispositions {
		if _, err := appendReviewFindingDisposition(ledger, d); err != nil {
			t.Fatal(err)
		}
	}
}
func (f *reviewCorrectionFixture) prove(t *testing.T) delivery.ReviewCorrection {
	t.Helper()
	result, err := f.workflow.ProveReviewCorrection(t.Context(), f.review.repository, []string{"0225-example"}, f.candidate, f.head)
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func assertReviewCorrectionRefused(t *testing.T, f *reviewCorrectionFixture, reason string) {
	t.Helper()
	result := f.prove(t)
	if result.Accepted || !strings.Contains(result.Reason, reason) {
		t.Fatalf("proof=%+v want refusal naming %q", result, reason)
	}
}
func TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord(t *testing.T) {
	t.Parallel()
	f := newReviewCorrectionFixture(t)
	if result := f.prove(t); !result.Accepted || result.Reason != "" {
		t.Fatalf("proof=%+v", result)
	}
	candidate := f.record
	candidate.HeadCommit = f.head
	plan, err := decideReviewLineage(t.Context(), &f.record, candidate, f.workflow.git)
	if err != nil || plan.Lineage.Round != 2 || plan.Lineage.PreviousHead != f.candidate {
		t.Fatalf("lineage=%+v err=%v", plan, err)
	}
}
func TestACorrectionOutsideTheArchivedSpecIsRefused(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"outside.txt", "docs/history/specs/0225-example-sibling/_prd.md", "docs/history/specs/another/_prd.md"} {
		t.Run(name, func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			f.head = commitReviewCorrectionFile(t, f.review, name, "outside\n")
			assertReviewCorrectionRefused(t, f, name)
		})
	}
}
func TestACorrectionWithAnUndisposedFindingIsRefused(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		change func(*testing.T, *reviewCorrectionFixture)
	}{
		{"missing", func(t *testing.T, f *reviewCorrectionFixture) { f.dispositions = f.dispositions[:1] }},
		{"duplicate", func(t *testing.T, f *reviewCorrectionFixture) {
			f.dispositions = append(f.dispositions, f.dispositions[0])
		}},
		{"empty evidence", func(t *testing.T, f *reviewCorrectionFixture) { f.dispositions[1].Evidence = " " }},
		{"wrong finding text", func(t *testing.T, f *reviewCorrectionFixture) { f.dispositions[0].Text = "another finding" }},
		{"ledger another head", func(t *testing.T, f *reviewCorrectionFixture) { f.dispositions[0].HeadCommit = f.head }},
		{"fix before candidate", func(t *testing.T, f *reviewCorrectionFixture) { f.dispositions[0].FixedBy = f.record.BaseCommit }},
		{"fix absent from head", func(t *testing.T, f *reviewCorrectionFixture) {
			f.dispositions[0].FixedBy = commitReviewCorrectionFile(t, f.review, correctionArchivePath, "later fix\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			tc.change(t, f)
			f.persist(t)
			assertReviewCorrectionRefused(t, f, "finding")
		})
	}
}
func TestACorrectionThatDoesNotDescendIsRefused(t *testing.T) {
	t.Parallel()
	f := newReviewCorrectionFixture(t)
	f.head = f.record.BaseCommit
	assertReviewCorrectionRefused(t, f, "does not descend")
}
func TestReviewCorrectionRequiresTheParkedReviewRecord(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"head", "repository", "outcome"} {
		t.Run(name, func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			switch name {
			case "head":
				f.record.HeadCommit = f.head
			case "repository":
				f.record.Repository += "-other"
			case "outcome":
				f.record.Outcome = reviewOutcomeReviewed
				f.record.Findings = ""
				f.record.FindingItems = nil
			}
			f.persist(t)
			assertReviewCorrectionRefused(t, f, "review record")
		})
	}
}

func TestRetryAdmitsAFixedRoundTwoParkOutsideDeclaredPaths(t *testing.T) {
	t.Parallel()
	f := newReviewCorrectionFixture(t)
	f.record.Lineage = &reviewLineage{Round: 2}
	f.head = commitReviewCorrectionFile(t, f.review, "outside.txt", "structural fix outside declared paths\n")
	for i := range f.dispositions {
		f.dispositions[i].Disposition = "fixed"
		f.dispositions[i].FixedBy = f.head
		f.dispositions[i].Evidence = ""
	}
	f.persist(t)
	runner := newSessionReviewRunner("No findings.")
	writeReviewCommandConfig(t, f.review.repository, "codex", f.review.artifactDir)
	setCommandEnvironmentForTest(t, f.review.homeDir, f.review.repository)
	withAgentRunner(t, runner)
	if result := f.prove(t); !result.Accepted || result.Reason != "" {
		t.Fatalf("proof=%+v", result)
	}
	f.review.baseCommit = f.record.BaseCommit
	code, record, stderr := runLineageReview(t, f.review)
	if code != exitOK || !record.Lineage.Correction || len(runner.runs) != 1 || !strings.Contains(runner.runs[0].Prompt, "structural fix outside declared paths") {
		t.Fatalf("review code=%d record=%+v calls=%d stderr=%q", code, record, len(runner.runs), stderr)
	}
}

func TestRetryKeepsTheBoundForADismissedRoundTwoFinding(t *testing.T) {
	t.Parallel()
	for _, lineage := range []reviewLineage{{Round: 2}, {Round: 2, Correction: true}, {Round: 1}} {
		name := "dismissed round two"
		if lineage.Correction {
			name = "fixed Correction Review"
		} else if lineage.Round == 1 {
			name = "fixed round one"
		}
		t.Run(name, func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			f.record.Lineage = &lineage
			// Declare a different path to prove that the refusal still uses the
			// archived Tasks' exact bound, including when declarations are readable.
			commitReviewCorrectionFile(t, f.review, "docs/history/specs/0225-example/task_01.md", "---\ntask: task_01\nspec: 0225-example\nstatus: completed\ntype: backend\n---\n\n# Task\n\n## Context\n\n- interface: declared.txt\n")
			f.record.HeadCommit = commitReviewCorrectionFile(t, f.review, correctionArchivePath, "parked round two\n")
			f.candidate = f.record.HeadCommit
			f.head = commitReviewCorrectionFile(t, f.review, "outside.txt", "outside fix\n")
			for i := range f.dispositions {
				f.dispositions[i].HeadCommit = f.candidate
				if i == 0 || lineage.Correction || lineage.Round == 1 {
					f.dispositions[i].Disposition = "fixed"
					f.dispositions[i].FixedBy = f.head
					f.dispositions[i].Evidence = ""
				}
			}
			f.persist(t)
			runner := newSessionReviewRunner("No findings.")
			withAgentRunner(t, runner)
			result := f.prove(t)
			want := "changed path outside.txt is outside the named archived Specs and the paths their Tasks declared"
			if result.Accepted || result.Reason != want || len(runner.runs) != 0 {
				t.Fatalf("proof=%+v calls=%d want=%q", result, len(runner.runs), want)
			}
		})
	}
}
