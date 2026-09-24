---
schema: spec-tasks/v1
spec: 0158-daemon-verification-and-access-readiness
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
      needs: [task_01, task_02, task_03, task_04, task_06, task_07]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Collect every independent Verification failure |
| task_02 | backend | Let a named Task repair a red repository gate |
| task_03 | backend | Prove access policy in readiness |
| task_04 | docs | Record the decisions |
| task_06 | backend | A repair entry that the gate still closes |
| task_07 | backend | Keep every failure and name a degraded access |
| task_05 | qa | Run the final QA gate |
