---
schema: spec-tasks/v1
spec: 0033-console-log-tool-summary-deduplication
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — Console Log tool-summary deduplication

| id      | title                                                     | type    | complexity | needs   |
| ------- | --------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Render one summary per tool call in plain-text output      | backend | high       | —       |
| task_02 | Preserve lossless Run evidence across deduplicated output | backend | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02
