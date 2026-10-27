---
schema: spec-tasks/v1
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — The Baseline setups stop sending agents to the Secondbrain inbox

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Every upstream setup and this repository carry the knowledge-workspace of b56fc5ff | backend | medium     | —       |
| task_02 | Run the final QA gate                                                         | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 imports the snapshots, restores this repository's
`knowledge-workspace`, regenerates the derived files and the Setup Manifest,
and moves the revision test constant with a new tree test; it proves them with
the whole `internal/baseline` and `skills` packages and the Setup Manifest
contract, so it needs no separate Integration Checkpoint. The gate follows
it. No other Spec is required.
