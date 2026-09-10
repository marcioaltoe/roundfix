---
schema: spec-tasks/v1
spec: 0132-a-grant-read-exactly-where-it-lives
qa: task_07
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
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_02, task_03, task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Tasks — A grant read exactly where it lives

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Characterize today's parse, resolution and discovery answers | test | medium | — |
| task_02 | Close the frontmatter on a complete marker line | backend | medium | task_01 |
| task_03 | Derive the record path from the resolved Spec Root | backend | high | task_01 |
| task_04 | Resolve a citation against the artifact that carries it | backend | high | task_01, task_03 |
| task_05 | Recognize the record naming already in use | backend | low | task_01 |
| task_06 | Read a record wherever the archive left it | test | medium | task_02, task_03, task_04, task_05 |
| task_07 | Run the final QA gate | qa | high | task_06 |

Waves: 1 → task_01 · 2 → task_02, task_03, task_05 · 3 → task_04 · 4 → task_06 · 5 → task_07.
