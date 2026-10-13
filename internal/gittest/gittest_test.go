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

func TestHardenedRepositoryNeverStartsAFilesystemMonitor(t *testing.T) {
	t.Parallel()
	global := filepath.Join(t.TempDir(), "global.gitconfig")
	if err := os.WriteFile(global, []byte("[core]\n fsmonitor = true\n untrackedCache = true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, bootstrap := range []string{"InitRepo", "Harden", "bare", "linked-worktree"} {
		t.Run(bootstrap, func(t *testing.T) {
			repository := filepath.Join(t.TempDir(), "repo")
			switch bootstrap {
			case "InitRepo":
				InitRepo(t, repository, "-b", "main")
			case "Harden":
				Run(t, "", "init", "-b", "main", repository)
				Harden(t, repository)
			case "bare":
				InitRepo(t, repository, "--bare")
			case "linked-worktree":
				base := filepath.Join(t.TempDir(), "base")
				InitRepo(t, base, "-b", "main")
				Run(t, base, "commit", "--allow-empty", "-m", "seed")
				Run(t, base, "worktree", "add", "-b", "linked", repository)
				for _, key := range []string{"core.fsmonitor", "core.untrackedCache"} {
					Run(t, repository, "config", "--unset-all", key)
				}
				Harden(t, repository)
			}
			// Read without ConfigArgs: its per-command override would hide whether
			// Harden wrote the local settings. Production Git gets this same config.
			tracePath := filepath.Join(t.TempDir(), "git-trace.json")
			git := func(args ...string) string {
				t.Helper()
				command := exec.Command("git", append([]string{"-C", repository}, args...)...)
				environment := IsolatedEnv()
				for index, entry := range environment {
					if strings.HasPrefix(entry, "GIT_CONFIG_GLOBAL=") {
						environment[index] = "GIT_CONFIG_GLOBAL=" + global
					}
				}
				command.Env = append(environment, "GIT_TRACE2_EVENT="+tracePath)
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v\n%s", args, err, output)
				}
				return string(output)
			}
			for _, key := range []string{"core.fsmonitor", "core.untrackedCache"} {
				if got := strings.TrimSpace(git("config", "--get", key)); got != "false" {
					t.Fatalf("%s = %q, want false", key, got)
				}
				if got := strings.TrimSpace(git("config", "--local", "--get", key)); got != "false" {
					t.Fatalf("local %s = %q, want false", key, got)
				}
			}
			if bootstrap == "bare" {
				return
			}
			identityArguments := []string{"-c", "user.name=Roundfix Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false"}
			commit := func(message string) { git(append(identityArguments, "commit", "-m", message)...) }
			file := filepath.Join(repository, "tracked.txt")
			if err := os.WriteFile(file, []byte("before\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			git("add", "--all")
			commit("before")
			if err := os.WriteFile(file, []byte("after!\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			git("add", "--all")
			commit("same size rewrite")
			if got := git("show", "HEAD:tracked.txt"); got != "after!\n" {
				t.Fatalf("committed content = %q", got)
			}
			trace, err := os.ReadFile(tracePath)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(trace), "fsmonitor--daemon") {
				t.Fatalf("hardened repository launched a filesystem monitor:\n%s", trace)
			}
		})
	}
}
