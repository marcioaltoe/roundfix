---
schema: spec-tasks/v1
spec: 0231-checks-that-hold-in-delivery
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_03]
---

# Tasks — Checks that hold in delivery

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The detach fixture tests count a process-group member that is already exiting as ended        | test    | low        | —                |
| task_02 | CI's pull request job tests the head merged with the current default branch and records it    | infra   | low        | —                |
| task_03 | The Delivery Queue re-runs a failed check that tested an older default branch instead of parking it | backend | medium     | task_02          |
| task_04 | Run the final QA gate                                                                         | qa      | high       | task_01, task_03 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

task_01 and task_02 share no file and run together. task_02 changes only the
Governed Path it is granted, so its commit is the tooling change alone;
task_03 follows it because the queue reads the annotation the workflow
writes. The gate follows both leaves.
