// Suite: auditing-binary age evidence.
// Invariant: evidence answers only where it is meaningful, and never reports current from absence.
// Boundary IN: a real Git repository and its declared Roundfix version.
// Boundary OUT: the staleness comparison itself, owned by internal/app.
package spec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/app"
)

func initEvidenceRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{
		{"init", "--initial-branch=main"},
		{"config", "user.email", "evidence@example.test"},
		{"config", "user.name", "Evidence"},
		// A developer with global commit signing on would otherwise have this
		// fixture prompt or fail; the repository disables it in every other
		// Git fixture for the same reason.
		{"config", "commit.gpgsign", "false"},
	} {
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	return root
}

func commitEvidenceFile(t *testing.T, root, name, body string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-m", "commit " + name}} {
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	command := exec.Command("git", "-C", root, "rev-parse", "HEAD")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("rev-parse HEAD: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func TestResolveAuditorEvidence(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("a build commit behind the audited tree is older", func(t *testing.T) {
		t.Parallel()
		root := initEvidenceRepo(t)
		first := commitEvidenceFile(t, root, "a.txt", "one\n")
		commitEvidenceFile(t, root, "b.txt", "two\n")

		evidence := ResolveAuditorEvidence(ctx, root, app.AuditingBinary{Commit: first})

		if evidence.Ancestry != app.AncestryOlder {
			t.Fatalf("ancestry = %v, want %v", evidence.Ancestry, app.AncestryOlder)
		}
	})

	t.Run("a build commit at the audited tree is not older", func(t *testing.T) {
		t.Parallel()
		root := initEvidenceRepo(t)
		head := commitEvidenceFile(t, root, "a.txt", "one\n")

		evidence := ResolveAuditorEvidence(ctx, root, app.AuditingBinary{Commit: head})

		if evidence.Ancestry != app.AncestryNotOlder {
			t.Fatalf("ancestry = %v, want %v", evidence.Ancestry, app.AncestryNotOlder)
		}
	})

	t.Run("the dirty marker still resolves its base commit", func(t *testing.T) {
		t.Parallel()
		root := initEvidenceRepo(t)
		first := commitEvidenceFile(t, root, "a.txt", "one\n")
		commitEvidenceFile(t, root, "b.txt", "two\n")

		evidence := ResolveAuditorEvidence(ctx, root, app.AuditingBinary{Commit: first + "-dirty"})

		if evidence.Ancestry != app.AncestryOlder {
			t.Fatalf("dirty-marked ancestry = %v, want %v", evidence.Ancestry, app.AncestryOlder)
		}
	})

	t.Run("a commit from another repository answers nothing", func(t *testing.T) {
		t.Parallel()
		// The ordinary fleet case: Roundfix audits an adopting repository, so
		// its build commit is not an object there. Comparing it against history
		// it never belonged to would be worse than reporting unknown.
		root := initEvidenceRepo(t)
		commitEvidenceFile(t, root, "a.txt", "one\n")

		evidence := ResolveAuditorEvidence(ctx, root, app.AuditingBinary{Commit: "0123456789abcdef0123456789abcdef01234567"})

		if evidence.Ancestry != app.AncestryUnknown {
			t.Fatalf("foreign-commit ancestry = %v, want %v", evidence.Ancestry, app.AncestryUnknown)
		}
		if evidence.TreeVersion != "" {
			t.Fatalf("foreign-commit tree version = %q, want empty", evidence.TreeVersion)
		}
	})

	t.Run("a released build falls back to the declared tree version", func(t *testing.T) {
		t.Parallel()
		root := initEvidenceRepo(t)
		manifest := filepath.Join(root, roundfixVersionManifest)
		if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
			t.Fatalf("create manifest directory: %v", err)
		}
		if err := os.WriteFile(manifest, []byte(`{"version":"9.9.9"}`), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}
		commitEvidenceFile(t, root, "a.txt", "one\n")

		evidence := ResolveAuditorEvidence(ctx, root, app.AuditingBinary{Version: "0.1.0"})

		if evidence.Ancestry != app.AncestryUnknown {
			t.Fatalf("released-build ancestry = %v, want %v", evidence.Ancestry, app.AncestryUnknown)
		}
		if evidence.TreeVersion != "9.9.9" {
			t.Fatalf("declared tree version = %q, want 9.9.9", evidence.TreeVersion)
		}
	})

	t.Run("a repository declaring no Roundfix version answers nothing", func(t *testing.T) {
		t.Parallel()
		root := initEvidenceRepo(t)
		commitEvidenceFile(t, root, "a.txt", "one\n")

		evidence := ResolveAuditorEvidence(ctx, root, app.AuditingBinary{Version: "0.1.0"})

		if evidence.Ancestry != app.AncestryUnknown || evidence.TreeVersion != "" {
			t.Fatalf("evidence = %+v, want no signal at all", evidence)
		}
	})
}
