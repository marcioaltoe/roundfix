---
schema: spec-tasks/v1
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
qa: task_04
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
---

# Tasks — QA evidence links resolve in the committed tree

| id      | title                                                                                                     | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The shared evidence link check refuses unresolved links in qa-report accept and settle                    | backend | medium     | —       |
| task_02 | The Daemon gives the QA gate one repair turn and refuses an unresolved evidence link at the commit         | backend | medium     | task_01 |
| task_03 | The qa-gate skill keeps evidence under names the repository does not ignore                                | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                                     | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 needs task_01 because it calls the shared check and both raise the
Roundfix skill version, re-record the owned-skill version record and the
Behavior Surface Record. task_03 needs task_02 because both re-record the
owned-skill version record. Each Task runs the suites that consume its changed
contract in its own Verification, so no separate Integration Checkpoint is
needed. No other Spec is required.
