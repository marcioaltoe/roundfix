---
schema: spec-tasks/v1
spec: 0311-the-suite-audit-and-no-test-hard-links-its-binary
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

# Tasks — The suite audit, and no test hard-links its binary

| id      | title                                                                   | type | complexity | needs   |
| ------- | ----------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The agent package adapters are copies and a fixture child makes none    | test | medium     | —       |
| task_02 | A repository rule refuses a hard link in any test                       | test | low        | task_01 |
| task_03 | Run the final QA gate                                                   | qa   | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 removes the only hard link the suite creates and proves the copies
and the leak with tests in `internal/agent`. task_02 follows it because its
repository rule fails while that hard link remains; it adds the rule, the
guide line and the glossary term. The gate follows both and runs the load
harness that reproduced the Run failures. No other Spec is required.
