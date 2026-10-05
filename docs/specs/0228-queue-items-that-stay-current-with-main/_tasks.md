---
schema: spec-tasks/v1
spec: 0228-queue-items-that-stay-current-with-main
qa: task_05
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
---

# Tasks — Queue items that stay current with main

| id      | title                                                                   | type    | complexity | needs                              |
| ------- | ----------------------------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | The owned-skill record command raises a colliding version               | backend | medium     | —                                  |
| task_02 | A derived merge resolves conflicts confined to declared version lines   | backend | high       | —                                  |
| task_03 | A review-only correction after archive returns the item to review round 2 | backend | medium     | task_02                            |
| task_04 | The skills, guides and version rule describe the three rules            | docs    | low        | task_01                            |
| task_06 | The real skill version hunk is merged and regenerated                   | backend | medium     | task_03                            |
| task_05 | Run the final QA gate                                                   | qa      | high       | task_01, task_02, task_03, task_04, task_06 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_06 · 4 → task_05

task_01 changes the record command and task_02 the derived merge, its
configuration and the repository declaration; they share no file and run
together. task_03 adds the retry proof and wires it in the delivery workflow
file task_02 changes, so it follows task_02. task_04 rewrites the skills and
guides and records both skill versions with the record command task_01
changes, so it follows task_01. task_06 is the corrective Task for QA finding
F1 of 2026-10-05: it changes the line-scoped resolver in the delivery workflow
file task_02 and task_03 change, and the configuration guide task_02 changes,
so it follows task_03. The gate follows all five.
