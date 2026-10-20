---
schema: spec-tasks/v1
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
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

# Tasks — GC sanitize reaches the default Artifact Roots

| id      | title                                                                 | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | One Artifact Root resolver serves the storage report and Run Retention | backend | medium     | —       |
| task_02 | GC sanitation classifies the Artifact Root each Run wrote to           | backend | high       | task_01 |
| task_03 | The storage reference and the glossary name the resolved Artifact Root | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                  | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The resolver comes first because sanitation groups Runs by its result. Both
code Tasks change `internal/store/journal.go`, `docs/user-guide/commands/gc.md`
and the Behavior Surface Record, so they run in series. The reference and the
glossary follow, because they describe both slices; the gate follows the last
of them. No Spec is listed in `requires`: the storage report code this plan
builds on merges before this plan (`_prd.md` → Prerequisites).
