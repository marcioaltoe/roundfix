---
schema: spec-tasks/v1
spec: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
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

# Tasks — A review whose reviewed head is gone opens a new lineage

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The review and deliver guides and references name a commit that is gone      | docs    | low        | —       |
| task_02 | A recorded commit proven gone restarts the lineage and is refused by name     | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                         | qa      | low        | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 runs, in its own Verification, every review, lineage, disposition,
dismissal, correction and deliver test of `./internal/cli`, which consume the
changed lineage decision and correction proof, so no separate Integration
Checkpoint is needed. The guides come first: task_01 changes text only, and
task_02, the code Task, names its guides through that dependency. No other
Spec is required; the queue delivers this Spec before 0371 and 0373.
