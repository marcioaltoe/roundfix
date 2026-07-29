---
schema: spec-tasks/v1
spec: 0054-tooling-task-and-verification-hygiene
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
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — Tooling task and verification environment hygiene

| id      | title                                                        | type    | complexity | needs                            |
| ------- | ------------------------------------------------------------ | ------- | ---------- | -------------------------------- |
| task_01 | Regenerate every derived digest from its canonical source     | test    | high       | —                                |
| task_02 | Ship the sanctioned regeneration target and a portable cache  | infra   | medium     | task_01                          |
| task_03 | Refuse executable files in Daemon commits                     | backend | medium     | —                                |
| task_04 | Prove the repository is green before a gate-bound Task starts  | backend | high       | —                                |
| task_05 | Document the choreography and the standing regeneration policy | docs    | medium     | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01, task_03, task_04 · 2 → task_02 · 3 → task_05
