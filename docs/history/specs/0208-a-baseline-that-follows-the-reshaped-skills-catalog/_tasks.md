---
schema: spec-tasks/v1
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — A Baseline that follows the reshaped skills catalog

| id      | title                                                                                | type    | complexity | needs                              |
| ------- | ------------------------------------------------------------------------------------ | ------- | ---------- | ---------------------------------- |
| task_01 | The setup snapshots take their upstream names, and the removed skills leave the catalog | backend | high       | —                                  |
| task_02 | External triage states its rules as clauses with force                               | backend | medium     | task_01                            |
| task_03 | Every managed repository requires the TypeSafe and README skills                     | backend | medium     | task_02                            |
| task_04 | This repository drops its `review` lock entry and tree                               | chore   | low        | task_03                            |
| task_05 | Run the final QA gate                                                                | qa      | high       | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
