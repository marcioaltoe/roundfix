---
schema: spec-tasks/v1
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
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
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A Run Worktree under a symlinked root is cleaned

| id      | title                                                                   | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A recorded worktree path is resolved within its Worktree Root           | backend | high       | —                |
| task_02 | Worktree paths compare through their longest existing prefix            | backend | medium     | task_01          |
| task_03 | Every consumer cleans under a symlinked root, and the reclaim is documented | backend | medium | task_02          |
| task_04 | The disposal lock tests wait on a synchronization point                 | test    | medium     | —                |
| task_05 | Run the final QA gate                                                   | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_04 · 2 → task_02 · 3 → task_03 · 4 → task_05

task_01 and task_02 both edit `internal/worktree/worktree.go`, so they run in
series. task_03 is the Integration Checkpoint for both: its Verification runs
the `internal/cli` and `internal/daemon` suites that consume the worktree
package's reconciliation and item cleanup. task_04 changes only
`internal/worktree/disposal.go` and `internal/worktree/disposal_test.go`, which
no other Task declares, so it runs beside task_01. The gate follows both
leaves.
