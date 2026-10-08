---
schema: spec-tasks/v1
spec: 0249-a-baseline-update-that-sanitizes-pending-history
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

# Tasks — A Baseline update that sanitizes pending history

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The glossary, the command references and the Roundfix Skill describe the Pending History      | docs    | medium     | —       |
| task_02 | A legacy unproven list of maps converts and every refusal prints on one line                  | backend | medium     | task_01 |
| task_03 | Baseline update plans and applies the Pending History in the change it plans                  | backend | high       | task_02 |
| task_04 | The upgrade notice names the Pending History                                                  | backend | low        | task_03 |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 describes ADR-0254 and the TechSpec in the glossary, the command
references and the Roundfix Skill first, so every code Task names its guide
through its dependencies. task_02 and task_03 both change
`internal/cli/history.go`, so they run in series, and task_03 prints the
one-line reasons task_02 introduces. task_04 reads the inventory and the
history helpers that task_03 leaves in the same package, so it follows
task_03. The gate follows the only leaf.
