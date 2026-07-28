package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
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
	mu        sync.Mutex
	snapshots [][]string
	calls     int
}

func (worktree *engineFakeWorktree) Snapshot(context.Context, string) ([]string, error) {
	worktree.mu.Lock()
	defer worktree.mu.Unlock()
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
	mu        sync.Mutex
	events    []runevent.RunEvent
	published chan runevent.RunEvent
}

func (sink *captureEventSink) Publish(_ context.Context, event runevent.RunEvent) error {
	sink.mu.Lock()
	sink.events = append(sink.events, event)
	published := sink.published
	sink.mu.Unlock()
	if published != nil {
		published <- event
	}
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

func (sink *captureEventSink) snapshot() []runevent.RunEvent {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	events := make([]runevent.RunEvent, len(sink.events))
	copy(events, sink.events)
	return events
}

func eventPayloadString(t *testing.T, event runevent.RunEvent, key string) string {
	t.Helper()
	payload := eventPayloadMap(t, event)
	value, _ := payload[key].(string)
	return value
}

func eventPayloadMap(t *testing.T, event runevent.RunEvent) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("decode event payload %s: %v", event.Payload, err)
	}
	return payload
}

const modelNotAdvertisedReasonForTest = `Agent Model "gpt-5.6-sol" not advertised by runtime "codex"; advertised: gpt-5.5, gpt-5.1`

func modelNotAdvertisedBatchErrorForTest() error {
	return &agent.BatchFailureError{
		ExitCode: 1,
		Reason:   "agent/protocol error",
		Err: &agent.ModelNotAdvertisedError{
			Runtime:    "codex",
			Model:      "gpt-5.6-sol",
			Advertised: []string{"gpt-5.5", "gpt-5.1"},
		},
	}
}

func assertCommentContains(t *testing.T, actions []engineSourceAction, sourceRef string, expected ...string) {
	t.Helper()
	for _, action := range actions {
		if action.Kind != "comment" || action.SourceRef != sourceRef || !action.Posted {
			continue
		}
		for _, fragment := range expected {
			if !strings.Contains(action.Body, fragment) {
				t.Fatalf("expected comment for %s to contain %q, got %q", sourceRef, fragment, action.Body)
			}
		}
		return
	}
	t.Fatalf("expected posted comment for %s, got actions %+v", sourceRef, actions)
}

type engineFakeRunner struct {
	calls                     *[]string
	status                    string
	statusBySourceRef         map[string]string
	terminalReason            string
	terminalReasonBySourceRef map[string]string
	duplicateOfBySourceRef    map[string]string
	err                       error
	errByCall                 []error
	result                    agent.ExecuteResult
	store                     *store.Store
	seen                      []string
	requests                  []agent.ExecuteRequest
}

func (runner *engineFakeRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *engineFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	*runner.calls = append(*runner.calls, "agent")
	callIndex := len(runner.requests)
	runner.requests = append(runner.requests, req)
	runner.seen = append(runner.seen, runStateForTest(runner.store, req.RunID))
	runErr := runner.err
	if callIndex < len(runner.errByCall) {
		runErr = runner.errByCall[callIndex]
	}
	if runErr != nil {
		if agent.IsStopError(runErr) && sink != nil {
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
		return agent.ExecuteResult{}, runErr
	}
	status := runner.status
	if status == "" {
		status = rounds.StatusResolved
	}
	for _, issue := range req.Batch.Issues {
		current, err := rounds.ParseIssue(issue.Path)
		if err != nil {
			return agent.ExecuteResult{}, err
		}
		nextStatus := status
		if runner.statusBySourceRef != nil {
			if override := runner.statusBySourceRef[current.SourceRef]; override != "" {
				nextStatus = override
			}
		}
		reason := runner.terminalReason
		if runner.terminalReasonBySourceRef != nil {
			if override := runner.terminalReasonBySourceRef[current.SourceRef]; override != "" {
				reason = override
			}
		}
		duplicateOf := ""
		if runner.duplicateOfBySourceRef != nil {
			duplicateOf = runner.duplicateOfBySourceRef[current.SourceRef]
		}
		if err := rounds.SetIssueStatus(issue.Path, nextStatus, duplicateOf, reason); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	result := runner.result
	if result.LogPath == "" {
		result.LogPath = req.LogPath
	}
	return result, nil
}

func (runner *engineFakeRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type engineFakeVerifier struct {
	calls  *[]string
	err    error
	script []error
	// failFirst fails only the first verification, so tests can prove a
	// failed Batch does not stop later Batches from verifying cleanly.
	failFirst bool
	store     *store.Store
	runID     string
	seen      []string
}

func (verifier *engineFakeVerifier) Verify(_ context.Context, req VerifyRequest) (VerifyResult, error) {
	*verifier.calls = append(*verifier.calls, "verify")
	verifier.seen = append(verifier.seen, runStateForTest(verifier.store, verifier.runID))
	if len(verifier.script) > 0 {
		err := verifier.script[0]
		verifier.script = verifier.script[1:]
		if err != nil {
			return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: err}
		}
		return VerifyResult{OutputPath: req.OutputPath}, nil
	}
	if verifier.failFirst && len(verifier.seen) == 1 {
		return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("verification failed")}
	}
	if verifier.err != nil {
		return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: verifier.err}
	}
	return VerifyResult{OutputPath: req.OutputPath}, nil
}

type engineInfrastructureVerifier struct {
	calls *[]string
	err   error
}

func (verifier *engineInfrastructureVerifier) Verify(_ context.Context, req VerifyRequest) (VerifyResult, error) {
	*verifier.calls = append(*verifier.calls, "verify")
	return VerifyResult{OutputPath: req.OutputPath}, verifier.err
}

