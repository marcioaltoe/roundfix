---
schema: spec-tasks/v1
spec: 0042-verification-capacity-and-daemon-task-settlement
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_06]
---

# Tasks — Verification Capacity and Daemon Task Settlement

| id      | title                                                          | type     | complexity | needs                                             |
| ------- | -------------------------------------------------------------- | -------- | ---------- | ------------------------------------------------- |
| task_01 | Configure independent Verification Capacity                    | backend  | high       | —                                                 |
| task_02 | Make the Daemon own Implement Task status                      | backend  | high       | task_01                                           |
| task_03 | Queue Task Verification and journal Waiting for Verification   | backend  | high       | task_01, task_02                                  |
| task_04 | Retry Temporary Verification Failure under exclusive capacity  | backend  | high       | task_03                                           |
| task_05 | Render per-Task Verification phases in the Live Run View       | frontend | medium     | task_01, task_02, task_03, task_04                |
| task_06 | Prove the integrated capacity and settlement contract          | test     | high       | task_01, task_02, task_03, task_04, task_05       |
| task_07 | Align operator docs and shipped Agent Skills                   | docs     | medium     | task_01, task_02, task_03, task_04, task_05, task_06 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_06 · 7 → task_07
