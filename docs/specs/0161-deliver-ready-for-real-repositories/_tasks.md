---
schema: spec-tasks/v1
spec: 0161-deliver-ready-for-real-repositories
qa: task_04
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
      needs: [task_01, task_02, task_03, task_05, task_06]
    - id: task_05
      file: task_05.md
      needs: [task_03]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Wait for checks that have not reported |
| task_02 | backend | Untracked, per-delivery item branches |
| task_03 | backend | Clean parks, crash-safe archive, stale owners |
| task_05 | backend | Accept the archive commit the real command makes |
| task_06 | backend | Park only what the item touched |
| task_04 | qa | Run the final QA gate |
