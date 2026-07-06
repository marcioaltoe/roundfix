---
name: write-tasks
description: Decompose a spec's PRD/TechSpec into a dependency-ordered task graph — vertical-slice task files plus a machine-parseable _tasks.md DAG manifest under docs/specs/<slug>/, gated by user approval of the breakdown before any file is written.
argument-hint: "<spec slug or path under docs/specs/>"
metadata:
  category: issue-decomposition
  tags: [issues, workflow, prd, agents]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Write Tasks

Turn `docs/specs/<slug>/_prd.md` (and `_techspec.md` when present) into the executable task graph: one `task_NN.md` per task plus `_tasks.md`, the manifest that `implement-task`, `implement-spec`, and any future daemon schedule from. The artifacts are the contract — a fresh agent session must be able to pick up one task file and build it with no other context than the spec folder.

## Preconditions

`$ARGUMENTS` names the spec. `_prd.md` must exist. If `_techspec.md` is missing on a feature with real architectural surface, call that out and recommend `write-techspec` first; proceed only if the user accepts the trade-off, compensating with deeper codebase exploration before decomposing.

## Ownership rules (what lives where, and why)

- **Dependencies live only in the `_tasks.md` graph.** Never in task files. Two copies of an edge is how graphs drift.
- **Status lives only in each task file's frontmatter.** `_tasks.md` owns topology, not progress — schedulers re-read task files for state.
- **Content references, never duplicates.** Task files point at PRD/TechSpec sections by name; a task that restates the spec goes stale the first time the spec is amended.

## Decomposition rules

- **Vertical slices.** Each task delivers a narrow but complete path through every layer it touches, demoable or verifiable on its own — a tracer bullet, not a layer ("all the schemas" is a wrong task; "expired imports retry and surface their status" is a right one).
- **Prefactoring first.** When a slice needs the ground prepared, make that its own leading task: make the change easy, then make the easy change.
- **Sized for one fresh session.** A task an agent can complete in a single sitting with a fresh context. More than ~7 subtasks or files means split it.
- **Tests embedded, never separated.** Every task's acceptance criteria include its own tests; a trailing "write the tests" task means the earlier tasks were never done.
- **Independently implementable.** Once its `needs` are completed, a task must require no other unfinished work — that's what allows parallel execution across worktrees later.

## Process

### 1. Derive the breakdown

Start from the TechSpec's Build Order (its dependency statements become graph edges). Map every PRD user story and core feature onto at least one task; an uncovered story is a hole to fix now.

### 2. Approval gate — the one mandatory human check

Present the breakdown as a table (`id | title | type | complexity | needs`) plus a one-line rationale for the slicing, and iterate until the user explicitly approves. Granularity and dependency correctness are where human judgment pays most — **write no files before approval.**

### 3. Write the files

From the templates in [references/task-template.md](references/task-template.md), write `_tasks.md` and every `task_NN.md` (numbered from `01` in topological order). Acceptance criteria must be independently verifiable — a criterion nobody can check is a wish, not a criterion. Include a `## Verification` section with the exact commands that prove the task done; `implement-task` runs them verbatim.

Durability applies here too: describe behavior and interfaces, not repo file paths (relative references within the spec folder are fine — the folder moves as a unit).

### 4. Validate the graph

Before reporting, verify mechanically — parse, don't eyeball:

- Every `graph.nodes[].file` exists and every task file has parseable frontmatter with `status: pending`.
- Every `needs` entry names an existing node id; no cycles (a topological order can be printed).
- Every PRD user story appears in some task's References.

Print the wave plan (wave 1 = no needs; wave N = needs met by earlier waves) as the execution preview.

### 5. Report

Reply with the file list and the wave plan.

## Anti-patterns

- Horizontal slicing (a schema task, an API task, a UI task) — every task should cut through the stack.
- Mega-tasks that "do the feature" — if it can't be verified alone, it isn't a task yet.
- Vague criteria ("works correctly", "handles errors") — name the observable behavior.
- Editing the graph and the task files out of sync — regenerate both from the approved breakdown.
- Writing files before the user approves the breakdown.

## References

- [references/task-template.md](references/task-template.md) — templates for `_tasks.md` and `task_NN.md`. Read it before writing any file.
