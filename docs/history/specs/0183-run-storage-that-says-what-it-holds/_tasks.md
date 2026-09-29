---
schema: spec-tasks/v1
spec: 0183-run-storage-that-says-what-it-holds
qa: task_05
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_02]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A retention prune reports only what it reclaimed |
| task_02 | backend | The retained-Run count ignores a vanished checkout |
| task_03 | backend | Doctor reports reclaimable Run storage |
| task_04 | backend | `gc sanitize` recognizes a pre-key default Artifact Root |
| task_05 | qa | Run the final QA gate |
| task_06 | backend | Only a missing checkout counts as vanished |
