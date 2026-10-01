---
schema: spec-tasks/v1
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
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
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The prompt carries the Delivery Conventions and a finding parks only when its anchor is in the diff |
| task_02 | backend | A sealed validator dismisses a convention restatement only inside its region |
| task_03 | backend | A second review reads the delta, and a third closes only on dispositions |
| task_04 | backend | Round 2 continues the round-1 session, and the skill describes the reviewer |
| task_05 | qa | Run the final QA gate |
| task_06 | backend | A ceiling record names no reviewer selection |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_05
