---
schema: spec-tasks/v1
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
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

# Tasks — The Delivery Queue waits for every check the head triggered

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | ADR-0369, its notices, CONTEXT.md, the guides and the Roundfix deliver reference state the rule | docs    | medium     | —       |
| task_02 | The GitHub boundary reads the workflow runs and the Expected Workflows of the head             | backend | medium     | task_01 |
| task_03 | The checking stage waits for the Head Check Set and parks naming each failed check             | backend | high       | task_02 |
| task_04 | The Turn Record lists the checks an item waits for, and deliver status prints them             | backend | medium     | task_03 |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The four implementation Tasks run in series. task_01 writes the decision and
the documents first, so the code Tasks implement what the guides already
state. task_02 needs task_01 for those documents. task_03 needs task_02 for
the `CheckReport` fields and the `WorkflowRun` type of
`internal/delivery/head_check_set.go`. task_04 needs task_03 for the pending list `checkCandidate` builds,
and both change `internal/delivery/engine.go`. Each implementation Task runs
the suites that consume its changed contract in its own Verification, so no
separate Integration Checkpoint is needed. No `requires`: Spec 0360 is
authored in parallel and delivers after this Spec.
