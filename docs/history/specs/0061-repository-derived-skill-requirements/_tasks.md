---
schema: spec-tasks/v1
spec: 0061-repository-derived-skill-requirements
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
      needs: [task_01, task_02, task_03]
---

# Tasks — Repository-derived skill requirements

| id      | title                                                         | type    | complexity | needs                    |
| ------- | ------------------------------------------------------------- | ------- | ---------- | ------------------------ |
| task_01 | Accept the external requirement instead of deciding it         | backend | low        | —                        |
| task_02 | Derive the requirement from the repository's Setup Manifest    | backend | medium     | task_01                  |
| task_03 | Report every missing skill with a per-skill install command    | backend | low        | task_01                  |
| task_04 | Align the Roundfix Skill pair and the derived digests          | docs    | low        | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04
