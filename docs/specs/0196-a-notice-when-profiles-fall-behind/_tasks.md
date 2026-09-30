---
schema: spec-tasks/v1
spec: 0196-a-notice-when-profiles-fall-behind
qa: task_05
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
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A profile can carry a Profile Deviation |
| task_02 | backend | `profiles check` compares the configuration with the recommendation |
| task_03 | backend | `profiles check --apply` adopts the recommendation |
| task_04 | backend | Doctor and `upgrade` carry the notice, and nothing else reaches the comparison |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
