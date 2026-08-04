---
schema: spec-tasks/v1
spec: 0076-force-stop-exit-proof
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

# Tasks — Force Stop exit proof

| id      | title                                            | type | complexity | needs   |
| ------- | ------------------------------------------------ | ---- | ---------- | ------- |
| task_01 | Keep the signal-ignoring helper alive            | test | medium     | —       |
| task_02 | Prove the escalation caused the exit             | test | medium     | task_01 |
| task_03 | Run the final QA gate                            | qa   | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The graph is a chain because the two halves of the defect cannot be verified
independently. Fixing the parent's ordering while the helper still dies of a
runtime deadlock makes the test pass more reliably and prove no more than it
does today; task_01 must land first so task_02's causation assertion has a live
process to assert about.

The QA gate is authored as task_03 and depends on task_02, the graph's only
non-QA leaf.
