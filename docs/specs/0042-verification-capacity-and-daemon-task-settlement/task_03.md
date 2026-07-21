---
task: task_03
spec: 0042-verification-capacity-and-daemon-task-settlement
status: pending
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

- [ ] Add the Task-cycle-owned fair shared/exclusive gate.
- [ ] Thread one gate through every Task Worktree plan copy.
- [ ] Journal waiting and acquired attempt evidence in stable order.
- [ ] Bound shared Verification at the configured capacity.
- [ ] Release before repair and reacquire for the final attempt.
- [ ] Add channel-coordinated overlap, fairness, cancellation, and race tests.

## Acceptance Criteria

- [ ] With Task Capacity `2` and Verification Capacity `1`, two Agent turns can
      overlap while the observed maximum active Verification count is `1`.
- [ ] With Verification Capacity `2`, two ready shared attempts overlap and
      both complete without permit loss.
- [ ] Every attempt records `waiting` before `started`, including an immediately
      available gate.
- [ ] A first deterministic failure releases capacity before its Agent repair,
      and attempt `2` queues and acquires again.
- [ ] An exclusive waiter begins only after active shared attempts drain and
      later shared waiters cannot bypass it.
- [ ] Cancellation while queued starts zero Verification commands, leaves no
      goroutine blocked, and lets a later independent test acquire full capacity.
- [ ] Tests coordinate on observable channels and counters rather than sleeps,
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
