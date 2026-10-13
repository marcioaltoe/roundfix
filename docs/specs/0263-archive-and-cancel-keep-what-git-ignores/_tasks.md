---
schema: spec-tasks/v1
spec: 0263-archive-and-cancel-keep-what-git-ignores
qa: task_04
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
---

# Tasks — Archive and cancel keep what Git ignores

| id      | title                                                                       | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Archive and cancel refuse Ignored Spec Files                                | backend | medium     | —       |
| task_02 | A failed folder removal keeps what was written and prints how to finish     | backend | medium     | task_01 |
| task_03 | The skills and the glossary describe the refusal and the recovery           | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                       | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 and task_02 both edit `internal/cli/archive.go`, `internal/cli/cancel.go`,
their help tests, both user guides and the Behavior Surface Record, and task_02
extends the helper file task_01 creates, so they run in series. task_03 writes
the skill text that describes both and runs after them. The gate follows the
only leaf.
