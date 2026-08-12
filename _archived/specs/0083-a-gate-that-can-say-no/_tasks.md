---
schema: spec-tasks/v1
spec: 0083-a-gate-that-can-say-no
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_01]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_03, task_04, task_05, task_06]
---

# Tasks — A gate that can say no

| id      | title                                                      | type  | complexity | needs                                    |
| ------- | ---------------------------------------------------------- | ----- | ---------- | ---------------------------------------- |
| task_01 | Take the wrapper off the gate and prove the gate can fail   | infra | high       | —                                        |
| task_02 | Give the coverage contract a semantic owner and make it true | test  | medium     | task_01                                  |
| task_03 | Stop the archived corpus counter from gating authoring      | test  | medium     | task_01                                  |
| task_04 | Measure the corpus sweep in a unit the machine cannot move  | test  | medium     | task_01, task_03                         |
| task_05 | Make the cancellation test wait on its milestone            | test  | medium     | task_01                                  |
| task_06 | Make the capacity test wait on its milestone                | test  | medium     | task_01                                  |
| task_07 | Run the final QA gate                                       | qa    | high       | task_02, task_03, task_04, task_05, task_06 |

Waves: 1 → task_01 · 2 → task_02, task_03, task_05, task_06 · 3 → task_04 · 4 → task_07

task_04 follows task_03 rather than joining it: both change
`internal/speccheck/constraints_characterization_test.go`, and two tasks editing
one file in the same wave integrate into a conflict.
