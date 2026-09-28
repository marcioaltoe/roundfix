---
schema: spec-tasks/v1
spec: 0174-operator-surfaces-that-tell-the-truth
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
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A schema refusal names its remedy and roundfix migrate performs it |
| task_02 | backend | A help token means help only where it is an argument |
| task_03 | backend | Settlement reads the front matter the derived Verification reads |
| task_04 | backend | The governed set covers every bounded file and its contract runs in a gate |
| task_06 | backend | The qa-report accept refusal test pins the new front-matter message |
| task_05 | qa | Run the final QA gate |
