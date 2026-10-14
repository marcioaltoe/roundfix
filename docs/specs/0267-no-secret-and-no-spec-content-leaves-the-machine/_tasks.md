---
schema: spec-tasks/v1
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — No secret and no Spec content leaves the machine

| id      | title                                                           | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | A QA Report that carries a secret is never committed            | backend | medium     | —                         |
| task_02 | No Agent Session or Verification command receives a judge key   | backend | medium     | —                         |
| task_03 | Every Roundfix push refuses a history that adds a secret        | backend | medium     | task_01, task_02          |
| task_04 | The skills, the QA prompt, the guides and the glossary describe it | docs | low        | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                           | qa      | medium     | task_04                   |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 owns `internal/secretscan` and the QA step; task_02 owns the judge
helper and the two environment builders. They share no declared file and run
together. task_03 needs task_01 for the package and task_02 because both edit
`internal/daemon/daemon.go`. task_04 documents what the three behavior Tasks
built and owns every skill, guide and glossary file, so it runs after them.
No Task depends on another Spec; Spec 0266, delivered first, changes none of
these files.
