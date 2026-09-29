---
schema: spec-tasks/v1
spec: 0177-runs-that-fit-their-budget-and-park-honestly
qa: task_06
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_07]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Instruction paths never make two Tasks collide |
| task_02 | backend | The Implement Run Budget renews at each Task settlement |
| task_03 | backend | A budget-stopped Run parks as a Run outcome |
| task_04 | backend | Carry-forward stages commits without repository hooks |
| task_05 | backend | Phrase checks survive wrapping and name what they missed |
| task_07 | backend | The QA settlement renews the budget before its report commit, and wrap-fix remediations are shell-safe |
| task_06 | qa | Run the final QA gate |
