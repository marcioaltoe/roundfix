---
schema: spec-tasks/v1
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
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

# Tasks — internal/cli tests stop reading the live documentation

| id      | title                                                               | type | complexity | needs   |
| ------- | ------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The repository-profile tests become documentation contracts         | test | medium     | —       |
| task_02 | The other live-document tests become documentation contracts        | test | medium     | task_01 |
| task_03 | Run the final QA gate                                               | qa   | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

Both implementation Tasks change
`internal/verifyselect/contracts_repository_test.go`, so task_02 follows
task_01 (the Wave collision rule). Each moves test code inside
`internal/cli` and proves it with the moved tests themselves and the
repository contract declarations test; neither needs a separate Integration
Checkpoint. The gate follows both and runs the cache and unreadable-trees
measurements. No other Spec is required.