type engineStopAfterCommandFailureVerifier struct {
	calls *[]string
	store *store.Store
	runID string
}

func (verifier *engineStopAfterCommandFailureVerifier) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	*verifier.calls = append(*verifier.calls, "verify")
	if err := verifier.store.RequestStop(ctx, verifier.runID); err != nil {
		return VerifyResult{OutputPath: req.OutputPath}, err
	}
	return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{
		Command:    req.Command,
		OutputPath: req.OutputPath,
		Err:        errors.New("verification failed before stop"),
	}
}

type engineFakeCommitter struct {
	calls       *[]string
	err         error
	afterCommit func(context.Context, CommitRequest) error
	messages    []string
	paths       [][]string
}

func (committer *engineFakeCommitter) Commit(ctx context.Context, req CommitRequest) error {
	*committer.calls = append(*committer.calls, "commit")
	committer.messages = append(committer.messages, req.Message)
	committer.paths = append(committer.paths, req.Paths)
	if committer.err != nil {
		return committer.err
	}
	if committer.afterCommit != nil {
		return committer.afterCommit(ctx, req)
	}
	return nil
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
	calls             *[]string
	afterResolve      func(context.Context, reviewsource.ResolveRequest) error
	afterIssueResolve func(context.Context, reviewsource.IssueResolveRequest) error
	errByAction       map[string]error
	requests          []reviewsource.ResolveRequest
	resolveRequests   []reviewsource.IssueResolveRequest
	commentRequests   []reviewsource.IssueCommentRequest
	actions           []engineSourceAction
	postedMarkers     map[string]bool
}

func (source *engineFakeSource) ResolveIssues(ctx context.Context, req reviewsource.ResolveRequest) error {
	*source.calls = append(*source.calls, "source")
	source.requests = append(source.requests, req)
	if source.afterResolve != nil {
		return source.afterResolve(ctx, req)
	}
	return nil
}

func (source *engineFakeSource) ResolveIssue(ctx context.Context, req reviewsource.IssueResolveRequest) error {
	*source.calls = append(*source.calls, "source")
	source.resolveRequests = append(source.resolveRequests, req)
	source.actions = append(source.actions, engineSourceAction{Kind: "resolve", SourceRef: req.SourceRef, Resolved: true})
	if source.afterIssueResolve != nil {
		return source.afterIssueResolve(ctx, req)
	}
	if source.afterResolve != nil {
		return source.afterResolve(ctx, reviewsource.ResolveRequest{
			Source:         req.Source,
			PRNumber:       req.PRNumber,
			BaseRepository: req.BaseRepository,
			Issues: []reviewsource.ResolvedIssue{{
				SourceRef: req.SourceRef,
			}},
		})
	}
	return source.errFor("resolve", req.SourceRef)
}

func (source *engineFakeSource) ReplyToIssue(_ context.Context, req reviewsource.IssueCommentRequest) (reviewsource.IssueCommentResult, error) {
	*source.calls = append(*source.calls, "source")
	source.commentRequests = append(source.commentRequests, req)
	if err := source.errFor("comment", req.SourceRef); err != nil {
		source.actions = append(source.actions, engineSourceAction{Kind: "comment", SourceRef: req.SourceRef, Marker: req.Marker, Body: req.Body})
		return reviewsource.IssueCommentResult{}, err
	}
	if source.postedMarkers == nil {
		source.postedMarkers = map[string]bool{}
	}
	if source.postedMarkers[req.Marker] {
		source.actions = append(source.actions, engineSourceAction{Kind: "comment", SourceRef: req.SourceRef, Marker: req.Marker, Body: req.Body})
		return reviewsource.IssueCommentResult{Skipped: true}, nil
	}
	source.postedMarkers[req.Marker] = true
	source.actions = append(source.actions, engineSourceAction{Kind: "comment", SourceRef: req.SourceRef, Marker: req.Marker, Body: req.Body, Posted: true})
	return reviewsource.IssueCommentResult{Posted: true}, nil
}

func (source *engineFakeSource) errFor(kind string, sourceRef string) error {
	if source.errByAction == nil {
		return nil
	}
	if err := source.errByAction[kind+":"+sourceRef]; err != nil {
		return err
	}
	return source.errByAction[kind]
}

type engineSourceAction struct {
	Kind      string
	SourceRef string
	Marker    string
	Body      string
	Posted    bool
	Resolved  bool
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
	return newEngineFixtureWithItems(t, []reviewsource.ReviewItem{
		reviewItemForTest("major: handle nil cache", "internal/cache/cache.go", 42, "thread:PRRT_1,comment:PRRC_1", "abc", "9001"),
	})
}

func reviewItemForTest(title string, file string, line int, sourceRef string, reviewHash string, reviewID string) reviewsource.ReviewItem {
	return reviewsource.ReviewItem{
		Title:                   title,
		File:                    file,
		Line:                    line,
		Severity:                "major",
		Author:                  "coderabbitai[bot]",
		Body:                    "review body",
		SourceRef:               sourceRef,
		ReviewHash:              reviewHash,
		SourceReviewID:          reviewID,
		SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
	}
}

