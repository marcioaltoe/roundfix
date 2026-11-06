---
schema: spec-tasks/v1
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
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

# Tasks — The build cache stays within a bound, and a full disk is named

| id      | title                                                          | type    | complexity | needs   |
| ------- | -------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Makefile holds its Go build cache to the Build Cache Bound | infra   | medium     | —       |
| task_02 | Doctor and the Delivery Queue warn below the Free-Space Floor  | backend | medium     | task_01 |
| task_03 | A full filesystem parks the item as a Disk-Full Park           | backend | medium     | task_02 |
| task_04 | Run the final QA gate                                          | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

Each Task documents its own behavior and writes its own `CONTEXT.md` entry,
so the three share the glossary and run in series. task_03 also uses
task_02's `inspectFreeSpace` and changes the engine, the workflow, the
`deliver` documents and the Roundfix skill version that task_02 changes.
Each implementation Task runs the suites that consume its changed contract in
its own Verification, so no separate Integration Checkpoint is needed. No
other Spec is required; Spec 0359 delivers first and edits other functions of
the shared files.
