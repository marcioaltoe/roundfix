---
schema: spec-tasks/v1
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
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
      needs: [task_01, task_02]
---

# Tasks — The Baseline reports a required skill that is not installed

| id      | title                                                                         | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The baseline update preview names a missing required skill                    | backend | medium     | —                |
| task_02 | skills check names a missing required skill of the repository it runs in      | backend | medium     | task_01          |
| task_03 | Run the final QA gate                                                         | qa      | high       | task_01, task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 calls the helper task_01 adds and writes the Roundfix skill, the
owned-skill version record and the Behavior Surface record that task_01 also
writes, so it needs task_01. The gate follows both. No other Spec is required.
