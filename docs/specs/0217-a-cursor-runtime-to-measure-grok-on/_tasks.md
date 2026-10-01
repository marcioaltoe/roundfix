---
schema: spec-tasks/v1
spec: 0217-a-cursor-runtime-to-measure-grok-on
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
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — A Cursor runtime to measure Grok on

| id      | title                                                                 | type    | complexity | needs                              |
| ------- | --------------------------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | `cursor` is an ACP Runtime whose selection names the advertised value | backend | high       | —                                  |
| task_02 | A Cursor selection without the maintainer's login is refused, never logged in | backend | medium | task_01                       |
| task_03 | The Roundfix Skill describes the Cursor runtime                       | docs    | low        | task_01, task_02                   |
| task_04 | Grok through Cursor, measured on two replayed Tasks                   | chore   | high       | task_01, task_02                   |
| task_05 | Run the final QA gate                                                 | qa      | high       | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05
