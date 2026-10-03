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
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03]
---

# Tasks — A Cursor runtime to measure Grok on

| id      | title                                                                 | type    | complexity | needs                              |
| ------- | --------------------------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | `cursor` is an ACP Runtime whose selection names the advertised value | backend | high       | —                                  |
| task_02 | A Cursor selection without the maintainer's login is refused, never logged in | backend | medium | task_01                       |
| task_03 | The Roundfix Skill describes the Cursor runtime                       | docs    | low        | task_01, task_02                   |
| task_05 | Run the final QA gate | qa | high | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_05

Task 04's measurement was run on the maintainer's machine by the operator on 2026-10-03, because the Agent sandbox blocks the Cursor API; its record is `measurement/grok-through-cursor.md`, its recording `internal/agent/testdata/cursor-session-recorded.json`, and its test `internal/agent/cursor_recorded_session_test.go`.
