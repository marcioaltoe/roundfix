---
schema: spec-tasks/v1
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Every contract reads a document's entry file with its companion files |
| task_02 | backend | The Roundfix skill is an entry file plus one reference per command family |
| task_03 | docs | The command reference is an index plus one file per command |
| task_04 | docs | A Task declares the one command file it changes |
| task_05 | qa | Run the final QA gate |
| task_06 | test | The workflow contract reads a skill with its references |
