---
schema: spec-tasks/v1
spec: 0067-derived-artifact-regeneration-boundary
qa: task_05
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
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Derived artifact regeneration boundary

| id      | title                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Declare one owner per derived path and prove it exhaustive | backend | high | —          |
| task_02 | Prove every declared step actually rewrites its artifacts | test | medium | task_01          |
| task_03 | Make the sanctioned command cover what it claims    | infra   | medium     | task_01          |
| task_04 | Say what a human must do when the command cannot    | backend | low        | task_02, task_03 |
| task_05 | Run the final QA gate                               | qa      | medium     | task_04          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05

task_01 leads because exhaustiveness is what stops the next occurrence. Three
artifacts have already inherited this ambiguity, and a fourth instance happened
on 2026-08-05 while authoring the Spec queue: the flag was guessed from a test
name and the guess was wrong.

task_03 is the one authorized tooling Task, bounded to exactly `Makefile` under
the 2026-08-02 grant, which names 0067 for "the regeneration step list and
derived path scan".

The QA gate is authored as task_05 and depends on task_04, the graph's only
non-QA leaf.
