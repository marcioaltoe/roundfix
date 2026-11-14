---
schema: spec-tasks/v1
spec: 0392-each-authoring-phase-runs-on-its-declared-model
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Each authoring phase runs on its declared model

| id      | title                                                                                             | type    | complexity | needs                     |
| ------- | ------------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | A repository declares its authoring models and installation applies the declaration                | backend | high       | —                         |
| task_02 | The checks report an unapplied Authoring Model Declaration                                         | backend | medium     | task_01                   |
| task_03 | A reported user-level skill copy carries the command that removes it                               | backend | high       | task_01, task_02          |
| task_04 | The authoring skills delegate a declared Authoring Phase, and the glossary and ADR-0386 name it    | docs    | medium     | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                                             | qa      | high       | task_04                   |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series. Tasks 01 to 03 share the `skills` command of the
CLI, its guide and the Behavior Surface Record, tasks 02 and 03 change the
same three reports, and task 04 describes what the first three built. Each
implementation Task runs the suites that consume its changed contracts in its
own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, how a repository's
Roundfix skills are installed and reported, in four implementation Tasks and
the QA gate, and each Task earns its verification cycle: the declaration and
its installation, the inspector and the three reports, the user-level notice
with the removal command, and the delegation text with the glossary. The
inspector stays its own Task instead of joining task 03 because it carries
three reports, four tests and three guides, and task 03 is already the
largest slice after task 01. One part of `_techspec.md` → API Contract 12
moved from task 04 to task 03 when the Tasks were prototyped: the `setup`
paragraph that names `roundfix skills remove`, because the documentation
contract refuses a command the Roundfix skill does not name, so task 03 could
not leave the repository Verification green without it.

No other Spec is required. Spec 0391 delivers first and Spec 0393 after this
one; `_prd.md` → Prerequisites names the shared files, and each Task anchors
its skill edits on headings and sentences and raises each skill version with
the record command on the tree it runs on.
