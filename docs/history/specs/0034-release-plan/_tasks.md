---
schema: spec-tasks/v1
spec: 0034-release-plan
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Release Plan

| id      | title                                             | type    | complexity | needs   |
| ------- | ------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Build the Release Plan domain and version calculator | backend | medium     | —       |
| task_02 | Classify commits and maintenance-only changes     | backend | high       | task_01 |
| task_03 | Resolve committed release ranges through local Git | backend | medium     | task_02 |
| task_04 | Expose the read-only Release Plan CLI contracts   | backend | high       | task_03 |
| task_05 | Require Release Plan in release guidance          | docs    | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
