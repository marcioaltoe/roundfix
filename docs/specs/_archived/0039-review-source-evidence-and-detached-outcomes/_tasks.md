---
schema: spec-tasks/v1
spec: 0039-review-source-evidence-and-detached-outcomes
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
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_05]
    - id: task_08
      file: task_08.md
      needs: [task_03, task_04, task_06, task_07]
    - id: task_09
      file: task_09.md
      needs: [task_08]
---

# Tasks — Review Source evidence and Detached Run outcomes

Cross-Spec prerequisite: Spec 0037 Terminal Outcome Integrity must be completed
before this graph starts because Stop Request interruption, registered-session
cleanup, and winner-only terminal publication are required boundaries.

| id      | title                                                           | type    | complexity | needs                                  |
| ------- | --------------------------------------------------------------- | ------- | ---------- | -------------------------------------- |
| task_01 | Model Review Source Evidence and transient failures             | backend | high       | —                                      |
| task_02 | Unify CodeRabbit evidence classification                        | backend | high       | task_01                                |
| task_03 | Settle Review Skipped without Round artifacts                  | backend | high       | task_02                                |
| task_04 | Retry transient Review Source failures and project waits        | backend | high       | task_01, task_02                       |
| task_05 | Report Review Issue knowledge and terminal context             | backend | high       | task_03, task_04                       |
| task_06 | Deliver notification receipts and Detached monitoring          | backend | high       | task_05                                |
| task_07 | Inherit evidence only for Daemon artifact commits              | backend | high       | task_02, task_05                       |
| task_08 | Align review-evidence docs and glossary                        | docs    | medium     | task_03, task_04, task_06, task_07     |
| task_09 | Align the protected Roundfix Skill pair                        | docs    | low        | task_08                                |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05 · 5 → task_06, task_07 · 6 → task_08 · 7 → task_09
