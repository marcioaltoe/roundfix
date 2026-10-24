---
schema: spec-tasks/v1
spec: 0316-an-agent-may-start-the-setup-skill
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

# Tasks — An Agent may start the setup skill

| id      | title                                                                  | type | complexity | needs   |
| ------- | ---------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The setup skill no longer turns off model invocation                   | docs | medium     | —       |
| task_02 | Run the final QA gate                                                  | qa   | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 changes the skill, its mirror, its version record, one guide
paragraph and the Behavior Surface record, and proves them with a new
contract test in `skills` and the existing skills, Baseline and documentation
contracts, so it needs no separate Integration Checkpoint. The gate follows
it. No other Spec is required.
