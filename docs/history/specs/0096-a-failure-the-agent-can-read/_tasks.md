---
schema: spec-tasks/v1
spec: 0096-a-failure-the-agent-can-read
qa: task_08
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
      needs: []
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_07
      file: task_07.md
      needs: []
    - id: task_09
      file: task_09.md
      needs: [task_03]
    - id: task_08
      file: task_08.md
      needs: [task_03, task_04, task_05, task_06, task_07, task_09]
---

# Tasks — A failure the Agent can read

| id      | title                                                    | type    | complexity | needs                                        |
| ------- | -------------------------------------------------------- | ------- | ---------- | -------------------------------------------- |
| task_01 | Say that the command produced nothing                      | backend | medium     | —                                            |
| task_02 | A signature that survives a timestamp                      | backend | medium     | —                                            |
| task_03 | Name a failure the loop has already seen                   | backend | high       | task_01, task_02                             |
| task_04 | Let the vacuity refusal name its offenders                  | backend | medium     | —                                            |
| task_05 | Say what the run budget bounds where it is set              | backend | low        | —                                            |
| task_06 | Name the surface a Task file was read from                  | backend | low        | —                                            |
| task_07 | Write the ceiling's exits into the authoring contract       | docs    | low        | —                                            |
| task_09 | Stop the failure cause being shadowed away                   | backend | low        | task_03                                      |
| task_08 | Run the final QA gate                                       | qa      | high       | task_03, task_04, task_05, task_06, task_07, task_09 |

Wave plan: `1 → task_01, task_02, task_04, task_05, task_06, task_07 · 2 → task_03 · 3 → task_08`.

task_07 is the only tooling Task and travels alone: its commit may touch one
authorized skill file plus its own Task file, and a commit mixing that with Go
fails the changed-path audit, as one did during Spec 0095. task_03 is the only
join and lands after the prompt and the signature it composes.
