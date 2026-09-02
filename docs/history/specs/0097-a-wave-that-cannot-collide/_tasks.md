---
schema: spec-tasks/v1
spec: 0097-a-wave-that-cannot-collide
qa: task_06
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — A wave that cannot collide

| id      | title                                          | type    | complexity | needs   |
| ------- | ---------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The collision rule over a Task Graph           | backend | high       | —       |
| task_02 | The checker reports a collision at authoring   | backend | medium     | task_01 |
| task_03 | The Run refuses before it dispatches           | backend | medium     | task_02 |
| task_04 | Bootstrap serialized across sibling worktrees  | backend | high       | task_03 |
| task_05 | A worktree failure in Roundfix's own words     | backend | medium     | task_04 |
| task_06 | QA gate                                        | qa      | medium     | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_06

The chain is serial, and this Spec is the one graph where that ordering is the
subject rather than a precaution. task_02 and task_03 are the two callers of one
rule, so they follow it. task_04 and task_05 both rewrite the worktree package.

This graph must not collide with itself. Every Task's Verification names files
in a package no sibling Task's Verification names, and the `needs` chain
serializes them regardless — which is the fix the Spec's own refusal will tell a
Supervisor to apply.

Each Task's work, references, and Verification live in its own `task_NN.md`.
