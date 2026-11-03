---
schema: spec-tasks/v1
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
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
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — The follow-up names what waits on someone else and where each item went

| id      | title                                                                                          | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The inbox source and the delivery plan skip hidden files and folders                           | backend | low        | —       |
| task_02 | The follow-up holds Waiting Items back, lists them with --all and reports a malformed declaration | backend | high       | task_01 |
| task_03 | The follow-up mode passes --all, writes the declaration and names each item's destination      | docs    | medium     | task_02 |
| task_04 | Run the final QA gate                                                                          | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 needs task_01 because both edit the follow-up guide and both copies of
its reference, `internal/cli/follow_up_test.go`, `CONTEXT.md`, the Roundfix
`SKILL.md` pair, the owned-skill version record and the Behavior Surface
Record. task_03 needs task_02 because it teaches the flag and the declaration
task_02 builds and edits the `SKILL.md` pair and the version record. Each Task
runs the suites that consume its changed contract in its own Verification, so
no separate Integration Checkpoint is needed. No other Spec is required.
