---
schema: spec-tasks/v1
spec: 0030-context-driven-agent-instructions
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
      needs: [task_03]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04, task_05]
---

# Tasks — Context-driven agent instructions

| id      | title                                                        | type    | complexity | needs                      |
| ------- | ------------------------------------------------------------ | ------- | ---------- | -------------------------- |
| task_01 | Define portable instruction profiles and module contracts     | backend | high       | —                          |
| task_02 | Build the read-only agent-instruction audit engine             | backend | high       | task_01                    |
| task_03 | Add safe apply and durable decision migrations                 | backend | high       | task_02                    |
| task_04 | Validate canonical skill setups and installed skills           | backend | medium     | task_02                    |
| task_05 | Add optional Secondbrain guidance and setup orchestration       | docs    | medium     | task_03                    |
| task_06 | Gate and synchronize the complete portable skill               | test    | medium     | task_03, task_04, task_05  |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05 · 5 → task_06
