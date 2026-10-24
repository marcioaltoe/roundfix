---
schema: spec-tasks/v1
spec: 0314-gc-trims-the-shared-build-cache
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

# Tasks — GC trims the shared build cache

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A trim engine removes only idle Go build cache entries                        | backend | medium     | —       |
| task_02 | GC trims each declared Shared Build Cache                                     | backend | high       | task_01 |
| task_03 | This repository declares its cache, and the documents describe the trim       | docs    | medium     | task_02 |
| task_04 | Run the final QA gate                                                         | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 calls task_01's engine; task_03 documents task_02's key and section,
and both re-record `docs/references/behavior-surfaces.json`, so the three
run in series. task_02's GC tests drive the engine through the command, so no
separate Integration Checkpoint is needed. No other Spec is required.
