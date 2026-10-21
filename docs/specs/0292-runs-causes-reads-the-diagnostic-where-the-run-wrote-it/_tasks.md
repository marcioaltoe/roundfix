---
schema: spec-tasks/v1
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
qa: task_03
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
---

# Tasks — Runs causes reads the diagnostic where the Run wrote it

| id      | title                                                              | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | Runs causes reads the diagnostic under the resolved Artifact Root  | backend | medium     | —       |
| task_02 | The runs reference and the glossary say where diagnostics are read | docs    | low        | task_01 |
| task_03 | Run the final QA gate                                              | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The code slice carries the `runs` guide and the Behavior Surface Record,
because a Task that changes a command names its guide. The reference and the
glossary follow, because they describe that behavior; the gate follows them.
The two Tasks share no declared file, but task_02 states what task_01 built,
so it needs it. No
Spec is listed in `requires`: the resolver this plan reuses is already on the
default branch (`_prd.md` → Prerequisites).
