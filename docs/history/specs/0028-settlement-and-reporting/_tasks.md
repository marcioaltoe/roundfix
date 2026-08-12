---
schema: spec-tasks/v1
spec: 0028-settlement-and-reporting
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_07
      file: task_07.md
      needs: []
    - id: task_08
      file: task_08.md
      needs: [task_02, task_03, task_04, task_05, task_06, task_07]
---

# Tasks — Settlement and Reporting Robustness

| id      | title                                                       | type    | complexity | needs                                                    |
| ------- | ----------------------------------------------------------- | ------- | ---------- | -------------------------------------------------------- |
| task_01 | Record the Run owner PID and prove process liveness         | backend | medium     | —                                                         |
| task_02 | Reclaim orphaned Active-Run locks on proven owner death     | backend | high       | task_01                                                   |
| task_03 | Normalize task status synonyms before validation            | backend | low        | —                                                         |
| task_04 | Warn on no-op Task commits                                  | backend | low        | —                                                         |
| task_05 | Return per-Task outcomes and report failure reasons         | backend | medium     | —                                                         |
| task_06 | List settle commit paths and warn on shared worktrees       | backend | medium     | —                                                         |
| task_07 | Diagnose missing adapter binaries at preflight              | backend | medium     | —                                                         |
| task_08 | Sync authoring guidance and the Roundfix Skill              | docs    | medium     | task_02, task_03, task_04, task_05, task_06, task_07      |

Waves: 1 → task_01, task_03, task_04, task_05, task_06, task_07 · 2 → task_02 · 3 → task_08
