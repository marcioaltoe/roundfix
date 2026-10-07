---
schema: spec-tasks/v1
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
qa: task_05
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
---

# Tasks — A sanitize that reads older folders and names its refusals

| id      | title                                                                                                     | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The glossary, the history reference and the Roundfix Skill describe Refused Units and failed-qa           | docs    | low        | —                         |
| task_02 | A Legacy Archive Folder's Task Graph is read leniently                                                    | backend | medium     | —                         |
| task_03 | A failed QA without an override becomes the failed-qa disposition                                         | backend | medium     | task_02                   |
| task_04 | The History Sanitize Command lists every Refused Unit and fills the batch past it                         | backend | high       | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                                                     | qa      | high       | task_04                   |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_02 and task_03 both change `internal/spec/archive_record.go`, so they run
in series. task_04 reads the tolerances task_02 adds and the disposition
task_03 adds, and follows task_01, which documents the output it builds. The
gate follows the only leaf.
