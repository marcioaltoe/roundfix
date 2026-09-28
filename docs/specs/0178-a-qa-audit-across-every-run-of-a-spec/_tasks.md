---
schema: spec-tasks/v1
spec: 0178-a-qa-audit-across-every-run-of-a-spec
qa: task_05
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Every Run's Task commits reach the audit |
| task_02 | backend | The audit table names each commit |
| task_03 | backend | Staleness is measured against the Delivery Base, and stale warns |
| task_04 | backend | The report names the user-flow binary |
| task_05 | qa | Run the final QA gate |
