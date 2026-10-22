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
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_03
      file: task_03.md
      needs: [task_04]
---

# Tasks — The queue's archive subject and pull request title fit 100 characters

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The deliver guide, the skill references and the glossary give the archive subject and the title | docs    | low        | —       |
| task_02 | The archive subject and the pull request title keep the slug while it fits                    | backend | medium     | task_01 |
| task_04 | The title test tells UTF-16 units from bytes                                                  | test    | low        | task_02 |
| task_03 | Run the final QA gate                                                                         | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_04 · 4 → task_03

task_01 writes the `deliver` guide, the Roundfix `deliver` reference and the
`archive-spec` skill for the surface task_02 changes, so task_02 follows it in
one dependency chain. task_02 changes the code and proves it in
`internal/cli`, the package that consumes every contract it changes, so it
needs no separate Integration Checkpoint. task_04 is corrective Task 1 of 2:
it answers finding F-01 of the QA Report of 2026-10-21, and the gate follows
it. No other Spec is required.
