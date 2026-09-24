---
schema: spec-tasks/v1
spec: 0163-baseline-decisions-and-regeneration
qa: task_07
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
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_06]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A mode change keeps the HTTP decision |
| task_02 | backend | Greenfield refuses before an unreachable step |
| task_03 | backend | Skill regeneration declares its outputs |
| task_04 | backend | Reconcile the lock on proven absence |
| task_05 | backend | The lock reconciliation command |
| task_06 | docs | Shipped skills describe the changes |
| task_07 | qa | Run the final QA gate |
