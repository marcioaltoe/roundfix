---
schema: spec-tasks/v1
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
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

# Tasks — A Spec commits the executables it names, settle says what it removes, and archive honours a declined QA gate

| id      | title                                                                                   | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A Task commit keeps the executables the Spec's authorization lists and judges each file | backend | high       | —       |
| task_02 | Archive accepts a Declared QA Decline without an override                               | backend | medium     | task_01 |
| task_03 | Settle names each worktree it removes                                                   | backend | medium     | task_02 |
| task_04 | Run the final QA gate                                                                   | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

Each Task is a vertical slice through one command's code and the documents
that describe it. All three write `CONTEXT.md`, the Roundfix `SKILL.md`
version pair, `skills/testdata/owned-skill-versions.json` and
`docs/references/behavior-surfaces.json`, and task_02 and task_03 both write
the Roundfix `spec-delivery` reference, so they run in series. No other Spec
is required.
