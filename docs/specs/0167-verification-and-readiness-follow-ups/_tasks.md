---
schema: spec-tasks/v1
spec: 0167-verification-and-readiness-follow-ups
qa: task_03
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
      needs: [task_01, task_02, task_04]
    - id: task_04
      file: task_04.md
      needs: [task_02]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The retry's verdict wins; completed repairs are history |
| task_02 | backend | Degraded access is printed |
| task_04 | backend | A quiet project init and a budget test without a race |
| task_03 | qa | Run the final QA gate |
