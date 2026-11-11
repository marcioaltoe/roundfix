---
schema: spec-tasks/v1
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
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

# Tasks — An ADR marked as not applying in Portuguese draws in no citing ADR

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The decision record, glossary, skills and guides name the Portuguese not-applicable token     | docs    | low        | —       |
| task_02 | The strict check reads the Portuguese not-applicable token                                    | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                         | qa      | low        | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 runs the whole `./internal/speccheck` package, which consumes the
seed rule and the fix text, in its own Verification, so no separate
Integration Checkpoint is needed. The docs come first: task_01 changes text
only, and task_02, the code Task, names its guides through that dependency.
No other Spec is required; the queue delivers this Spec after 0373 to 0378.
