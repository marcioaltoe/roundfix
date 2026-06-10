package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

type engineFixture struct {
	store       *store.Store
	run         store.Run
	artifactDir string
	gitRoot     string
	issuePaths  []string
	calls       *[]string
	sink        *captureEventSink
	progress    *bytes.Buffer
	worktree    *engineFakeWorktree
}

// engineFakeWorktree returns scripted snapshots in call order, then keeps
// returning the last one.
type engineFakeWorktree struct {
	snapshots [][]string
	calls     int
}

func (worktree *engineFakeWorktree) Snapshot(context.Context, string) ([]string, error) {
	index := worktree.calls
	worktree.calls++
	if len(worktree.snapshots) == 0 {
		return nil, nil
	}
	if index >= len(worktree.snapshots) {
		index = len(worktree.snapshots) - 1
	}
	return worktree.snapshots[index], nil
}

type captureEventSink struct {
	mu     sync.Mutex
	events []runevent.RunEvent
}

func (sink *captureEventSink) Publish(_ context.Context, event runevent.RunEvent) error {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	sink.events = append(sink.events, event)
	return nil
}

func (sink *captureEventSink) kinds() []runevent.Kind {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	kinds := make([]runevent.Kind, 0, len(sink.events))
	for _, event := range sink.events {
		kinds = append(kinds, event.Kind)
	}
	return kinds
}

type engineFakeRunner struct {
	calls  *[]string
	status string
	err    error
	store  *store.Store
	seen   []string
}

func (runner *engineFakeRunner) Probe(context.Context, agent.RuntimeSpec) error { return nil }

func (runner *engineFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	*runner.calls = append(*runner.calls, "agent")
	runner.seen = append(runner.seen, runStateForTest(runner.store, req.RunID))
	if runner.err != nil {
		if agent.IsStopError(runner.err) && sink != nil {
			// Mirror the real runner: a stopped Agent publishes its status
			// event before returning.
			payload, _ := json.Marshal(struct {
				Status string `json:"status"`
			}{Status: "stopped"})
			_ = sink.Publish(ctx, runevent.RunEvent{
				RunID:   req.RunID,
				Batch:   req.Batch.Number,
				Source:  runevent.SourceAgent,
				Kind:    runevent.KindAgentStatus,
				Summary: "SESSION STOPPED\n",
				Payload: payload,
			})
		}
		return agent.ExecuteResult{}, runner.err
	}
	status := runner.status
	if status == "" {
		status = rounds.StatusResolved
	}
	for _, issue := range req.Batch.Issues {
		if err := rounds.SetIssueStatus(issue.Path, status, ""); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

type engineFakeVerifier struct {
	calls *[]string
	err   error
	store *store.Store
	runID string
	seen  []string
}

func (verifier *engineFakeVerifier) Verify(context.Context, VerifyRequest) error {
	*verifier.calls = append(*verifier.calls, "verify")
	verifier.seen = append(verifier.seen, runStateForTest(verifier.store, verifier.runID))
	return verifier.err
}

type engineFakeCommitter struct {
	calls    *[]string
	err      error
	messages []string
	paths    [][]string
}

func (committer *engineFakeCommitter) Commit(_ context.Context, req CommitRequest) error {
	*committer.calls = append(*committer.calls, "commit")
	committer.messages = append(committer.messages, req.Message)
	committer.paths = append(committer.paths, req.Paths)
	return committer.err
}

type engineFakePusher struct {
	calls   *[]string
	remotes []string
}

func (pusher *engineFakePusher) Push(_ context.Context, req PushRequest) error {
	*pusher.calls = append(*pusher.calls, "push")
	pusher.remotes = append(pusher.remotes, req.Remote+" HEAD:"+req.Branch)
	return nil
}

type engineFakeSource struct {
	calls    *[]string
	requests []reviewsource.ResolveRequest
}

func (source *engineFakeSource) ResolveIssues(_ context.Context, req reviewsource.ResolveRequest) error {
	*source.calls = append(*source.calls, "source")
	source.requests = append(source.requests, req)
	return nil
}

func runStateForTest(runStore *store.Store, runID string) string {
	run, found, err := runStore.Run(context.Background(), runID)
	if err != nil || !found {
		return "unknown"
	}
	return run.State
}

func newEngineFixture(t *testing.T) *engineFixture {
	t.Helper()
	ctx := context.Background()
	homeDir := t.TempDir()
	gitRoot := t.TempDir()
	artifactDir := filepath.Join(gitRoot, ".roundfix")

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = runStore.Close() })

	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        gitRoot,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    artifactDir,
	})
	if err != nil {
		t.Fatalf("create run: %v", err)
	}

	persisted, err := rounds.PersistRound(ctx, rounds.PersistRequest{
		ArtifactDir:    artifactDir,
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		HeadSHA:        "abc123",
		Round:          1,
		CreatedAt:      time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Items: []reviewsource.ReviewItem{
			{
				Title:                   "major: handle nil cache",
				File:                    "internal/cache/cache.go",
				Line:                    42,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "review body",
				SourceRef:               "thread:PRRT_1,comment:PRRC_1",
				ReviewHash:              "abc",
				SourceReviewID:          "9001",
				SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
			},
		},
	})
	if err != nil {
		t.Fatalf("persist round: %v", err)
	}

	calls := []string{}
	return &engineFixture{
		store:       runStore,
		run:         run,
		artifactDir: artifactDir,
		gitRoot:     gitRoot,
		issuePaths:  persisted.IssuePaths,
		calls:       &calls,
		sink:        &captureEventSink{},
		progress:    &bytes.Buffer{},
		worktree:    &engineFakeWorktree{snapshots: [][]string{nil, {"src/agent-change.go"}}},
	}
}

