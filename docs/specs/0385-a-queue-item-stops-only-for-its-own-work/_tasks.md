---
schema: spec-tasks/v1
spec: 0385-a-queue-item-stops-only-for-its-own-work
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

# Tasks — A queue item stops only for its own work

| id      | title                                                                                  | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Queue Owner logs its capacity wait and keeps its record when it stops with an item left | backend | high       | —       |
| task_02 | A review finding's anchor may end a sentence                                           | backend | medium     | task_01 |
| task_03 | The review validator's sealed session carries the reviewed repository's marker         | backend | medium     | task_02 |
| task_04 | The Roundfix deliver and review references describe the three changes                  | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                  | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series. Tasks 01, 02 and 03 change different code, but each
re-records `docs/references/behavior-surfaces.json` after its own guide edit,
and tasks 02 and 03 both edit `docs/user-guide/commands/review.md`; task_04
mirrors all three guide texts into the Roundfix skill's references and raises
its version once. Each code Task runs the suites that consume its changed
surfaces in its own Verification, so no separate Integration Checkpoint is
needed. No other Spec is required; Specs 0383 and 0384 deliver first and edit
other sections of the shared files named in `_prd.md` → Prerequisites.
