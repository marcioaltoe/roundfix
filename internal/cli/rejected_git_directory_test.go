// Suite: CLI behavior under a rejected Git directory
// Invariant: rejected .git directories mark no project root and cause no Project Config write.
// Boundary IN: CLI streams and exit codes, skills root resolution, and temporary homes and repositories.
// Boundary OUT: Git entry validation shapes, covered by internal/config.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/gittest"
)

func TestRejectedGitDirectoryTranscripts(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		args   []string
		stdout string
		stderr string
		code   int
	}{
		{
			name:   "Surface Transcript 1",
			args:   []string{"runs", "list"},
			stdout: "No Runs found.\n",
			code:   0,
		},
		{
			name:   "Surface Transcript 2",
			args:   []string{"reconcile"},
			stderr: "roundfix: reconcile failed: reconcile requires running inside a Git repository\nRun 'roundfix reconcile --help' for usage.\n",
			code:   2,
		},
		{
			name:   "Surface Transcript 3",
			args:   []string{"init", "--scope", "project"},
			stderr: "Init failed\n\nReason:\n  project init requires a Git root; use --scope user outside a repository\n\nUsage:\n  Run 'roundfix init --help' for usage.\n",
			code:   2,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			homeDir := t.TempDir()
			drive := filepath.Join(homeDir, "drive")
			mustMkdir(t, filepath.Join(homeDir, ".git", "gk"))
			mustMkdir(t, filepath.Join(homeDir, ".git", "info"))
			mustMkdir(t, drive)
			setCommandEnvironmentForTest(t, homeDir, drive)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, test.args, &stdout, &stderr)
			if code != test.code || stdout.String() != test.stdout || stderr.String() != test.stderr {
				t.Errorf("transcript: exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q stderr=%q",
					code, stdout.String(), stderr.String(), test.code, test.stdout, test.stderr)
			}
			if _, err := os.Stat(filepath.Join(homeDir, ".roundfixrc.yml")); !os.IsNotExist(err) {
				t.Fatalf("Project Config must not exist: stat error=%v", err)
			}
		})
	}
}

func TestSkillsProjectRootSkipsARejectedGitDirectory(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	drive := filepath.Join(homeDir, "drive")
	mustMkdir(t, filepath.Join(homeDir, ".git", "gk"))
	mustMkdir(t, filepath.Join(homeDir, ".git", "info"))
	mustMkdir(t, drive)

	root, err := defaultResolveSkillsProjectRoot(context.Background(), drive)
	const refusal = "project skill install requires running inside a git repository or passing --dir"
	if root != "" || err == nil || err.Error() != refusal {
		t.Fatalf("rejected .git: root=%q error=%v; want empty root and error=%q", root, err, refusal)
	}

	repoDir := filepath.Join(drive, "repo")
	gittest.InitRepo(t, repoDir)
	subDir := filepath.Join(repoDir, "sub")
	mustMkdir(t, subDir)
	root, err = defaultResolveSkillsProjectRoot(context.Background(), subDir)
	if err != nil || root != repoDir {
		t.Fatalf("nested repository: root=%q error=%v; want root=%q and no error", root, err, repoDir)
	}
}
