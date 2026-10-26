---
schema: spec-tasks/v1
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
qa: task_04
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
      needs: [task_01, task_02, task_03]
---

# Tasks — The status, the Park Class table and write-tasks say what to do next

| id      | title                                                                  | type    | complexity | needs                     |
| ------- | ---------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The follow-up status names the manual release classification           | backend | medium     | —                         |
| task_02 | The Park Class table lists blockers and states the secret-found action | docs    | medium     | task_01                   |
| task_03 | write-tasks tells authors to declare a shared test file in every Task  | docs    | low        | task_02                   |
| task_04 | Run the final QA gate                                                  | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 and task_02 both raise the `roundfix` skill version, edit
`skills/testdata/owned-skill-versions.json` and re-record
`docs/references/behavior-surfaces.json`, so task_02 needs task_01. task_03
raises the `write-tasks` version in the same record, so it needs task_02. The
gate follows all three. No other Spec is required.
