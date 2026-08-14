---
schema: spec-tasks/v1
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
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
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_03, task_04, task_05, task_06]
---

# Tasks — Context-Driven Baseline coverage and Repository Skill Set restoration

Cross-Spec integration gate: Spec 0036 Task 01's lock compatibility fixture
must agree with the restoration adapter before `restore-skills` writes can be
enabled. Doctor behavior remains owned by Spec 0036.

| id      | title                                             | type    | complexity | needs                                       |
| ------- | ------------------------------------------------- | ------- | ---------- | ------------------------------------------- |
| task_01 | Enforce versioned setup asset contracts           | backend | high       | —                                           |
| task_02 | Generate complete Context-Driven Baseline guidance | docs    | high       | task_01                                     |
| task_03 | Reject unresolved Decision Plan references        | backend | high       | task_01, task_02                            |
| task_04 | Authorize exact setup Change Plans                | backend | high       | task_01, task_03                            |
| task_05 | Prove portable Repository Skill Set snapshots     | backend | high       | task_01, task_04                            |
| task_06 | Restore external Repository Skill Set members     | backend | high       | task_04, task_05                            |
| task_07 | Synchronize setup workflow guidance               | docs    | medium     | task_02, task_03, task_04, task_05, task_06 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_06 · 7 → task_07
