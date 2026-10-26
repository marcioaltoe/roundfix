---
schema: spec-tasks/v1
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
qa: task_03
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
      needs: [task_01, task_02]
---

# Tasks — The group probe waits for a zombie, and the Hard Link rule sees all forms

| id      | title                                                                  | type | complexity | needs            |
| ------- | ---------------------------------------------------------------------- | ---- | ---------- | ---------------- |
| task_01 | The unreaped-group test waits for its member to be a zombie            | test | medium     | —                |
| task_02 | The Hard Link rule reports every reference in every module test file   | test | medium     | —                |
| task_03 | Run the final QA gate                                                  | qa   | high       | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03

task_01 changes only the `internal/cli` detach fixture test files and
task_02 only `internal/testfixture/hard_link_rule_test.go`; they share no
file, so they run in the same Wave. The gate follows both and runs the load
harness the planning change recorded in the audit reference. No other Spec
is required.
