---
schema: spec-tasks/v1
spec: 0156-a-delivery-loop-that-outlives-the-session
qa: task_06
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
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_07]
    - id: task_07
      file: task_07.md
      needs: [task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The queue store |
| task_02 | backend | The pull request boundary |
| task_03 | backend | The delivery engine |
| task_04 | backend | The command family |
| task_05 | docs | Describe the queue in the shipped skill and the guide |
| task_06 | qa | Run the final QA gate |
| task_07 | backend | A migration ladder that applies from every earlier version |
