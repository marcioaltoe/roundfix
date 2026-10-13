---
schema: spec-tasks/v1
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
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
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A killed child names its signal, and a cancelled Run starts nothing

| id      | title                                                                 | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The deliver reference, the user guide and the glossary name a Signal Exit | docs    | low        | —                |
| task_02 | The Delivery Queue names the signal that ended an item binary         | backend | medium     | task_01          |
| task_03 | Test helpers that run a child Roundfix binary name its Signal Exit    | test    | medium     | task_02          |
| task_04 | A cancelled Run starts no Verification command after an immediate grant | backend | medium     | —                |
| task_05 | Run the final QA gate                                                 | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_04 · 2 → task_02 · 3 → task_03 · 4 → task_05

No two Tasks share a declared file. task_02 depends on task_01 because the
Task that changes the `deliver` surface needs the Task that documents it in
its ancestry; task_03 depends on task_02 because its helper imports the
`internal/childexit` package task_02 creates. task_04 touches only
`internal/daemon` and runs beside task_01.