func newEngineFixtureWithItems(t *testing.T, items []reviewsource.ReviewItem) *engineFixture {
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
		Items:          items,
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
		Session:      agent.SessionRefForRun(fixture.run.ID, fixture.gitRoot),
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
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusResolved {
		t.Fatalf("expected resolved Review Issue, got %q", issue.Status)
	}
	if issue.TerminalReason != "" {
		t.Fatalf("expected resolved Review Issue without terminal reason, got %q", issue.TerminalReason)
	}
	if committer.messages[0] != BatchCommitMessage(1) {
		t.Fatalf("expected Batch commit message, got %q", committer.messages[0])
	}
	if source.resolveRequests[0].SourceRef != "thread:PRRT_1,comment:PRRC_1" {
		t.Fatalf("expected Source Reference forwarded, got %+v", source.resolveRequests[0])
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

func TestPerWorkAgentSessionReviewUsesReviewProfile(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &selectionLifecycleRunner{}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryReview: selectionProfileForTest(selectionForTest("claude", "review-model", "medium"), selectionForTest("codex", "review-fallback", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	result, err := engine.ResolveCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("ResolveCycle returned error: %v", err)
	}
	if result.Remaining != 0 {
		t.Fatalf("expected review cycle to resolve all issues, got %+v", result)
	}
	if got := lifecycleRequestSummary(runner.runRequests()); strings.Join(got, "\n") != "roundfix-"+fixture.run.ID+"-review-001|claude|review-model" {
		t.Fatalf("expected review profile session and selection, got %v", got)
	}
	if got := strings.Join(runner.closedSessions(), "\n"); got != "roundfix-"+fixture.run.ID+"-review-001" {
		t.Fatalf("expected review session closed once, got %q", got)
	}
}

func TestResolveCycleReportsReviewAgentSessionCloseFailure(t *testing.T) {
	fixture := newEngineFixture(t)
	closeErr := errors.New("close failed")
	runner := &selectionLifecycleRunner{closeErr: closeErr}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryReview: selectionProfileForTest(selectionForTest("codex", "review-model", "high"), selectionForTest("codex", "review-fallback", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	_, err := engine.ResolveCycle(context.Background(), plan)

	if !errors.Is(err, closeErr) {
		t.Fatalf("ResolveCycle error = %v, want close failure", err)
	}
	if !strings.Contains(err.Error(), "close Agent Session") || !strings.Contains(err.Error(), "Batch 001") {
		t.Fatalf("close failure lacks context: %v", err)
	}
	if got := strings.Join(runner.closedSessions(), "\n"); got != "roundfix-"+fixture.run.ID+"-review-001" {
		t.Fatalf("expected review session close attempted once, got %q", got)
	}
}

func TestResolveCyclePropagatesSettledIssueOutcomesIndividually(t *testing.T) {
	const (
		resolvedRef   = "thread:PRRT_resolved,comment:PRRC_resolved"
		invalidRef    = "thread:PRRT_invalid,comment:PRRC_invalid"
		duplicatedRef = "thread:PRRT_duplicate,comment:PRRC_duplicate"
		failedRef     = "thread:PRRT_failed,comment:PRRC_failed"
	)
	fixture := newEngineFixtureWithItems(t, []reviewsource.ReviewItem{
		reviewItemForTest("major: resolved", "internal/one.go", 1, resolvedRef, "hash-resolved", "9001"),
		reviewItemForTest("major: invalid", "internal/two.go", 2, invalidRef, "hash-invalid", "9002"),
		reviewItemForTest("major: duplicated", "internal/three.go", 3, duplicatedRef, "hash-duplicated", "9003"),
		reviewItemForTest("major: failed", "internal/four.go", 4, failedRef, "hash-failed", "9004"),
	})
	runner := &engineFakeRunner{
		calls: fixture.calls,
		store: fixture.store,
		statusBySourceRef: map[string]string{
			resolvedRef:   rounds.StatusResolved,
			invalidRef:    rounds.StatusInvalid,
			duplicatedRef: rounds.StatusDuplicated,
			failedRef:     rounds.StatusFailed,
		},
		terminalReasonBySourceRef: map[string]string{
			invalidRef: "invalid: generated file",
			failedRef:  "Verification failed: command make verify exit status 7; diagnostics: /tmp/roundfix.log",
		},
		duplicateOfBySourceRef: map[string]string{
			duplicatedRef: "thread:PRRT_resolved",
		},
	}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, runner, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if len(result.Batches) != 1 || result.Batches[0].ResolvedSourceThreads != 3 {
		t.Fatalf("expected three resolved source threads, got %+v", result.Batches)
	}
	gotPrefix := []string{}
	for index, action := range source.actions {
		if index == 6 {
			break
		}
		gotPrefix = append(gotPrefix, action.Kind+":"+action.SourceRef)
	}
	wantPrefix := []string{
		"resolve:" + resolvedRef,
		"comment:" + invalidRef,
		"resolve:" + invalidRef,
		"comment:" + duplicatedRef,
		"resolve:" + duplicatedRef,
		"comment:" + failedRef,
	}
	if strings.Join(gotPrefix, "|") != strings.Join(wantPrefix, "|") {
		t.Fatalf("expected per-status source action order %v, got %v", wantPrefix, gotPrefix)
	}
	for _, action := range source.actions {
		if action.SourceRef == failedRef && action.Kind == "resolve" {
			t.Fatalf("expected failed Review Issue to stay open, got source actions %+v", source.actions)
		}
	}
	assertCommentContains(t, source.actions, invalidRef, "roundfix:outcome", "invalid: generated file")
	assertCommentContains(t, source.actions, duplicatedRef, "roundfix:outcome", "Canonical Review Issue: thread:PRRT_resolved")
	assertCommentContains(t, source.actions, failedRef, "roundfix:outcome", "make verify", "exit status 7", "later Round retries")
	calls := strings.Join(*fixture.calls, ">")
	sourceIndex := strings.Index(calls, "source")
	if sourceIndex < 0 {
		t.Fatalf("expected source propagation call, got %q", calls)
	}
	if strings.Index(calls, "verify") > sourceIndex || strings.Index(calls, "commit") > sourceIndex {
		t.Fatalf("expected source propagation after verification and commit, got %q", calls)
	}
}

func TestResolveCycleOutcomeCommentsAreIdempotent(t *testing.T) {
	const sourceRef = "thread:PRRT_invalid,comment:PRRC_invalid"
	fixture := newEngineFixtureWithItems(t, []reviewsource.ReviewItem{
		reviewItemForTest("major: invalid", "internal/invalid.go", 5, sourceRef, "hash-invalid", "9001"),
	})
	if err := rounds.SetIssueStatus(fixture.issuePaths[0], rounds.StatusInvalid, "", "invalid: not actionable"); err != nil {
		t.Fatalf("set issue status: %v", err)
	}
	issue, err := rounds.ParseIssue(fixture.issuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)
	plan := fixture.plan()

	if _, err := engine.propagateSourceIssue(context.Background(), plan, 1, issue, "batch", nil); err != nil {
		t.Fatalf("first propagation: %v", err)
	}
	if _, err := engine.propagateSourceIssue(context.Background(), plan, 1, issue, "batch", nil); err != nil {
		t.Fatalf("second propagation: %v", err)
	}

	posted := 0
	for _, action := range source.actions {
		if action.Kind == "comment" && action.Posted {
			posted++
		}
	}
	if posted != 1 {
		t.Fatalf("expected one posted idempotent comment, got %d actions %+v", posted, source.actions)
	}
}

func TestResolveCycleSourcePropagationFailureContinues(t *testing.T) {
	const (
		firstRef  = "thread:PRRT_first,comment:PRRC_first"
		secondRef = "thread:PRRT_second,comment:PRRC_second"
	)
	fixture := newEngineFixtureWithItems(t, []reviewsource.ReviewItem{
		reviewItemForTest("major: first", "internal/first.go", 1, firstRef, "hash-first", "9001"),
		reviewItemForTest("major: second", "internal/second.go", 2, secondRef, "hash-second", "9002"),
	})
	runner := &engineFakeRunner{
		calls: fixture.calls,
		store: fixture.store,
		statusBySourceRef: map[string]string{
			firstRef:  rounds.StatusInvalid,
			secondRef: rounds.StatusInvalid,
		},
		terminalReasonBySourceRef: map[string]string{
			firstRef:  "invalid: first",
			secondRef: "invalid: second",
		},
	}
	source := &engineFakeSource{
		calls: fixture.calls,
		errByAction: map[string]error{
			"comment:" + firstRef: errors.New("review source unavailable"),
		},
	}
	engine := fixture.engine(t, runner, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle should continue after source propagation failure: %v", err)
	}
	if len(result.Batches) != 1 || !result.Batches[0].Failed || result.Batches[0].ResolvedSourceThreads != 1 || result.Remaining != 1 {
		t.Fatalf("expected retryable source propagation failure and one resolved source thread, got result=%+v batches=%+v", result, result.Batches)
	}
	if len(source.resolveRequests) != 1 || source.resolveRequests[0].SourceRef != secondRef {
		t.Fatalf("expected only the second issue resolved after first comment failed, got %+v", source.resolveRequests)
	}
	first, err := rounds.ParseIssue(fixture.issuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	if first.Status != rounds.StatusFailed || !strings.Contains(first.TerminalReason, "Review Source propagation failed during invalid") {
		t.Fatalf("expected first issue failed for retry, got status=%q reason=%q", first.Status, first.TerminalReason)
	}
	assertCommentContains(t, source.actions, secondRef, "invalid: second")
	if !strings.Contains(fixture.progress.String(), "Review Source propagation failed") {
		t.Fatalf("expected propagation failure reported, got %q", fixture.progress.String())
	}
	foundFailureEvent := false
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonSourceResolution) {
		if strings.Contains(eventPayloadString(t, event, "error"), "review source unavailable") {
			foundFailureEvent = true
		}
	}
	if !foundFailureEvent {
		t.Fatalf("expected propagation failure event, got %+v", fixture.sink.snapshot())
	}
}

func TestOutcomeCommentBodyUsesPublicDuplicateReference(t *testing.T) {
	canonicalPath := filepath.Join(t.TempDir(), "issue_001.md")
	content := `---
source: coderabbit
pr: "123"
round: 1
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: owner/project
head_branch: feature/review
source_ref: "thread:PRRT_public,comment:PRRC_public"
---

# Issue 001
`
	if err := os.WriteFile(canonicalPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write canonical issue: %v", err)
	}

	body := outcomeCommentBody("run_1", rounds.Issue{
		Status:      rounds.StatusDuplicated,
		DuplicateOf: canonicalPath,
	}, "duplicated", "<!-- marker -->")

	if !strings.Contains(body, "Canonical Review Issue: thread:PRRT_public,comment:PRRC_public") {
		t.Fatalf("expected public Source Reference in comment body, got %q", body)
	}
	if strings.Contains(body, canonicalPath) {
		t.Fatalf("comment body leaked local artifact path %q: %q", canonicalPath, body)
	}
}

func TestOutcomeCommentBodySanitizesPublicReasons(t *testing.T) {
	body := outcomeCommentBody("run_1", rounds.Issue{
		Status:         rounds.StatusFailed,
		TerminalReason: "Verification failed:\nsee /Users/alice/dev/roundfix/.roundfix/verification.log and ping @team for details",
	}, "failed", "<!-- marker -->")

	for _, forbidden := range []string{"/Users/alice", ".roundfix/verification.log", "@team"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("comment body leaked %q: %q", forbidden, body)
		}
	}
	for _, want := range []string{"Reason: Verification failed: see <path>", "＠team"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected sanitized comment body to contain %q, got %q", want, body)
		}
	}
}

func TestResolveCycleRunEndLeavesUnresolvedIssuesCommented(t *testing.T) {
	fixture := newEngineFixture(t)
	source := &engineFakeSource{calls: fixture.calls}
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store, err: errors.New("agent crashed")}
	engine := fixture.engine(t, runner, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if result.Remaining != 1 {
		t.Fatalf("expected failed issue to remain open for later Round, got %d", result.Remaining)
	}
	assertCommentContains(t, source.actions, "thread:PRRT_1,comment:PRRC_1", "Agent failed", "later Round retries")
	if len(source.resolveRequests) != 0 {
		t.Fatalf("expected failed issue to remain unresolved on Review Source, got %+v", source.resolveRequests)
	}
}

func TestResolveCycleJournalsTransportAnomalyBeforeVerification(t *testing.T) {
	fixture := newEngineFixture(t)
	const anomaly = "acpx exited with exit code 1 after parsed session/prompt result\n--- acpx stderr tail ---\nMessage buffer exceeded 10485760 bytes"
	runner := &engineFakeRunner{
		calls:  fixture.calls,
		store:  fixture.store,
		result: agent.ExecuteResult{TransportAnomaly: anomaly},
	}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if len(result.Batches) != 1 || result.Batches[0].Failed {
		t.Fatalf("expected anomaly Batch to proceed through verification, got %+v", result.Batches)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit>source" {
		t.Fatalf("expected anomaly to preserve resolve flow, got %q", got)
	}
	events := fixture.sink.snapshot()
	anomalyIndex := -1
	verificationIndex := -1
	for index, event := range events {
		if event.Kind == runevent.KindDaemonBatch && eventPayloadString(t, event, "anomaly") == anomaly {
			anomalyIndex = index
			if event.Source != runevent.SourceDaemon || event.Batch != 1 {
				t.Fatalf("expected daemon Batch anomaly event, got %+v", event)
			}
		}
		if event.Kind == runevent.KindDaemonVerification && strings.Contains(string(event.Payload), `"phase":"started"`) {
			verificationIndex = index
		}
	}
	if anomalyIndex < 0 {
		t.Fatalf("expected transport anomaly in Run Event Journal, got %+v", events)
	}
	if verificationIndex < 0 {
		t.Fatalf("expected verification start event, got %+v", events)
	}
	if anomalyIndex > verificationIndex {
		t.Fatalf("expected anomaly event before verification, anomaly index %d verification index %d", anomalyIndex, verificationIndex)
	}
}

func TestResolveCycleKeepsAgentProvidedTerminalReason(t *testing.T) {
	fixture := newEngineFixture(t)
	const reason = "invalid: reviewer asked for generated code"
	runner := &engineFakeRunner{
		calls:          fixture.calls,
		store:          fixture.store,
		status:         rounds.StatusInvalid,
		terminalReason: reason,
	}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if len(result.Batches) != 1 || result.Batches[0].Failed {
		t.Fatalf("expected invalid triage to complete the Batch, got %+v", result.Batches)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusInvalid {
		t.Fatalf("expected invalid Review Issue, got %q", issue.Status)
	}
	if issue.TerminalReason != reason {
		t.Fatalf("expected agent terminal reason %q, got %q", reason, issue.TerminalReason)
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

func TestResolveCycleAgentFailureFailsBatchAndContinues(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store, err: errors.New("agent crashed")}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected Agent failure to fail only the Batch, got cycle error %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); strings.Contains(got, "verify") || strings.Contains(got, "commit") {
		t.Fatalf("expected no verification or commit after Agent failure, got %q", got)
	}
	if len(source.commentRequests) == 0 {
		t.Fatal("expected failed Review Issue propagation comment after Agent failure")
	}
	if len(result.Batches) != 1 || !result.Batches[0].Failed {
		t.Fatalf("expected failed Batch outcome, got %+v", result.Batches)
	}
	if result.Batches[0].FailureReason != "Agent failed: agent crashed" {
		t.Fatalf("expected Agent failure reason, got %q", result.Batches[0].FailureReason)
	}
	if result.Remaining != 1 {
		t.Fatalf("expected failed issue to stay Unresolved, got remaining %d", result.Remaining)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed Batch issue status, got %q", issue.Status)
	}
	for _, expected := range []string{"Agent failed", "agent crashed"} {
		if !strings.Contains(issue.TerminalReason, expected) {
			t.Fatalf("expected terminal reason to contain %q, got %q", expected, issue.TerminalReason)
		}
	}
	if strings.Contains(issue.TerminalReason, "\n") {
		t.Fatalf("expected single-line terminal reason, got %q", issue.TerminalReason)
	}
}

func TestResolveCycleModelNotAdvertisedFailureSettlesReviewIssueReason(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store, err: modelNotAdvertisedBatchErrorForTest()}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected model rejection to fail only the Batch, got cycle error %v", err)
	}
	if len(result.Batches) != 1 || !result.Batches[0].Failed {
		t.Fatalf("expected failed Batch outcome, got %+v", result.Batches)
	}
	if result.Batches[0].FailureReason != modelNotAdvertisedReasonForTest {
		t.Fatalf("expected model rejection Batch reason %q, got %q", modelNotAdvertisedReasonForTest, result.Batches[0].FailureReason)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed Batch issue status, got %q", issue.Status)
	}
	if issue.TerminalReason != modelNotAdvertisedReasonForTest {
		t.Fatalf("expected model rejection terminal reason %q, got %q", modelNotAdvertisedReasonForTest, issue.TerminalReason)
	}
	assertCommentContains(t, source.actions, issue.SourceRef, modelNotAdvertisedReasonForTest)

	var journaled bool
	for _, event := range fixture.sink.snapshot() {
		if event.Kind != runevent.KindDaemonBatch || event.Batch != 1 {
			continue
		}
		payload := eventPayloadMap(t, event)
		if payload["phase"] == "failed" && payload["error"] == modelNotAdvertisedReasonForTest {
			journaled = true
		}
	}
	if !journaled {
		t.Fatalf("expected failed Batch event to journal reason %q", modelNotAdvertisedReasonForTest)
	}
}

func TestResolveCycleVerificationFailureFailsBatchAndContinues(t *testing.T) {
	fixture := newEngineFixture(t)
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID, err: errors.New("exit status 7")}
	committer := &engineFakeCommitter{calls: fixture.calls}
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected verification failure to fail only the Batch, got cycle error %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify>source>source" {
		t.Fatalf("expected one repair attempt, failed-thread comment, and run-end idempotency check after final failed verification, got %q", got)
	}
	if len(source.resolveRequests) != 0 {
		t.Fatalf("expected failed Review Issue to remain unresolved, got resolves %+v", source.resolveRequests)
	}
	if len(result.Batches) != 1 || !result.Batches[0].Failed {
		t.Fatalf("expected failed Batch outcome, got %+v", result.Batches)
	}
	if !strings.Contains(result.Batches[0].FailureReason, "exit status 7") {
		t.Fatalf("expected verification failure reason, got %q", result.Batches[0].FailureReason)
	}
	if result.Remaining != 1 {
		t.Fatalf("expected failed issue to stay Unresolved, got remaining %d", result.Remaining)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed Batch issue status, got %q", issue.Status)
	}
	finalOutputPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 2)
	for _, expected := range []string{"Verification failed", "make verify", "exit status 7", finalOutputPath} {
		if !strings.Contains(issue.TerminalReason, expected) {
			t.Fatalf("expected terminal reason to contain %q, got %q", expected, issue.TerminalReason)
		}
	}
	if strings.Contains(issue.TerminalReason, "\n") {
		t.Fatalf("expected single-line terminal reason, got %q", issue.TerminalReason)
	}
	if len(runner.requests) != 2 {
		t.Fatalf("expected initial Agent request plus one Verification Feedback request, got %d", len(runner.requests))
	}
	if !strings.Contains(runner.requests[1].Prompt, "Verification Feedback") {
		t.Fatalf("expected second Agent request to be Verification Feedback, got:\n%s", runner.requests[1].Prompt)
	}
}

