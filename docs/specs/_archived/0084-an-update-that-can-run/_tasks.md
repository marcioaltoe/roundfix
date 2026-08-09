---
schema: spec-tasks/v1
spec: 0084-an-update-that-can-run
qa: task_10
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
    - id: task_09
      file: task_09.md
      needs: [task_03, task_05, task_06]
    - id: task_10
      file: task_10.md
      needs: [task_04, task_08, task_09]
---

# Tasks — A Baseline Command update that runs on the repositories that already exist

| id      | title                                                        | type    | complexity | needs                     |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ------------------------- |
| task_01 | Classify a managed region instead of blocking on it            | backend | high       | —                         |
| task_02 | Name the lines a refresh removes                               | backend | medium     | task_01                   |
| task_03 | Show the maintainer what the refresh replaces                  | backend | medium     | task_02                   |
| task_04 | Prove the update converges on a second run                     | test    | medium     | task_02                   |
| task_05 | Make an unresolved Baseline Profile diagnose itself            | backend | low        | —                         |
| task_06 | Emit the fourteen structural clauses again                     | backend | high       | —                         |
| task_07 | Seat the consultation, outside-evidence, and glossary clauses  | backend | medium     | task_06                   |
| task_08 | State the two obligations where a Task author reads them       | docs    | low        | task_07                   |
| task_09 | Sweep a fleet the Spec did not author                          | test    | high       | task_03, task_05, task_06 |
| task_10 | Run the final QA gate                                          | qa      | high       | task_04, task_08, task_09 |

Waves: 1 → task_01, task_05, task_06 · 2 → task_02, task_07 · 3 → task_03, task_04, task_08 · 4 → task_09 · 5 → task_10
