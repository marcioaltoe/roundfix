---
schema: spec-tasks/v1
spec: 0048-context-driven-project-decisions-and-spec-constraints
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_03, task_04, task_05, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
---

# Tasks — Context-Driven project decisions and Spec constraints

External gate: the newest QA Report for Spec 0047 must carry
`verdict: pass` before task_01 can enter `in_progress`.

| id      | title                                  | type    | complexity | needs                                        |
| ------- | -------------------------------------- | ------- | ---------- | -------------------------------------------- |
| task_01 | Model typed project decisions          | backend | high       | —                                            |
| task_02 | Derive Better Auth HTTP policy         | backend | high       | task_01                                      |
| task_03 | Render confirmed project guidance      | backend | medium     | task_02                                      |
| task_04 | Protect tooling authority              | backend | medium     | task_01                                      |
| task_05 | Add Project Constraints to authoring   | docs    | high       | task_03, task_04                             |
| task_06 | Enforce Project Constraints downstream | docs    | high       | task_05                                      |
| task_07 | Synchronize project-decision assets    | infra   | high       | task_02, task_03, task_04, task_05, task_06 |
| task_08 | Prove project-constraint journeys      | test    | high       | task_07                                      |

Waves: external QA gate → 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05 · 5 → task_06 · 6 → task_07 · 7 → task_08
