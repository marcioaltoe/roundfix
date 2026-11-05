---
schema: spec-tasks/v1
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
qa: task_05
requires: [0352-only-a-git-entry-that-git-accepts-marks-a-git-root]
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

# Tasks — A queue item shows its turn, and the pre-PR review gets its own capacity

| id      | title                                                                                       | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | ADR-0365, CONTEXT.md, the Roundfix deliver reference and the guides state the behavior       | docs    | medium     | —       |
| task_02 | delivery.review_concurrency sets the Review Capacity, and the Limits line prints it          | backend | medium     | task_01 |
| task_03 | The Queue Owner reviews items in review slots beside its one item, and keeps the start order | backend | high       | task_02 |
| task_04 | The owner records each item's Queue Turn, and deliver status prints the Turn lines           | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                                       | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The four implementation Tasks run in series. task_01 writes the decision and
the documents first, so every Task that changes the `deliver` command has its
guide among its dependencies and the later code Tasks implement what the
guides already state. task_02 needs task_01 for those guides. task_03 needs
task_02 for the `Delivery.ReviewConcurrency` field it wires into the engine.
task_04 needs task_03 because it records the review slots task_03 creates
through the engine methods task_03 adds, and shares
`internal/delivery/engine.go` and `internal/cli/deliver_workflow.go` with it;
it also changes `internal/cli/deliver.go`, which task_02 changes. Each
implementation Task runs the existing suites that consume its changed
contract in its own Verification, so no separate Integration Checkpoint is
needed. Spec 0352 is required: task_02's Project Config test marks its
directory with `gittest.MarkRepository`, which 0352 adds.
