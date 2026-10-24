---
schema: spec-tasks/v1
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
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

# Tasks — The CI job runs the documentation contracts once

| id      | title                                                  | type  | complexity | needs   |
| ------- | ------------------------------------------------------ | ----- | ---------- | ------- |
| task_01 | A pull request runs the documentation contracts once   | infra | medium     | —       |
| task_02 | Run the final QA gate                                  | qa    | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 proves its Makefile and workflow change with `make` runs over stub
tools in temporary directories and a workflow parse; it needs no separate
Integration Checkpoint. The gate follows it. No other Spec is required; the
audit reference it cites comes from Spec 0311's planning change.
