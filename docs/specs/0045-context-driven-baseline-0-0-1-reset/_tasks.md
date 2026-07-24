---
schema: spec-tasks/v1
spec: 0045-context-driven-baseline-0-0-1-reset
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_05, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_05, task_07]
    - id: task_09
      file: task_09.md
      needs: [task_03, task_08]
    - id: task_10
      file: task_10.md
      needs: [task_01]
    - id: task_11
      file: task_11.md
      needs: [task_09, task_10]
    - id: task_12
      file: task_12.md
      needs: [task_08, task_09, task_11]
    - id: task_13
      file: task_13.md
      needs: [task_08, task_09, task_10, task_11, task_12]
---

# Tasks — Context-Driven Baseline 0.0.1 reset

| id      | title                                               | type    | complexity | needs                                          |
| ------- | --------------------------------------------------- | ------- | ---------- | ---------------------------------------------- |
| task_01 | Enforce the 0.0.1 Source Baseline contract          | backend | high       | —                                              |
| task_02 | Publish the project-agnostic governed corpus         | docs    | high       | task_01                                        |
| task_03 | Render exact skill activation bundles               | backend | high       | task_01, task_02                               |
| task_04 | Evaluate local Repository Capability evidence       | backend | high       | task_01                                        |
| task_05 | Deliver the Standard TypeScript Monorepo Profile     | backend | high       | task_02, task_03, task_04                      |
| task_06 | Inventory incompatible Source Baselines              | backend | high       | task_01                                        |
| task_07 | Resolve individual Readoption dispositions           | backend | high       | task_02, task_05, task_06                      |
| task_08 | Apply Baseline Readoption through one Change Plan    | backend | high       | task_05, task_07                               |
| task_09 | Align maintained profiles with the Repository Skill Set | backend | high    | task_03, task_08                               |
| task_10 | Plan the release history reset                       | backend | high       | task_01                                        |
| task_11 | Align the Roundfix 0.0.1 distribution                | infra   | high       | task_09, task_10                               |
| task_12 | Prove complete profile journeys                      | test    | high       | task_08, task_09, task_11                      |
| task_13 | Document the 0.0.1 operating contract                | docs    | medium     | task_08, task_09, task_10, task_11, task_12    |

Waves: 1 → task_01 · 2 → task_02, task_04, task_06, task_10 · 3 → task_03 · 4 → task_05 · 5 → task_07 · 6 → task_08 · 7 → task_09 · 8 → task_11 · 9 → task_12 · 10 → task_13
