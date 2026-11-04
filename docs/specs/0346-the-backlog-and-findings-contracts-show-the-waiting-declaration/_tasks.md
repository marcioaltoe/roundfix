---
schema: spec-tasks/v1
spec: 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — The Backlog and Findings contracts show the Waiting Declaration

| id      | title                                                                         | type  | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ----- | ---------- | ------- |
| task_01 | The Backlog and Findings contracts show the Waiting Declaration commented     | infra | medium     | —       |
| task_02 | Run the final QA gate                                                         | qa    | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02

One implementation Task: the module, its derived files, this repository's
managed guide and the user guide's Findings template are held equal by
repository tests, so they change in one settlement. task_01 runs the suites
that consume the changed contract in its own Verification, so no separate
Integration Checkpoint is needed. No other Spec is required; cycle AX
delivers this Spec last, and task_01 regenerates every shared derived file
by command on its own base.
