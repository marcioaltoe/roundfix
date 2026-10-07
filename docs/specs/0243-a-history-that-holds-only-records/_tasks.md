---
schema: spec-tasks/v1
spec: 0243-a-history-that-holds-only-records
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A history that holds only records

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A Legacy Archive Folder converts to an Archive Record that Git still backs                    | backend | high       | —                |
| task_02 | Retired Findings and Backlog Entries reduce, and retired reviews and handoffs leave           | backend | medium     | —                |
| task_03 | roundfix history sanitize plans the history and applies one batch after the History Full Tag | backend | high       | task_01, task_02 |
| task_04 | The skill, a Baseline clause, the glossary and the export contract describe the sanitize      | docs    | medium     | —                |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_02, task_04 · 2 → task_03 · 3 → task_05

task_01 and task_02 add separate files to `internal/spec` and share no
declared file, so they run in parallel. task_03 composes both into the
command, so it follows them. task_04 changes only guidance, the Baseline
module and its derived files, `CHANGELOG.md` and a documentation test. It
shares no declared file with task_01, task_02 or task_03, so it runs in the
first wave. The gate follows task_03 and task_04.
