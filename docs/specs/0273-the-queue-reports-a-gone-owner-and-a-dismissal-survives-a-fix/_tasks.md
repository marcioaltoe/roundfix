---
schema: spec-tasks/v1
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
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

# Tasks — The queue reports a gone owner, and a dismissal survives a fix

| id      | title                                                                               | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Deliver status names a recorded Queue Owner that is not running                      | backend | medium     | —       |
| task_02 | A deliver command in an item worktree names the checkout that owns the queue          | backend | medium     | task_01 |
| task_03 | A dismissal is accepted after a fix that leaves the finding's anchor unchanged        | backend | high       | task_02 |
| task_04 | The Roundfix skill and the glossary describe the owner report, the refusal and the dismissal | docs    | low        | task_03 |
| task_05 | Run the final QA gate                                                               | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 and task_02 both edit `internal/cli/deliver.go` and the `deliver` user
guide, and task_01, task_02 and task_03 each re-record
`docs/references/behavior-surfaces.json` after their guide edit, so they run in
series. task_04 describes all three in the skill and the glossary and runs
after them. The gate follows the only leaf.
