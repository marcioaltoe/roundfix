---
schema: spec-tasks/v1
spec: 0011-storage-lifecycle
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — Review Artifacts, Run Logs, and Spec Archiving

| id      | title                                                             | type    | complexity | needs                            |
| ------- | ----------------------------------------------------------------- | ------- | ---------- | -------------------------------- |
| task_01 | Review artifact location resolver and spec association            | backend | medium     | —                                |
| task_02 | Opt-in per-Batch agent logs                                       | backend | low        | —                                |
| task_03 | Archive Command with all-completed and QA-passed precondition     | backend | medium     | —                                |
| task_04 | Review Issue title hygiene, status-poll dedup, merge-readiness note| backend | low        | —                                |
| task_05 | Docs and skill sync                                               | docs    | low        | task_01, task_02, task_03, task_04|

Waves: 1 → task_01, task_02, task_03, task_04 · 2 → task_05
