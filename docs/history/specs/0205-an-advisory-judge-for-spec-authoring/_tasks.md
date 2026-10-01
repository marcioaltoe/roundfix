---
schema: spec-tasks/v1
spec: 0205-an-advisory-judge-for-spec-authoring
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

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The judge plans a Spec's judgments from its own artifacts, as measured |
| task_02 | backend | The judge asks, records every call and stops at the monthly ceiling |
| task_03 | backend | `roundfix spec judge` prints what the judge raises |
| task_04 | docs | The authoring skills run the judge and answer what it raises |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
