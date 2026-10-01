---
schema: spec-tasks/v1
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
qa: task_04
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
      needs: [task_01, task_02, task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A Claim Receipt is proved, and a held attribution without one is a gap |
| task_02 | backend | A TechSpec's Surface Transcripts are declared, well-formed and traced |
| task_03 | docs | The authoring skills teach receipts and concrete contracts |
| task_04 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
