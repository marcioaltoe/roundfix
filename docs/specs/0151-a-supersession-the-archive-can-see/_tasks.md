---
schema: spec-tasks/v1
spec: 0151-a-supersession-the-archive-can-see
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
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_03]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03, task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The supersession record and the command that writes it |
| task_02 | backend | Archive's second accepted proof |
| task_03 | docs | Describe both in the shipped skill and the guide |
| task_05 | backend | A deliverer check that checks what it claims |
| task_04 | qa | Run the final QA gate |
