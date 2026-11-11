---
schema: spec-tasks/v1
spec: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
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

# Tasks — A push whose tests pass runs every contract under a measured budget

| id      | title                                                                         | type  | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ----- | ---------- | ------- |
| task_01 | The push gate records its suite-time verdict and judges it after the contracts | infra | medium     | —       |
| task_02 | Run the final QA gate                                                         | qa    | low        | task_01 |

Waves: 1 → task_01 · 2 → task_02

One implementation Task: the new test reads both the `Makefile` and the
workflow, and the setup manifest must be regenerated from the final
`Makefile` bytes, so splitting them would only add a Task whose Verification
depends on the other. task_01's Verification runs the whole
`./internal/verifyselect` package, which holds every test that consumes the
`Makefile` targets and the workflow, so no separate Integration Checkpoint is
needed. No other Spec is required.
