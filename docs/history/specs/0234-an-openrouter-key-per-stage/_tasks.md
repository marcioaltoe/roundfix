---
schema: spec-tasks/v1
spec: 0234-an-openrouter-key-per-stage
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

# Tasks — An OpenRouter key per stage

| id      | title                                                                                          | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The judge reads its own OpenRouter key first and names the variable it used                    | backend | medium     | —       |
| task_02 | The Doctor Command lists each stage's keys, and the guides and skills name the stage keys      | backend | low        | task_01 |
| task_03 | Implementation on an open model reads its own OpenRouter key first and records the variable    | backend | low        | task_02 |
| task_04 | Run the final QA gate                                                                          | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

Each implementation Task raises an owned skill's version or follows one that
does, so they share `skills/testdata/owned-skill-versions.json` and run in
series. task_03 changes the helper and spend record that Spec 0233 adds, so
the operator queues this Spec after 0233; the graph cannot express a
dependency on another Spec. The gate follows the only leaf.
