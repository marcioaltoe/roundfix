---
schema: spec-tasks/v1
spec: 0250-a-run-database-that-keeps-only-recent-runs
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A Run Database that keeps only recent Runs

| id      | title                                                                                   | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Run Database removes a terminal Run whole and compacts incrementally                 | data    | high       | —                |
| task_02 | roundfix gc applies Run Retention and reports the Runs, rows and bytes it removes        | backend | high       | task_01          |
| task_03 | The Run Retention Sweep runs once a day at the start of Runs and Delivery Queues         | backend | high       | task_02          |
| task_04 | The glossary and the lifecycle policy describe Run Retention                             | docs    | low        | task_01          |
| task_05 | Run the final QA gate                                                                    | qa      | high       | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05

task_02 builds the sweep on task_01's store API, and task_03 calls task_02's
sweep from the start paths and shares `internal/cli/cli.go` with it, so the
three run in series. task_04 edits the lifecycle policy prose after task_01
edits its table, and declares only `CONTEXT.md` and that policy, so it shares
no declared file with task_02 and runs beside it. The gate follows both leaves.
