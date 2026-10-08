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
