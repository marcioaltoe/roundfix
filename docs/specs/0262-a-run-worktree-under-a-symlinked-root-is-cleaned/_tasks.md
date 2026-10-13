---
schema: spec-tasks/v1
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
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

# Tasks — A Run Worktree under a symlinked root is cleaned

| id      | title                                                                   | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A recorded worktree path is resolved within its Worktree Root           | backend | high       | —                |
| task_02 | Worktree paths compare through their longest existing prefix            | backend | medium     | task_01          |
| task_03 | Every consumer cleans under a symlinked root, and the reclaim is documented | backend | medium | task_02          |
| task_04 | Run the final QA gate                                                   | qa      | high       | task_03          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 and task_02 both edit `internal/worktree/worktree.go`, so they run in
series. task_03 is the Integration Checkpoint for both: its Verification runs
the `internal/cli` and `internal/daemon` suites that consume the worktree
package's reconciliation and item cleanup. The gate follows task_03. The
disposal lock tests are not part of this Spec: Spec 0259's corrective task_06
fixes them.
