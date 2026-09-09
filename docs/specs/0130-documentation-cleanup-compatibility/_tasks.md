---
schema: spec-tasks/v1
spec: 0130-documentation-cleanup-compatibility
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — Documentation cleanup compatibility

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Preserve authorization and regeneration after cleanup | test | high | — |
| task_02 | Verify the cleanup compatibility boundary | qa | medium | task_01 |

Waves: 1 → task_01 · 2 → task_02.
