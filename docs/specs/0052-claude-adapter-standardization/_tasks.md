---
schema: spec-tasks/v1
spec: 0052-claude-adapter-standardization
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_02]
    - id: task_07
      file: task_07.md
      needs: [task_04, task_05, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
---

# Tasks — Official Claude adapter and opaque model identifiers

| id      | title                                                           | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Prove official Claude adapter lineage                           | backend | high       | —                         |
| task_02 | Make advertised model identifiers opaque                        | backend | high       | —                         |
| task_03 | Default the frontend profile to the proven Claude tuple         | backend | low        | task_02                   |
| task_04 | Report multi-runtime Adapter Readiness in Doctor                | backend | medium     | task_01, task_03          |
| task_05 | Migrate stale Claude overrides through Setup                    | backend | high       | task_01                   |
| task_06 | Explain the Preflight fallback boundary                         | backend | low        | task_02                   |
| task_07 | Align adapter docs and user guide                               | docs    | medium     | task_04, task_05, task_06 |
| task_08 | Align the protected Roundfix Skill pair and derived digest pins | docs    | low        | task_07                   |

Waves: 1 → task_01, task_02 · 2 → task_03, task_05, task_06 · 3 → task_04 · 4 → task_07 · 5 → task_08
