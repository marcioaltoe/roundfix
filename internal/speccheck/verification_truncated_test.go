package speccheck_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"roundfix/internal/speccheck"
)

func checkTruncatedBullet(t *testing.T, status, bullet string) []speccheck.Finding {
	t.Helper()
	root := t.TempDir()
	const slug = "truncated"
	prefix := "docs/specs/" + slug + "/"
	writeCitationFixtureFile(t, root, prefix+"_prd.md", "---\nstatus: active\n---\n# Fixture\n")
	writeCitationFixtureFile(t, root, prefix+"_tasks.md", "---\nschema: spec-tasks/v1\nspec: truncated\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n---\n")
	writeCitationFixtureFile(t, root, prefix+"task_01.md", fmt.Sprintf("---\ntask: task_01\nspec: truncated\nstatus: %s\ntype: backend\ncomplexity: low\n---\n\n## Verification\n\n%s\n", status, bullet))
	result, err := speccheck.Check(filepath.Join(root, "docs/specs"), root, slug)
	if err != nil {
		t.Fatal(err)
	}
	return findingsWithCode(result, speccheck.CodeVerifyTruncated)
}

func TestVerificationSpanEndingInBackslashIsTruncated(t *testing.T) {
	t.Parallel()
	bullet := "- `grep -qF '| \\`' file`"
	for _, status := range []string{"pending", "in_progress", "failed"} {
		t.Run(status, func(t *testing.T) {
			findings := checkTruncatedBullet(t, status, bullet+"\n"+bullet)
			if len(findings) != 2 {
				t.Fatalf("findings = %+v", findings)
			}
			for i, f := range findings {
				if f.Severity != speccheck.SeverityError || f.Where[0].Path != "docs/specs/truncated/task_01.md" || f.Where[0].Line != 11+i {
					t.Fatalf("finding = %+v", f)
				}
			}
		})
	}
	if findings := checkTruncatedBullet(t, "completed", bullet); len(findings) != 0 {
		t.Fatalf("completed findings = %+v", findings)
	}
}

func TestVerificationLineWithAStrayBacktickIsTruncated(t *testing.T) {
	t.Parallel()
	if findings := checkTruncatedBullet(t, "pending", "- `test -f effect` — stray `"); len(findings) != 1 {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCleanVerificationBulletIsNotTruncated(t *testing.T) {
	t.Parallel()
	for _, bullet := range []string{"- `test -f effect`", "- `test -f effect` — expected `exit 1`", "- prose without a code span"} {
		t.Run(bullet, func(t *testing.T) {
			// The prose-only bullet accompanies a valid command so Spec loading succeeds.
			findings := checkTruncatedBullet(t, "pending", "- `test -f other-effect`\n"+bullet)
			if len(findings) != 0 {
				t.Fatalf("findings = %+v", findings)
			}
		})
	}
}
