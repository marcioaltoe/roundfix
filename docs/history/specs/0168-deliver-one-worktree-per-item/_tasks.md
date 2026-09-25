---
schema: spec-tasks/v1
spec: 0168-deliver-one-worktree-per-item
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
      needs: [task_01, task_02, task_03, task_04, task_06, task_07, task_08]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Record each item's worktree |
| task_02 | backend | Run every item in its own worktree |
| task_03 | backend | Park, resume and merge without the restore machinery |
| task_04 | docs | Describe item worktrees in the shipped skill and the guide |
| task_06 | backend | Provision every item worktree, and let migrated merged items rest |
| task_07 | backend | A half-removed item worktree never blocks the queue |
| task_08 | backend | Remove only a truly half-removed item worktree |
| task_05 | qa | Run the final QA gate |
