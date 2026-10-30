---
schema: spec-tasks/v1
spec: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
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

# Tasks — The Delivery Queue regenerates the derived records of an archived Spec

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The coverage record writer names its command, and the regeneration contract models an archived item | backend | medium     | —       |
| task_02 | The queue merges an archived item's derived-record conflict and refuses it without the sanction    | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                              | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 follows task_01 because its resolution test passes only once the
coverage record writer names its command. task_02 also re-records the coverage
record after it adds two top-level tests, and that re-record needs the same
fix. The two Tasks share no file. Each runs the suites that consume its
change, so no separate Integration Checkpoint is needed. No other Spec is
required.
