---
schema: spec-tasks/v1
spec: 0113-a-refused-gate-writes-its-refusal-once
qa: task_06
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
    - id: task_07
      file: task_07.md
      needs: [task_04]
    - id: task_05
      file: task_05.md
      needs: [task_07]
    - id: task_06
      file: task_06.md
      needs: [task_05, task_07]
---

# Tasks — Gate Refusal Report Shape

| id      | title                                   | type    | complexity | needs             |
| ------- | --------------------------------------- | ------- | ---------- | ----------------- |
| task_01 | Write terminal row on precondition      | backend | high       | —                 |
| task_02 | Detect precondition failure             | backend | high       | task_01           |
| task_03 | Store precondition metadata             | backend | medium     | task_02           |
| task_04 | Update mechanical stage validation      | backend | high       | task_03           |
| task_07 | Parse results from the Results table    | backend | high       | task_04           |
| task_05 | Read newest report only                 | backend | high       | task_07           |
| task_06 | QA gate                                 | qa      | medium     | task_05, task_07  |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_07 · 6 → task_05 · 7 → task_06

task_05 depends on task_07 for edit locality, not for behavior: both rewrite the
row-collection path in `internal/speccheck/mechanical.go` and its test. Run as
siblings on 2026-08-26 they produced an integration conflict on
`mechanical_test.go` that discarded a passing task_05. Two Tasks that rewrite the
same region are not independently implementable, whatever their logic says, so
the edge is the fix rather than a merge strategy.

Each Task's work, references, and Verification live in its own `task_NN.md`.
This file owns topology only; a second copy of a Verification command is how the
two drift, which is what the 2026-08-26 measurement below was caused by.

task_07 carries Core Feature 4, added to the PRD after the defect was measured
live on Spec 0098 on 2026-08-26: the shape detector parses every markdown table
in a report as its Results matrix, so three cells of an evidence table written
in prose each became a separate blocker on the run that carried the fix.
