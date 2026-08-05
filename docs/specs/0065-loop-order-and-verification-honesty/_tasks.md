---
schema: spec-tasks/v1
spec: 0065-loop-order-and-verification-honesty
qa: task_06
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Loop order and verification honesty

| id      | title                                          | type    | complexity | needs            |
| ------- | ---------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | State the loop order once                       | chore   | medium     | —                |
| task_02 | Refuse a Verification that cannot fail          | backend | high       | —                |
| task_03 | Refuse a Task that contradicts itself           | backend | medium     | task_02          |
| task_04 | Check that the order restatements agree         | backend | low        | task_01          |
| task_05 | Synchronise the authoring and CLI Skills        | chore   | low        | task_03, task_04 |
| task_06 | Run the final QA gate                           | qa      | medium     | task_05          |

Waves: 1 → task_01 · task_02 · 2 → task_03 · task_04 · 3 → task_05 ·
4 → task_06

task_04 depends on task_01 deliberately, and the dependency is not cosmetic.
`SC-LOOP-ORDER-DIVERGENT` fails while the three statements disagree, and
`spec check` runs inside `make verify`. Authoring the check before the sources
agree would leave the repository red between two Tasks — and the Daemon runs
the configured Verification command as a precondition, so the Task meant to
repair that state would be settled without ever starting. Spec 0075 lost a Run
to exactly that shape on 2026-08-05.

task_01 and task_02 are independent and form the first wave. task_03 shares the
requirement parsing task_02 introduces.

Two Tasks carry authorized tooling scope. task_01 edits
`internal/baseline/assets/modules/autonomous-work.json` under the 2026-08-04
confirmation of the boundary the 2026-08-02 grant deferred. task_05 edits the
`write-tasks` and `roundfix` Skills with their mirrors under the same records.

The QA gate is authored as task_06 and depends on task_05, the graph's only
non-QA leaf.
