---
schema: spec-tasks/v1
spec: 0037-terminal-outcome-integrity
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Tasks — Terminal outcome integrity

| id      | title                                                    | type    | complexity | needs                                  |
| ------- | -------------------------------------------------------- | ------- | ---------- | -------------------------------------- |
| task_01 | Guard terminal completion and reconciliation             | data    | high       | —                                      |
| task_02 | Target only registered Agent Sessions during cleanup     | backend | medium     | task_01                                |
| task_03 | Prove owner exit before Force Stop completion            | backend | high       | task_01, task_02                       |
| task_04 | Interrupt every Review Source wait on Stop Request       | backend | medium     | task_01                                |
| task_05 | Publish only the winning terminal outcome                | backend | high       | task_01, task_02, task_03, task_04     |
| task_06 | Align terminal-outcome operator guidance                 | docs    | medium     | task_03, task_04, task_05              |
| task_07 | Align the protected Roundfix Skill pair                  | docs    | low        | task_06                                |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05 · 5 → task_06 · 6 → task_07
