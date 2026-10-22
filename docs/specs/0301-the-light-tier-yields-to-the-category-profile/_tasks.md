---
schema: spec-tasks/v1
spec: 0301-the-light-tier-yields-to-the-category-profile
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

# Tasks — The Light Tier yields to the category profile

| id      | title                                                                                                  | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The guides, the skill references and the glossary describe the skip, the yield and the Doctor line     | docs    | low        | —       |
| task_02 | A light Task skips or yields the Light Tier when the light session cannot serve it                      | backend | medium     | task_01 |
| task_03 | Doctor warns when the Light Tier cannot serve full access                                              | backend | medium     | task_02 |
| task_04 | Run the final QA gate                                                                                  | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 writes the configuration and Doctor guides, the Roundfix `runtime` and
`setup` references and the glossary for the surfaces task_02 and task_03
change, and re-records the Behavior Surface Record after its last guide edit,
so the code Tasks follow it in one dependency chain. task_02 adds the
`internal/lighttier` helper and changes the Daemon, and proves both in
`internal/lighttier`, `internal/daemon` and, through `roundfix implement`,
`internal/cli`, the packages that consume every contract it changes. task_03
calls that helper from Doctor, so it follows task_02. The gate follows all
three. No other Spec is required.
