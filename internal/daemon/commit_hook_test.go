package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func TestClassifyCommitHookRefusal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		output  string
		hook    string
		refused bool
	}{
		{
			name:    "husky pre-commit banner names the hook",
			output:  "src/parser.go:12: function exceeds the 80-line limit\nhusky - pre-commit script failed (code 1)",
			hook:    "pre-commit",
			refused: true,
		},
		{
			name:    "lefthook banner names the hook",
			output:  "lefthook: pre-commit: lint failed",
			hook:    "pre-commit",
			refused: true,
		},
		{
			name:    "commit-msg refusal names its own hook",
			output:  "commit-msg hook refused: message does not follow Conventional Commits",
			hook:    "commit-msg",
			refused: true,
		},
		{
			name:    "prepare-commit-msg wins over its commit-msg suffix",
			output:  "prepare-commit-msg hook failed",
			hook:    "prepare-commit-msg",
			refused: true,
		},
		{
			name:    "a hook git could not spawn is still a refusal",
			output:  "error: cannot spawn .git/hooks/pre-merge-commit: Permission denied",
			hook:    "pre-merge-commit",
			refused: true,
		},
		{
			name:    "runner banner without a hook name still classifies",
			output:  "lefthook ❯ lint ❯ failed",
			refused: true,
		},
		{
			name:   "nothing staged is a plain git failure",
			output: "On branch main\nnothing to commit, working tree clean",
		},
		{
			name:   "empty message is a plain git failure",
			output: "Aborting commit due to empty commit message.",
		},
		{
			name:   "unmerged files are a plain git failure",
			output: "error: Committing is not possible because you have unmerged files.",
		},
		{
			name:   "signing failure is a plain git failure",
			output: "error: gpg failed to sign the data\nfatal: failed to write commit object",
		},
		{
			name:   "no output classifies nothing",
			output: "   \n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			hook, refused := ClassifyCommitHookRefusal(test.output)
			if refused != test.refused {
				t.Fatalf("ClassifyCommitHookRefusal(%q) refused = %t, want %t", test.output, refused, test.refused)
			}
			if hook != test.hook {
				t.Fatalf("ClassifyCommitHookRefusal(%q) hook = %q, want %q", test.output, hook, test.hook)
			}
		})
	}
}

func TestHookRefusalErrorNamesTheHookAndUnwraps(t *testing.T) {
	t.Parallel()
	cause := errors.New("exit status 1")
	refusal := &HookRefusalError{Hook: "pre-commit", ExitCode: 1, Output: "line too long", Err: cause}

	if !errors.Is(refusal, cause) {
		t.Fatalf("expected the git failure to stay unwrappable, got %v", refusal)
	}
	if got := refusal.Error(); !strings.Contains(got, "pre-commit") || !strings.Contains(got, "exit status 1") || !strings.Contains(got, "line too long") {
		t.Fatalf("expected hook, exit status and output in the message, got %q", got)
	}
	unnamed := &HookRefusalError{ExitCode: 1}
	if got := unnamed.HookName(); got != "commit" {
		t.Fatalf("expected an unnamed hook to fall back to %q, got %q", "commit", got)
	}
}

func TestHookOutputExcerptBoundsOutputAtItsHead(t *testing.T) {
	t.Parallel()
	if got := hookOutputExcerpt("  finding\n", hookRefusalOutputMax); got != "finding" {
		t.Fatalf("expected trimmed output, got %q", got)
	}
	if got := hookOutputExcerpt("   ", hookRefusalOutputMax); got != "" {
		t.Fatalf("expected blank output to excerpt to empty, got %q", got)
	}
	long := "objection: " + strings.Repeat("x", hookRefusalOutputMax*2)
	got := hookOutputExcerpt(long, hookRefusalOutputMax)
	if !strings.HasPrefix(got, "objection: ") {
		t.Fatalf("expected the hook's objection kept at the head, got %q", got)
	}
	if !strings.HasSuffix(got, " ...") {
		t.Fatalf("expected the cut marked, got %q", got)
	}
	if len(got) > hookRefusalOutputMax+len(" ...") {
		t.Fatalf("expected the excerpt bounded to %d bytes, got %d", hookRefusalOutputMax, len(got))
	}
}

// TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged drives the real
// committer against a real repository whose pre-commit hook refuses, which is
// the shape of every measured Run death: verified work, a refused commit.
func TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged(t *testing.T) {
	t.Parallel()
	repoDir := newHookRepoForTest(t)
	writeCommitHookForTest(t, repoDir, "pre-commit",
		"echo 'src/parser.go:12: function exceeds the 80-line limit' >&2\n"+
			"echo 'husky - pre-commit script failed (code 1)' >&2\n"+
			"exit 1\n")
	mustWriteForTest(t, filepath.Join(repoDir, "src.txt"), "verified work\n")

	err := GitCommitter{}.Commit(context.Background(), CommitRequest{
		WorkDir: repoDir,
		Message: "feat: add the parser",
		Paths:   []string{"src.txt"},
	})

	var refusal *HookRefusalError
	if !errors.As(err, &refusal) {
		t.Fatalf("expected a classified hook refusal, got %v", err)
	}
	if refusal.Hook != "pre-commit" {
		t.Fatalf("expected the refusing hook named, got %q", refusal.Hook)
	}
	if refusal.ExitCode != 1 {
		t.Fatalf("expected the hook exit code recorded, got %d", refusal.ExitCode)
	}
	if !strings.Contains(refusal.Output, "function exceeds the 80-line limit") {
		t.Fatalf("expected the hook's objection recorded, got %q", refusal.Output)
	}
	staged := runGitForTest(t, repoDir, "diff", "--cached", "--name-only")
	if !strings.Contains(staged, "src.txt") {
		t.Fatalf("expected the verified work left staged for recovery, got %q", staged)
	}
	if got := strings.TrimSpace(runGitForTest(t, repoDir, "rev-list", "--count", "HEAD")); got != "1" {
		t.Fatalf("expected no commit beyond the base, got %q commits", got)
	}
}

// TestGitCommitterKeepsPlainGitFailureUnclassified holds the other half of
// the contract: a refusal git itself raised is not a hook refusal, so it
// keeps ending the Run as an infrastructure failure.
func TestGitCommitterKeepsPlainGitFailureUnclassified(t *testing.T) {
	t.Parallel()
	repoDir := newHookRepoForTest(t)

	err := GitCommitter{}.Commit(context.Background(), CommitRequest{
		WorkDir: repoDir,
		Message: "feat: commit nothing",
		Paths:   []string{"base.txt"},
	})

	if err == nil {
		t.Fatal("expected git to refuse a commit with nothing staged")
	}
	var refusal *HookRefusalError
	if errors.As(err, &refusal) {
		t.Fatalf("expected a plain git failure to stay unclassified, got %+v", refusal)
	}
}

