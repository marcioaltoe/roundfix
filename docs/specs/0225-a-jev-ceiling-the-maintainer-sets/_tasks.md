---
schema: spec-tasks/v1
spec: 0225-a-jev-ceiling-the-maintainer-sets
qa: task_03
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
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The guides and the Roundfix Skill describe the configured Jev ceiling |
| task_02 | backend | The judge and the Jev Router gate read the User Config Jev ceiling |
| task_03 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 changes the `spec judge` command and the configuration model, whose
guides and skill references task_01 writes, so it follows task_01; they share
no declared file. The gate follows both.
