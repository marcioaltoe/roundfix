---
schema: spec-tasks/v1
spec: 0029-launch-and-recovery-fixes
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: [task_01, task_03, task_04, task_05]
---

# Tasks — Launch and Recovery Fixes

| id      | title                                                                  | type    | complexity | needs                                |
| ------- | ---------------------------------------------------------------------- | ------- | ---------- | ------------------------------------ |
| task_01 | Split the Detached Run handshake into liveness and Run-creation phases | backend | high       | —                                    |
| task_02 | Classify Agent Model not-advertised rejections from acpx stderr        | backend | medium     | —                                    |
| task_03 | Settle Batches with the model-rejection reason instead of protocol error | backend | medium   | task_02                              |
| task_04 | Report the effective Agent Model probe in the Doctor Command           | backend | low        | task_02                              |
| task_05 | Prefer failed-Task surfaces in settle and report the chosen surface    | backend | medium     | —                                    |
| task_06 | Sync the Roundfix Skill and docs to the shipped behavior               | docs    | medium     | task_01, task_03, task_04, task_05   |

Waves: 1 → task_01, task_02, task_05 · 2 → task_03, task_04 · 3 → task_06
