---
schema: spec-tasks/v1
spec: 0046-public-context-driven-baseline-command
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
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_05, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_04]
    - id: task_09
      file: task_09.md
      needs: [task_07, task_08]
    - id: task_10
      file: task_10.md
      needs: [task_07]
    - id: task_11
      file: task_11.md
      needs: [task_09, task_10]
    - id: task_12
      file: task_12.md
      needs: [task_11]
    - id: task_13
      file: task_13.md
      needs: [task_09]
    - id: task_14
      file: task_14.md
      needs: [task_09]
    - id: task_15
      file: task_15.md
      needs: [task_12, task_13, task_14]
    - id: task_16
      file: task_16.md
      needs: [task_15]
    - id: task_17
      file: task_17.md
      needs: [task_16]
---

# Tasks — Public Context-Driven Baseline Command

| id      | title                                                  | type    | complexity | needs                            |
| ------- | ------------------------------------------------------ | ------- | ---------- | -------------------------------- |
| task_01 | Freeze the Python compatibility corpus                 | test    | high       | —                                |
| task_02 | Establish the embedded Baseline catalog authority      | backend | high       | task_01                          |
| task_03 | Manage repository-owned Baseline Profiles              | backend | medium     | task_02                          |
| task_04 | Inspect bounded Git repository state                    | backend | high       | task_03                          |
| task_05 | Plan root-instruction preservation                      | backend | high       | task_04                          |
| task_06 | Resolve profile alignment decisions                     | backend | high       | task_04                          |
| task_07 | Emit portable Baseline Plans                            | backend | high       | task_05, task_06                 |
| task_08 | Recover interrupted file transactions                   | backend | high       | task_04                          |
| task_09 | Apply approved Baseline Plans                           | backend | high       | task_07, task_08                 |
| task_10 | Classify root instructions through sealed ACP proposals | backend | high       | task_07                          |
| task_11 | Guide human Baseline workflows                          | backend | high       | task_09, task_10                 |
| task_12 | Recalculate rejected plans from scoped proposals        | backend | high       | task_11                          |
| task_13 | Restore Repository Skill Sets through Baseline           | backend | high       | task_09                          |
| task_14 | Synchronize canonical Baseline assets through Go         | backend | high       | task_09                          |
| task_15 | Document the public Baseline operating contract          | docs    | high       | task_12, task_13, task_14        |
| task_16 | Cut over to the Go Baseline authority                    | backend | high       | task_15                          |
| task_17 | Prove release-ready Baseline journeys                    | test    | high       | task_16                          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05, task_06, task_08 · 6 → task_07 · 7 → task_09, task_10 · 8 → task_11, task_13, task_14 · 9 → task_12 · 10 → task_15 · 11 → task_16 · 12 → task_17
