---
schema: spec-tasks/v1
spec: 0047-context-driven-guidance-composition
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
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_04, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_06, task_07]
    - id: task_09
      file: task_09.md
      needs: [task_07, task_08]
---

# Tasks — Context-Driven guidance composition

| id      | title                                       | type    | complexity | needs                    |
| ------- | ------------------------------------------- | ------- | ---------- | ------------------------ |
| task_01 | Render the Instruction Hierarchy            | backend | high       | —                        |
| task_02 | Complete ADR and Findings contracts         | backend | medium     | task_01                  |
| task_03 | Segment preserved rules without byte loss   | backend | high       | task_01                  |
| task_04 | Distribute rules to semantic owners         | backend | high       | task_03                  |
| task_05 | Plan repository-owned Profile adaptations   | backend | high       | task_01                  |
| task_06 | Guide public Profile divergence resolution  | backend | high       | task_05                  |
| task_07 | Synchronize composed Baseline assets        | infra   | high       | task_02, task_04, task_06 |
| task_08 | Document the composed Baseline workflow     | docs    | medium     | task_06, task_07         |
| task_09 | Prove composed Baseline journeys            | test    | high       | task_07, task_08         |

Waves: 1 → task_01 · 2 → task_02, task_03, task_05 · 3 → task_04, task_06 · 4 → task_07 · 5 → task_08 · 6 → task_09
