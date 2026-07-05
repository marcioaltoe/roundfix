---
schema: spec-tasks/v1
spec: 0006-acpx-run-robustness
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — ACPX Run Robustness

| id      | title                                                            | type    | complexity | needs                     |
| ------- | ---------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Result-over-exit classification with anomaly journaling          | backend | medium     | —                         |
| task_02 | Settle Command for failed-but-done Tasks                         | backend | high       | —                         |
| task_03 | acpx message-buffer mitigation and upstream evidence             | infra   | low        | —                         |
| task_04 | Docs and skill sync                                              | docs    | low        | task_01, task_02, task_03 |

Waves: 1 → task_01, task_02, task_03 · 2 → task_04
