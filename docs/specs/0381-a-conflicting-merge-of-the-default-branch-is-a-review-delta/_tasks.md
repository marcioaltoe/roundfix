---
schema: spec-tasks/v1
spec: 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
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

# Tasks — A conflicting merge of the default branch is a review delta

| id      | title                                                                                     | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The guides, references, glossary and ADR notices describe the Conflicted Carry             | docs    | low        | —       |
| task_02 | A conflicting merge is measured from its Conflicted Carry and keeps the Reviewer Lineage   | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                     | qa      | low        | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 runs, in its own Verification, every review, lineage, disposition,
dismissal, correction, deliver and segment test of `./internal/cli`, which
consume the changed measurement and prompts, so no separate Integration
Checkpoint is needed. The guides come first: task_01 changes text only, and
task_02, the code Task, names its guides through that dependency. No other
Spec is required; the queue delivers Specs 0377, 0376, 0378, 0379 and 0380
before this one.
