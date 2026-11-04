---
schema: spec-tasks/v1
spec: 0347-read-only-run-commands-work-outside-a-repository
qa: task_05
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
---

# Tasks — Read-only Run commands work outside a repository

| id      | title                                                                                                | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The worktree-location refusal names its Git root and Load can return it as data                       | backend | medium     | —       |
| task_02 | The Read-Only Run Commands load without the worktree-location refusal and warn instead                | backend | medium     | task_01 |
| task_03 | Runs List outside a repository lists every repository                                                 | backend | medium     | task_02 |
| task_04 | The Roundfix skill, the configuration guide and the glossary state the Read-Only Run Command contract | docs    | low        | task_03 |
| task_05 | Run the final QA gate                                                                                 | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 needs task_01 for `WorktreeLocationOptional`. task_03 changes
`runs.go` and `cli.go` after task_02 and re-records the Behavior Surface
Record task_02 re-records. task_04 documents what tasks 02 and 03 built and
re-records the same Behavior Surface Record. Each Task runs the suites that
consume its changed contract in its own Verification, so no separate
Integration Checkpoint is needed. No other Spec is required.
