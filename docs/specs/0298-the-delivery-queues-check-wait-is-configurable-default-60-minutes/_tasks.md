---
schema: spec-tasks/v1
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
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

# Tasks — The Delivery Queue's check wait is configurable, default 60 minutes

| id      | title                                                                              | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The guides, the Roundfix reference and the glossary describe the Check Wait        | docs    | low        | —       |
| task_02 | The Check Wait reads delivery.check_timeout, 60 minutes by default                 | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                              | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 writes the configuration guide, the `deliver` guide and the Roundfix
`deliver` reference for the surface task_02 changes, so task_02 follows it in
one dependency chain. task_02 changes the code in `internal/config`,
`internal/delivery` and `internal/cli` and proves each contract in the package
that consumes it, the wiring included, so it needs no separate Integration
Checkpoint. The gate follows both. No other Spec is required.
