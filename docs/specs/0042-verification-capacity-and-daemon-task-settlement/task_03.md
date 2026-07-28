---
task: task_03
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
type: backend
complexity: high
---

# Task 03: Queue Task Verification and journal Waiting for Verification

## Overview

Introduce the Task-cycle-owned Verification gate and run every Task attempt
through cancellation-aware capacity acquisition. Concurrent Agent work must
remain possible while the journal gives every attempt deterministic waiting
then started ordering and Verification Feedback runs without holding capacity.

## Requirements

1. MUST create one bounded Verification gate per Task cycle and share it across
   every Task Worktree worker in that cycle.
2. MUST support fair shared and exclusive acquisition with no process-global,
   repository-global, filesystem, or external coordination service.
3. MUST emit Waiting for Verification before every acquisition and emit
   Verification started only after acquisition succeeds.
4. MUST hold one shared capacity unit for the complete command sequence in one
   numbered attempt and release it on pass, failure, cancellation, or error.
5. MUST release capacity before invoking Verification Feedback and reacquire it
   for the final numbered attempt.
6. MUST honor context cancellation while queued without starting a command,
   leaking a permit, blocking a worker goroutine, or assigning a false terminal
   Task verdict.
7. MUST prevent queued exclusive acquisition from starvation by later shared
   attempts and introduce no dependency beyond the standard library.

## Subtasks

- [x] Add the Task-cycle-owned fair shared/exclusive gate.
- [x] Thread one gate through every Task Worktree plan copy.
- [x] Journal waiting and acquired attempt evidence in stable order.
- [x] Bound shared Verification at the configured capacity.
- [x] Release before repair and reacquire for the final attempt.
- [x] Add channel-coordinated overlap, fairness, cancellation, and race tests.

## Acceptance Criteria

- [x] With Task Capacity `2` and Verification Capacity `1`, two Agent turns can
      overlap while the observed maximum active Verification count is `1`.
- [x] With Verification Capacity `2`, two ready shared attempts overlap and
      both complete without permit loss.
- [x] Every attempt records `waiting` before `started`, including an immediately
      available gate.
- [x] A first deterministic failure releases capacity before its Agent repair,
      and attempt `2` queues and acquires again.
- [x] An exclusive waiter begins only after active shared attempts drain and
      later shared waiters cannot bypass it.
- [x] Cancellation while queued starts zero Verification commands, leaves no
      goroutine blocked, and lets a later independent test acquire full capacity.
- [x] Tests coordinate on observable channels and counters rather than sleeps,
      private method calls, or production-only hooks.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-concurrency/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`

## Verification

- `rtk go test ./internal/daemon -run 'TestTaskCycle.*(VerificationCapacit|WaitingForVerification|RepairReacquires)' -count=1` — expected: shared attempt bounds, event ordering, and release/reacquire behavior pass without timing sleeps.
- `rtk go test ./internal/daemon -run 'TestVerificationGate.*(Exclusive|Fair|Cancel|Release)' -count=20` — expected: exclusive fairness, cancellation, and permit restoration remain deterministic across repeated runs.
- `rtk go test ./internal/runevent -run 'Test.*Verification.*Waiting' -count=1` — expected: waiting is accepted and projected as an additive `roundfix-events/v1` Verification phase.
- `rtk go test -race ./internal/daemon ./internal/runevent -run 'Test(TaskCycle.*(VerificationCapacit|WaitingForVerification|RepairReacquires)|VerificationGate.*(Exclusive|Fair|Cancel|Release)|.*Verification.*Waiting)' -count=1` — expected: queue, worker, event, and cancellation behavior is race-free.

## References

- `_prd.md` → Goals 1 and 4–6; User Stories 1–3 and 8; Core Features 1–2, 5, 8–9; Success Metrics.
- `_techspec.md` → Implementation Design: Interfaces and Data Models; Testing Approach; Build Order 3.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → per-cycle capacity, event ordering, and fairness decision.
- `../../adr/0038-daemon-allows-one-verification-repair.md` → release around the preserved Agent repair.

## Result

Implemented one fair, cancellation-aware Verification gate per Task cycle.
Every Task Worktree plan copy shares that gate. Shared attempts hold one unit
across their full command sequence, publish `waiting` before acquisition and
`started` after acquisition, and release capacity before Verification Feedback
or any final numbered attempt.

Verification evidence:

- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test ./internal/daemon -run 'TestTaskCycle.*(VerificationCapacit|WaitingForVerification|RepairReacquires)' -count=1` — passed: 5 tests.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test ./internal/daemon -run 'TestVerificationGate.*(Exclusive|Fair|Cancel|Release)' -count=20` — passed: 40 repeated tests.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test ./internal/runevent -run 'Test.*Verification.*Waiting' -count=1` — passed: 1 test.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test -race ./internal/daemon ./internal/runevent -run 'Test(TaskCycle.*(VerificationCapacit|WaitingForVerification|RepairReacquires)|VerificationGate.*(Exclusive|Fair|Cancel|Release)|.*Verification.*Waiting)' -count=1` — passed: 8 tests.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test -race ./internal/daemon ./internal/runevent -count=1` — passed: 164 tests.
- `rtk make verify` — passed outside the sandbox: 2,704 repository tests, 4 Skill tests, the Roundfix Skill check, and the build. The first sandboxed run reached 2,699 passes but its five `/bin/ps`-dependent owner-process tests were denied by the environment; the exact unsandboxed rerun passed.

Acceptance evidence:

- `TestTaskCycleVerificationCapacityOneBoundsConcurrentTaskWorktrees` observed two overlapping Agent turns and a maximum of one active Verification attempt.
- `TestTaskCycleVerificationCapacityTwoOverlapsReadyAttemptsWithoutPermitLoss` observed two concurrent shared attempts and two completed Tasks.
- `TestTaskCycleWaitingForVerificationPrecedesStartedWhenImmediatelyAvailable`, `TestTaskCycleRepairReacquiresVerificationCapacityAfterFeedback`, and `TestVerificationWaitingEventProjectsAdditivePhase` proved stable `waiting` then `started` ordering and the additive `roundfix-events/v1` projection.
- `TestTaskCycleRepairReacquiresVerificationCapacityAfterFeedback` proved the failed first attempt released capacity before Verification Feedback, another Task verified during repair, and attempt 2 queued and reacquired.
- `TestVerificationGateExclusiveWaiterBlocksLaterShared` proved active shared attempts drain before exclusive acquisition and later shared arrivals cannot bypass the exclusive waiter.
- `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement` and `TestVerificationGateCancelExclusiveWaiterRestoresFullSharedCapacity` proved queued cancellation starts no command, assigns no false terminal Task status, returns the worker, and restores full capacity.
- The concurrency tests use channels, mutex-protected counters, and observable Run Events. They contain no sleeps, gate-state inspection, or production-only test hooks.

Follow-up: Task 04 owns Temporary Verification Failure classification and the
production path that selects exclusive acquisition. This task supplies and
proves the exclusive gate contract without implementing that later slice.
