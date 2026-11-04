---
schema: spec-tasks/v1
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
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

# Tasks — A Task Context path the Task deletes is not unresolved

| id      | title                                                                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A Task Context path tracked at HEAD or the Delivery Base resolves and the fix stops saying Create | backend | medium     | —                |
| task_02 | The guides list a removed file only under deletes and state the base rule                         | docs    | low        | —                |
| task_03 | Run the final QA gate                                                                              | qa      | medium     | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03

task_01 and task_02 share no declared file and neither reads what the other
writes: task_02's test reads only guide text, so they run in one Wave. Each
Task runs the suites that consume its changed contract in its own
Verification, so no separate Integration Checkpoint is needed. No other Spec
is required.
