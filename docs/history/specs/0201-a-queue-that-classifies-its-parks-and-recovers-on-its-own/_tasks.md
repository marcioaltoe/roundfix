---
schema: spec-tasks/v1
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — A queue that classifies its parks and recovers on its own

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Every park names its class and next command, and a check that failed outside the change is re-run once | backend | high | — |
| task_02 | A Spec names its prerequisites and the owner waits for them | backend | high | task_01 |
| task_03 | A retry resumes an item the operator archived after an environment-only partial | backend | high | task_02 |
| task_04 | A conflicting Pull Request parks at once and a derived conflict is resolved by regeneration | backend | high | task_03 |
| task_05 | Run the final QA gate | qa | high | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
