---
schema: spec-tasks/v1
spec: 0009-parallel-scheduling
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
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04, task_05]
---

# Tasks — Parallel Scheduling

| id      | title                                                             | type     | complexity | needs                     |
| ------- | ----------------------------------------------------------------- | -------- | ---------- | ------------------------- |
| task_01 | Worktree config: concurrency, location hierarchy, slug paths      | data     | medium     | —                         |
| task_02 | Task Worktree lifecycle and the serialized integration queue      | backend  | high       | —                         |
| task_03 | Ready-set scheduler in TaskCycle                                  | backend  | high       | task_01, task_02          |
| task_04 | Concurrency-correct surfaces: Work Queue, header, report order    | frontend | medium     | task_03                   |
| task_05 | Settle over Task Worktrees and empty-debris reap wiring           | backend  | medium     | task_02, task_03          |
| task_06 | Docs and skill sync                                               | docs     | low        | task_03, task_04, task_05 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04, task_05 · 4 → task_06
