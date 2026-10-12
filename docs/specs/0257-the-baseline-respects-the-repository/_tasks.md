---
schema: spec-tasks/v1
spec: 0257-the-baseline-respects-the-repository
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

# Tasks — The Baseline respects the repository's own choices

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A repository records its Better Auth base path                                | backend | high       | —       |
| task_02 | A Baseline update records a declared decision                                 | backend | high       | task_01 |
| task_03 | A skill that differs from its Setup Snapshot is kept, not called trailing     | backend | high       | task_02 |
| task_04 | The Baseline states how to handle a diverging vendored skill                  | backend | medium     | task_03 |
| task_05 | Run the final QA gate                                                         | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The four implementation Tasks run in series (Wave collision). task_02's
end-to-end test declares the base path task_01 adds, and both re-record the
Behavior Surface Record. task_02 and task_03 both change
`internal/cli/baseline_update.go`, the update guide, the `roundfix` skill's
Baseline reference and the owned-skill version record. task_03 and task_04
both change that reference and record, and task_04 shares the catalog's
derived files and `docs/agents/setup-context.json` with task_01. The gate
follows the only leaf.
