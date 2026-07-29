---
name: implement-task
description: Execute one task file from docs/specs/<slug>/ end-to-end — ground in the PRD/TechSpec, implement only the slice, and hand back evidence under the execution mode's settlement contract. Starts immediately when assigned, with no confirmation prompt. Use when the user says "run task_03", "execute the next task", "pick up a task", or when an implementation loop or daemon assigns a task file from a spec folder.
metadata:
  category: implementation
  tags: [workflow, coding, agents, testing]
  version: 0.0.2
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Implement task

Build exactly one Task from a Spec folder. The Task file defines the slice,
acceptance criteria, and Verification contract. Standalone execution owns the
full lifecycle; a Roundfix Daemon-assigned turn hands back
implementation-ready work for Daemon Verification and settlement.

**Being assigned the Task — by the user, a loop, or a Daemon — is the
authorization to start.** Begin immediately: no greeting, no plan-approval
question, no waiting for a "go". Contradictory requirements, or a stale
standalone `in_progress` status, are the only reasons to pause before
implementing.

## Execution modes

Standalone execution and a Roundfix Daemon-assigned turn share the same Task
slice, Result, and evidence requirements, but they have different settlement
owners.

- In standalone execution, the Agent updates Task status, runs every command
  in `## Verification`, settles the Task from fresh evidence, and commits only
  when the user has authorized a commit.
- In a Roundfix Daemon-assigned turn, the Daemon is the sole Task-status
  writer. The Agent must not edit status, run the declared `## Verification`
  commands, claim a terminal verdict, or commit. It may run focused checks,
  records implementation and focused-check evidence in `## Result`, and hands
  back implementation-ready work. The Daemon then runs the complete declared
  Verification verbatim before settlement.
- A Daemon Verification command failure releases Verification Capacity before
  one Verification Feedback repair turn in the same Agent Session. The Agent
  repairs the implementation, updates `## Result`, runs focused checks when
  useful, and hands back again without running declared Verification or
  settling status.

For a Daemon-assigned Task that arrives `in_progress`, start the Task fresh.
The work-target lock proves that no live Agent owns it. Do not ask how to
resume, and do not normalize or settle its status.

## 1. Load

Read the assigned `task_NN.md`, then run the Project Constraint preflight
before changing the Task status to `in_progress` in standalone execution:

1. A completed or archived legacy Spec is exempt from forced constraint
   backfill. Leave its artifacts byte-identical and stop rather than
   re-executing or rewriting it.
2. For an active Spec, require complete `Project Constraints` sections in the
   PRD and every present TechSpec. Each must account for identifier strategy,
   authentication and HTTP, active ADR obligations, and tooling authority as
   applicable or not applicable with a reason, with an operative
   `docs/agents/` source path for every row. Stop without changing Task status
   when any row, reason, or source is missing.
3. Determine whether the assigned Task creates, edits, renames, moves, or
   deletes repository-tooling configuration, scripts, ignore files, plugin
   declarations, or version pins. Tooling authorization is not implied by Task
   assignment, setup approval, or a generic implementation request.
4. Before a tooling Task can start, require express maintainer authorization
   and an exact bounded repository-relative file list in the PRD and every
   present TechSpec. The mutation allowlist is those exact paths plus the
   assigned Task file itself. Missing or conflicting authorization leaves the
   Task pending and stops execution.
5. Run `git status --short` and capture the pre-existing changed paths before
   work begins. They remain user-owned and cannot be counted as this Task's
   changes.

After the preflight, standalone execution checks `status`:

- `pending` → set `status: in_progress` and proceed.
- `in_progress` → a previous session may have died mid-task. Inspect the
  worktree and git log for partial work, report what you find, and ask how to
  proceed rather than double-building.
- `completed` / `failed` → stop and report; re-running is an explicit human
  decision.

In a Daemon-assigned turn, treat status as Daemon-owned and proceed without
editing it.

- Read the sections of `_prd.md` and `_techspec.md` named in the task's References, plus `CONTEXT.md` (use its vocabulary in code names and test names) and any referenced ADRs.
- **If requirements contradict each other** — task vs spec, spec vs ADR — stop and report the conflict. Guessing buries a spec bug inside an implementation.

