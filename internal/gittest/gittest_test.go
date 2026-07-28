package gittest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitRepoDisablesAutomaticMaintenance(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "worktree", args: []string{"-b", "main"}},
		{name: "bare", args: []string{"--bare"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			repoDir := filepath.Join(t.TempDir(), "repo")
			InitRepo(t, repoDir, testCase.args...)

			want := map[string]string{
				"maintenance.auto":       "false",
				"gc.auto":                "0",
				"maintenance.autoDetach": "false",
			}
			for key, value := range want {
				got := strings.TrimSpace(Run(t, repoDir, "config", "--local", "--get", key))
				if got != value {
					t.Fatalf("%s = %q, want %q", key, got, value)
				}
			}
		})
	}
}

func TestInitRepoKeepsCommitFromSpawningDetachedMaintenance(t *testing.T) {
	repoDir := filepath.Join(t.TempDir(), "repo")
	InitRepo(t, repoDir, "-b", "main")
	if err := os.WriteFile(filepath.Join(repoDir, "tracked.txt"), []byte("seed\n"), 0o644); err != nil {
		t.Fatalf("write tracked file: %v", err)
	}
	Run(t, repoDir, "add", "tracked.txt")

	tracePath := filepath.Join(t.TempDir(), "trace2.json")
	args := append([]string{"-C", repoDir}, ConfigArgs()...)
	cmd := exec.Command("git", append(args, "commit", "-q", "-m", "seed")...)
	cmd.Env = append(IsolatedEnv(), "GIT_TRACE2_EVENT="+tracePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, output)
	}

	trace, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatalf("read Trace2 events: %v", err)
	}
	if strings.Contains(string(trace), "maintenance") {
		t.Fatalf("commit started automatic maintenance:\n%s", trace)
	}
}
