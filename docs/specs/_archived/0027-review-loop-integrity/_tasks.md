---
schema: spec-tasks/v1
spec: 0027-review-loop-integrity
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02]
    - id: task_08
      file: task_08.md
      needs: [task_01]
    - id: task_09
      file: task_09.md
      needs: [task_03, task_07]
    - id: task_10
      file: task_10.md
      needs: [task_01, task_07]
    - id: task_11
      file: task_11.md
      needs: [task_05, task_06, task_08, task_09, task_10]
---

# Tasks — Review Loop Integrity

| id      | title                                                              | type    | complexity | needs                                            |
| ------- | ------------------------------------------------------------------ | ------- | ---------- | ------------------------------------------------ |
| task_01 | Add the CleanUnverified terminal state and exit code 3             | backend | low        | —                                                |
| task_02 | Persist a terminal reason on Review Issue artifacts                | backend | low        | —                                                |
| task_03 | Add thread-reply and PR-comment mutations to the Review Source client | backend | medium  | —                                                |
| task_04 | Enumerate pending Run Branch work and fast-forward-integrate it    | backend | medium     | —                                                |
| task_05 | Enforce the Branch Integrity Preflight on fetch, resolve, and watch | backend | high      | task_03, task_04                                 |
| task_06 | Execute review Runs in the user's checkout                         | backend | high       | task_05                                          |
| task_07 | Record terminal reasons from the engine's settle paths             | backend | medium     | task_02                                          |
| task_08 | Confirm Merge-Ready through a grace window or end Clean Unverified | backend | medium     | task_01                                          |
| task_09 | Propagate per-issue outcomes with Outcome Comments                 | backend | high       | task_03, task_07                                 |
| task_10 | Split the report into per-Run and cumulative counts with reasons   | backend | medium     | task_01, task_07                                 |
| task_11 | Sync the Roundfix Skill, glossary, and command docs to the new contract | docs | medium    | task_05, task_06, task_08, task_09, task_10      |

Waves: 1 → task_01, task_02, task_03, task_04 · 2 → task_05, task_07, task_08 · 3 → task_06, task_09, task_10 · 4 → task_11
