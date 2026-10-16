---
schema: spec-tasks/v1
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
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

# Tasks — The F6 sender is named, and a rebuilt item binary is replaced

| id      | title                                                                                     | type | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The item binary tests report how the probe ended, and the diagnosis records the F6 sender | test | medium     | —       |
| task_02 | A rebuilt item binary replaces the file the previous build left                           | test | medium     | task_01 |
| task_03 | Run the final QA gate                                                                     | qa   | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 and task_02 both edit `internal/cli/deliver_item_binary_test.go`, and
the diagnosis of task_01 must run while the build still copies in place, so
task_02 follows it. The gate follows the only leaf.
