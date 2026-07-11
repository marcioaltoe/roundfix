---
schema: spec-tasks/v1
spec: 0026-model-fallback-guardrail
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
      needs: [task_02, task_03]
---

# Tasks — Model Fallback Guardrail

| id      | title                                                     | type    | complexity | needs            |
| ------- | --------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Prove a Fallback Selection with disposable sessions       | backend | medium     | —                |
| task_02 | Report the fallback on non-interactive selection failures | backend | medium     | task_01          |
| task_03 | Confirm the fallback interactively before Run creation    | backend | medium     | task_01, task_02 |
| task_04 | Ship fallback guardrail guidance                          | docs    | low        | task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
