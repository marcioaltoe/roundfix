// Boundary: real Git history and public review/dispose commands in temporary
// symlinked worktrees. The reviewer is a local fake; no remote is configured.
package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

const declaredCorrectionSlug = "0225-example"

func declaredCorrectionTask(context string, recorded bool) string {
	body := "---\ntask: task_01\nstatus: completed\ntype: backend\n---\n\n# Task 01: Review correction\n\n## Context\n\n" + context
	if recorded {
		body += "\n## Recorded paths\n\n- `review.txt`\n"
	}
	return body
}

func newDeclaredCorrectionFixture(t *testing.T, form, taskContent string, missingSource bool) (reviewCommandFixture, string) {
	t.Helper()
	fixture, _, _ := newSymlinkedReviewCheckoutFixture(t, true)
	folder := "docs/history/specs/" + declaredCorrectionSlug
	if form == "archive-record" {
		source := "docs/specs/" + declaredCorrectionSlug
		commitReviewCorrectionFile(t, fixture, source+"/_prd.md", "# Example\n")
		revision := commitReviewCorrectionFile(t, fixture, source+"/task_01.md", taskContent)
		if missingSource {
			revision = strings.Repeat("a", 40)
		}
		gittest.Run(t, fixture.repository, "rm", "-r", "--", source, folder)
		recordPath := writeReaderRecord(t, fixture.repository, spec.ArchiveRecord{
			Spec: declaredCorrectionSlug, Title: "Review correction", Source: source,
			SourceRevision: revision, Disposition: spec.ArchivePass, PullRequest: "123",
		})
		gittest.Run(t, fixture.repository, "add", "--", recordPath)
		gittest.Run(t, fixture.repository, "commit", "-m", "docs: archive example")
	} else {
		commitReviewCorrectionFile(t, fixture, folder+"/task_01.md", taskContent)
	}
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	return fixture, folder + "/task_01.md"
}

func TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview(t *testing.T) {
	t.Parallel()
	for _, form := range []string{"archive-record", "legacy-folder"} {
		t.Run(form, func(t *testing.T) {
			t.Parallel()
			fixture, _ := newDeclaredCorrectionFixture(t, form, declaredCorrectionTask("- interface: review.txt\n", false), false)
			record := reviewSymlinkedCheckout(t, fixture)
			head := commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
			disposeSymlinkedReviewFinding(t, fixture, head)
			result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
			if !result.Accepted || result.Reason != "" {
				t.Fatalf("proof=%+v, want accepted declared code correction", result)
			}
		})
	}
}

func TestARecordedPathCountsAsDeclared(t *testing.T) {
	t.Parallel()
	for _, form := range []string{"archive-record", "legacy-folder"} {
		t.Run(form, func(t *testing.T) {
			t.Parallel()
			fixture, _ := newDeclaredCorrectionFixture(t, form, declaredCorrectionTask("", true), false)
			record := reviewSymlinkedCheckout(t, fixture)
			head := commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
			disposeSymlinkedReviewFinding(t, fixture, head)
			result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
			if !result.Accepted || result.Reason != "" {
				t.Fatalf("proof=%+v, want accepted recorded path", result)
			}
		})
	}
}

func TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath(t *testing.T) {
	t.Parallel()
	fixture, _ := newDeclaredCorrectionFixture(t, "archive-record", declaredCorrectionTask("- interface: review.txt\n", false), false)
	record := reviewSymlinkedCheckout(t, fixture)
	commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
	head := commitReviewCorrectionFile(t, fixture, "other.txt", "undeclared code\n")
	disposeSymlinkedReviewFinding(t, fixture, head)
	result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
	want := "changed path other.txt is outside the named archived Specs and the paths their Tasks declared"
	if result.Accepted || result.Reason != want {
		t.Fatalf("proof=%+v, want refusal %q", result, want)
	}
}

