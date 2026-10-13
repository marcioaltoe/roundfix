---
schema: spec-tasks/v1
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Tests that wait on a child do not race the machine

| id      | title                                                                     | type | complexity | needs                     |
| ------- | ------------------------------------------------------------------------- | ---- | ---------- | ------------------------- |
| task_01 | The detach and other CLI tests bound by real time only what they assert   | test | high       | —                         |
| task_02 | The agent tests stop racing the real cancel grace and cleanup timeout     | test | high       | —                         |
| task_03 | The daemon, worktree, store and spec tests stop racing their work         | test | medium     | —                         |
| task_04 | The Real-Time Budget rule holds each test file to a ceiling               | test | medium     | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                     | qa   | high       | task_04                   |

Waves: 1 → task_01, task_02, task_03 · 2 → task_04 · 3 → task_05

task_01, task_02 and task_03 share no declared file and run side by side.
task_04 measures each file's ceiling on the tree the first three leave, so it
follows all of them. The gate follows task_04, the only non-QA leaf, and runs
its load windows on the tree every Task built.
