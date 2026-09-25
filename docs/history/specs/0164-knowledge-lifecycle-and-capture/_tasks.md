---
schema: spec-tasks/v1
spec: 0164-knowledge-lifecycle-and-capture
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | Supersede ADR-0123 with ADR-0163 |
| task_02 | backend | Review retirement reads recorded evidence |
| task_03 | backend | A terminal record closes without a Spec |
| task_04 | docs | Every retired family has one documented home |
| task_05 | docs | Capture names its publication owner |
| task_06 | qa | Run the final QA gate |
