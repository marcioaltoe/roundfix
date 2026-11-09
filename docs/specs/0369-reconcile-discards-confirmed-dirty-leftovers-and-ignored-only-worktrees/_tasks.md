---
schema: spec-tasks/v1
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
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

# Tasks — Reconcile discards confirmed dirty leftovers and releases ignored-only worktrees

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A confirmed Dirty Discard removes an ended Run's leftovers after recording them                      | backend | high       | —       |
| task_02 | An unregistered Run Worktree holding only ignored files is released, and empty directories go with it | backend | medium     | task_01 |
| task_03 | The amended ADRs and the Roundfix skill's reconcile reference describe both                            | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                                | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 and task_02 both change `internal/worktree/worktree.go`,
`internal/cli/reconcile.go`, the `reconcile` user guide, `CONTEXT.md` and the
Behavior Surface Record, so they run in series; task_03 documents the
delivered behavior of both and regenerates the skill version record and the
Behavior Surface Record on its own base. Each Task runs the suites that
consume its contract in its own Verification, so no separate Integration
Checkpoint is needed. No other Spec is required; the queue delivers 0366,
0367 and 0368 first.
