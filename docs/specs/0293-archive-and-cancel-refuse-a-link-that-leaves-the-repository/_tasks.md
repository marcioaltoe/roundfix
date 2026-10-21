---
schema: spec-tasks/v1
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
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
      needs: [task_03]
---

# Tasks — Archive and cancel refuse a link that leaves the repository

| id      | title                                                                    | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The resolver refuses a link that leaves the repository                   | backend | medium     | —       |
| task_02 | Archive and cancel refuse it through the command and the guides say so   | backend | medium     | task_01 |
| task_03 | The Roundfix archive reference and ADR-0269 state the refusal            | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                    | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 changes the resolver and the package functions; task_02 proves the
two commands end to end through the resolver and documents the refusal in the
guides, so it is task_01's Integration Checkpoint. task_03 repeats the guide
wording in the Roundfix skill and records the narrowing in ADR-0269, so it
follows task_02. The gate follows the last of them. No Spec is listed in
`requires`.
