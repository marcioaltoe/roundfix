---
title: "Daemon Run engine: resolve cycle extraction with explicit dependencies"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 5
  - 6
  - 7
  - 10
  - 11
  - 12
  - 14
  - 15
  - 16
blocked_by:
  - 03-runevent-seam-runner-and-writersink.md
---

# Daemon Run engine: resolve cycle extraction with explicit dependencies

## Parent

.scratch/roundfix-event-journal/04-daemon-run-engine-prd.md

## What to build

Structural core of the Daemon Run engine PRD: a Run engine in the daemon module executing one resolve cycle over a
validated plan (per Batch: Agent → verify issue statuses → commit when auto-commit → resolve Review Source threads
for resolved/invalid), returning per-Batch outcomes and remaining Unresolved Review Issue count. Final Push is a
separate explicit engine operation; gating (no Unresolved Review Issues, auto-push on) stays in the caller per
ADR 0001. The engine updates intermediate Run states; Run creation and terminal completion stay with the caller.
Constructor takes an explicit dependencies struct (runner, verifier, committer, pusher, Review Source, rounds/store,
clock, Run Event sink); CLI package globals for orchestration die. The resolve command uses the engine. Stop
semantics move with the orchestration: no new Batch/verify/commit/push/source mutation after cancellation; stop
published to the sink; Agent worktree changes preserved.

## Acceptance criteria

- [x] Engine proves resolve→verify→commit→source-resolution contract with fake collaborators and a temp store; no terminal, no network
- [x] Final Push is a separate operation; caller-side gating decides when it runs
- [x] Intermediate Run states updated during the cycle; terminal completion remains caller-owned
- [x] CLI orchestration globals removed; resolve command delegates to the engine
- [x] Stop tests: no post-stop daemon actions, stop event reaches sink, worktree preserved
- [x] Existing CLI orchestration tests and fakes migrate to engine tests
- [x] Batch failure policy unchanged: verification/Agent failure fails the Run

## Blocked by

- 03-runevent-seam-runner-and-writersink.md

## Comments

**2026-06-10 (agent):** Built `daemon.Engine` (`internal/daemon/engine.go`): `NewEngine(Dependencies)` validates an
explicit struct (Runner, Verifier, Committer, Pusher, Source resolver, Runs store, Sink, clock, Progress writer);
`ResolveCycle(ctx, CyclePlan)` executes per Batch agent → terminal-status validation → verification →
duplicate-marking → commit (auto-commit) → Review Source thread resolution, returning per-Batch outcomes
(committed, resolved threads) and the remaining Unresolved Review Issue count; `FinalPush` is a separate explicit
operation with gating left in the caller (ADR 0001). Intermediate Run states: store gained non-terminal
`ResolvingWithAgent`/`Verifying`/`Pushing` constants and `UpdateRunState` (rejects terminal states and missing
Runs); engine sets them at phase boundaries; terminal completion stays caller-owned. Stop semantics: at any daemon
boundary a canceled context halts the cycle, publishes a `daemon.status` stopped event (new `KindDaemonStatus` in
runevent, full taxonomy deferred to PRD 6) via `context.WithoutCancel`; a stop during the Agent surfaces the
runner's own stop event and the Batch is not marked failed (worktree preserved). The five CLI orchestration globals
(`runAgentRuntime`, `runVerificationGate`, `createBatchCommit`, `resolveReviewSourceIssues`, `runFinalPush`) are
deleted; one `newEngineCollaborators` factory remains as the CLI dispatch-test seam; resolve delegates to the
engine via `executeResolveCycle`, and the watch resolver calls the same path interim (unit 09 finishes the watch
migration). Engine tests (`internal/daemon/engine_test.go`) use fake collaborators plus a real temp Run Database
and prove: call order `agent>verify>commit>source`, no push during a cycle, explicit FinalPush sets `Pushing`,
agent/verification failures fail the cycle and mark the Batch failed, stop-before-Batch publishes the daemon stop
event with zero collaborator calls, and stop-during-Agent preserves issue state. All existing CLI text-contract
tests pass unchanged. Verification: `rtk go vet ./...` clean, `rtk go test ./...` 181 passed in 15 packages,
`rtk go test -race` on daemon+cli 77 passed, `rtk go run ./cmd/roundfix --help` green.
