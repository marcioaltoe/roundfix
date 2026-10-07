---
schema: spec-tasks/v1
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
qa: task_05
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A formatted QA report and a reopen for late dependencies

| id      | title                                                                                                  | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | The glossary, the guides and the Roundfix Skill describe the Format Command and the Late Dependency     | docs    | low        | —                |
| task_02 | The Project Config reads the verification.format command                                                | backend | low        | —                |
| task_03 | The QA step formats its QA directory with the Format Command and reverts a failed run                  | backend | high       | task_01, task_02 |
| task_04 | Reopen returns a completed QA gate to pending over a Late Dependency                                   | backend | medium     | task_01          |
| task_05 | Run the final QA gate                                                                                  | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05

No two Tasks share a declared file. task_03 reads the configuration value
task_02 adds, and both behavior Tasks follow task_01, which writes the guides
and the Roundfix Skill references for the surfaces they change. The gate
follows both leaves.
