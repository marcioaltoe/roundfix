---
name: implement-task
description: Execute one task file from docs/specs/<slug>/ end-to-end — ground in the PRD/TechSpec, implement the slice, verify every acceptance criterion with fresh evidence, record the result, update status, and commit. Starts immediately when assigned, with no confirmation prompt. Use when the user says "run task_03", "execute the next task", "pick up a task", or when an implementation loop or daemon assigns a task file from a spec folder.
metadata:
  category: implementation
  tags: [workflow, coding, agents, testing]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Implement Task

Build exactly one task from a spec folder, end to end: load → plan → implement → **verify** → record → commit. The task file is the contract; the verification gate at the end is what makes `status: completed` mean something to the scheduler that reads it next.

**Being assigned the task — by the user, a loop, or a daemon — is the authorization to start.** Begin immediately: no greeting, no plan-approval question, no waiting for a "go". The anomalies called out below (stale `in_progress` status, contradictory requirements) are the only reasons to pause before implementing.

## 1. Load

- Read the assigned `task_NN.md`. Check `status`:
  - `pending` → set `status: in_progress` and proceed.
  - `in_progress` → a previous session may have died mid-task. Inspect the worktree and git log for partial work, report what you find, and ask how to proceed rather than double-building.
  - `completed` / `failed` → stop and report; re-running is an explicit human decision.
- Read the sections of `_prd.md` and `_techspec.md` named in the task's References, plus `CONTEXT.md` (use its vocabulary in code names and test names) and any referenced ADRs.
- **If requirements contradict each other** — task vs spec, spec vs ADR — stop and report the conflict. Guessing buries a spec bug inside an implementation.

## 2. Plan

- Build a checklist from Requirements + Acceptance Criteria.
- Capture the pre-change signal that proves the task is _not yet_ done — the failing test, the missing behavior, the 404. Without a red starting point you cannot show your change is what turned things green.

## 3. Implement

- Stay inside the slice. The PRD's Non-Goals and the task's scope are walls, not suggestions — work that belongs to another task goes in a follow-up note, not in this diff.
- Tests first at the seams the TechSpec names (they are pre-agreed; a new seam needs the user's sign-off). Typecheck and run the focused tests frequently; save the full suite for the gate.
- Root cause only — no lint/type suppressions, no swallowed errors, no timing hacks. A workaround closes the task and opens a bug.

## 4. Verify — the gate

Evidence before status, always in this order:

1. Run every command in the task's `## Verification` section verbatim.
2. Run the repository's verify pipeline (`make verify`, or the build/lint/typecheck/test equivalents this repo documents).
3. Walk Acceptance Criteria one by one: each needs fresh evidence from this session — a command output, a test name that passes, an observed behavior. A green suite is not evidence for a criterion the suite doesn't cover.

A narrow verification never supports a broad claim. If any check fails after honest root-cause attempts: set `status: failed`, record what was tried in `## Result`, and report — a loud failure the scheduler can retry beats a quiet "mostly done".

## 5. Record

Append a `## Result` section to the task file: what changed (described by behavior, not file lists), commands run with outcomes, evidence per acceptance criterion, and any follow-ups discovered. Tick the Subtasks and Acceptance Criteria checkboxes that the evidence supports. Set `status: completed`.

Never touch `_tasks.md` — it owns graph topology, not progress.

## 6. Commit

- Stage only this task's files (`git status --short` first; unrelated changes stay out).
- One commit per task, Conventional Commits format, task id in the body for traceability (`spec: <slug> / task_02`).
- Never push and never open a PR from inside a task — publishing is a separate, explicit action.

## Anti-patterns

- Marking `completed` on "should work" — the status field is read by machines that won't double-check.
- Weakening, skipping, or deleting a failing test to get to green.
- Fixing "one more thing" spotted along the way — that's a follow-up, not scope.
- Re-verifying with stale output from earlier in the session; evidence must postdate the last edit.
- Editing other tasks' files or statuses.
