---
schema: spec-tasks/v1
spec: 0116-a-verdict-that-states-its-own-scope
qa: task_06
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
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_05]
    - id: task_06
      file: task_06.md
      needs: [task_07]
---

# Tasks — A verdict that states its own scope

| id      | title                                        | type    | complexity | needs   |
| ------- | -------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The auditing binary and its staleness         | backend | medium     | —       |
| task_02 | A QA Report names its auditor                 | backend | medium     | task_01 |
| task_03 | The verdict line states the probe's coverage  | backend | medium     | task_02 |
| task_04 | The glossary names the Auditing Binary        | docs    | low        | task_03 |
| task_05 | The authoring skills name the probing check   | docs    | medium     | task_04 |
| task_07 | The QA gate keeps the non-probing check       | docs    | low        | task_05 |
| task_06 | QA gate                                       | qa      | medium     | task_07 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_07 · 7 → task_06

The chain is serial by two rules, not by logic.

Edit locality: task_02 and task_03 both reach `internal/cli`, and two Tasks that
rewrite one region are not independently implementable whatever their logic
says — run as siblings on 2026-08-26, that shape produced an integration
conflict that discarded a Task which had already passed its Verification.

Documentation after behavior: task_04 and task_05 describe delivered behavior,
so every behavior Task precedes them. Spec 0118 ordered its documentation
Tasks before a corrective Task changed the rule they described, and shipped two
surfaces that contradicted the code until its QA gate caught it as F-02. This
graph applies that lesson before the fact.

task_05 is the only tooling Task. Its bounded paths are the four canonical
authoring skills and their four generated copies, authorized at
`docs/workflow/authorizations/2026-08-30-the-authoring-skills-name-the-probing-check.md`.

task_07 is corrective, from this Spec's own first Run: task_05 put the probing
form into the QA gate's precondition, where the probe's question has no true
answer and every completed Task reports vacuous. It is numbered after the gate
but ordered before it; the graph, not the number, carries the topology. One
corrective slot of the two remains unused.

Each Task's work, references, and Verification live in its own `task_NN.md`.