func TestResolveCycleVerificationFailureRepairsSameSessionAndAvoidsDuplicateBatchStart(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID, failFirst: true}
	committer := &engineFakeCommitter{calls: fixture.calls}
	source := &engineFakeSource{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, source)
	plan := fixture.plan()

	result, err := engine.ResolveCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify>commit>source" {
		t.Fatalf("expected initial Agent, attempt 1, repair, final attempt, commit, source; got %q", got)
	}
	if len(result.Batches) != 1 || result.Batches[0].Failed || !result.Batches[0].Committed || result.Batches[0].ResolvedSourceThreads != 1 {
		t.Fatalf("expected repaired Batch to settle successfully, got %+v", result.Batches)
	}
	if len(runner.requests) != 2 {
		t.Fatalf("expected initial prompt plus one Verification Feedback prompt, got %d", len(runner.requests))
	}
	if runner.requests[0].Session != plan.Session || runner.requests[1].Session != plan.Session {
		t.Fatalf("expected repair to reuse SessionRef %#v, got %#v then %#v", plan.Session, runner.requests[0].Session, runner.requests[1].Session)
	}
	repairPrompt := runner.requests[1].Prompt
	expectedPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1)
	for _, expected := range []string{
		"Verification Feedback",
		"Work Item: Batch 001",
		"Failed command: make verify",
		"Diagnostic artifact: " + expectedPath,
		"verification failed",
	} {
		if !strings.Contains(repairPrompt, expected) {
			t.Fatalf("expected repair prompt to contain %q, got:\n%s", expected, repairPrompt)
		}
	}
	startedBatchEvents := 0
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonBatch) {
		if eventPayloadString(t, event, "phase") == "started" {
			startedBatchEvents++
		}
	}
	if startedBatchEvents != 1 {
		t.Fatalf("expected exactly one Batch-start boundary, got %d", startedBatchEvents)
	}
	verificationEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification)
	verdicts := []string{}
	for _, event := range verificationEvents {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == string(runevent.VerificationPhaseVerdict) {
			verdicts = append(verdicts, fmt.Sprintf("%v:%v", payload["attempt"], payload["verdict"]))
		}
	}
	if got := strings.Join(verdicts, "|"); got != "1:failed|2:passed" {
		t.Fatalf("expected failed attempt 1 and passed attempt 2 verdicts, got %s", got)
	}
}

