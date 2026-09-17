---
schema: spec-tasks/v1
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
qa: task_04
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
    - id: task_05
      file: task_05.md
      needs: [task_03]
    - id: task_04
      file: task_04.md
      needs: [task_05]
---

# Tasks — A commit that carries the work and nothing else

| id      | title                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Ask the index before refusing a file for its mode    | backend | medium     | —       |
| task_02 | Make a lost output fail the Task                     | backend | medium     | task_01 |
| task_03 | Keep the Verification's writes out of the QA commit  | backend | medium     | task_02 |
| task_05 | Restore the permission-variant coverage the rule kept true | test | low | task_03 |
| task_04 | Run the final QA gate                                | qa      | high       | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_05 · 5 → task_04
