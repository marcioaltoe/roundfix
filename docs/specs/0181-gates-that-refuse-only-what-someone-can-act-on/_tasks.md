---
schema: spec-tasks/v1
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
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
| task_01 | backend | Every Task commit records the paths the Task changed without declaring them |
| task_02 | backend | The pre-PR Pull Request row never decides a qualifying partial |
| task_03 | backend | A related-ADR gap opens only for ADRs that predate the Spec |
| task_04 | docs | The skills, the guides and the glossary describe the three gates as they now behave |
| task_05 | qa | Run the final QA gate |
| task_06 | backend | Citation checks read only what the Spec's authors wrote |