func TestResolveCycleVerificationFailureRetainsDiagnosticsWithoutStreamingOutput(t *testing.T) {
	fixture := newEngineFixture(t)
	plan := fixture.plan()
	plan.Verification = `printf '\117\125\124\120\125\124\137\102\117\104\131'; exit 7`
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, ExecVerifier{}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("expected command failure to fail only the Batch, got %v", err)
	}
	if len(result.Batches) != 1 || !result.Batches[0].Failed {
		t.Fatalf("expected failed Batch outcome, got %+v", result.Batches)
	}
	outputPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1)
	finalOutputPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 2)
	content, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatalf("read verification artifact: %v", readErr)
	}
	if string(content) != "OUTPUT_BODY" {
		t.Fatalf("expected failed command output retained, got %q", string(content))
	}
	finalContent, readErr := os.ReadFile(finalOutputPath)
	if readErr != nil {
		t.Fatalf("read final verification artifact: %v", readErr)
	}
	if string(finalContent) != "OUTPUT_BODY" {
		t.Fatalf("expected final failed command output retained, got %q", string(finalContent))
	}
	progress := fixture.progress.String()
	if !strings.Contains(progress, "Verification failed (attempt 1); diagnostics: "+outputPath) {
		t.Fatalf("expected verdict summary with diagnostic path, got %q", progress)
	}
	if !strings.Contains(progress, "Verification failed (attempt 2); diagnostics: "+finalOutputPath) {
		t.Fatalf("expected final verdict summary with diagnostic path, got %q", progress)
	}
	if strings.Contains(progress, "OUTPUT_BODY") {
		t.Fatalf("expected progress to omit raw command output, got %q", progress)
	}
	verificationEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification)
	phases := []string{}
	verdicts := 0
	for _, event := range verificationEvents {
		if strings.Contains(event.Summary, "OUTPUT_BODY") || strings.Contains(string(event.Payload), "OUTPUT_BODY") {
			t.Fatalf("expected Run Events to omit output body, got summary=%q payload=%s", event.Summary, event.Payload)
		}
		payload := eventPayloadMap(t, event)
		phase, _ := payload["phase"].(string)
		phases = append(phases, phase)
		if phase == string(runevent.VerificationPhaseVerdict) {
			verdicts++
			if payload["verdict"] != string(runevent.VerificationVerdictFailed) {
				t.Fatalf("expected failed aggregate verdict, got %v", payload["verdict"])
			}
			attempt := int(payload["attempt"].(float64))
			expectedPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, attempt)
			if payload["diagnostic_path"] != expectedPath {
				t.Fatalf("expected verdict diagnostic path %q, got %v", expectedPath, payload["diagnostic_path"])
			}
		}
		if phase == string(runevent.VerificationPhaseFailed) {
			if _, ok := payload["attempt"].(float64); !ok {
				t.Fatalf("expected attempt metadata, got %s", event.Payload)
			}
			attempt := int(payload["attempt"].(float64))
			expectedPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, attempt)
			if payload["diagnostic_path"] != expectedPath {
				t.Fatalf("expected failed command diagnostic path %q, got %v", expectedPath, payload["diagnostic_path"])
			}
		}
	}
	if got := strings.Join(phases, "|"); got != "started|failed|verdict|started|failed|verdict" {
		t.Fatalf("expected command phases and one aggregate verdict per attempt, got %s", got)
	}
	if verdicts != 2 {
		t.Fatalf("expected exactly one aggregate verdict per attempt, got %d in %+v", verdicts, verificationEvents)
	}
}

