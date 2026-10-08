---
schema: spec-tasks/v1
spec: 0254-a-faster-make-test
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A faster make test

| id      | title                                                                         | type | complexity | needs            |
| ------- | ----------------------------------------------------------------------------- | ---- | ---------- | ---------------- |
| task_01 | internal/cli runs its tests in parallel under the Parallel Test Package rule  | test | high       | —                |
| task_02 | Six more packages run their tests in parallel                                 | test | high       | task_01          |
| task_03 | The sequential residue in internal/cli shrinks and pinned history reads less  | test | medium     | task_02          |
| task_04 | One go test runs every selected Repository Contract Test                      | infra | medium    | task_01          |
| task_05 | Run the final QA gate                                                         | qa   | high       | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05

task_01, task_02 and task_03 share
`internal/testfixture/parallel_tests_test.go`, and task_01 and task_03 both
edit `internal/cli` test files, so they run in series. task_04 touches only
`internal/verifyselect`; it follows task_01 so that its gate runs the
converted `internal/cli`. The conversions change more test files than a Task
may declare (the Context bound is 50 entries): each Task declares its Governed
Paths and the files it changes beyond the inserted line, and the Daemon
records the remaining test files of its named packages under
`## Recorded paths`. The gate follows both leaves.
