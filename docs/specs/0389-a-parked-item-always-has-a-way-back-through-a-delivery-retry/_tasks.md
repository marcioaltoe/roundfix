---
schema: spec-tasks/v1
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
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

# Tasks — A parked item always has a way back through a Delivery Retry

| id      | title                                                                                 | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A Delivery Retry requeues an item parked before it has an Item Branch                  | backend | medium     | —       |
| task_02 | The delivery authorization of an archived Spec is read without a Spec Root directory   | backend | medium     | task_01 |
| task_03 | A Pull Request that holds an earlier candidate returns the item to the gate            | backend | medium     | task_02 |
| task_04 | The Roundfix deliver reference describes the three changes                             | docs    | low        | task_03 |
| task_05 | Run the final QA gate                                                                 | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series. Tasks 01, 02 and 03 each edit
`docs/user-guide/commands/deliver.md` and re-record
`docs/references/behavior-surfaces.json` after their own guide edit, and tasks
01 and 03 both change `internal/delivery/engine.go`; task_04 mirrors the four
guide texts into the Roundfix skill's `deliver` reference and raises its
version once. Each code Task runs the suites that consume its changed surfaces
in its own Verification, so no separate Integration Checkpoint is needed. No
other Spec is required; Spec 0388 delivers first and edits other parts of the
shared files named in `_prd.md` → Prerequisites.
