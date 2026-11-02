---
schema: spec-tasks/v1
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
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

# Tasks — An approved Baseline update deletes the Retired Skills

| id      | title                                                                                     | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | An approved Baseline update deletes the Retired Skills it proves unedited                 | backend | high       | —       |
| task_02 | The glossary, the guide and the Roundfix skill say what an approved update deletes        | docs    | medium     | task_01 |
| task_03 | Run the final QA gate                                                                     | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 needs task_01 because the documents describe the behavior task_01
ships, and ADR-0246 is retired only once the code follows ADR-0350. The two
share no declared file. Each Task runs the suites that consume its changed
contract in its own Verification, so no separate Integration Checkpoint is
needed. No other Spec is required.
