package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func TestSettleRefusesAnUnlistedExecutableInANewDirectory(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "backend", status: string(spec.StatusFailed), verification: []string{"true"}}})
	writeSettleExecutable(t, repo, "newdir/helper.sh")
	beforeHead := gitSettleOutput(t, repo, "rev-parse", "HEAD")
	beforeStatus := gitSettleOutput(t, repo, "status", "--porcelain=v1", "--untracked-files=all")
	beforeTask := mustRead(t, implementTaskPath(repo, "task_01"))
	var out, diag bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &out, &diag)
	refusal := "refuse Settle commit: executable file newdir/helper.sh (mode 0755) is not valid Work Item output unless the Spec's _authorization.md lists the path under executables"
	if code != exitPreflight || !strings.Contains(diag.String(), refusal) || out.Len() != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
	if got := gitSettleOutput(t, repo, "rev-parse", "HEAD"); got != beforeHead {
		t.Fatal("refusal moved HEAD")
	}
	if got := gitSettleOutput(t, repo, "status", "--porcelain=v1", "--untracked-files=all"); got != beforeStatus {
		t.Fatalf("refusal changed status: %s", got)
	}
	if got := mustRead(t, implementTaskPath(repo, "task_01")); got != beforeTask {
		t.Fatal("refusal changed Task")
	}
}

func TestSettleCommitsAListedExecutableWithItsMode(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "backend", status: string(spec.StatusFailed), verification: []string{"test -x scripts/launch.sh"}}})
	path := filepath.Join(repo, "docs", "specs", implementTestSlug, "_authorization.md")
	mustWrite(t, path, strings.Replace(mustRead(t, path), "---\n", "---\nexecutables:\n  - scripts/launch.sh\n", 1))
	gitImplement(t, repo, "add", path)
	gitImplement(t, repo, "commit", "-m", "authorize launcher")
	writeSettleExecutable(t, repo, "scripts/launch.sh")
	assertSettleCommitsExecutable(t, repo, "scripts/launch.sh")
}

func TestSettleCommitsATrackedFileMadeExecutable(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "backend", status: string(spec.StatusFailed), verification: []string{"test -x helper.sh"}}})
	mustWrite(t, filepath.Join(repo, "helper.sh"), "#!/bin/sh\nexit 0\n")
	gitImplement(t, repo, "add", "helper.sh")
	gitImplement(t, repo, "commit", "-m", "track script")
	writeSettleExecutable(t, repo, "helper.sh")
	assertSettleCommitsExecutable(t, repo, "helper.sh")
}

func writeSettleExecutable(t *testing.T, repo, path string) {
	t.Helper()
	gittest.AppendConfig(t, repo, "[core]\n\tfileMode = true\n")
	absolute := filepath.Join(repo, filepath.FromSlash(path))
	mustMkdir(t, filepath.Dir(absolute))
	mustWrite(t, absolute, "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(absolute, 0o755); err != nil {
		t.Fatal(err)
	}
}

func assertSettleCommitsExecutable(t *testing.T, repo, path string) {
	t.Helper()
	var out, diag bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &out, &diag)
	if code != exitOK {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
	if tree := gitSettleOutput(t, repo, "ls-tree", "HEAD", "--", path); !strings.HasPrefix(tree, "100755 blob ") || !strings.HasSuffix(tree, "\t"+path+"\n") {
		t.Fatalf("tree = %q", tree)
	}
	if task := mustRead(t, implementTaskPath(repo, "task_01")); !strings.Contains(task, "status: completed") {
		t.Fatalf("Task = %s", task)
	}
}
