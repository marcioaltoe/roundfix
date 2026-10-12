---
schema: spec-tasks/v1
spec: 0256-an-archive-record-that-survives-a-squash-merge
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_04]
---

# Tasks — An Archive Record that survives a squash merge

| id      | title                                                                                                  | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | The skills, guides and glossary describe the record's pull request, its recovery and the archive warning | docs    | medium     | —                |
| task_02 | Archive-root readers read only real records and recover a squash-dropped source revision                | backend | high       | task_01          |
| task_03 | The Archive Command records the pull request and warns about a squash-dropped revision                  | backend | medium     | task_01          |
| task_04 | The Delivery Queue records the pull request in a Pull Request Record Commit                              | backend | high       | task_03          |
| task_05 | Run the final QA gate                                                                                  | qa      | high       | task_02, task_04, task_06 |
| task_06 | The runs guidance admits Source Revision Recovery and the delivery texts match the TechSpec (corrective, QA-001/QA-002 of 2026-10-12) | docs | low | task_04 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_06 · 5 → task_05

task_01 writes every skill, guide and glossary text first, so each code Task
depends on the Task that documents its surface, and the owned skills' versions
rise once. task_01 and task_03 both declare
`docs/references/behavior-surfaces.json`, and task_03 follows task_01. task_02
and task_03 share no declared file, so they run in one wave. task_04 runs the
Archive Command's new flag through the step executable, so it needs task_03.
The gate follows both leaves. No Task depends on another active Spec.