func (fixture *engineFixture) plan() CyclePlan {
	issues := make([]rounds.Issue, 0, len(fixture.issuePaths))
	for _, path := range fixture.issuePaths {
		issues = append(issues, rounds.Issue{Path: path})
	}
	return CyclePlan{
		RunID:        fixture.run.ID,
		GitRoot:      fixture.gitRoot,
		ArtifactDir:  fixture.artifactDir,
		SourceName:   reviewsource.SourceCodeRabbit,
		AgentName:    "codex",
		Runtime:      agent.RuntimeSpec{ID: "codex", DisplayName: "Codex"},
		Verification: "make verify",
		AutoCommit:   true,
		PullRequest: PullRequestRef{
			Number:         "123",
			BaseRepository: "owner/project",
			HeadRepository: "owner/project",
			HeadBranch:     "feature/review",
		},
		Batches:     []rounds.Batch{{Number: 1, Issues: issues}},
		TotalIssues: len(issues),
	}
}

func (fixture *engineFixture) engine(t *testing.T, runner agent.Runner, verifier Verifier, committer Committer, pusher Pusher, source ReviewSourceResolver) *Engine {
	t.Helper()
	engine, err := NewEngine(Dependencies{
		Runner:    runner,
		Verifier:  verifier,
		Committer: committer,
		Pusher:    pusher,
		Source:    source,
		Runs:      fixture.store,
		Worktree:  fixture.worktree,
		Sink:      fixture.sink,
		Progress:  fixture.progress,
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	return engine
}

func TestResolveCycleExecutesResolveVerifyCommitSourceContract(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	committer := &engineFakeCommitter{calls: fixture.calls}
	pusher := &engineFakePusher{calls: fixture.calls}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, pusher, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit>source" {
		t.Fatalf("expected resolve>verify>commit>source contract, got %q", got)
	}
	if len(result.Batches) != 1 {
		t.Fatalf("expected one Batch outcome, got %+v", result.Batches)
	}
	outcome := result.Batches[0]
	if !outcome.Committed || outcome.ResolvedSourceThreads != 1 || outcome.Issues != 1 {
		t.Fatalf("unexpected Batch outcome: %+v", outcome)
	}
	if result.Remaining != 0 {
		t.Fatalf("expected no remaining Unresolved Review Issues, got %d", result.Remaining)
	}
	if committer.messages[0] != BatchCommitMessage(1) {
		t.Fatalf("expected Batch commit message, got %q", committer.messages[0])
	}
	if source.requests[0].Issues[0].SourceRef != "thread:PRRT_1,comment:PRRC_1" {
		t.Fatalf("expected Source Reference forwarded, got %+v", source.requests[0])
	}
	// Intermediate Run states observed by collaborators during the cycle.
	if runner.seen[0] != store.StateResolvingWithAgent {
		t.Fatalf("expected ResolvingWithAgent during Agent run, got %q", runner.seen[0])
	}
	if verifier.seen[0] != store.StateVerifying {
		t.Fatalf("expected Verifying during verification, got %q", verifier.seen[0])
	}
	// Terminal completion stays caller-owned: the Run is still non-terminal.
	if state := runStateForTest(fixture.store, fixture.run.ID); store.IsTerminalState(state) {
		t.Fatalf("expected non-terminal Run after cycle, got %q", state)
	}
	// Final Push is a separate operation: the cycle never pushes.
	if len(pusher.remotes) != 0 {
		t.Fatalf("expected no push during cycle, got %v", pusher.remotes)
	}
}

func TestFinalPushIsASeparateExplicitOperation(t *testing.T) {
	fixture := newEngineFixture(t)
	pusher := &engineFakePusher{calls: fixture.calls}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, pusher, &engineFakeSource{calls: fixture.calls})

	err := engine.FinalPush(context.Background(), FinalPushRequest{
		RunID:   fixture.run.ID,
		WorkDir: fixture.gitRoot,
		Remote:  "origin",
		Branch:  "feature/review",
	})

	if err != nil {
		t.Fatalf("final push: %v", err)
	}
	if len(pusher.remotes) != 1 || pusher.remotes[0] != "origin HEAD:feature/review" {
		t.Fatalf("expected explicit push, got %v", pusher.remotes)
	}
	if state := runStateForTest(fixture.store, fixture.run.ID); state != store.StatePushing {
		t.Fatalf("expected Pushing state during Final Push, got %q", state)
	}
}

