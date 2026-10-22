---
schema: spec-tasks/v1
spec: 0300-reconcile-releases-a-stale-run-owner-record
qa: task_03
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
---

# Tasks — Reconcile releases a stale Run owner record

| id      | title                                                                                 | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The reconcile guide, the skill reference and the glossary describe the Run owner record | docs    | low        | —       |
| task_02 | Reconcile releases a Run owner record whose owner identity mismatch is proven          | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                 | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 writes the `reconcile` guide, the Roundfix `reconcile` reference and
the glossary for the surface task_02 changes, so task_02 follows it in one
dependency chain. task_02 changes the store and the command and proves both in
`internal/store` and `internal/cli`, the packages that consume every contract
it changes, so it needs no separate Integration Checkpoint. The gate follows
both. No other Spec is required.
