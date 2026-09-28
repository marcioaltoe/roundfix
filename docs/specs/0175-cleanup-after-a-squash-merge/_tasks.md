---
schema: spec-tasks/v1
spec: 0175-cleanup-after-a-squash-merge
qa: task_07
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
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_01, task_02, task_03, task_04, task_05, task_06]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Prove a merged Spec's Run against the merged head |
| task_02 | backend | Key legacy Runs from their Run Worktree |
| task_03 | backend | Reconcile reads the Delivery Queue merge record |
| task_04 | backend | Release a merged Spec's Runs after a delivery merge |
| task_05 | backend | Own and sweep carry-forward staging worktrees |
| task_06 | backend | Protect nested worktrees and lock every worktree administration |
| task_07 | qa | Run the final QA gate |
