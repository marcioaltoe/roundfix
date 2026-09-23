---
schema: spec-tasks/v1
spec: 0154-evidence-when-the-target-branch-is-gone
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Ask the default branch when the target is gone |
| task_02 | backend | Preserved reasons that name the missing proof |
| task_03 | docs | Describe the fallback in the shipped skill |
| task_04 | qa | Run the final QA gate |
