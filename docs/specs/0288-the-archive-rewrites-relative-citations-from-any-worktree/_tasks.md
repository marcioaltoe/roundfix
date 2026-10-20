---
schema: spec-tasks/v1
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
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

# Tasks — The archive rewrites relative citations from any worktree

| id      | title                                                                  | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The archive and cancel build repository paths from one root spelling   | backend | medium     | —       |
| task_02 | Archive and cancel work from every spelling of the checkout            | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                  | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 adds the resolver and fixes the package functions; task_02 calls the
resolver from the two commands and proves them end to end from every
spelling of a linked worktree, so it is task_01's Integration Checkpoint. The
gate follows it. No other Spec is required.
