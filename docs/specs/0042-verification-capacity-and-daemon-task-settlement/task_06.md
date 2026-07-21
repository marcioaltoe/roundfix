---
task: task_06
spec: 0042-verification-capacity-and-daemon-task-settlement
status: pending
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

- [ ] Map macro invariants to the existing CLI/daemon integration suites.
- [ ] Build a disposable concurrent Implement fixture with real Verification.
- [ ] Cover serialized success and Agent-status normalization.
- [ ] Cover temporary retry success and retained evidence.
- [ ] Add repeated-temporary, deterministic, and queued-cancellation companions.
- [ ] Prove journal, Task files, worktrees, stdout/stderr, and exit outcomes.
- [ ] Run the integrated matrix under the race detector.

## Acceptance Criteria

- [ ] The disposable flow observes two simultaneous Agent turns and never more
      than one active Verification attempt at capacity `1`.
- [ ] Journal replay orders waiting before started for every attempt and shows
      one exclusive retry with no second retry or hidden Agent repair.
- [ ] A premature Agent terminal status is normalized, the real gate runs, and
      only a Daemon verdict determines final Task status.
- [ ] A one-time exit `75` flow ends Clean and retains distinct initial/retry
      evidence; repeated `75` ends non-clean and preserves the Task Worktree.
- [ ] A non-`75` failure still follows exactly one Agent repair and numbered
      attempt `2`, proving the old deterministic policy was not generalized.
- [ ] Cancellation while queued starts no child process and leaves resumable
      status without a leaked worker or permit.
- [ ] Normal Implement stdout and exit code remain byte-compatible; waiting,
      retry, and diagnostic progress remains on stderr or the Run Event Stream.
- [ ] Focused, repeated, race, and build gates pass from a clean disposable
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
- interface: `internal/runevent/stream_test.go`
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
