---
schema: spec-tasks/v1
spec: 0140-a-spec-traces-the-promises-it-makes
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
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_05
      file: task_05.md
      needs: [task_07]
---

# Tasks — A Spec traces the promises it makes

| id      | title                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Read a declared promise and refuse an undeclared one | backend | medium     | —       |
| task_02 | Trace a declared promise to its Task                 | backend | medium     | task_01 |
| task_03 | Characterize the coined codes and the corpus         | test    | medium     | task_02 |
| task_04 | State the declaration rule where artifacts are written | docs  | medium     | task_03 |
| task_06 | Let the operational fixtures declare their promises   | test    | low        | task_04 |
| task_07 | Let the Daemon fixtures and the archive characterization catch up | test | medium | task_06 |
| task_05 | Run the final QA gate                                | qa      | high       | task_07 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_07 · 7 → task_05
