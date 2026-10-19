---
schema: spec-tasks/v1
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
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
    - id: task_03
      file: task_03.md
      needs: [task_04]
---

# Tasks — The Codex full-access mode is the one the adapter offers

| id      | title                                                                                   | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A full-access request selects the mode the Agent Session advertises                     | backend | high       | —       |
| task_02 | Describe the full-access selection and the sandboxed browser driver                     | docs    | medium     | task_01 |
| task_04 | The proof test tells the advertised mode from the first candidate                       | test    | low        | task_02 |
| task_03 | Run the final QA gate                                                                   | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_04 · 4 → task_03

task_02 documents the behavior task_01 builds, so it follows it. The gate
follows task_02. task_04 is corrective Task 1 of 2: it answers finding F1 of
the QA Report of 2026-10-19 and the gate follows it. No other active Spec is required: Specs 0282 and 0283 share
no file with this one.
