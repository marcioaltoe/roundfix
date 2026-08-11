---
schema: spec-tasks/v1
spec: 0038-terminal-run-worktree-reconciliation
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Tasks — Terminal Run Worktree reconciliation

Cross-Spec prerequisite: Spec 0037 Terminal Outcome Integrity must be completed
before this graph starts because guarded Integration Pending reconciliation is
part of the safety boundary.

| id      | title                                                         | type    | complexity | needs                          |
| ------- | ------------------------------------------------------------- | ------- | ---------- | ------------------------------ |
| task_01 | Classify terminal Run Worktrees with positive proof           | backend | high       | —                              |
| task_02 | Apply stale-proof cleanup and migrate terminal reaping         | backend | high       | task_01                        |
| task_03 | Reconcile Integration Pending with durable evidence           | data    | high       | task_01, task_02               |
| task_04 | Deliver the Reconcile Command and apply contract              | backend | high       | task_01, task_02, task_03      |
| task_05 | Surface retained terminal Run Worktrees in Runs List          | backend | medium     | task_01                        |
| task_06 | Align reconciliation docs and glossary                       | docs    | medium     | task_04, task_05               |
| task_07 | Align the protected Roundfix Skill pair                      | docs    | low        | task_06                        |

Waves: 1 → task_01 · 2 → task_02, task_05 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_07
