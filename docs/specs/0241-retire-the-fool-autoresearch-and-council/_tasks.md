---
schema: spec-tasks/v1
spec: 0241-retire-the-fool-autoresearch-and-council
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01, task_05]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03, task_05]
    - id: task_05
      file: task_05.md
      needs: [task_01]
---

# Tasks — Retire the-fool, autoresearch and council

| id      | title                                                                                   | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Baseline retires council and the-fool and the asset sync keeps them out             | backend | high       | —                         |
| task_02 | Roundfix stops shipping council and this repository drops the three skills              | chore   | medium     | task_01, task_05          |
| task_03 | Baseline update lists the retired copies an adopter still holds                         | backend | medium     | task_02                   |
| task_04 | Run the final QA gate                                                                   | qa      | high       | task_01, task_02, task_03, task_05 |
| task_05 | The special-entry fixture binds its socket under a short root                           | test    | low        | task_01                   |

Waves: 1 → task_01 · 2 → task_05 · 3 → task_02 · 4 → task_03 · 5 → task_04

task_02 changes the doctor counts only after task_01 has removed both skills
from the module, and both refresh `docs/agents/setup-context.json`. task_03
reads task_01's Retired Skill set and shares
`skills/testdata/owned-skill-versions.json` with task_02, so the three run in
series. task_05 is a corrective Task: it moves the special-entry fixture of
`skills/repository_test.go` under a short root so its socket binds on macOS,
and it runs before task_02, which edits the same file. The gate follows every
Task.