// TestTaskCycleHookRefusalKeepsTaskCompletedAndNamesRecovery is the Daemon
// boundary: the Task settled completed before the commit, so a hook refusal
// records itself and names the recovery instead of failing the Task.
func TestTaskCycleHookRefusalKeepsTaskCompletedAndNamesRecovery(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Add the parser"}})
	fixture.worktree.snapshots = [][]string{nil, {"src/parser.go"}}
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	committer := &hookRefusingCommitter{refusal: &HookRefusalError{
		Hook:     "pre-commit",
		ExitCode: 1,
		Output:   "src/parser.go:12: function exceeds the 80-line limit\nhusky - pre-commit script failed (code 1)",
		Err:      errors.New("git commit failed: exit status 1"),
	}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	var refusal *HookRefusalError
	if !errors.As(err, &refusal) {
		t.Fatalf("expected the hook refusal to reach the caller, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 0 {
		t.Fatalf("expected no settled Task counted after the refusal, got %+v", result)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected the verified Task to stay completed, got %q", got)
	}
	if got := committer.attempts(); got != 1 {
		t.Fatalf("expected exactly one commit attempt, got %d", got)
	}

	event := hookRefusedEvent(t, fixture)
	if event.ReviewIssue != "task_01" {
		t.Fatalf("expected the Task id on the refusal event, got %q", event.ReviewIssue)
	}
	if event.RunID != fixture.run.ID {
		t.Fatalf("expected the Run id on the refusal event, got %q", event.RunID)
	}
	payload := eventPayloadMap(t, event)
	if payload["decision"] != "refused" {
		t.Fatalf("expected a refused commit decision, got %v", payload["decision"])
	}
	if payload["hook"] != "pre-commit" {
		t.Fatalf("expected the refusing hook recorded, got %v", payload["hook"])
	}
	if payload["status"] != string(spec.StatusCompleted) {
		t.Fatalf("expected the Task recorded completed, got %v", payload["status"])
	}
	if got, ok := payload["exit_code"].(float64); !ok || int(got) != 1 {
		t.Fatalf("expected the hook exit code recorded, got %v", payload["exit_code"])
	}
	output, _ := payload["output"].(string)
	if !strings.Contains(output, "function exceeds the 80-line limit") {
		t.Fatalf("expected the hook output recorded, got %q", output)
	}
	recovery := fmt.Sprintf("roundfix settle --spec %s --task task_01", taskCycleSlug)
	if payload["recovery"] != recovery {
		t.Fatalf("expected the recovery command recorded, got %v", payload["recovery"])
	}
	progress := fixture.progress.String()
	if !strings.Contains(progress, recovery) {
		t.Fatalf("expected the recovery command reported in the Run's output, got %q", progress)
	}
	if !strings.Contains(progress, "function exceeds the 80-line limit") {
		t.Fatalf("expected what the hook objected to reported, got %q", progress)
	}
}

// hookRefusedEvent finds the single commit event classified hook_refused.
// The classification is matched as the literal the Run Event carries, not
// through the constant, so renaming the constant cannot silently change the
// vocabulary a Supervisor's journal is read with.
func hookRefusedEvent(t *testing.T, fixture *taskCycleFixture) runevent.RunEvent {
	t.Helper()
	var matched []runevent.RunEvent
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit) {
		if eventPayloadString(t, event, "classification") == "hook_refused" {
			matched = append(matched, event)
		}
	}
	if len(matched) != 1 {
		t.Fatalf("expected one commit event classified hook_refused, got %d", len(matched))
	}
	return matched[0]
}

type hookRefusingCommitter struct {
	mu       sync.Mutex
	calls    int
	refusal  *HookRefusalError
	requests []CommitRequest
}

func (committer *hookRefusingCommitter) Commit(_ context.Context, req CommitRequest) error {
	committer.mu.Lock()
	defer committer.mu.Unlock()
	committer.calls++
	committer.requests = append(committer.requests, req)
	return committer.refusal
}

func (committer *hookRefusingCommitter) attempts() int {
	committer.mu.Lock()
	defer committer.mu.Unlock()
	return committer.calls
}

// initHookRepoForTest turns an existing directory into a repository with a
// repository-local hooks path, so a global core.hooksPath cannot decide
// whether the test's hook runs.
func initHookRepoForTest(t *testing.T, repoDir string) {
	t.Helper()
	gittest.InitRepo(t, repoDir, "-b", "main")
	hooksDir := filepath.Join(repoDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatalf("create hooks directory: %v", err)
	}
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = test@example.com\n[commit]\n\tgpgsign = false\n[core]\n\thooksPath = "+hooksDir+"\n")
}

// newHookRepoForTest builds a disposable repository with one base commit and
// a repository-local hooks path.
func newHookRepoForTest(t *testing.T) string {
	t.Helper()
	repoDir := t.TempDir()
	initHookRepoForTest(t, repoDir)
	mustWriteForTest(t, filepath.Join(repoDir, "base.txt"), "base\n")
	runGitForTest(t, repoDir, "add", "base.txt")
	runGitForTest(t, repoDir, "commit", "-q", "-m", "init")
	return repoDir
}

func writeCommitHookForTest(t *testing.T, repoDir string, hook string, body string) {
	t.Helper()
	path := filepath.Join(repoDir, ".git", "hooks", hook)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatalf("write %s hook: %v", hook, err)
	}
}
