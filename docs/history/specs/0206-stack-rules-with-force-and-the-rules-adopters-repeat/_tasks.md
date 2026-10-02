---
schema: spec-tasks/v1
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — Stack rules with force, and the rules adopters repeat

| id      | title                                                        | type    | complexity | needs                              |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ---------------------------------- |
| task_01 | Recurring rules become core clauses, and lint warnings block | backend | high       | —                                  |
| task_02 | Go, CLI and TUI rules carry force                            | backend | medium     | task_01                            |
| task_03 | Rust rules carry force, with the typed-error policy          | backend | medium     | task_02                            |
| task_04 | Recurring Spec-workflow and TypeScript rules become clauses  | backend | medium     | task_03                            |
| task_05 | Run the final QA gate                                        | qa      | high       | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
