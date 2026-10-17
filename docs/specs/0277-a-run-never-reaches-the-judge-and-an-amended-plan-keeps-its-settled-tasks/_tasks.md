---
schema: spec-tasks/v1
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04]
---

# Tasks — A Run never reaches the judge, and an amended plan keeps its settled Tasks

| id      | title                                                                                                   | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | No Agent or Verification of a Run can make the judge send, and the TypeSafe key stays out of Runs        | backend | high       | —                |
| task_02 | The bootstrap test reads only its own release, and every test compiles for Windows                       | test    | high       | —                |
| task_03 | Carry-forward keeps a settled Task whose Binding Text an amendment left unchanged                        | backend | high       | task_01          |
| task_04 | The history cleanup reads every ADR reference form adopters use                                          | backend | medium     | task_03          |
| task_05 | Run the final QA gate                                                                                    | qa      | high       | task_02, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01, task_03 and task_04 each raise the Roundfix skill version and
re-record `skills/testdata/owned-skill-versions.json` and the Behavior Surface
Record, and task_01 and task_04 both edit the `history` guide and the Roundfix
`archive` reference, so the three run in series. task_02 changes only tests and
the coverage record, shares no file with them, and runs beside task_01. The
gate follows both leaves.
