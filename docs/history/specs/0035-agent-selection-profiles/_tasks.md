---
schema: spec-tasks/v1
spec: 0035-agent-selection-profiles
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
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
    - id: task_09
      file: task_09.md
      needs: [task_08]
    - id: task_10
      file: task_10.md
      needs: [task_09]
---

# Tasks — Agent Selection Profiles

| id      | title                                                        | type     | complexity | needs            |
| ------- | ------------------------------------------------------------ | -------- | ---------- | ---------------- |
| task_01 | Resolve atomic Agent Selection Profiles                      | backend  | high       | —                |
| task_02 | Enforce author-declared Task Types                           | backend  | medium     | task_01          |
| task_03 | Show advisory profile recommendations                        | backend  | high       | task_01          |
| task_04 | Configure and validate complete profiles                     | backend  | high       | task_01, task_03 |
| task_05 | Prove relevant profiles before operational Runs              | backend  | high       | task_02, task_04 |
| task_06 | Route each work action through an owned Agent Session         | backend  | high       | task_05          |
| task_07 | Persist Agent Selection attempts and events                  | data     | high       | task_06          |
| task_08 | Project per-work selection state in the Live Run surfaces    | frontend | high       | task_07          |
| task_09 | Synchronize profile guidance and owned skills                | docs     | medium     | task_08          |
| task_10 | Prove mixed-profile fallback behavior end to end             | test     | high       | task_09          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05 · 5 → task_06 · 6 → task_07 · 7 → task_08 · 8 → task_09 · 9 → task_10
