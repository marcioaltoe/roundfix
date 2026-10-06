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
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — Retire the-fool, autoresearch and council

| id      | title                                                                                   | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Baseline retires council and the-fool and the asset sync keeps them out             | backend | high       | —                         |
| task_02 | Roundfix stops shipping council and this repository drops the three skills              | chore   | medium     | task_01                   |
| task_03 | Baseline update lists the retired copies an adopter still holds                         | backend | medium     | task_02                   |
| task_04 | Run the final QA gate                                                                   | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 changes the doctor counts only after task_01 has removed both skills
from the module, and both refresh `docs/agents/setup-context.json`. task_03
reads task_01's Retired Skill set and shares
`skills/testdata/owned-skill-versions.json` with task_02, so the three run in
series. The gate follows every Task.
