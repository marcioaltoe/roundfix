---
schema: spec-tasks/v1
spec: 0171-a-deterministic-suite-under-load
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Deadline-bound waits that watch the work |
| task_02 | backend | Implement tests wait on events |
| task_03 | backend | Daemon and worktree tests wait on events |
| task_04 | backend | No live release lookup in tests |
| task_05 | backend | An adapter fixture that survives a relative path |
| task_06 | qa | Run the final QA gate |
