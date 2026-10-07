---
schema: spec-tasks/v1
spec: 0247-a-run-gate-that-runs-the-repository-contracts
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — A Run gate that runs the repository contracts

| id      | title                                                                                   | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | verify-select discovers Repository Contract Tests and selects them by Contract Relevance | backend | medium     | —                         |
| task_02 | make verify-changed runs the selected contracts, and every contract declares its relevance | backend | medium     | task_01                   |
| task_03 | The glossary and the repository rules describe Contract Relevance                       | docs    | low        | —                         |
| task_04 | Run the final QA gate                                                                   | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04

task_02 uses the selector that task_01 builds and shares
`internal/verifyselect` with it, so they run in series. task_03 shares no file
with either and runs beside task_01. The gate follows every Task.
