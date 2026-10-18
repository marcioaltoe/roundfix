// Boundary: real Git objects in disposable repositories; no provider or network.
package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/secretscan"
)

type secretRewriteFixture struct {
	review                               reviewCommandFixture
	base, candidate, head, secret, clean string
}

func newSecretRewriteFixture(t *testing.T) *secretRewriteFixture {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.PersistIdentity(t, repo)
	f := &secretRewriteFixture{review: reviewCommandFixture{repository: repo}}
	f.base = commitReviewCorrectionFile(t, f.review, "keep.txt", "unchanged\n")
	f.secret = "first\nSERVICE_" + "TOKEN=" + strings.Repeat("a1", 12) + "\nlast\n"
	f.clean = "first\nSERVICE_" + "TOKEN=" + strings.Repeat("a", 24) + "\nlast\n"
	if matches := secretscan.ScanContent("named.txt", []byte(f.secret)); len(matches) != 1 || matches[0].Line != 2 {
		t.Fatal("fixture must hold one Secret Match on line 2")
	}
	commitReviewCorrectionFile(t, f.review, "named.txt", f.secret)
	f.candidate = commitReviewCorrectionFile(t, f.review, "archive.txt", "archive preserved\n")
	// Recreate the first secret commit from its parent and replay the archive.
	gittest.Run(t, repo, "checkout", "-b", "fix/secret-rewrite", f.base)
	commitReviewCorrectionFile(t, f.review, "named.txt", f.clean)
	f.head = commitReviewCorrectionFile(t, f.review, "archive.txt", "archive preserved\n")
	return f
}

func (f *secretRewriteFixture) prove(t *testing.T, named ...string) delivery.SecretRewrite {
	t.Helper()
	workflow := &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}}
	proof, err := workflow.ProveSecretRewrite(t.Context(), f.review.repository, f.base, f.candidate, f.head, named)
	if err != nil {
		t.Fatal(err)
	}
	return proof
}

func TestSecretRewriteAcceptsReplacedMatchedLines(t *testing.T) {
	t.Parallel()
	f := newSecretRewriteFixture(t)
	// Local edits and the index are not evidence; only the two committed trees are.
	mustWrite(t, filepath.Join(f.review.repository, "named.txt"), "local change\n")
	gittest.Run(t, f.review.repository, "add", "named.txt")
	before := gittest.Run(t, f.review.repository, "status", "--porcelain=v1")
	proof := f.prove(t, "named.txt")
	if !proof.Accepted || proof.Reason != "" {
		t.Fatalf("proof=%+v", proof)
	}
	after := gittest.Run(t, f.review.repository, "status", "--porcelain=v1")
	if after != before {
		t.Fatal("proof changed the index or worktree")
	}
}

func TestSecretRewriteRefusesAnyOtherChange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		change func(*testing.T, *secretRewriteFixture)
		reason string
	}{
		{"Run start", func(t *testing.T, f *secretRewriteFixture) { f.base = f.candidate }, "item head does not descend from the Run start "},
		{"unnamed path", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "keep.txt", "other change\n")
		}, "the rewrite changes keep.txt, which the blocker does not name"},
		{"added path", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "added.txt", "new\n")
		}, "the rewrite adds, removes or changes the mode of added.txt"},
		{"removed path", func(t *testing.T, f *secretRewriteFixture) {
			gittest.Run(t, f.review.repository, "rm", "named.txt")
			gittest.Run(t, f.review.repository, "commit", "-m", "test: remove named path")
			f.head = strings.TrimSpace(gittest.Run(t, f.review.repository, "rev-parse", "HEAD"))
		}, "the rewrite adds, removes or changes the mode of named.txt"},
		{"mode change", func(t *testing.T, f *secretRewriteFixture) {
			// Update the index mode so this case works on Windows too.
			gittest.Run(t, f.review.repository, "update-index", "--chmod=+x", "named.txt")
			gittest.Run(t, f.review.repository, "commit", "-m", "test: change mode")
			f.head = strings.TrimSpace(gittest.Run(t, f.review.repository, "rev-parse", "HEAD"))
		}, "the rewrite adds, removes or changes the mode of named.txt"},
		{"line count", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "named.txt", f.clean+"extra\n")
		}, "the rewrite changes the line count of named.txt"},
		{"trailing newline", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "named.txt", strings.TrimSuffix(f.clean, "\n"))
		}, "the rewrite changes the line count of named.txt"},
		{"nonmatched line", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "named.txt", strings.Replace(f.clean, "first", "other", 1))
		}, "the rewrite changes line 1 of named.txt, which holds no Secret Match at the parked candidate"},
		{"retained match", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "named.txt", strings.ReplaceAll(f.secret, "a1", "b2"))
		}, "the rewrite keeps a Secret Match at line 2 of named.txt"},
		{"unchanged named file", func(t *testing.T, f *secretRewriteFixture) {
			f.head = commitReviewCorrectionFile(t, f.review, "named.txt", f.secret)
		}, "the rewrite keeps a Secret Match at line 2 of named.txt"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newSecretRewriteFixture(t)
			tc.change(t, f)
			named := []string{"named.txt"}
			if tc.name == "added path" {
				named = append(named, "added.txt")
			}
			proof := f.prove(t, named...)
			want := tc.reason
			if tc.name == "Run start" {
				want += f.base
			}
			if proof.Accepted || proof.Reason != want {
				t.Fatalf("proof=%+v want=%q", proof, want)
			}
			if strings.Contains(proof.Reason, strings.Repeat("a1", 12)) {
				t.Fatal("refusal leaked a value")
			}
		})
	}
}
