---
schema: spec-tasks/v1
spec: 0094-one-history-root-under-docs
qa: task_09
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02]
    - id: task_08
      file: task_08.md
      needs: [task_02]
    - id: task_10
      file: task_10.md
      needs: [task_07]
    - id: task_11
      file: task_11.md
      needs: [task_10]
    - id: task_12
      file: task_12.md
      needs: [task_06]
    - id: task_13
      file: task_13.md
      needs: [task_06]
    - id: task_14
      file: task_14.md
      needs: [task_13]
    - id: task_15
      file: task_15.md
      needs: [task_07]
    - id: task_09
      file: task_09.md
      needs: [task_06, task_07, task_08, task_10, task_11, task_12, task_13, task_14, task_15]
---

# Tasks — One history root under `docs/`

| id      | title                                                       | type    | complexity | needs                     |
| ------- | ----------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Classify a retired decision apart from a pending one         | backend | low        | —                         |
| task_02 | Resolve one history root and move this repository into it    | backend | medium     | —                         |
| task_03 | Decide a review's liveness from local Git alone              | backend | medium     | task_02                   |
| task_04 | Discover retired material sitting on an older layout         | backend | medium     | task_01, task_02, task_03 |
| task_05 | Plan a history relocation as a ledger of identities          | backend | medium     | task_04                   |
| task_06 | Apply and roll back a relocation inside the transaction      | backend | high       | task_05                   |
| task_07 | State the history location in every carrier that names it    | docs    | medium     | task_02                   |
| task_08 | Prove review reaches no history tree                         | test    | low        | task_02                   |
| task_10 | Repoint the pinned archive-path assertion                     | test    | low        | task_07                   |
| task_11 | Repoint the derived-artifact test's archive path              | test    | low        | task_10                   |
| task_12 | Admit a real fleet archive into the migration                 | backend | medium     | task_06                   |
| task_13 | Report what the migration moved and what it kept              | backend | medium     | task_06                   |
| task_14 | Refuse a collision and stop calling the tree current          | backend | high       | task_13                   |
| task_15 | Correct the two carriers that still name the old destination  | docs    | low        | task_07                   |
| task_09 | Run the final QA gate                                        | qa      | high       | task_06…task_08, task_10…task_15 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_07, task_08 · 3 → task_04, task_10 · 4 → task_05, task_11 · 5 → task_06 · 6 → task_12, task_13 · 7 → task_14, task_15 · 8 → task_09