func TestResolveCycleAgentFailureFailsBatchAndRun(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store, err: errors.New("agent crashed")}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	_, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err == nil || !strings.Contains(err.Error(), "agent crashed") {
		t.Fatalf("expected Agent failure to fail the cycle, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent" {
		t.Fatalf("expected no daemon actions after Agent failure, got %q", got)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed Batch issue status, got %q", issue.Status)
	}
}

func TestResolveCycleVerificationFailureFailsBatchAndRun(t *testing.T) {
	fixture := newEngineFixture(t)
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID, err: errors.New("verification failed")}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, verifier, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	_, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err == nil || !strings.Contains(err.Error(), "verification failed") {
		t.Fatalf("expected verification failure to fail the cycle, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify" {
		t.Fatalf("expected no commit or source mutation after failed verification, got %q", got)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed Batch issue status, got %q", issue.Status)
	}
}

func TestResolveCycleStopBeforeBatchPublishesStopAndDoesNothing(t *testing.T) {
	fixture := newEngineFixture(t)
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := engine.ResolveCycle(ctx, fixture.plan())

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled cycle, got %v", err)
	}
	if len(*fixture.calls) != 0 {
		t.Fatalf("expected no daemon actions after Stop Request, got %v", *fixture.calls)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindDaemonStatus {
		t.Fatalf("expected daemon stop event published to the sink, got %v", kinds)
	}
}

func TestResolveCycleStopDuringAgentPreservesWorktreeAndHaltsDaemonActions(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store, err: agent.StopError{Err: context.Canceled}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	_, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if !agent.IsStopError(err) {
		t.Fatalf("expected StopError surfaced, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent" {
		t.Fatalf("expected no verification, commit, push, or source mutation after stop, got %q", got)
	}
	// Worktree preserved: the Batch is not marked failed and nothing was
	// committed.
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status == rounds.StatusFailed {
		t.Fatal("expected stopped Batch to preserve issue state, got failed")
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindAgentStatus {
		t.Fatalf("expected the runner's stop event to reach the sink, got %v", kinds)
	}
	for _, kind := range kinds {
		switch kind {
		case runevent.KindDaemonVerification, runevent.KindDaemonCommit, runevent.KindDaemonPush, runevent.KindDaemonSourceResolution:
			t.Fatalf("expected no unsafe daemon events after stop, got %v", kinds)
		}
	}
}

func TestNewEngineRequiresExplicitDependencies(t *testing.T) {
	_, err := NewEngine(Dependencies{})

	if err == nil {
		t.Fatal("expected missing dependencies error")
	}
	for _, expected := range []string{"Runner", "Verifier", "Committer", "Pusher", "Source", "Runs", "Worktree"} {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("expected %q in missing dependency error, got %v", expected, err)
		}
	}
}

type publishingFakeRunner struct {
	calls *[]string
}

func (runner *publishingFakeRunner) Probe(context.Context, agent.RuntimeSpec) error { return nil }

func (runner *publishingFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	*runner.calls = append(*runner.calls, "agent")
	// Mirror the real runner: a critical sink failure surfaces from
	// publication and fails the run.
	if err := sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Summary: "output",
		Payload: []byte(`{"text":"output"}`),
	}); err != nil {
		return agent.ExecuteResult{}, fmt.Errorf("publish Run Events: %w", err)
	}
	for _, issue := range req.Batch.Issues {
		if err := rounds.SetIssueStatus(issue.Path, rounds.StatusResolved, ""); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{}, nil
}

type failingCriticalSink struct{}

func (failingCriticalSink) Publish(context.Context, runevent.RunEvent) error {
	return errors.New("journal append failed")
}

func TestResolveCycleFailsRunWhenCriticalJournalSinkFails(t *testing.T) {
	fixture := newEngineFixture(t)
	fanout := runevent.NewFanout([]runevent.Sink{failingCriticalSink{}}, nil)
	defer fanout.Close()
	engine, err := NewEngine(Dependencies{
		Runner:    &publishingFakeRunner{calls: fixture.calls},
		Verifier:  &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID},
		Committer: &engineFakeCommitter{calls: fixture.calls},
		Pusher:    &engineFakePusher{calls: fixture.calls},
		Source:    &engineFakeSource{calls: fixture.calls},
		Runs:      fixture.store,
		Worktree:  fixture.worktree,
		Sink:      fanout,
		Progress:  fixture.progress,
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	_, cycleErr := engine.ResolveCycle(context.Background(), fixture.plan())

	if cycleErr == nil || !strings.Contains(cycleErr.Error(), "journal append failed") {
		t.Fatalf("expected journal append failure to fail the Run, got %v", cycleErr)
	}
	// Publication is non-optional: the selection event fails before any
	// daemon action, so nothing beyond (at most) the Agent ran.
	if got := strings.Join(*fixture.calls, ">"); strings.Contains(got, "verify") || strings.Contains(got, "commit") || strings.Contains(got, "source") {
		t.Fatalf("expected cycle halted after publish failure, got %q", got)
	}
}

func TestResolveCycleStagesOnlyAgentTouchedPaths(t *testing.T) {
	fixture := newEngineFixture(t)
	issuePath := fixture.issuePaths[0]
	// user-wip.txt is dirty before the Batch starts — pre-existing work or
	// a mid-Run user edit — and must never reach the Batch commit.
	fixture.worktree.snapshots = [][]string{
		{"user-wip.txt"},
		{"user-wip.txt", "src/fixed.go", issuePath},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if len(committer.paths) != 1 {
		t.Fatalf("expected one commit, got %v", committer.paths)
	}
	staged := strings.Join(committer.paths[0], "|")
	if staged != issuePath+"|src/fixed.go" {
		t.Fatalf("expected only Agent-touched paths staged (issue file + code), got %q", staged)
	}
	if strings.Contains(staged, "user-wip.txt") {
		t.Fatal("expected pre-existing user change kept out of the Batch commit")
	}
	if !result.Batches[0].Committed || result.Batches[0].CommitSkipped {
		t.Fatalf("expected committed outcome, got %+v", result.Batches[0])
	}
}

func TestResolveCycleSkipsCommitForTriageOnlyBatch(t *testing.T) {
	fixture := newEngineFixture(t)
	// Identical snapshots: the Agent triaged without changing the worktree.
	fixture.worktree.snapshots = [][]string{{"user-wip.txt"}, {"user-wip.txt"}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store, status: rounds.StatusInvalid}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected triage-only Batch to succeed, got %v", err)
	}
	if len(committer.paths) != 0 {
		t.Fatalf("expected no commit for triage-only Batch, got %v", committer.paths)
	}
	outcome := result.Batches[0]
	if outcome.Committed || !outcome.CommitSkipped {
		t.Fatalf("expected commit-skip outcome, got %+v", outcome)
	}
	if !strings.Contains(fixture.progress.String(), "Batch commit skipped: Batch 001 made no worktree changes.") {
		t.Fatalf("expected skip surfaced in command output, got %q", fixture.progress.String())
	}
}