func TestDeclaredPathsAreReadAtTheRecordedRevision(t *testing.T) {
	t.Parallel()
	for _, form := range []string{"archive-record", "legacy-folder"} {
		t.Run(form, func(t *testing.T) {
			t.Parallel()
			for _, state := range []string{"correction-commit", "working-tree"} {
				t.Run(state, func(t *testing.T) {
					t.Parallel()
					fixture, taskPath := newDeclaredCorrectionFixture(t, form, declaredCorrectionTask("- interface: review.txt\n", false), false)
					record := reviewSymlinkedCheckout(t, fixture)
					widened := declaredCorrectionTask("- interface: review.txt\n- interface: other.txt\n", false)
					if state == "correction-commit" {
						commitReviewCorrectionFile(t, fixture, taskPath, widened)
					} else {
						mustMkdir(t, filepath.Dir(filepath.Join(fixture.repository, taskPath)))
						mustWrite(t, filepath.Join(fixture.repository, taskPath), widened)
						if form == "archive-record" {
							sourceTask := filepath.Join(fixture.repository, "docs/specs/", declaredCorrectionSlug, "task_01.md")
							mustMkdir(t, filepath.Dir(sourceTask))
							mustWrite(t, sourceTask, widened)
						}
					}
					commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
					head := commitReviewCorrectionFile(t, fixture, "other.txt", "undeclared code\n")
					disposeSymlinkedReviewFinding(t, fixture, head)
					result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
					want := "changed path other.txt is outside the named archived Specs and the paths their Tasks declared"
					if result.Accepted || result.Reason != want {
						t.Fatalf("proof=%+v, want original declared set and refusal %q", result, want)
					}
				})
			}
		})
	}
}

func TestUnreadableDeclaredPathsRefuseTheCorrection(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, form, task string
		missing          bool
	}{
		{"missing-source-revision", "archive-record", declaredCorrectionTask("- interface: review.txt\n", false), true},
		{"unparsable-record-task", "archive-record", "# Missing Task frontmatter\n", false},
		{"unparsable-folder-task", "legacy-folder", "# Missing Task frontmatter\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture, _ := newDeclaredCorrectionFixture(t, tc.form, tc.task, tc.missing)
			if remotes := strings.TrimSpace(gittest.Run(t, fixture.repository, "remote")); remotes != "" {
				t.Fatalf("unexpected remotes %q", remotes)
			}
			record := reviewSymlinkedCheckout(t, fixture)
			head := commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
			disposeSymlinkedReviewFinding(t, fixture, head)
			result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
			prefix := "declared paths of archived Spec " + declaredCorrectionSlug + " are unavailable: "
			if result.Accepted || !strings.HasPrefix(result.Reason, prefix) || len(result.Reason) == len(prefix) {
				t.Fatalf("proof=%+v, want unavailable refusal with a cause", result)
			}
		})
	}
}

func TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared(t *testing.T) {
	t.Parallel()
	for _, form := range []string{"archive-record", "legacy-folder"} {
		t.Run(form, func(t *testing.T) {
			t.Parallel()
			task := declaredCorrectionTask("", true) + "\n## Carry-forward provenance\n\n- Source Run: synthetic-run\n- Source commit: abc123\n"
			fixture, _ := newDeclaredCorrectionFixture(t, form, task, false)
			record := reviewSymlinkedCheckout(t, fixture)
			head := commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
			disposeSymlinkedReviewFinding(t, fixture, head)
			result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
			if !result.Accepted || result.Reason != "" {
				t.Fatalf("proof=%+v, want accepted carried recorded path", result)
			}
			head = commitReviewCorrectionFile(t, fixture, "other.txt", "undeclared code\n")
			result = proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
			want := "changed path other.txt is outside the named archived Specs and the paths their Tasks declared"
			if result.Accepted || result.Reason != want {
				t.Fatalf("proof=%+v, want refusal %q", result, want)
			}
		})
	}
}
