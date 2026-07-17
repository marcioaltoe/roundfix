---
schema: spec-tasks/v1
spec: 0032-deterministic-agent-session-cancellation
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

# Tasks — Deterministic Agent Session cancellation

| id      | title                                       | type    | complexity | needs   |
| ------- | ------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Create deterministic cancellation controls  | test    | medium     | —       |
| task_02 | Enforce deterministic forced-close ordering | backend | high       | task_01 |
| task_03 | Prove cooperative cancellation stability    | test    | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03
