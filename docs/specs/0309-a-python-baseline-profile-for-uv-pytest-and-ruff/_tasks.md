---
schema: spec-tasks/v1
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
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

# Tasks — A Python Baseline Profile for uv, pytest and ruff

| id      | title                                                                                                  | type    | complexity | needs                     |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------------------------- |
| task_01 | The asset sync imports the upstream python preset at the commit every setup names                      | backend | high       | —                         |
| task_02 | The python profile governs uv, ruff and pytest and dispatches the Python skills                        | backend | high       | task_01                   |
| task_03 | The user guide lists the Python profile and the skills it dispatches                                   | docs    | low        | task_02                   |
| task_04 | Run the final QA gate                                                                                  | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 imports the `python` setup at the commit Spec 0317 gave the other
upstream snapshots, and proves it in `internal/baseline`; the operator
queues this Spec only after Spec 0317 merges. task_02's profile names that setup, and both Tasks
regenerate the same derived files and change
`internal/baseline/derived_ownership_test.go` and
`docs/agents/setup-context.json`, so it needs task_01. task_03 describes the
profile task_02 ships and re-records the Behavior Surface Record after its
guide edit, so it needs task_02. The gate follows every Task. No other Spec
is required.
