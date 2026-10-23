---
schema: spec-tasks/v1
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
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

# Tasks — The documentation contracts cost seconds per active Spec

| id      | title                                                           | type    | complexity | needs   |
| ------- | --------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The removed-history citation pass looks each token up in a set  | backend | medium     | —       |
| task_02 | The corpus tests share one Check per active Spec                | test    | medium     | task_01 |
| task_03 | Run the final QA gate                                           | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 makes every Check cheaper in `internal/spec` and proves it with the
`internal/spec`, `internal/speccheck` and `internal/cli` tests that consume
the History Citations it reports, so it needs no separate Integration
Checkpoint. task_02 follows it so that its corpus tests and the package
timings run on the cheaper Check; it changes only test code in
`internal/docscontract`. The gate follows both. No other Spec is required.
