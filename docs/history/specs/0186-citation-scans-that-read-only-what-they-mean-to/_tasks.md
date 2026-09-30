---
schema: spec-tasks/v1
spec: 0186-citation-scans-that-read-only-what-they-mean-to
qa: task_03
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
---

# Task Graph

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | The citation projection applies to the Spec's own Task files only | backend | low | — |
| task_02 | The Relocation Citation scan reads through a repository root | backend | medium | — |
| task_03 | Run the final QA gate | qa | medium | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03
