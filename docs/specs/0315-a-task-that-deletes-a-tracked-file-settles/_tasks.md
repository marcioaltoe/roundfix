---
schema: spec-tasks/v1
spec: 0315-a-task-that-deletes-a-tracked-file-settles
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

# Tasks — A Task that deletes a tracked file settles

| id      | title                                                                    | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | Tracked-file readers skip a path the working tree lacks                  | backend | medium     | —       |
| task_02 | A gated Task whose declared deletion is only on disk settles             | test    | medium     | task_01 |
| task_03 | Run the final QA gate                                                    | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 changes the three readers and proves them in `internal/spec`. task_02
proves the fix through the Daemon's real settlement path; its `disk-only`
subtest passes only with task_01's readers, so it follows task_01 and serves
as the Integration Checkpoint between `internal/spec` and `internal/daemon`.
The gate follows both. No other Spec is required.
