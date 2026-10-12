---
schema: spec-tasks/v1
spec: 0261-a-cancelled-spec-leaves-a-record
requires: [0256-an-archive-record-that-survives-a-squash-merge]
qa: task_05
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A cancelled Spec leaves a record

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A cancelled Archive Record keeps the reason and the Tasks, and Source Return hands back sources | backend | high       | —                |
| task_02 | The docs layout guide, the skills, the user guides and the glossary describe cancellation      | docs    | high       | —                |
| task_03 | The Cancel Command cancels an active Spec and refuses while work holds it                      | backend | high       | task_01, task_02 |
| task_04 | A cancelled Spec meets no prerequisite, supersedes nothing and is never retried                | backend | medium     | task_01, task_02 |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05

task_01 owns the record and the cut in `internal/spec`, and task_02 owns every
text, so they share no declared file and run together. task_02 writes the
guides first, so the two code Tasks depend on the Task that documents their
surfaces, and the owned skills' versions rise once. task_03 needs task_01's
`Cancel` and task_02's routing row, user guide and coverage entries; both
task_02 and task_03 declare `docs/references/skill-coverage.json` and
`docs/references/behavior-surfaces.json`, and task_03 follows task_02. task_04
needs task_01's disposition to build its fixtures and task_02's supersede and
deliver guide texts. task_03 and task_04 share no declared file. The gate
follows both leaves. Spec 0256 is a declared prerequisite in `requires`; no
Verification here runs a test it creates.
