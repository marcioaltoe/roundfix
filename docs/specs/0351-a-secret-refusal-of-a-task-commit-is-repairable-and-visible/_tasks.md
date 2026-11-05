---
schema: spec-tasks/v1
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
qa: task_05
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A secret refusal of a Task commit is repairable and visible

| id      | title                                                                                 | type    | complexity | needs                     |
| ------- | ------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Task commit's Secret Scan runs as the last Settlement Check of every attempt       | backend | high       | —                         |
| task_02 | The Causes report lists each secret refusal as a secret-scan item                     | backend | medium     | task_01                   |
| task_03 | A run-unresolved park names its failed Tasks, their reasons and the skipped Tasks      | backend | medium     | —                         |
| task_04 | The skill, the guides, the glossary and ADR-0361 state the new behavior               | docs    | medium     | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                                 | qa      | medium     | task_04                   |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05

task_02 needs task_01 because it reads `runevent.SettlementCheckSecretScan`,
which task_01 adds. task_03 shares no file with task_01 or task_02, so it runs
in the first Wave. task_04 documents what the three make true. Each
implementation Task runs the existing suites that consume its changed contract
in its own Verification, so no separate Integration Checkpoint is needed. No
other Spec is required.
