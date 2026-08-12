---
schema: spec-tasks/v1
spec: 0005-tui-cockpit
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_02]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Tasks — TUI Cockpit Redesign

| id      | title                                                            | type     | complexity | needs                     |
| ------- | ---------------------------------------------------------------- | -------- | ---------- | ------------------------- |
| task_01 | Renderer decomposition, rendering-neutral                        | frontend | medium     | —                         |
| task_02 | Two-pane base layout and phase row                               | frontend | medium     | task_01                   |
| task_03 | Work Queue upgrade: separators, markers, severity, totals footer | frontend | medium     | task_02                   |
| task_04 | Timeline grouping by Batch and event kind                        | frontend | medium     | task_02                   |
| task_05 | Modal detail for both Work Item kinds                            | frontend | high       | task_02                   |
| task_06 | Footer states and responsive fallback                            | frontend | low        | task_03, task_04, task_05 |
| task_07 | Attach parity pass and docs/skill sync                           | docs     | low        | task_06                   |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04, task_05 · 4 → task_06 · 5 → task_07
