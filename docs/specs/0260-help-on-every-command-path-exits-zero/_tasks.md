---
schema: spec-tasks/v1
spec: 0260-help-on-every-command-path-exits-zero
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

# Tasks — Help on every command path exits zero

| id      | title                                  | type    | complexity | needs   |
| ------- | -------------------------------------- | ------- | ---------- | ------- |
| task_01 | Help on every command path exits zero  | backend | medium     | —       |
| task_02 | Run the final QA gate                  | qa      | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02

One implementation Task carries the whole fix: the two dispatchers, the tests
that guard the class, the contract helper and the skill text change together,
because each Verification needs the others (the sweep and the contract fail
until both dispatchers answer help, and the skills rule needs the skill files
in the Task that changes their sources). The gate follows it.
