---
schema: spec-tasks/v1
spec: 0236-baseline-update-with-a-repository-profile
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
      needs: [task_02]
---

# Tasks — Baseline update with a repository profile

| id      | title                                                                                     | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The skill readers resolve a repository profile as baseline update does                    | backend | medium     | —       |
| task_02 | Baseline update, Doctor and skill restore work end to end with a repository profile       | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                     | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02's regression tests run the real commands over task_01's loader, so
they follow it. The gate follows the only leaf.
