---
schema: spec-tasks/v1
spec: 0095-a-verification-that-ran-before-anyone-believed-it
qa: task_09
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
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: []
    - id: task_08
      file: task_08.md
      needs: [task_03]
    - id: task_09
      file: task_09.md
      needs: [task_03, task_06, task_07, task_08]
---

# Tasks — A Verification that ran before anyone believed it

| id      | title                                                        | type    | complexity | needs                              |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ---------------------------------- |
| task_01 | Extract the prober the Daemon already runs                    | backend | medium     | —                                  |
| task_02 | Give the check a tree of its own to run in                    | backend | medium     | —                                  |
| task_03 | Execute every authored command before any Run                 | backend | high       | task_01, task_02                   |
| task_04 | Refuse a reversed exit condition                              | backend | medium     | —                                  |
| task_05 | Refuse a Verification that reaches outside the repository     | backend | medium     | task_04                            |
| task_06 | Restore the vacuity refusal and account what it finds         | backend | medium     | task_05                            |
| task_07 | Let a Task declare the file it creates                        | backend | low        | —                                  |
| task_08 | Write the exit-zero rule where an author reads it             | docs    | low        | task_03                            |
| task_09 | Run the final QA gate                                         | qa      | high       | task_03, task_06, task_07, task_08 |

Waves: 1 → task_01, task_02, task_04, task_07 · 2 → task_03, task_05 · 3 → task_06, task_08 · 4 → task_09
