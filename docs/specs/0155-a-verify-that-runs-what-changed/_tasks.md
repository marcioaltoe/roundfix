---
schema: spec-tasks/v1
spec: 0155-a-verify-that-runs-what-changed
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
      needs: [task_01, task_02]
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
| task_01 | backend | The selector and its classification |
| task_02 | test | A partition that covers every test once |
| task_03 | infra | The selective target and the Daemon's gate |
| task_04 | infra | Pull requests run the selective gate |
| task_05 | qa | Run the final QA gate |
