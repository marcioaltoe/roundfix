---
schema: spec-tasks/v1
spec: 0013-codex-runtime-hygiene
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — Codex Runtime Hygiene

| id      | title                                                        | type    | complexity | needs            |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Extract shared read-only health checks from the Setup Command| backend | medium     | —                |
| task_02 | Codex hygiene inspector (quarantine and notarization)        | backend | medium     | —                |
| task_03 | Doctor Command                                               | backend | medium     | task_01, task_02 |
| task_04 | Verified-clean codex on the codex-acp spawn path             | backend | medium     | task_02          |
| task_05 | Docs and skill sync                                          | docs    | low        | task_03, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05
