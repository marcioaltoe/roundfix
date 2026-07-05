---
schema: spec-tasks/v1
spec: 0001-implement-command
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
      needs: [task_01, task_02, task_03]
    - id: task_06
      file: task_06.md
      needs: [task_02, task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_06]
    - id: task_09
      file: task_09.md
      needs: [task_06]
    - id: task_10
      file: task_10.md
      needs: [task_07, task_08]
---

# Tasks — Implement Command

| id      | title                                                                    | type     | complexity | needs                     |
| ------- | ------------------------------------------------------------------------ | -------- | ---------- | ------------------------- |
| task_01 | Spec contract parser: discovery, Task Graph, task files, QA verdict      | backend  | high       | —                         |
| task_02 | Run Database v4: work-target locks, implement Kind, working-tree query   | data     | medium     | —                         |
| task_03 | Task and QA prompt builders mirroring implement-task and qa-gate         | backend  | low        | —                         |
| task_04 | Preflight default-branch detection and veto message                      | backend  | low        | —                         |
| task_05 | Engine TaskCycle: per-Task agent, verification, settle, commit, failure policy | backend | high  | task_01, task_02, task_03 |
| task_06 | Implement Command CLI: flags, Preflight Validation, resume, output, exit codes | backend | high  | task_02, task_04, task_05 |
| task_07 | Opt-in QA gate: --qa flag through verdict settling and QA Report commit  | backend  | medium     | task_06                   |
| task_08 | Interactive Input Spec picker                                            | frontend | medium     | task_06                   |
| task_09 | Live Run View renders Tasks as Work Items; Attach parity                 | frontend | medium     | task_06                   |
| task_10 | Roundfix skill and docs update for the Implement Command                 | docs     | low        | task_07, task_08          |

Waves: 1 → task_01, task_02, task_03, task_04 · 2 → task_05 · 3 → task_06 · 4 → task_07, task_08, task_09 · 5 → task_10
