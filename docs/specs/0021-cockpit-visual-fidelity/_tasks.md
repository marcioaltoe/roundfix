---
schema: spec-tasks/v1
spec: 0021-cockpit-visual-fidelity
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
---

# Tasks — Cockpit Visual Fidelity

| id      | title                                                          | type     | complexity | needs                     |
| ------- | -------------------------------------------------------------- | -------- | ---------- | ------------------------- |
| task_01 | Style tokens and color-mode resolution                         | frontend | medium     | —                         |
| task_02 | Header, Phase Row, and Work Queue fidelity                     | frontend | high       | task_01                   |
| task_03 | Timeline fidelity: groups, gutter, bounded summaries           | frontend | high       | task_01                   |
| task_04 | Detail Modal fidelity and pane empty states                    | frontend | medium     | task_01                   |
| task_05 | Docs and skill sync for the cockpit visual contract            | docs     | low        | task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_03, task_04 · 3 → task_05
