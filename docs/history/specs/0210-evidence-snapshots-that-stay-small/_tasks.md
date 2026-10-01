---
schema: spec-tasks/v1
spec: 0210-evidence-snapshots-that-stay-small
qa: task_03
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
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | An Evidence Snapshot records one digest per declared input and carries by it |
| task_02 | docs | The guide describes one digest per declared input |
| task_03 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03
