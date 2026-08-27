---
schema: spec-tasks/v1
spec: 0118-a-task-proved-once-does-not-run-twice
qa: task_07
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
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_08
      file: task_08.md
      needs: [task_06]
    - id: task_09
      file: task_09.md
      needs: [task_08]
    - id: task_07
      file: task_07.md
      needs: [task_09]
---

# Tasks — A Task proved once does not run twice

| id      | title                                       | type    | complexity | needs   |
| ------- | ------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Carry-forward accepts an Unresolved Run     | backend | medium     | —       |
| task_02 | A Spec-scoped carry-forward query           | backend | medium     | task_01 |
| task_03 | Preflight refuses to re-execute proved work | backend | high       | task_02 |
| task_04 | The glossary names the accepted outcomes    | docs    | low        | task_03 |
| task_05 | Document both command contracts             | docs    | medium     | task_04 |
| task_06 | The skill ships with the CLI change         | docs    | medium     | task_05 |
| task_08 | A carry-forward that can carry a serial graph | backend | high     | task_06 |
| task_09 | The command help states the outcomes it accepts | backend | low    | task_08 |
| task_07 | QA gate                                     | qa      | medium     | task_09 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_06 · 7 → task_08 · 8 → task_09 · 9 → task_07

The chain is serial by edit locality, not by logic. task_01 and task_02 both
rewrite the reconcile command path, so they cannot be siblings. task_04,
task_05, and task_06 touch three provably disjoint files and could run as a
wave; they are serialized anyway because each reads the delivered code to
describe it, and a Task that documents behavior its sibling is still writing
documents a draft. Two Tasks that rewrite one region are not independently
implementable whatever their logic says: run as siblings on 2026-08-26, that
shape produced an integration conflict that discarded a Task which had already
passed its Verification.

task_06 is the only tooling Task. Its bounded path is
`.agents/skills/roundfix/SKILL.md`, authorized at
`docs/workflow/authorizations/2026-08-27-carry-forward-reaches-an-unresolved-run.md`;
the generated copy under `skills/` is sanctioned fallout of `make skills-sync`.

task_08 is corrective, added after this Spec's own first Run ended Unresolved:
carry-forward refused all six verified Tasks because their declared inputs were
compared against a checkout that had not received the earlier carries yet. It
is numbered after the gate but ordered before it; the graph, not the number,
carries the topology. task_09 is the second, from QA finding F-01: the reconcile help still called carry-forward a stopped-Run act. Both corrective slots are now used; a third finding requires a sanctioned exit, not a third patch.

Each Task's work, references, and Verification live in its own `task_NN.md`.