func TestResolveCycleVerificationInfrastructureErrorHaltsWithoutFailedSettlement(t *testing.T) {
	fixture := newEngineFixture(t)
	infraErr := errors.New("artifact filesystem unavailable")
	verifier := &engineInfrastructureVerifier{calls: fixture.calls, err: infraErr}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, verifier, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if !errors.Is(err, infraErr) {
		t.Fatalf("expected infrastructure error identity preserved, got %v", err)
	}
	if len(result.Batches) != 0 {
		t.Fatalf("expected no Batch settlement on infrastructure error, got %+v", result.Batches)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify" {
		t.Fatalf("expected cycle to halt after verification infrastructure error, got %q", got)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status == rounds.StatusFailed {
		t.Fatal("expected infrastructure error not to mark Review Issue failed")
	}
	verificationEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification)
	verdicts := 0
	for _, event := range verificationEvents {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == string(runevent.VerificationPhaseVerdict) {
			verdicts++
			if payload["verdict"] != string(runevent.VerificationVerdictFailed) {
				t.Fatalf("expected failed infrastructure verdict, got %v", payload["verdict"])
			}
		}
	}
	if verdicts != 1 {
		t.Fatalf("expected one aggregate verdict before halting, got %d", verdicts)
	}
}

func TestResolveCycleStopRequestAfterAttemptOneFailureDoesNotRepair(t *testing.T) {
	fixture := newEngineFixture(t)
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store}
	verifier := &engineStopAfterCommandFailureVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify" {
		t.Fatalf("expected Stop Request to prevent Verification Feedback repair, got %q", got)
	}
	if len(runner.requests) != 1 {
		t.Fatalf("expected no repair Agent request after Stop Request, got %d", len(runner.requests))
	}
	if len(result.Batches) != 0 {
		t.Fatalf("expected no Batch settlement after Stop Request, got %+v", result.Batches)
	}
	issue, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse issue: %v", parseErr)
	}
	if issue.Status == rounds.StatusFailed {
		t.Fatal("expected stopped Batch not to be marked failed")
	}
}

