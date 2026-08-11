---
schema: spec-tasks/v1
spec: 0080-cheap-detectors-run-before-the-gate
qa: task_08
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_09
      file: task_09.md
      needs: [task_03]
    - id: task_10
      file: task_10.md
      needs: [task_09]
    - id: task_11
      file: task_11.md
      needs: [task_10]
    - id: task_08
      file: task_08.md
      needs: [task_05, task_07, task_09, task_10, task_11]
---

# Tasks — Cheap detectors run before the gate

| id      | title                                                | type    | complexity | needs            |
| ------- | ---------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Detect the facts the gate proves by hand              | backend | high       | —                |
| task_02 | Tell the gate prompt what changed                     | backend | medium     | —                |
| task_03 | Run the mechanical stage before the Agent Session     | backend | high       | task_01, task_02 |
| task_04 | Let a report row declare its evidence inputs          | chore   | medium     | task_02          |
| task_05 | Carry a row forward only on unmoved evidence          | backend | high       | task_03, task_04 |
| task_06 | Declare the two-tier verification contract            | chore   | high       | task_04          |
| task_07 | Re-record the expectation task_06 invalidates         | test    | low        | task_06          |
| task_09 | Let the QA harness complete the report the Daemon seeded | test    | medium     | task_03          |
| task_10 | Give the prompt fixture the Git precondition the stage now needs | test | low | task_09 |
| task_11 | Let the authorization detector read the fallout it already sanctions | backend | medium | task_10 |
| task_08 | Run the final QA gate                                 | qa      | medium     | task_05, task_07, task_09, task_10, task_11 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05, task_06 ·
4 → task_07, task_09 · 5 → task_10 · 6 → task_11 · 7 → task_08

Three edges carry reasoning worth stating.

**task_02 has no dependency**, though an earlier draft of the TechSpec gave it
one. It edits the agent prompt surface while task_01 edits the consistency
checker, so the two are file-disjoint and start together. The TechSpec was
amended rather than silently contradicted.

**task_06 follows task_04 rather than task_01.** Both are tooling Tasks under
the same authorization and both regenerate derived artifacts under
`internal/baseline/assets/**`; running them in one wave invites an integration
conflict for no gain.

**task_07 exists because Spec 0079 paid for it.** task_06's clauses grow the
Source Baseline, staleing the maintained compatibility expectation. A
consequent fix folded into an authorized tooling commit fails the
tooling-authority gate — `docs/agents/specific-repository.md` says so, and the
0079 gate refused exactly that shape. Authoring it as its own node applies the
lesson at authoring time instead of paying a gate round to relearn it.

The parallel waves are file-disjoint by declared Context: wave 1 pairs
`internal/speccheck` with `internal/agent`; wave 2 pairs the Daemon's QA step
with the qa-gate skill; wave 3 pairs the carry-forward resolver with the
Baseline modules and Makefile.

task_04 and task_06 are the authorized tooling Tasks, bounded by
`docs/workflow/authorizations/2026-08-06-proof-cost.md`.

The QA gate is authored as task_08 and depends on task_05 and task_07, the
graph's only non-QA leaves.
