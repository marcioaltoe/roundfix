---
schema: spec-tasks/v1
spec: 0365-the-capability-projection-reads-grouped-select-options
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

# Tasks — The capability projection reads grouped select options

| id      | title                                                                                     | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The capability projection reads a grouped select option as its groups' values             | backend | medium     | —       |
| task_02 | The Exact Agent Selection Proof passes over the recorded codex-acp 2.6.0 payload          | test    | medium     | task_01 |
| task_03 | Run the final QA gate                                                                     | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 reads the fixtures and the projection task_01 adds, so it runs
after it. Each Task runs the suites that consume its contract in its own
Verification, so no separate Integration Checkpoint is needed. No other Spec
is required; this Spec delivers before 0366 and 0367, and none of them
changes `internal/agent`.
