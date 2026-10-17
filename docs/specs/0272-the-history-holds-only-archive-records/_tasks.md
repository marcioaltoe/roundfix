---
schema: spec-tasks/v1
spec: 0272-the-history-holds-only-archive-records
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_05
      file: task_05.md
      needs: [task_06]
---

# Tasks — The history holds only Archive Records

| id      | title                                                                                  | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A citation of a removed folder file becomes a Git Locator in the record and its copies | backend | high       | —       |
| task_02 | The history commands remove every retired family and rewrite what cites it             | backend | high       | task_01 |
| task_03 | The check holds the History Root to Archive Records and links to files that exist      | backend | high       | task_02 |
| task_04 | The Baseline states retirement once and this repository's history leaves               | docs    | high       | task_03 |
| task_06 | The baseline update plan names every file it will remove                               | backend | low        | task_04 |
| task_05 | Run the final QA gate                                                                  | qa      | high       | task_06 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_05

task_02 removes files through the Git Locator and the Citation Rewrite that
task_01 adds. task_03 adds a unit to the same history commands and the same
user guide pages task_02 edits, and its check reads what task_02 removes.
task_04 refreshes this repository's guide, which turns task_03's check on, and
then removes the history that check would refuse, so both happen in one
commit. The gate follows the only leaf.

On 2026-10-17 task_04 also took the fix that lets the history citation scan
skip this repository's tracked `CLAUDE.md` link. It was first planned as a
corrective task_06, which could not settle alone: once the scan stops
refusing the link, every Managed Refresh of the tree applies the pending
history, so `TestRegenerationIsDeclared` reports the rewrites as undeclared
until the cleanup lands in the same commit.

task_06 is the corrective Task added on 2026-10-17 after the first QA gate:
the `roundfix baseline update` plan counted each unit's files but did not name
them, which User Story 3 asks for.
