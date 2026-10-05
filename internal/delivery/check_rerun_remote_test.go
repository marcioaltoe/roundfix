package delivery

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitHubCLIInspectsAgainstTheDeliveryRemote(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.test/repo\n\ngo 1.26\n"), 0600); err != nil {
		t.Fatal(err)
	}
	script := newScriptedCommandRunner(t,
		commandStep{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":1}`}},
		commandStep{name: "gh", args: []string{"api", "repos/{owner}/{repo}/check-runs/7/annotations?per_page=100"}, result: CommandResult{Stdout: `[]`}},
		commandStep{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{Stdout: "FAIL\texample.test/repo/internal/other\t1s"}},
		commandStep{name: "git", args: []string{"symbolic-ref", "refs/remotes/upstream/HEAD"}, result: CommandResult{Stdout: "refs/remotes/upstream/trunk"}},
		commandStep{name: "git", args: []string{"fetch", "upstream", "+refs/heads/trunk:refs/remotes/upstream/trunk"}},
		commandStep{name: "git", args: []string{"merge-base", "refs/remotes/upstream/trunk", "candidate"}, result: CommandResult{Stdout: "common"}},
		commandStep{name: "git", args: []string{"diff", "--name-only", "--no-renames", "-z", "common", "candidate"}, result: CommandResult{Stdout: "internal/change/file.go\x00"}},
	)
	client := GitHubCLI{Runner: recoveryCommandRunner{dir, script}}
	failure, err := client.InspectFailedCheck(t.Context(), dir, "upstream", "candidate", PullRequestCheck{Link: recoveryLink})
	if err != nil || !failure.OutsideChange {
		t.Fatalf("failure=%+v err=%v", failure, err)
	}
}

func TestGitHubCLIRefusesAnUnsafeDeliveryRemote(t *testing.T) {
	for _, remote := range []string{"", "-u", "a/b", "a b"} {
		script := newScriptedCommandRunner(t)
		_, err := (GitHubCLI{Runner: script}).InspectFailedCheck(t.Context(), "/repo", remote, "candidate", PullRequestCheck{Link: recoveryLink})
		if err == nil || !strings.Contains(err.Error(), "invalid delivery remote") {
			t.Fatalf("remote %q err=%v", remote, err)
		}
	}
}
