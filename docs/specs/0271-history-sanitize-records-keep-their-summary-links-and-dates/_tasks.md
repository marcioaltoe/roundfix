---
schema: spec-tasks/v1
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
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
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — History sanitize records keep their summary, links and dates

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A legacy record's summary drops block quote markers and falls back to a section | backend | medium     | —       |
| task_02 | A legacy record states its override gaps, adopted sources and creation date   | backend | medium     | task_01 |
| task_03 | Legacy Link Rebase points the outcome and promoted copies at surviving files  | backend | high       | task_02 |
| task_04 | The history guide, the Roundfix skill and the glossary describe the readings  | docs    | low        | task_03 |
| task_05 | Run the final QA gate                                                         | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01, task_02 and task_03 all edit `BuildArchiveRecord` in
`internal/spec/archive_record.go`, and task_02 and task_03 both edit
`PlanLegacyConversion` and its request in `internal/spec/history_sanitize.go`,
so they run in series. task_04 describes all three and runs after them. The
gate follows the only leaf.
