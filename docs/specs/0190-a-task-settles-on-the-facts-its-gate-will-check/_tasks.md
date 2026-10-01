---
schema: spec-tasks/v1
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
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
      needs: [task_01, task_02, task_03]
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
| task_01 | backend | A Task of a gated graph runs the repository Verification when it settles |
| task_02 | backend | A Verification attempt audits the commit the Daemon is about to create |
| task_03 | backend | A Task is refused the Spec Consistency findings it introduced |
| task_04 | docs | The skills, the commands guide and the glossary describe Settlement Checks |
| task_05 | qa | Run the final QA gate |
| task_06 | chore | The two skills task_04 changed declare new versions |
