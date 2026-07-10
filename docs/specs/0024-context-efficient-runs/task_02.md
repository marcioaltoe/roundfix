---
task: task_02
spec: 0024-context-efficient-runs
status: pending
type: backend
complexity: high
---

# Task 02: Allow one same-session Verification repair

## Overview

Complete the Daemon-owned Verification lifecycle by returning one failed
verdict to the same Agent Session and running one final authoritative attempt.
The slice covers both Task and review Batch Work Items while preserving
independent Task scheduling and explicit failed-prerequisite settlement.

## Requirements

1. MUST remove the requirement for Agents to run authoritative configured Verification from initial Task and review prompts.
2. MUST send one Verification Feedback prompt only for a typed attempt-1 command failure and reuse the same Agent Session.
3. MUST include the failed command, wrapped failure, and diagnostic artifact path in feedback without embedding log output.
4. MUST rerun the complete Verification command sequence after repair and settle from the final verdict without another repair.
5. MUST treat repair Agent errors, Stop Requests, and Verification infrastructure errors under existing infrastructure policies rather than as command failures.
6. MUST preserve clean Agent-authored failed settlement for missing credentials and continue ready independent Tasks while blocking dependents.
7. MUST avoid publishing a second Batch-start boundary for the repair turn.

## Subtasks

- [ ] Revise initial Task and review prompt contracts for Daemon-owned Verification.
- [ ] Add the path-based Verification Feedback prompt builder.
- [ ] Add the one-repair state transition to Task execution.
- [ ] Add the same transition to review Batch execution.
- [ ] Preserve Stop, infrastructure, prerequisite, and scheduler policies.
- [ ] Cover all first/final verdict combinations and session identity.

## Acceptance Criteria

- [ ] Attempt 1 pass settles normally and sends zero Verification output or feedback messages to the Agent.
- [ ] Attempt 1 failure sends exactly one feedback prompt to the same SessionRef and attempt 2 reruns every configured command.
- [ ] Attempt 2 pass settles completed/resolved; attempt 2 failure settles under the existing failed Work Item policy.
- [ ] No path performs a third Verification attempt or second repair prompt.
- [ ] Infrastructure failure or Stop Request never enters the repair loop.
- [ ] A Task that records a missing credential as failed does not prevent an independent ready Task from starting, while its dependents remain blocked.
- [ ] Review repair produces no duplicate Batch-start event.

## Verification

- `rtk go test ./internal/agent ./internal/daemon` - expected: prompt, same-session repair, final settlement, stop, infrastructure, and scheduler tests pass.
- `rtk go test -race ./internal/daemon` - expected: scheduler and repair flows pass under the race detector.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-concurrency/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`

## References

`_prd.md` -> User Stories 1, 2, 6; Core Features 1-4; User Experience. `_techspec.md` -> System Architecture: Verification flow; API Contracts: repair prompts; Build Order 2. ADR-0014; ADR-0038.
