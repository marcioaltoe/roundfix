---
schema: spec-tasks/v1
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
qa: task_03
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_03
      file: task_03.md
      needs: [task_05]
---

# Tasks — internal/cli tests stop reading the live documentation

| id      | title                                                               | type | complexity | needs   |
| ------- | ------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The repository-profile tests become documentation contracts         | test | medium     | —       |
| task_02 | The other live-document tests become documentation contracts        | test | medium     | task_01 |
| task_04 | The upgrade characterization test runs in a fixture workspace       | test | low        | task_02 |
| task_05 | The suite guard fingerprints the repository through Git             | backend | medium  | task_04 |
| task_03 | Run the final QA gate                                               | qa   | high       | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_04 · 4 → task_05 · 5 → task_03

task_01 and task_02 both change
`internal/verifyselect/contracts_repository_test.go`, so task_02 follows
task_01 (the Wave collision rule). Each moves test code inside
`internal/cli` and proves it with the moved tests themselves and the
repository contract declarations test; neither needs a separate Integration
Checkpoint. task_04 is the corrective Task for F-01 of the first Run's QA
gate. It changes one test file that neither earlier Task declares, and it
follows task_02 so that it runs alone in its Wave on the tree both moves
produced. task_05 is the second and last corrective Task, from the
maintainer's answer of 2026-10-25 (about 11:00 BRT): the suite guard
fingerprints through Git. It shares no file with the earlier Tasks and
follows task_04, so it runs alone in Wave 4. The gate follows it and runs the
cache, test input log, unreadable-trees and guard sabotage measurements. No other Spec is required.
