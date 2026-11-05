---
schema: spec-tasks/v1
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
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

# Tasks — runs list reads as a table, and the Run Browser draws, at a terminal

| id      | title                                                                                              | type     | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | -------- | ---------- | ------- |
| task_01 | The Run Browser and both cockpits draw on the terminal, and the Run Browser never ends silently     | frontend | medium     | —       |
| task_02 | runs list prints a table at a terminal and keeps its stable line everywhere else                    | backend  | high       | task_01 |
| task_03 | attach and runs show accept the short Run id, and ADR-0363 records the two formats                  | backend  | medium     | task_02 |
| task_04 | Run the final QA gate                                                                               | qa       | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The three implementation Tasks run in series. task_02 needs task_01 because
the table uses `RunRepositoryName`, which task_01 adds, and both change the
Roundfix `runs` reference, its mirror, the `SKILL.md` pair, the version
record, the `runs` guide and the Behavior Surface Record. task_03 needs
task_02 for the same shared documents, `internal/cli/cli.go` and
`internal/cli/cli_test.go`. Each implementation Task runs the existing suites
that consume its changed contract in its own Verification, so no separate
Integration Checkpoint is needed. No other Spec is required; Spec 0352 is
delivered first and shares only documents and records, which each Task
re-records.
