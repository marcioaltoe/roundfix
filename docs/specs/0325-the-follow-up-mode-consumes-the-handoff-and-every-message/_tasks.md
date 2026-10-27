---
schema: spec-tasks/v1
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
qa: task_03
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02, task_04]
---

# Tasks — The follow-up mode consumes the handoff and every message

| id      | title                                                                    | type | complexity | needs            |
| ------- | ------------------------------------------------------------------------ | ---- | ---------- | ---------------- |
| task_01 | The follow-up mode takes the handoff and settles every handed-in message | docs | medium     | —                |
| task_02 | The glossary defines Session Handoff                                     | docs | low        | —                |
| task_04 | The modes test catches the restored "without consuming" sentence       | test | low        | task_01          |
| task_03 | Run the final QA gate                                                    | qa   | high       | task_02, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_04 · 3 → task_03

task_01 and task_02 share no file: task_01 changes the `roundfix` skill, its
mirror, its version record and `skills/roundfix_modes_test.go`, and task_02
changes only `CONTEXT.md`. task_04 is corrective Task 1 of 2: it answers
findings F-01 and F-02 of the QA Report of 2026-10-26 in the test task_01
wrote, and the gate follows it and task_02. No other Spec is required.