## 2. Plan

- Build a checklist from Requirements + Acceptance Criteria.
- Capture the pre-change signal that proves the task is _not yet_ done — the failing test, the missing behavior, the 404. Without a red starting point you cannot show your change is what turned things green.

## 3. Implement

- Stay inside the slice. The PRD's Non-Goals and the task's scope are walls, not suggestions — work that belongs to another task goes in a follow-up note, not in this diff.
- Tests first at the seams the TechSpec names (they are pre-agreed; a new seam needs the user's sign-off). Typecheck and run the focused tests frequently; save the full suite for the gate.
- Root cause only — no lint/type suppressions, no swallowed errors, no timing hacks. A workaround closes the task and opens a bug.
- For an authorized tooling Task, compare the target path with the mutation
  allowlist before every mutation. Never edit `_tasks.md` or any other Task
  file. If a required path is absent from the authorization, stop and request
  a revised Spec instead of widening the list.

## 4. Standalone verify — the gate

In standalone execution, evidence comes before terminal status in this order:

1. Run every command in the Task's `## Verification` section verbatim.
2. Run the repository's verify pipeline (`make verify`, or the build/lint/typecheck/test equivalents this repo documents) when the current execution mode requires local completion evidence.
3. Walk Acceptance Criteria one by one: each needs fresh evidence from this session — a command output, a test name that passes, an observed behavior. A green suite is not evidence for a criterion the suite doesn't cover.

For every authorized tooling Task, run a changed-file postflight after the last
edit. Compare the captured baseline with `git status --short` and
`git diff --name-only`, accounting for staged, unstaged, and untracked paths.
Every newly changed path must be either one exact authorized path or the
assigned Task file. If any other path appears, make no further mutation, set
`status: failed`, record the out-of-scope paths in `## Result`, and leave the
worktree intact for recovery.

A narrow verification never supports a broad claim. If any standalone check
fails after honest root-cause attempts: set `status: failed`, record what was
tried in `## Result`, and report — a loud failure the scheduler can retry
beats a quiet "mostly done".

## 5. Standalone record

Append a `## Result` section to the task file: what changed (described by behavior, not file lists), commands run with outcomes, evidence per acceptance criterion, and any follow-ups discovered. Tick the Subtasks and Acceptance Criteria checkboxes that the evidence supports. Set `status: completed`.

Never touch `_tasks.md` — it owns graph topology, not progress.

## 6. Standalone commit

- Commit only when the user has explicitly authorized it.
- Stage only this task's files (`git status --short` first; unrelated changes stay out).
- One commit per task, Conventional Commits format, task id in the body for traceability (`spec: <slug> / task_02`).
- Never push and never open a PR from inside a task — publishing is a separate, explicit action.

## 7. Daemon-assigned handoff

1. Read the assigned Task, PRD, TechSpec, `CONTEXT.md`, active ADRs, and bounded
   context paths before editing.
2. Implement only the assigned slice. Never edit `_tasks.md`, another Task
   file, or a path outside an authorized tooling allowlist.
3. Run focused implementation checks while working when useful. Do not run any
   command from the Task's `## Verification` section.
4. Append or update `## Result` with the implementation and focused-check
   evidence for every acceptance criterion. Do not use Daemon Verification
   evidence that has not run yet.
5. Hand back implementation-ready work without editing Task status or claiming
   `completed`, `failed`, passing Verification, commit readiness, or delivery.
6. Never commit, push, or open a pull request. The Daemon runs declared
   Verification, writes terminal Task status, and owns the Task commit.

## Anti-patterns

- Marking `completed` on "should work" — the status field is read by machines that won't double-check.
- Weakening, skipping, or deleting a failing test to get to green.
- Fixing "one more thing" spotted along the way — that's a follow-up, not scope.
- Re-verifying with stale output from earlier in the session; evidence must postdate the last edit.
- Editing other tasks' files or statuses.
- In a Daemon-assigned turn, running declared Verification, editing Task
  status, claiming a terminal verdict, or committing.
