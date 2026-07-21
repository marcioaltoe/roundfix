---
task: task_02
spec: 0042-verification-capacity-and-daemon-task-settlement
status: pending
type: backend
complexity: high
---

# Task 02: Make the Daemon own Implement Task status

## Overview

Move every Implement Task status transition to the Daemon and redefine the
Agent turn as an implementation-ready handoff. A premature Agent-authored
`completed` or `failed` value must be normalized to `in_progress`, preserve
Result evidence, and continue into authoritative Daemon Verification.

## Requirements

1. MUST make the Daemon write `in_progress` before the initial Task Agent is
   invoked and remain the sole writer of terminal Task status during Implement.
2. MUST instruct initial and Verification Feedback Agents not to edit Task
   status, run declared `## Verification` commands, or claim Task completion.
3. MUST allow focused implementation checks and Agent-authored Result evidence
   as part of the implementation-ready handoff.
4. MUST reload the Task after every Agent turn, preserve valid non-status
   changes, and normalize any status value back to the Daemon's current
   `in_progress` state.
5. MUST run Daemon Verification after a successful Agent handoff regardless of
   whether the Agent wrote `completed`, `failed`, `pending`, or `in_progress`.
6. MUST keep genuine Agent execution, unreadable Task artifact,
   infrastructure, and Stop Request paths distinct; only the Daemon may decide
   whether they settle failed or remain resumable.
7. MUST preserve one same-Session Verification Feedback turn, failed Task
   Worktree retention, independent Task continuation, dependency blocking,
   and the Settle Command contract.

## Subtasks

- [ ] Move initial `in_progress` persistence to the Daemon start boundary.
- [ ] Replace Agent status/verdict instructions with implementation-ready handoff.
- [ ] Normalize reloaded Agent status without losing Result content.
- [ ] Remove the Agent-authored failed-status Verification bypass.
- [ ] Update the repository-owned implement-task workflow contract.
- [ ] Replace bypass tests with authoritative Verification and failure-path cases.

## Acceptance Criteria

- [ ] The Task file is `in_progress` before the first Agent request and the
      journal's Task-start event agrees with it.
- [ ] Agent-authored `completed` and `failed` values both reach Daemon
      Verification; neither directly settles the Task.
- [ ] Result prose written by the Agent survives status normalization and is
      committed with a passing Task.
- [ ] A genuine Agent runtime failure starts no Verification command and is
      settled only through the Daemon's existing failure policy.
- [ ] A Stop Request after Agent work invents no `completed` or `failed` Task
      status and preserves resumable evidence.
- [ ] Prompt and authorial skill tests forbid declared Verification/status
      authorship while still allowing focused checks.
- [ ] Existing repair, dependency, worktree retention, and Settle Command
      regression tests pass.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/spec_prompt_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/spec/spec.go`

## Verification

- `rtk go test ./internal/agent -run 'TestBuild(Task|VerificationRepair)Prompt' -count=1` — expected: prompts require implementation-ready handoff, prohibit status and declared Verification authorship, and preserve focused-check guidance.
- `rtk go test ./internal/daemon -run 'TestTaskCycle.*(AgentStatus|DaemonStatus|Verification|AgentFailure|Stop)' -count=1` — expected: every Agent status variant follows the Daemon-owned settlement matrix and genuine failure/cancellation paths stay distinct.
- `rtk go test ./internal/cli -run 'Test(Settle|Implement).*TaskStatus' -count=1` — expected: Implement ownership changes do not alter Settle Command recovery.
- `rtk go test -race ./internal/agent ./internal/daemon -run 'Test(BuildTaskPrompt|BuildVerificationRepairPrompt|TaskCycle.*(AgentStatus|DaemonStatus|Verification|AgentFailure|Stop))' -count=1` — expected: prompt and status ownership behavior is race-free.

## References

- `_prd.md` → Goals 2–3; User Stories 4–5 and 8; Core Features 3–5 and 9–10; Success Metrics.
- `_techspec.md` → System Architecture; Implementation Design: Data Models; Testing Approach; Build Order 2.
- `../../adr/0057-daemon-exclusively-owns-implement-task-status.md` → exclusive Daemon status authorship and superseded bypass.
- `../../adr/0014-daemon-runs-task-verification-and-settles-status.md` → preserved Daemon Verification and settlement boundary.
