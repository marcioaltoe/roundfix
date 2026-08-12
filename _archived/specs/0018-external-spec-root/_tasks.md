---
schema: spec-tasks/v1
spec: 0018-external-spec-root
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
      needs: [task_02, task_03]
---

# Tasks — External Spec Root

| id      | title                                                             | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | specs.root config with resolution and the external predicate      | backend | low        | —                |
| task_02 | Spec Root threading through every Spec consumer                   | backend | high       | task_01          |
| task_03 | Commit boundary: staging guard and settle without a commit        | backend | medium     | task_02          |
| task_04 | Docs and skill sync for the External Spec Root surface            | docs    | low        | task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