func TestResolveCycleContinuesToNextBatchAfterFailedBatch(t *testing.T) {
	fixture := newEngineFixtureWithItems(t, []reviewsource.ReviewItem{
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
		{
			Title:                   "major: close response body",
			File:                    "internal/http/client.go",
			Line:                    17,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "second review body",
			SourceRef:               "thread:PRRT_2,comment:PRRC_2",
			ReviewHash:              "def",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		},
	})
	// One before-snapshot for each Batch, then the after-snapshot of the
	// second Batch: its commit must exclude the first Batch's residue.
	fixture.worktree = &engineFakeWorktree{snapshots: [][]string{
		nil,
		{"src/batch-one-residue.go"},
		{"src/batch-one-residue.go", "src/batch-two-change.go"},
	}}
	verifier := &engineFakeVerifier{
		calls: fixture.calls,
		store: fixture.store,
		runID: fixture.run.ID,
		script: []error{
			errors.New("attempt 1 failed"),
			errors.New("attempt 2 failed"),
			nil,
		},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	source := &engineFakeSource{calls: fixture.calls}
	plan := fixture.plan()
	plan.Session = agent.SessionRefForRun(fixture.run.ID, fixture.gitRoot)
	plan.Batches = []rounds.Batch{
		{Number: 1, Issues: []rounds.Issue{{Path: fixture.issuePaths[0]}}},
		{Number: 2, Issues: []rounds.Issue{{Path: fixture.issuePaths[1]}}},
	}
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("resolve cycle: %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify>source>agent>verify>commit>source>source" {
		t.Fatalf("expected second Batch to run after the first failed, got %q", got)
	}
	if len(result.Batches) != 2 {
		t.Fatalf("expected two Batch outcomes, got %+v", result.Batches)
	}
	if !result.Batches[0].Failed || result.Batches[0].Committed || result.Batches[0].ResolvedSourceThreads != 0 {
		t.Fatalf("expected first Batch failed without commit or thread resolution, got %+v", result.Batches[0])
	}
	if result.Batches[1].Failed || !result.Batches[1].Committed || result.Batches[1].ResolvedSourceThreads != 1 {
		t.Fatalf("expected second Batch committed with one resolved thread, got %+v", result.Batches[1])
	}
	if result.Remaining != 1 {
		t.Fatalf("expected only the failed issue to remain Unresolved, got %d", result.Remaining)
	}
	first, parseErr := rounds.ParseIssue(fixture.issuePaths[0])
	if parseErr != nil {
		t.Fatalf("parse first issue: %v", parseErr)
	}
	if first.Status != rounds.StatusFailed {
		t.Fatalf("expected first issue failed, got %q", first.Status)
	}
	second, parseErr := rounds.ParseIssue(fixture.issuePaths[1])
	if parseErr != nil {
		t.Fatalf("parse second issue: %v", parseErr)
	}
	if second.Status != rounds.StatusResolved {
		t.Fatalf("expected second issue resolved, got %q", second.Status)
	}
	if len(committer.paths) != 1 || strings.Join(committer.paths[0], ",") != "src/batch-two-change.go" {
		t.Fatalf("expected second Batch commit to exclude first Batch residue, got %v", committer.paths)
	}
	if len(runner.requests) != 3 {
		t.Fatalf("expected initial+repair for first Batch and initial for second Batch, got %d", len(runner.requests))
	}
	for _, req := range runner.requests {
		if req.Session != plan.Session {
			t.Fatalf("expected shared Agent Session %#v, got %#v", plan.Session, req.Session)
		}
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

func TestResolveCycleStopRequestBeforeBatchPublishesStopAndDoesNothing(t *testing.T) {
	fixture := newEngineFixture(t)
	if err := fixture.store.RequestStop(context.Background(), fixture.run.ID); err != nil {
		t.Fatalf("request Stop: %v", err)
	}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})

	result, err := engine.ResolveCycle(context.Background(), fixture.plan())

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested, got %v", err)
	}
	if len(*fixture.calls) != 0 {
		t.Fatalf("expected no daemon actions after Stop Request, got %v", *fixture.calls)
	}
	if result.Remaining != fixture.plan().TotalIssues {
		t.Fatalf("expected all issues to remain before the first Batch, got %+v", result)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindDaemonStatus {
		t.Fatalf("expected daemon stop event published to the sink, got %v", kinds)
	}
}

func TestResolveCycleStopRequestAfterBatchSettlementHaltsBeforeNextBatch(t *testing.T) {
	fixture := newEngineFixtureWithItems(t, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first",
			File:                    "internal/first.go",
			Line:                    10,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "first",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second",
			File:                    "internal/second.go",
			Line:                    20,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "second",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 10, 12, 1, 0, 0, time.UTC),
		},
	})
	plan := fixture.plan()
	plan.Batches = []rounds.Batch{
		{Number: 1, Issues: []rounds.Issue{{Path: fixture.issuePaths[0]}}},
		{Number: 2, Issues: []rounds.Issue{{Path: fixture.issuePaths[1]}}},
	}
	plan.TotalIssues = 2
	source := &engineFakeSource{
		calls: fixture.calls,
		afterResolve: func(context.Context, reviewsource.ResolveRequest) error {
			return fixture.store.RequestStop(context.Background(), fixture.run.ID)
		},
	}
	engine := fixture.engine(t, &engineFakeRunner{calls: fixture.calls, store: fixture.store}, &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}, &engineFakeCommitter{calls: fixture.calls}, &engineFakePusher{calls: fixture.calls}, source)

	result, err := engine.ResolveCycle(context.Background(), plan)

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit>source" {
		t.Fatalf("expected first Batch to verify, commit, and resolve source before stop, got %q", got)
	}
	if len(result.Batches) != 1 || !result.Batches[0].Committed || result.Remaining != 1 {
		t.Fatalf("expected first Batch settled and one issue remaining, got %+v", result)
	}
	first, err := rounds.ParseIssue(fixture.issuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	second, err := rounds.ParseIssue(fixture.issuePaths[1])
	if err != nil {
		t.Fatalf("parse second issue: %v", err)
	}
	if first.Status != rounds.StatusResolved {
		t.Fatalf("expected first issue resolved before stop, got %q", first.Status)
	}
	if second.Status != rounds.StatusPending {
		t.Fatalf("expected second issue left pending, got %q", second.Status)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindDaemonStatus {
		t.Fatalf("expected daemon stop event after first Batch settlement, got %v", kinds)
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

func (runner *publishingFakeRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

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
		if err := rounds.SetIssueStatus(issue.Path, rounds.StatusResolved, "", ""); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{}, nil
}

func (runner *publishingFakeRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
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
