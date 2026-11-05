---
schema: spec-tasks/v1
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
qa: task_05
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
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A large authored data file is reviewed apart, not omitted

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Project Config declares the Review Data Paths                                                            | backend | medium     | —       |
| task_02 | Review Data Paths leave the code diff and are planned into Data Review Segments                          | backend | medium     | task_01 |
| task_03 | The review reads the Review Data Paths apart, in its session                                             | backend | high       | task_02 |
| task_04 | The review guide, the Roundfix skill, the glossary and ADR-0282 describe Data Review Segments            | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                                    | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 matches with `ReviewsApart` from task_01, task_03 sends what task_02
plans, and task_04 documents all three. Each implementation Task runs the
suites that consume its changed contract in its own Verification, so no
separate Integration Checkpoint is needed. No other Spec is required; Spec
0356, delivered first, changes neighbouring functions of the same files.
