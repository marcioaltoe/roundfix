---
schema: spec-tasks/v1
spec: 0310-the-python-cli-profile-writes-its-starter-files
requires: [0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff]
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

# Tasks — The Python CLI profile writes its Starter Files

| id      | title                                                                           | type    | complexity | needs                     |
| ------- | ------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The python module declares its Starter Files and the catalog validates them     | backend | high       | —                         |
| task_02 | Starter Files are previewed and written only where absent on a confirmed digest | backend | high       | task_01                   |
| task_03 | roundfix baseline starters previews and writes a profile's Starter Files        | backend | high       | task_02                   |
| task_04 | Run the final QA gate                                                           | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 adds the declarations and the Python starters to the catalog. task_02
builds the preview and the confirmed write over task_01's `StarterFiles`.
task_03 exposes task_02's operation as a command and documents it, and is the
integration checkpoint for task_02's payload. The chain is serial because each
Task consumes the previous one's code. The Spec requires Spec 0309, whose
`python` module and `python-cli` profile task_01 extends and every
Verification command exercises.
