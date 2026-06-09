# Tasks playbook

Task lists, individual task files, and vertical-slice tracker issues. A task is a prompt for its executor — human or agent: a clear problem, complete acceptance criteria, explicit boundaries, and the minimum context needed. The specification is the program; vague input produces vague work.

## Decomposition principles

- **Slice vertically.** Each task cuts a thin, complete path through every layer (schema + API + UI + tests) and is demoable on its own. Never "build all the models, then all the APIs".
- **Small and independent.** Each task is implementable in one focused session once its dependencies are met, leaves the system in a working state, and embeds its own tests — never a separate "write tests" task.
- **Map the dependency graph.** Name what blocks what; no circular dependencies. Publish or order blockers first. Front-load the riskiest task so it fails fast.
- **What, not how.** Describe what must be accomplished and how to verify it. Implementation detail lives in the spec — reference its sections by name instead of duplicating them.
- **Specific test cases.** "POST /jobs/done with unknown job ID returns 404" — not "test the happy path".
- **Mark the human gates.** Label each slice HITL (needs a human decision or design review) or AFK (implementable without one). Prefer AFK; isolate the HITL parts into their own slices.

A task is too large when any of these hold — split it:

- The title needs an "and".
- It touches more than ~5–7 files or spans two independent subsystems.
- It needs more than 7 subtasks or more than ~10 acceptance criteria.
- It would take more than one focused session.

## Master task list

Common convention: a `_tasks.md` beside the PRD/spec (leading underscore marks meta documents), with individual files named `task_01.md`, `task_02.md`, ... zero-padded.

```markdown
# <feature> — Task List

## Tasks

| #  | Title   | Status  | Complexity | Dependencies |
| -- | ------- | ------- | ---------- | ------------ |
| 01 | <title> | pending | medium     | —            |
| 02 | <title> | pending | low        | task_01      |
```

## Task file skeleton

YAML frontmatter first, then the body. Field values:

- `status`: `pending` | `in_progress` | `completed`
- `type`: `frontend` | `backend` | `docs` | `test` | `infra` | `refactor` | `chore` | `bugfix` (or the project's own taxonomy)
- `complexity`: `low` (single file, straightforward logic) | `medium` (2–4 files, maybe one new interface) | `high` (5+ files, new subsystem, multiple integration points or concurrency) | `critical` (cross-cutting, high regression risk, needs coordination)
- `dependencies`: list of task file names (`[task_01, task_02]`) or `[]`

```markdown
---
status: pending
title: <task title — must match the H1>
type: backend
complexity: medium
dependencies: [task_01]
---

# Task N: <title>

## Overview

<2–3 sentences: what this task accomplishes and why it matters>

<requirements>
1. MUST <specific, checkable requirement>
2. SHOULD <specific, checkable requirement>
</requirements>

## Subtasks

- [ ] N.1 <what, not how>
- [ ] N.2 <what, not how>

## Implementation Details

<reference the spec section by name for patterns; do not restate it>

### Relevant Files

- `<path>` — <why it matters here>

### Dependent Files

- `<path>` — <how it is affected>

### Related ADRs

- [ADR-NNN: <title>](../adrs/adr-NNN.md) — <relevance>  <!-- omit section if none apply -->

## Deliverables

- <concrete output>
- Unit tests for <behavior> (required)
- Integration tests for <flow> (required)

## Tests

Unit tests:
- [ ] <specific case: exact input, condition, expected result>

Integration tests:
- [ ] <specific case>

All tests must pass; meet the project's coverage target.

## Success Criteria

- All tests passing
- <measurable outcome>
```

Gates: 3–7 subtasks; tests embedded with specific case descriptions; every requirement checkable; spec referenced, not duplicated.

## Vertical-slice tracker issue

For publishing decomposed work to an issue tracker. Title is action-forward: component, then behavior. Avoid file paths and code — they go stale. Exception: a snippet that encodes a decision more precisely than prose (state machine, reducer, schema, type shape); trim it to the decision-rich parts.

```markdown
## Parent

<reference to the parent issue — omit if none>

## What to build

<concise end-to-end description of the slice's behavior, not layer-by-layer implementation>

## Acceptance criteria

- [ ] <observable criterion>
- [ ] <observable criterion>

## Blocked by

<reference to the blocking issue, or "None — can start immediately">
```

Publish in dependency order so "Blocked by" can cite real issue identifiers. Use the project's domain vocabulary, and respect existing ADRs in the area being touched. Add scope constraints when needed to prevent yak-shaving ("no new dependencies", "keep p95 under 200ms").

## Checkpoints

For plans longer than a few tasks, insert an explicit checkpoint every 2–3 tasks:

```markdown
## Checkpoint: after tasks 1–3

- [ ] All tests pass
- [ ] Application builds without errors
- [ ] Core flow works end to end
- [ ] Review with a human before proceeding
```

## Quality gates

1. Every task independently implementable once dependencies are met; no circular dependencies.
2. Every task slices vertically and leaves the system working.
3. Tests embedded in every task, with named specific cases.
4. Spec and PRD referenced by section, never restated.
5. Sizing red flags clear on every task (no "and" titles, ≤7 subtasks, one subsystem).
6. Dependency order published blockers-first; riskiest work early.
