---
schema: spec-tasks/v1
spec: 0172-a-qa-gate-that-tells-the-truth
qa: task_06
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
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_07]
    - id: task_07
      file: task_07.md
      needs: [task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A report starts pending and a hollow pass is refused |
| task_02 | backend | The allocator writes where the reader looks |
| task_03 | backend | A temporary retry keeps the deterministic failure |
| task_04 | backend | An unobserved Verification is classified |
| task_05 | backend | The event stream survives |
| task_07 | backend | The journal consumer corpus replays the new event stream signature |
| task_06 | qa | Run the final QA gate |
