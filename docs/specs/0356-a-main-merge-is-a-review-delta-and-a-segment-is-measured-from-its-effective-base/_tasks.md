---
schema: spec-tasks/v1
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
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

# Tasks — A main merge is a review delta, and a segment is measured from its Effective Base

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Each Review Segment is measured from its Effective Base                                                  | backend | medium     | —       |
| task_02 | A conflict-free merge of the default branch keeps the Reviewer Lineage                                   | backend | medium     | task_01 |
| task_03 | deliver status and deliver retry name a review that restarts at round 1                                  | backend | medium     | task_02 |
| task_04 | The guides, the Roundfix skill, the glossary and four ADRs state the Effective Base and lineage rules    | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                                    | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 calls `reviewMeasuredFrom` from task_01, task_03 reads the `Restart`
task_02 sets, and task_04 documents all three. Each implementation Task runs
the whole `./internal/cli` suite, which consumes its changed contract, in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required.
