---
schema: spec-tasks/v1
spec: 0117-a-window-the-preflight-owns
qa: task_06
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
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — A window the Preflight owns

| id      | title                                    | type    | complexity | needs   |
| ------- | ---------------------------------------- | ------- | ---------- | ------- |
| task_01 | Store the Run Window at schema 13        | data    | high       | —       |
| task_02 | The window command                       | backend | medium     | task_01 |
| task_03 | Preflight refuses a closed window        | backend | high       | task_02 |
| task_04 | Report a Run that may cross the window   | backend | medium     | task_03 |
| task_05 | Document the term and the two bounds     | docs    | medium     | task_04 |
| task_06 | QA gate                                  | qa      | medium     | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_06

The chain is serial by edit locality, not by logic. task_03 and task_04 both
rewrite `internal/cli/implement.go`, and task_02 and task_03 both rewrite the
command surface in `internal/cli/cli.go`. Two Tasks that rewrite one region are
not independently implementable whatever their logic says: run as siblings on
2026-08-26, that shape produced an integration conflict that discarded a Task
which had already passed its Verification.

Each Task's work, references, and Verification live in its own `task_NN.md`.
