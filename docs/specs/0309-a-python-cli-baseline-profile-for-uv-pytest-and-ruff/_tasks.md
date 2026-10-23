---
schema: spec-tasks/v1
spec: 0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff
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
      needs: [task_01, task_02, task_03]
---

# Tasks — A Python CLI Baseline Profile for uv, pytest and ruff

| id      | title                                                                                          | type    | complexity | needs                     |
| ------- | ---------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | A Common Setup Snapshot keeps only the skills its components share                             | backend | medium     | —                         |
| task_02 | The python-cli profile governs uv, ruff and pytest on the common setup of Go and Rust          | backend | high       | task_01                   |
| task_03 | The user guide lists the Python CLI profile and the glossary defines the Common Setup Snapshot | docs    | low        | task_02                   |
| task_04 | Run the final QA gate                                                                          | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 changes only the setup derivation, its validator and the asset sync,
and proves them in `internal/baseline`, the package that consumes them.
task_02 adds the `python` common snapshot, which the catalog accepts only once
task_01's validator knows the kind, and both Tasks change
`internal/baseline/assets_sync_test.go`, so it needs task_01. task_03
describes the profile task_02 ships and re-records the Behavior Surface Record
after its guide edit, so it needs task_02. The gate follows every Task. No
other Spec is required.
