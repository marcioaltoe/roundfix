---
schema: spec-tasks/v1
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
requires: [0374-a-machine-verification-capacity-bounds-verifications-across-repositories]
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
      needs: [task_06]
    - id: task_06
      file: task_06.md
      needs: [task_04]
---

# Tasks — A Verification Core Budget keeps concurrent Verifications light

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Verification Core Budget is computed and a Verification command receives it                         | backend | medium     | —       |
| task_02 | Every Verification and Agent Session of an Implement Run receives the Verification Core Budget          | backend | high       | task_01 |
| task_03 | The repository gate and the Linux Lane receive the budget, and doctor reports it                        | backend | high       | task_02 |
| task_04 | The glossary defines the Verification Core Budget                                                       | docs    | low        | task_03 |
| task_05 | Run the final QA gate                                                                                   | qa      | medium     | task_06 |
| task_06 | The Linux Lane builds the budget entries from the corebudget names                                       | backend | low        | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_05

task_06 is corrective Task 1 of 2, added on 2026-11-11 for finding F01 of the QA Report of that day; the QA gate now needs it.

The Tasks run in series. task_02 and task_03 call the package task_01 adds,
share the Roundfix skill version and the Behavior Surface Record, and task_03
calls task_02's budget helper; task_04 defines the term the
guides of tasks 02 and 03 use. Each implementation Task runs the suites that
consume its changed contract in its own Verification, so no separate
Integration Checkpoint is needed. Spec 0374 is required: it adds
`runs.max_verifying`, the doctor line and the configuration section this
Spec places its own after; Spec 0373 delivers before both and edits other
sections of the shared files named in `_prd.md` → Prerequisites.
