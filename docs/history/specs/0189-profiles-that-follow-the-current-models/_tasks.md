---
schema: spec-tasks/v1
spec: 0189-profiles-that-follow-the-current-models
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The Model Catalog, the picker efforts and the adapter floors follow today's adapters |
| task_02 | backend | The recommendation is a dated Recommended Profile |
| task_03 | backend | Built-in selections are the Recommended Profile, and `gpt-5.5` leaves production code |
| task_04 | docs | The model reference states the shipped snapshot |
| task_05 | qa | Run the final QA gate |
| task_06 | backend | An override equal to a fallback swaps places with the preferred |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05
