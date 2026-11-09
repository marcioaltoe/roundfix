---
schema: spec-tasks/v1
spec: 0368-the-disk-full-tests-read-no-real-host-path
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

# Tasks — The disk-full tests read no real host path

| id      | title                                                                                | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The disk-full tests use their own temporary paths and the lookup walks past a file   | backend | medium     | —       |
| task_02 | Run the final QA gate                                                                | qa      | low        | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 runs, in its own Verification, every test of `./internal/cli` that
consumes the changed lookup, so no separate Integration Checkpoint is needed.
No other Spec is required; the queue delivers 0366 and 0367 first.
