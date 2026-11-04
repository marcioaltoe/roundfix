---
schema: spec-tasks/v1
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
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

# Tasks — The vendored catalog skills follow the catalog's main through a pinned bump

| id      | title                                                                                  | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Catalog Skill Bump pins the vendored catalog skills to one catalog commit          | infra   | medium     | —       |
| task_02 | The vendored catalog skills and every upstream setup follow catalog commit c8b2e85     | backend | high       | task_01 |
| task_03 | The glossary and the repository rules name the Catalog Skill Bump                      | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                  | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 needs task_01 because it runs the target task_01 adds and both
rewrite `docs/agents/setup-context.json`. task_03 needs task_02 because it
names the rule task_02's contract enforces; it shares no file with the others
and could run beside task_02, but the chain keeps the glossary after the
behavior it describes. Each implementation Task runs the suites that consume
its change, so no separate Integration Checkpoint is needed. The gate follows
task_03. No other Spec is required; the queue delivers this Spec after Spec
0349.
