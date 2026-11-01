---
schema: spec-tasks/v1
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
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

# Tasks — The Spec Check leaves the QA gate unprobed and prints its expected Verification

| id      | title                                                                                               | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Spec Check leaves the QA gate unprobed through the Daemon's QA predicate                        | backend | medium     | —       |
| task_02 | The authored-QA-Verification finding prints the expected block, and write-tasks teaches copying it  | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                               | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 needs task_01 because both edit the spec command guide, the Roundfix
spec reference and its `SKILL.md` pair, `internal/speccheck/report.go`, the
owned-skill version record and the Behavior Surface Record. Each Task runs the
suites that consume its changed contract in its own Verification, so no
separate Integration Checkpoint is needed. No other Spec is required.
