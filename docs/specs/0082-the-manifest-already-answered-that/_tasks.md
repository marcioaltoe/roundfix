---
schema: spec-tasks/v1
spec: 0082-the-manifest-already-answered-that
qa: task_08
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
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_02, task_03]
    - id: task_07
      file: task_07.md
      needs: [task_05, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
---

# Tasks — Baseline update

| id      | title                                                        | type    | complexity | needs            |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Pin the adoption path with a characterization corpus         | test    | medium     | —                |
| task_02 | Refresh managed regions without asking who owns the prose    | backend | high       | task_01          |
| task_03 | Read the Setup Manifest back into plan inputs                | backend | medium     | task_01          |
| task_04 | Refresh one repository from its own manifest                 | backend | high       | task_02, task_03 |
| task_05 | Carry the repository's skills with its guidance              | backend | medium     | task_04          |
| task_06 | Ask only what the manifest does not already answer           | backend | medium     | task_02, task_03 |
| task_07 | Teach the update path to the docs and the owned skills       | docs    | medium     | task_05, task_06 |
| task_08 | Run the final QA gate                                        | qa      | high       | task_07          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04, task_06 · 4 → task_05 · 5 → task_07 · 6 → task_08
