---
schema: spec-tasks/v1
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The Daemon records an Evidence Snapshot for every row that can carry |
| task_02 | backend | A later pass carries unmoved rows and names why every other row re-runs |
| task_03 | backend | The next pass reads a failed pass's QA Report from its Run Branch |
| task_04 | docs | The qa-gate skill, the QA prompt and the guide teach the carry |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
