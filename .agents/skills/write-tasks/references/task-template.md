# Task artifact templates

Two templates: the `_tasks.md` manifest and the per-task `task_NN.md` file. Guidance appears as `<!-- comments -->`; delete every comment from final files.

## `_tasks.md` — the DAG manifest

Machine-parseable. Dependencies live **only** here; status lives **only** in task files. The body table is a human-readable projection of the frontmatter graph — regenerate it whenever the graph changes.

```markdown
---
schema: spec-tasks/v1
spec: <feature-slug>
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
---

# Tasks — <feature name>

| id      | title                            | type    | complexity | needs            |
| ------- | -------------------------------- | ------- | ---------- | ---------------- |
| task_01 | <title>                          | chore   | low        | —                |
| task_02 | <title>                          | backend | medium     | task_01          |
| task_03 | <title>                          | frontend| medium     | task_01          |
| task_04 | <title>                          | backend | high       | task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04
```

## `task_NN.md` — one task

A fresh agent session must be able to build this with no context beyond the spec folder. Reference the PRD/TechSpec by section name; never restate them.

```markdown
---
task: task_02
spec: <feature-slug>
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend # frontend | backend | data | infra | docs | test | chore
complexity: medium # low | medium | high
---

# Task 02: <imperative title in glossary vocabulary>

## Overview

<!-- 2-4 sentences: the slice this task delivers and how it's demoable/verifiable on its own. -->

## Requirements

<!-- Numbered. MUST for the contract, SHOULD for preferences. Behavior and interfaces, not file paths. -->

1. MUST ...
2. SHOULD ...

## Subtasks

<!-- 3-7 checkboxes, WHAT not HOW. More than 7 means the task should have been split. -->

- [ ] ...

## Acceptance Criteria

<!-- Independently verifiable checkboxes — each one checkable by a command, a test, or an observable behavior.
     qa-gate re-validates these against the running app at the end of the spec. -->

- [ ] ...

## Context

<!-- Optional. Include only task-specific path references that help a fresh Agent orient.
     Use at most 50 unique clean repository-relative paths, labeled exactly as:
     - instruction: `.agents/skills/<skill>/SKILL.md`
     - interface: `internal/package/file.go`
     The Daemon includes these paths in the 200-path Spec Context Bundle before prior changed files. -->

- instruction: `<path>`
- interface: `<path>`

## Verification

<!-- The exact commands that prove this task done, with what to expect from each.
     The Daemon runs these verbatim after the Agent turn and will not settle the task completed until they pass. -->

- `<command>` — expected: ...

## References

<!-- PRD/TechSpec sections and ADRs this task implements, by name:
     `_prd.md` → User Stories 3-4; `_techspec.md` → Build Order 2, Interfaces: ImportScheduler; ADR-0012. -->
```

## Result section (appended by `implement-task`, not by `write-tasks`)

On completion, `implement-task` appends a `## Result` section: what changed (described by behavior), commands run with outcomes, evidence per acceptance criterion, and follow-ups. `write-tasks` never creates this section.
