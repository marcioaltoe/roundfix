---
schema: spec-tasks/v1
spec: 0173-a-delivery-queue-that-recovers
qa: task_07
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
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_06]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Carry-forward follows the Run's integration order |
| task_02 | backend | Project Config is committed on authority or refused aloud |
| task_03 | backend | The engine can retry a parked item |
| task_04 | backend | The delivery workflow recovers an item's Run |
| task_05 | backend | `roundfix deliver retry` returns a parked item to its owner |
| task_06 | backend | A timed-out profile proof is retried once and called temporary |
| task_07 | qa | Run the final QA gate |
