---
schema: spec-tasks/v1
spec: 0308-recorded-paths-are-read-before-carry-forward-provenance
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

# Tasks — Recorded paths are read before carry-forward provenance

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A carried Task's recorded paths are read and replaced before its provenance    | backend | medium     | —       |
| task_02 | Run the final QA gate                                                         | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 changes the one reader and writer of the recorded section and proves
it in `internal/spec` and in the packages that consume it (`internal/cli`,
`internal/worktree`, `internal/daemon`, `internal/speccheck`), so it needs no
separate Integration Checkpoint. The gate follows it. No other Spec is
required.
