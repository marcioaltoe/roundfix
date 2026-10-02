---
schema: spec-tasks/v1
spec: 0216-baseline-wording-left-after-the-stack-wave
qa: task_04
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
      needs: [task_01, task_02, task_03]
---

# Tasks — Baseline wording left after the stack wave

| id      | title                                                                 | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The backend bucket clause names what it forbids                       | backend | medium     | —                         |
| task_02 | The backend and frontend guides name their declared workspace         | backend | medium     | task_01                   |
| task_03 | The skill guide tells the Agent to ask for a skill only a person can start | backend | medium | task_02                   |
| task_04 | Run the final QA gate                                                 | qa      | medium     | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
