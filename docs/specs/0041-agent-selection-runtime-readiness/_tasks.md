---
schema: spec-tasks/v1
spec: 0041-agent-selection-runtime-readiness
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
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_04]
    - id: task_08
      file: task_08.md
      needs: [task_04]
    - id: task_09
      file: task_09.md
      needs: [task_05, task_06, task_07, task_08]
---

# Tasks — Agent Selection Runtime Readiness

| id      | title                                                    | type    | complexity | needs                                  |
| ------- | -------------------------------------------------------- | ------- | ---------- | -------------------------------------- |
| task_01 | Prove official Codex adapter readiness                   | backend | high       | —                                      |
| task_02 | Expose bounded ACP selection capabilities                | backend | high       | task_01                                |
| task_03 | Prove exact advertised Agent Selections                  | backend | high       | task_02                                |
| task_04 | Centralize Agent Selection Profile readiness             | backend | high       | task_03                                |
| task_05 | Enforce complete one-Run Agent Selection overrides       | backend | medium     | task_04                                |
| task_06 | Generate and persist only proven Setup profiles          | backend | high       | task_04                                |
| task_07 | Prove profile configuration before persistence           | backend | medium     | task_04                                |
| task_08 | Diagnose effective Agent Selection Profiles in Doctor    | backend | high       | task_04                                |
| task_09 | Synchronize Agent Selection readiness guidance           | docs    | medium     | task_05, task_06, task_07, task_08     |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05, task_06, task_07, task_08 · 6 → task_09
