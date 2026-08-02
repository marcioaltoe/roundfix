---
schema: spec-tasks/v1
spec: 0071-verification-cost
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04]
    - id: task_07
      file: task_07.md
      needs: [task_03, task_04, task_06]
---

# Tasks — Verification cost

| id      | title                                            | type    | complexity | needs                     |
| ------- | ------------------------------------------------ | ------- | ---------- | ------------------------- |
| task_01 | Record which tests the suite executes             | test    | medium     | —                         |
| task_02 | Let the CLI package take its environment          | backend | high       | task_01                   |
| task_03 | Run the CLI tests in parallel                     | test    | high       | task_02                   |
| task_04 | Free the Baseline package from process state      | backend | high       | task_01                   |
| task_05 | Stop charging every Task for the whole suite      | docs    | medium     | task_01                   |
| task_06 | Assert a suite-time budget                        | infra   | medium     | task_03, task_04          |
| task_07 | Publish the before-and-after                      | docs    | low        | task_03, task_04, task_06 |

Waves: 1 → task_01 · 2 → task_02, task_04, task_05 · 3 → task_03 · 4 → task_06 ·
5 → task_07

Wave 2 fans out because its three Tasks touch disjoint files: the CLI package,
the Baseline package, and the Task files plus the authoring skill.

**Only tasks 02 and 03 can move the headline number.** The baseline shows
packages already overlap under `go test ./...`, so the suite can never finish
faster than its slowest package, and `internal/cli` alone is 113.2s of 136.9s.
task_04 reduces the sum and helps a single-package run; it cannot move the
floor.

task_01 changes no behavior and lands first. It records which test functions
the suite executes, so "coverage is unchanged" becomes an assertion instead of
a claim. The timing baseline is already committed under `baseline/` and is not
re-derived — comparing against a re-measured "before" would make the comparison
circular.
