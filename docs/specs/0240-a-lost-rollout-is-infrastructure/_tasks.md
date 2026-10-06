---
schema: spec-tasks/v1
spec: 0240-a-lost-rollout-is-infrastructure
qa: task_04
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
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — A lost rollout is infrastructure

| id      | title                                                                                                   | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The glossary, the guides and the Roundfix Skill describe the Lost Rollout, its recovery and its park    | docs    | low        | —                |
| task_02 | The ACPX Runner names a Lost Rollout by its code and phrase                                              | backend | medium     | —                |
| task_03 | The Daemon recovers a Lost Rollout and the queue parks an unrecovered one without counting its retry    | backend | high       | task_01, task_02 |
| task_04 | Run the final QA gate                                                                                   | qa      | high       | task_03          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

task_01 and task_02 share no declared file. task_03 reads the runner's
`DescribeLostRollout` from task_02 and changes the `deliver` surface whose
guide task_01 writes, and it produces and consumes the `rollout_lost` Run
Event in one Task, so the event contract is proved end to end. The gate
follows the only leaf.
