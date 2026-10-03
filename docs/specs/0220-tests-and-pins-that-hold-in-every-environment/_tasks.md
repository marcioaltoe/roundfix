---
schema: spec-tasks/v1
spec: 0220-tests-and-pins-that-hold-in-every-environment
qa: task_03
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
      needs: [task_01, task_02]
---

# Tasks — Tests and pins that hold in every environment

| id      | title                                                                    | type  | complexity | needs            |
| ------- | ------------------------------------------------------------------------ | ----- | ---------- | ---------------- |
| task_01 | A fixture's process group is proved ended by its live members            | test  | medium     | —                |
| task_02 | Restore the trailing skills and hold them to the lock and the snapshot   | chore | medium     | —                |
| task_03 | Run the final QA gate                                                    | qa    | high       | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03
