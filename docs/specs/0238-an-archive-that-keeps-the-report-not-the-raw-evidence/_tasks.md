---
schema: spec-tasks/v1
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
qa: task_05
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_04]
---

# Tasks — An archive that keeps the report, not the raw evidence

| id      | title                                                                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | No test or record reads archived QA evidence                                                       | test    | low        | —                |
| task_02 | A normal archive drops the raw QA evidence, keeps a manifest, and the queue accepts it as exact     | backend | high       | —                |
| task_03 | The Archive Command cuts one archived Spec on request, as a dry run unless applied                  | backend | medium     | task_02          |
| task_04 | The skills and the Spec workflow guides state the evidence cut                                     | docs    | medium     | task_03          |
| task_05 | Run the final QA gate                                                                              | qa      | high       | task_01, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 and task_02 share no file. task_03 builds on task_02's planner and
shares the Archive Command and its user guide with it. task_04 documents
task_03's command and runs the derived-file regeneration alone. The gate
depends on both leaves.
