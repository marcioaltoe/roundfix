---
task: task_06
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
type: test
complexity: high
---

# Task 06: Prove the integrated capacity and settlement contract

## Overview

Exercise the shipped Implement flow across configuration, concurrent Task
Agents, real shell Verification, exit-75 retry, journal replay, Live Run View,
Daemon settlement, and terminal CLI behavior. This is the macro acceptance
slice complementing the focused invariants embedded in Tasks 01–05, not a
replacement for them or a mock-only confidence layer.

## Requirements

1. MUST state the invariant, owning layer, and existing canonical suite before
   adding each integration case and extend existing suites where they own the
   behavior.
2. MUST use disposable repositories and real shell Verification processes;
   fakes may isolate only the ACP Runtime or other external boundaries already
   represented by production interfaces.
3. MUST prove Task Capacity `2` permits Agent overlap while Verification
   Capacity `1` limits active gates to one and records waiting before started.
4. MUST prove an Agent-authored terminal status cannot bypass Daemon
   Verification and terminal status is written only by Daemon settlement.
5. MUST prove a project wrapper returning `75` once receives one exclusive
   retry, preserves both diagnostic identities, runs no Agent repair, and can
   finish Clean.
6. MUST include negative companions for repeated exit `75`, deterministic
   non-`75` failure, and cancellation while queued.
7. MUST assert journal and resulting Task/worktree state rather than values
   copied from fake setup, and MUST preserve stdout/stderr and top-level exit
   contracts.
8. MUST add no production-only test hook, sleep-based synchronization, new test
   dependency, weakened assertion, quarantine, or generic retry.

## Subtasks

- [x] Map macro invariants to the existing CLI/daemon integration suites.
- [x] Build a disposable concurrent Implement fixture with real Verification.
- [x] Cover serialized success and Agent-status normalization.
- [x] Cover temporary retry success and retained evidence.
- [x] Add repeated-temporary, deterministic, and queued-cancellation companions.
- [x] Prove journal, Task files, worktrees, stdout/stderr, and exit outcomes.
- [x] Run the integrated matrix under the race detector.

## Acceptance Criteria

- [x] The disposable flow observes two simultaneous Agent turns and never more
      than one active Verification attempt at capacity `1`.
- [x] Journal replay orders waiting before started for every attempt and shows
      one exclusive retry with no second retry or hidden Agent repair.
- [x] A premature Agent terminal status is normalized, the real gate runs, and
      only a Daemon verdict determines final Task status.
- [x] A one-time exit `75` flow ends Clean and retains distinct initial/retry
      evidence; repeated `75` ends non-clean and preserves the Task Worktree.
- [x] A non-`75` failure still follows exactly one Agent repair and numbered
      attempt `2`, proving the old deterministic policy was not generalized.
- [x] Cancellation while queued starts no child process and leaves resumable
      status without a leaked worker or permit.
- [x] Normal Implement stdout and exit code remain byte-compatible; waiting,
      retry, and diagnostic progress remains on stderr or the Run Event Stream.
- [x] Focused, repeated, race, and build gates pass from a clean disposable
      fixture without network access or assumed machine services.

## Context

- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/daemon_test.go`
- interface: `internal/runevent/event_test.go`
- interface: `internal/tui/cockpit_test.go`
- interface: `internal/worktree/worktree_test.go`

## Verification

- `rtk go test ./internal/cli -run 'TestRunImplement.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation)' -count=1` — expected: disposable CLI flows prove serialized capacity, Daemon status, temporary retry, and cancellation outcomes.
- `rtk go test ./internal/daemon ./internal/runevent ./internal/tui -run 'Test.*(IntegratedVerificationCapacit|TemporaryVerificationFlow|WaitingForVerificationReplay)' -count=1` — expected: cross-package state, journal, and rendered evidence agree.
- `rtk go test ./internal/cli ./internal/daemon -run 'Test.*(VerificationCapacit|TemporaryVerification|QueuedCancellation)' -count=20` — expected: condition-coordinated integration cases remain deterministic without sleeps.
- `rtk go test -race ./internal/cli ./internal/daemon ./internal/runevent ./internal/tui -run 'Test.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation|WaitingForVerification)' -count=1` — expected: the integrated user flows have no data race, blocked worker, or permit leak.
- `rtk go build -buildvcs=false ./...` — expected: the tested production surface builds without test-only seams.

## References

- `_prd.md` → all Goals, User Stories, Core Features, User Experience, and Success Metrics.
- `_techspec.md` → Coverage Map; Integration Points; Testing Approach; Build Order 7.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → capacity and exclusive-retry invariants.
- `../../adr/0057-daemon-exclusively-owns-implement-task-status.md` → status and Verification-handoff invariants.

## Result

### What changed

- `internal/cli/implement_test.go` extends the canonical disposable Implement
  suite with real shell Verification cases for capacity `2:1`, Agent-authored
  terminal-status normalization, temporary exit `75`, repeated exit `75`,
  deterministic repair, and cancellation while queued. Coordination uses Agent
  channels, named pipes, and journal-condition polling; only the ACP Runner is
  faked.
- `internal/daemon/task_engine_test.go`, `internal/runevent/event_test.go`, and
  `internal/tui/cockpit_test.go` align the existing owning tests with the
  integrated verification matrix names, so the declared cross-package command
  exercises capacity, temporary retry, journal replay, and rendered phase
  projection.
- `internal/cli/implement.go` now projects an authoritative terminal
  `daemon.TaskOutcome` into the closing report. The repeated-`75` fixture
  exposed that a preserved failed Task Worktree could otherwise be reported as
  `skipped` from the Run Worktree's stale pre-settlement Task file. The focused
  regression leaves that file pending and proves the Daemon's failed outcome
  controls the report and counts.

### Acceptance criteria

1. **Independent capacities** —
   `TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow`
   observes two simultaneous Agent turns, blocks real shell children on named
   pipes, and derives a maximum of one active Verification from the journal and
   serialized shell log.
2. **Journal ordering and retry budget** — the capacity flow asserts each
   attempt's `waiting` cursor precedes `started`.
   `TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair`
   observes one exclusive retry and no Verification Feedback;
   `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`
   observes exactly two temporary failures, one retry start, and no retry
   number greater than one.
3. **Daemon-owned settlement** — the capacity flow has the Agent write both
   `completed` and `failed`, then proves both real gates execute and both
   resulting Task files are `completed` after Daemon settlement.
4. **Temporary outcomes and evidence** — the one-time `75` case finishes
   Clean, retains the initial diagnostic, and proves the successful retry uses
   its distinct diagnostic identity. The repeated-`75` case finishes
   Unresolved, retains distinct initial/retry contents, and preserves the
   failed Task Worktree and branch with `status: failed`.
5. **Deterministic policy** —
   `TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair`
   makes a real shell command exit `42`, observes waiting attempts `[1 2]`,
   exactly one Verification Feedback event, exactly two Agent calls, and no
   exclusive-retry metadata.
6. **Queued cancellation** —
   `TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks`
   waits until one gate is active and one queued, cancels, and proves the queued
   marker never exists, both Task Worktrees remain `in_progress`, no queued
   settlement is journaled, the worker set returns, and the Run lock is
   released.
7. **Terminal contract** — success, repaired success, and Stopped cases assert
   exact stdout bytes and exit codes. Run identity, terminal progress,
   diagnostics, waiting, and retry evidence are asserted on stderr or the Run
   Event Journal and forbidden from normal stdout.
8. **Hermetic verification** — every macro fixture creates a disposable Git
   repository and uses local shell processes only. Focused, repeated, race,
   formatting, diff, and build checks passed without network access or machine
   services.

### Verification evidence

- `rtk go test ./internal/cli -run 'TestRunImplement.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation)' -count=1` →
  `Go test: 12 passed in 1 packages`
- `rtk go test ./internal/daemon ./internal/runevent ./internal/tui -run 'Test.*(IntegratedVerificationCapacit|TemporaryVerificationFlow|WaitingForVerificationReplay)' -count=1` →
  `Go test: 4 passed in 3 packages`
- `rtk go test ./internal/cli ./internal/daemon -run 'Test.*(VerificationCapacit|TemporaryVerification|QueuedCancellation)' -count=20` →
  `Go test: 400 passed in 2 packages`
- `rtk go test -race ./internal/cli ./internal/daemon ./internal/runevent ./internal/tui -run 'Test.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation|WaitingForVerification)' -count=1` →
  `Go test: 29 passed in 4 packages`
- `rtk go build -buildvcs=false ./...` → passed
- `rtk gofmt -d internal/cli/implement.go internal/cli/implement_test.go internal/daemon/task_engine_test.go internal/runevent/event_test.go internal/tui/cockpit_test.go` →
  no diff
- `rtk git -c core.fsmonitor=false diff --check` → passed

All Go gates used
`GOCACHE=/private/tmp/roundfix-task06-gocache` because the sandbox denies writes
to the default user Go build cache. No production test hook, dependency,
sleep-based synchronization, quarantine, generic retry, commit, push, or pull
request was added.
