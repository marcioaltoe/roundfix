---
schema: spec-tasks/v1
spec: 0209-sources-that-share-a-context-share-a-spec
qa: task_05
requires: [0205-an-advisory-judge-for-spec-authoring]
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Sources that share a context share a Spec, as Baseline clauses with force |
| task_02 | backend | The judge asks whether an adopted and an open source share a Spec |
| task_03 | backend | `roundfix spec judge` prints and counts grouping suggestions |
| task_04 | docs | The authoring skills teach grouping and extending, and answer a suggestion |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 and task_02 share no file. task_03 and task_04 share
`skills/testdata/owned-skill-versions.json`, so they run in series.
