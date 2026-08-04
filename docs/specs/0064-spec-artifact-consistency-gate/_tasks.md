---
schema: spec-tasks/v1
spec: 0064-spec-artifact-consistency-gate
qa: task_09
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_04, task_05]
    - id: task_08
      file: task_08.md
      needs: [task_06, task_07]
    - id: task_10
      file: task_10.md
      needs: [task_08]
    - id: task_11
      file: task_11.md
      needs: [task_10]
    - id: task_09
      file: task_09.md
      needs: [task_11]
---

# Tasks — Spec artifact consistency gate

| id      | title                                              | type    | complexity | needs                     |
| ------- | -------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Report a Constraint contradiction with both sides   | backend | high       | —                         |
| task_02 | Report ADR and coverage gaps from citations         | backend | high       | task_01                   |
| task_03 | Report undocumented emitted vocabulary              | backend | medium     | task_01                   |
| task_04 | Expose the check as `roundfix spec check`           | backend | medium     | task_01, task_02, task_03 |
| task_05 | Replay the four findings that motivated this Spec   | test    | high       | task_02, task_03          |
| task_06 | Name the check in the glossary                      | docs    | low        | task_04                   |
| task_07 | Bring this repository's own Specs to a clean report | docs    | medium     | task_04, task_05          |
| task_08 | Wire the check into the local gate                  | infra   | low        | task_06, task_07          |
| task_10 | Prove the sweep budget under conditions that mean something | test | medium | task_08                |
| task_11 | Run the budget proof as its own gate step           | infra   | low        | task_10                   |
| task_09 | Run the final QA gate                               | qa      | high       | task_11                   |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04, task_05 ·
4 → task_06, task_07 · 5 → task_08 · 6 → task_10 · 7 → task_11 · 8 → task_09

task_01 is a tracer bullet, not a scaffold: it delivers the report model and
the first detector family together, so the very first settled Task reports a
real contradiction with a file and line on each side. It leads because the
constraint and tooling-authority family is the defect that cost four gate
executions closing Spec 0072, and it is the reason this Spec leads the queue.

Wave 2 fans out because the ADR/coverage detectors and the vocabulary contract
read disjoint artifacts and share only the report model.

task_08 is the one authorized tooling Task. Its bounded file is exactly
`Makefile`, per the Tooling authority row in `_prd.md` and `_techspec.md`. It
runs last among the implementation Tasks by design: wiring a fail-closed gate
before task_07 has brought the repository's own Specs to a clean report would
turn `make verify` red for every contributor.

The QA gate is authored as task_09 and depends on the graph's only non-QA
leaf, which transitively covers every other node.

## Corrective Tasks from the 2026-08-03 gate

The gate reported one finding, F-001: the corpus sweep's sub-second budget was
asserted as wall clock inside `go test -parallel 16 ./...`, where it measures
the runner's load rather than the check. The same sweep measured 0.64 s alone
and 3.0–4.0 s under the suite, so `make verify` exited 2 for a product that
meets its budget, and twelve of fifteen QA rows were blocked behind it.

The defect is in this Spec's authoring, not in the implementation: the
TechSpec's Testing Approach asked for a measurement it never gave execution
conditions. That section is amended, and the repair is split into two Tasks
because it crosses the tooling boundary — task_10 owns the test change and
must not touch `Makefile`, while task_11 is the authorized tooling Task bounded
to exactly `Makefile`. Folding them together would fail the tooling-authority
audit the gate performs against each Task commit.

Two corrective Tasks is the ceiling `docs/agents/autonomous-work.md` sets before
the decomposition itself must be re-examined. These two close one finding
together rather than turning discovery into a serial chain of gate cycles.
Appending them invalidates the 2026-08-03 gate result by construction, which is
why task_09 returns to `pending`.
