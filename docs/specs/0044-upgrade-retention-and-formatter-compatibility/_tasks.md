---
schema: spec-tasks/v1
spec: 0044-upgrade-retention-and-formatter-compatibility
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
      needs: [task_03]
    - id: task_06
      file: task_06.md
      needs: [task_03]
    - id: task_07
      file: task_07.md
      needs: [task_03]
    - id: task_08
      file: task_08.md
      needs: [task_04, task_05, task_06, task_07]
    - id: task_09
      file: task_09.md
      needs: [task_08]
---

# Tasks — Context-Driven Baseline upgrade retention and formatter compatibility

| id      | title                                          | type    | complexity | needs                                  |
| ------- | ---------------------------------------------- | ------- | ---------- | -------------------------------------- |
| task_01 | Establish upgrade compatibility contracts      | backend | high       | —                                      |
| task_02 | Restore operational workflow clauses           | docs    | high       | task_01                                |
| task_03 | Complete portable hard-rule coverage           | docs    | high       | task_02                                |
| task_04 | Render one dispatch entry per installed skill  | backend | high       | task_03                                |
| task_05 | Enforce the Upgrade Retention Contract         | backend | high       | task_03                                |
| task_06 | Create the Repository-Owned Extension          | backend | medium     | task_03                                |
| task_07 | Report uncovered instruction delegation       | backend | medium     | task_03                                |
| task_08 | Guarantee Formatter-Stable Output              | backend | high       | task_04, task_05, task_06, task_07     |
| task_09 | Synchronize the shipped setup contract         | docs    | medium     | task_08                                |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04, task_05, task_06, task_07 · 5 → task_08 · 6 → task_09
