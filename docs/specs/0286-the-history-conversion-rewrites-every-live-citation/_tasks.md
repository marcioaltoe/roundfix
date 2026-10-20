---
schema: spec-tasks/v1
spec: 0286-the-history-conversion-rewrites-every-live-citation
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The history conversion rewrites every live citation

| id      | title                                                                                   | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The history plan rewrites every citation in every Citation Source                       | backend | high       | —                |
| task_02 | A lineage line names a retired ADR by its identifier and dates its supersession         | backend | medium     | task_01          |
| task_03 | A pending Task's target survives, and a broken link is resolved by its file name        | backend | high       | task_01, task_02 |
| task_04 | The check reports a removed history citation, and the guides describe the conversion    | backend | high       | task_01, task_03 |
| task_05 | Run the final QA gate                                                                   | qa      | high       | task_04          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 and task_02 both change `internal/spec/history_entries.go`. task_03
adds to the file task_01 creates and asserts Surface Transcript 1, which needs
both earlier slices. task_04 reads the helpers task_01 and task_03 add, shares
`internal/speccheck/history_retained.go` with task_03, and carries every guide,
reference, glossary and skill-version change. So the four run in series, and
the gate follows the last of them. No other Spec is required.
