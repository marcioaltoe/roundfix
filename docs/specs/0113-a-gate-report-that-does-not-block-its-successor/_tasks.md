---
schema: spec-tasks/v1
spec: 0113-a-gate-report-that-does-not-block-its-successor
qa: task_07
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_08
      file: task_08.md
      needs: [task_06]
    - id: task_09
      file: task_09.md
      needs: [task_08]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_04, task_06, task_08, task_09]
---

# Tasks — A gate report that does not block its successor

| id      | title                                                      | type    | complexity | needs                     |
| ------- | ---------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Record the refusal as the report's terminal row              | backend | medium     | —                         |
| task_02 | Prove a refused report does not block its successor          | backend | medium     | task_01                   |
| task_03 | Name the literal the blocked cause requires                  | backend | low        | —                         |
| task_04 | Report one cause once                                        | backend | medium     | task_03                   |
| task_05 | Let a Spec's own coined term reach its gate                  | backend | medium     | —                         |
| task_06 | Perform the repairs the Task names                           | backend | high       | task_05                   |
| task_08 | Wire the gate's new inputs into the request it receives      | backend | high       | task_06                   |
| task_09 | A refusal row records a refusal, and the precondition breaks no journey | backend | high | task_08 |
| task_07 | Run the final QA gate                                        | qa      | high       | task_02, task_04, task_06, task_08, task_09 |

Wave plan: `1 → task_01, task_03, task_05 · 2 → task_02, task_04, task_06 · 3 → task_07`.

task_02 follows task_01 rather than running beside it because it proves the
successor case against the report task_01 teaches the writer to produce; task_04
follows task_03 because both report on the same parse result and would otherwise
edit the same detector; task_06 follows task_05 because the vocabulary update is
the first repair the gate can actually perform.
