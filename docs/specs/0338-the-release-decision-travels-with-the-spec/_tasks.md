---
schema: spec-tasks/v1
spec: 0338-the-release-decision-travels-with-the-spec
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

# Tasks — The release decision travels with the Spec

| id      | title                                                                                                 | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The PRD's Delivery Type reaches the Archive Record and the queue's pull request title                 | backend | medium     | —       |
| task_02 | An authorization record's Release Ceiling answers the release plan's approval                         | backend | high       | task_01 |
| task_03 | The Baseline's release and authorization clauses name the Release Ceiling                             | infra   | medium     | task_02 |
| task_04 | The authoring skills fill the Delivery Type and ask the release question                              | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                                 | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 needs task_01 because both edit the Roundfix skill's `SKILL.md` pair,
the owned-skill version record, `CONTEXT.md` and the Behavior Surface Record.
task_03 needs task_02 because its clauses describe the keys and the answer
task_02 builds. task_04 needs task_03 because it teaches the same keys and
edits the owned-skill version record. Each Task runs the suites that consume
its changed contract in its own Verification, so no separate Integration
Checkpoint is needed. No other Spec is required.
