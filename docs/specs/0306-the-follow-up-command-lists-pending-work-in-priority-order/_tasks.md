---
schema: spec-tasks/v1
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
qa: task_05
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The Follow-up Command lists a repository's pending work in priority order

| id      | title                                                                                                                         | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The command guide, the follow-up reference, the configuration guide and the glossary describe the Follow-up Command           | docs    | medium     | —                |
| task_02 | The judge answers four yes/no questions per Follow-up Item, after a secret scan                                               | backend | medium     | —                |
| task_03 | roundfix follow-up lists every pending item in Follow-up Priority order beside the work in progress                           | backend | high       | task_01, task_02 |
| task_04 | roundfix follow-up --status shows the work in progress, the pending Specs, the next release and what waits on the maintainer | backend | medium     | task_03          |
| task_05 | Run the final QA gate                                                                                                         | qa      | high       | task_04          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 documents the whole command first, as Spec 0261 did for the Cancel
Command, so the code Tasks that follow it satisfy the command documentation
contracts from their first commit; it adds only the guide's Skill Coverage
Map entry, because the `command follow-up` entry needs the command. task_02
changes only `internal/judge` and shares no file with task_01, so they run
together. task_03 is task_02's Integration Checkpoint: it is the judge
contract's only caller. task_03 and task_04 share `internal/cli/cli.go`,
`internal/cli/follow_up.go`, the Skill Coverage Map and the Behavior Surface
Record, so they run in series. No other Spec is required.
