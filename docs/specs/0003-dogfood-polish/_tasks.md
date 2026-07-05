---
schema: spec-tasks/v1
spec: 0003-dogfood-polish
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
      needs: []
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_07
      file: task_07.md
      needs: []
    - id: task_08
      file: task_08.md
      needs: []
    - id: task_09
      file: task_09.md
      needs: [task_01, task_05, task_06]
---

# Tasks — Dogfood Polish

| id      | title                                                        | type     | complexity | needs                     |
| ------- | ------------------------------------------------------------ | -------- | ---------- | ------------------------- |
| task_01 | Lowercase task commit subjects; unscope the QA commit        | backend  | low        | —                         |
| task_02 | Remove Budget and Round lines from the implement header      | backend  | low        | —                         |
| task_03 | Hermetic git test helpers across packages                    | test     | medium     | —                         |
| task_04 | Bounded stderr tail in agent infrastructure errors           | backend  | low        | —                         |
| task_05 | Stop by spec target                                          | backend  | low        | —                         |
| task_06 | QA field in implement Interactive Input                      | frontend | low        | —                         |
| task_07 | Spec-Run agent logs under the Artifact Directory             | backend  | low        | —                         |
| task_08 | Spec discovery diagnostics in the picker                     | backend  | low        | —                         |
| task_09 | Docs and skill sync                                          | docs     | low        | task_01, task_05, task_06 |

Waves: 1 → task_01, task_02, task_03, task_04, task_05, task_06, task_07, task_08 · 2 → task_09
