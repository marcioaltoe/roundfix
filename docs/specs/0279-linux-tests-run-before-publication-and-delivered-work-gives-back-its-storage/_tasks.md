---
schema: spec-tasks/v1
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
qa: task_04
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
---

# Tasks — Linux tests run before publication, and delivered work gives back its storage

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Delivery Queue gate runs the changed Go packages in a Linux Lane                          | backend | high       | —       |
| task_02 | A delivered Spec's leftovers beneath its declared paths are superseded, and a kept dirty Run names why | backend | high       | task_01 |
| task_03 | Merged items give back their Item Binaries under a written retention policy                   | backend | high       | task_02 |
| task_04 | Run the final QA gate                                                                         | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

Each implementation Task raises the Roundfix skill version and re-records
`skills/testdata/owned-skill-versions.json`, the Behavior Surface Record and
`CONTEXT.md`, and task_01 and task_03 both change the `deliver` guide, its
reference and `internal/cli/deliver_workflow.go`, so the three run in series.
The gate follows the last one.
