---
schema: spec-tasks/v1
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A queued Spec is revalidated on its starting main before its first Run |
| task_02 | backend | A Delivery Queue records its limits and enforces them |
| task_03 | backend | `roundfix deliver plan` shows what is approved to run, and start refuses what is not |
| task_04 | backend | `deliver start` takes explicit limits and `deliver status` presents one Pending Question |
| task_05 | docs | The implement-spec entry point hands implementation to Roundfix |
| task_06 | qa | Run the final QA gate |
