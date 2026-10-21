---
schema: spec-tasks/v1
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
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
      needs: [task_02]
---

# Tasks — The queue's archive subject and pull request title fit 100 characters

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The deliver guide, the skill references and the glossary give the archive subject and the title | docs    | low        | —       |
| task_02 | The archive subject and the pull request title keep the slug while it fits                    | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                         | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 writes the `deliver` guide, the Roundfix `deliver` reference and the
`archive-spec` skill for the surface task_02 changes, so task_02 follows it in
one dependency chain. task_02 changes the code and proves it in
`internal/cli`, the package that consumes every contract it changes, so it
needs no separate Integration Checkpoint. The gate follows both. No other Spec
is required.
