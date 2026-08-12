---
schema: spec-tasks/v1
spec: 0012-npm-distribution
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
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03, task_06]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04, task_06]
---

# Tasks — npm Distribution and Skill Bundle

| id      | title                                                            | type    | complexity | needs                             |
| ------- | ---------------------------------------------------------------- | ------- | ---------- | --------------------------------- |
| task_01 | Platform mapping table and per-platform binary package scaffolding| infra   | medium     | —                                 |
| task_02 | Launcher package and pass-through bin shim                       | infra   | medium     | task_01                           |
| task_03 | Upgrade-asset compatibility test pinning the asset-name scheme   | test    | low        | task_01                           |
| task_06 | Roundfix skill bundle: embed, sync, check, install, list         | backend | high       | —                                 |
| task_04 | Tag-triggered release workflow                                   | infra   | high       | task_01, task_02, task_03, task_06|
| task_05 | Docs (install, release runbook, skill bundle) and skill sync     | docs    | low        | task_02, task_04, task_06         |

Waves: 1 → task_01, task_06 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05
