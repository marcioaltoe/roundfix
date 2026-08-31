---
schema: spec-tasks/v1
spec: 0105-the-gates-own-economics
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
    - id: task_07
      file: task_07.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_07]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — The gate's own economics

| id      | title                                          | type    | complexity | needs   |
| ------- | ---------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Roundfix owns the QA Task's Verification       | backend | high       | —       |
| task_02 | A finding blocks the rows it names             | backend | medium     | task_01 |
| task_03 | The citation parser reads the written forms    | backend | medium     | task_02 |
| task_07 | The derived Verification passes the checker     | backend | medium     | task_03 |
| task_05 | Characterization, and the Pull Request row     | docs    | high       | task_07 |
| task_06 | QA gate                                        | qa      | medium     | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_07 · 5 → task_05 · 6 → task_06

The chain is serial by two rules, not by logic.

Edit locality: task_02 and task_03 both rewrite `internal/speccheck`, and
task_01 reaches the checker through the finding it adds. Two Tasks that rewrite
one region are not independently implementable whatever their logic says — run
as siblings on 2026-08-26, that shape produced an integration conflict that
discarded a Task which had already passed its Verification.

Documentation after behavior: task_05 describes delivered behavior, so every
behavior Task precedes it. Spec 0118 ordered its documentation Tasks first and
shipped two surfaces contradicting the code until its gate caught it; Spec 0116
repeated the lesson from the other side, applying a rule uniformly to a skill
whose stage could not satisfy it.

The bounded authorization is not a Task. It landed as its own commit during
authoring, before any commit that edits a skill, which is the choreography the
tooling-authority rule requires. Authoring it as a Task would have made its
Verification vacuous the moment the record existed, and the record has to exist
before the graph can name it.

task_07 is corrective, from this Spec's own first Run: the command
`SC-QA-VERIFICATION-AUTHORED` demands is refused by `SC-VERIFY-NON-HERMETIC`,
which reads an awk regex as an external path, so no `qa` Task could be authored.
It is numbered after the gate but ordered before it; the graph, not the number,
carries the topology. One corrective slot of the two remains unused.

Each Task's work, references, and Verification live in its own `task_NN.md`.
