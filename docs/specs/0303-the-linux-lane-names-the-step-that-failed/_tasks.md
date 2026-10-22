---
schema: spec-tasks/v1
spec: 0303-the-linux-lane-names-the-step-that-failed
qa: task_03
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
---

# Tasks — The Linux Lane names the step that failed

| id      | title                                                                              | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The deliver guide, the skill reference and the glossary name a failed go list       | docs    | low        | —       |
| task_02 | The lane script marks a failed go list and the lane names it with its Park Class    | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                              | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 writes the `deliver` guide, the Roundfix `deliver` reference and the
glossary for the outcome task_02 changes, so task_02 follows it in one
dependency chain. task_02 changes the lane script and its classifier and
proves them in `internal/cli` and `internal/delivery`, the packages that
consume every contract it changes, so it needs no separate Integration
Checkpoint. The gate follows both. No other Spec is required.
