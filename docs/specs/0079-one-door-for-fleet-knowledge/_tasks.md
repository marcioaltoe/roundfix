---
schema: spec-tasks/v1
spec: 0079-one-door-for-fleet-knowledge
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
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_08
      file: task_08.md
      needs: [task_02, task_06]
    - id: task_09
      file: task_09.md
      needs: [task_06, task_08]
    - id: task_07
      file: task_07.md
      needs: [task_04, task_06, task_08, task_09]
---

# Tasks — One door for fleet knowledge

| id      | title                                            | type    | complexity | needs            |
| ------- | ------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Name the door's vocabulary in the glossary        | docs    | low        | —                |
| task_02 | Carry the contracts and the permission to the guides | chore | high       | task_01          |
| task_03 | Sweep the legacy findings behind their first rollups | docs  | medium     | task_02          |
| task_04 | Refuse what the findings contract forbids         | backend | high       | task_02, task_03 |
| task_05 | Pilot the door with this Spec's own material      | docs    | medium     | task_02          |
| task_06 | Teach the authorial skills to start at the door   | chore   | medium     | task_05          |
| task_08 | Re-record the maintained Source Baseline expectation | test  | low        | task_02, task_06 |
| task_09 | Teach the guides what the operator was promised   | chore   | medium     | task_06, task_08 |
| task_07 | Run the final QA gate                             | qa      | medium     | task_04, task_06, task_08, task_09 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_05 · 4 → task_04,
task_06 · 5 → task_08 · 6 → task_09 · 7 → task_07

task_08 was appended after the first QA gate run settled `fail`: the clauses
task_02 and task_06 authored grew the Source Baseline, and the maintained
compatibility fixture's declared expectation still named the pre-Spec count,
so the repository gate was red and the gate blocked seven rows behind that
one finding (F-001). Appending it returns task_07 to `pending` — a gate
result cannot certify a graph that grew after it ran.

task_09 followed from the second gate run, whose only red row was F-002: the
PRD promised four operator behaviors — oldest-first Triage, an empty inbox
read as rest, a live-work-only findings directory read as health, and a
Rollup with no open members as a closure candidate — and no task owned
carrying them into a durable instruction. It also owns the fallout its own
clauses cause in the maintained Source Baseline expectation, so the gate that
follows it never inherits a red repository.

Two orderings in this graph are load-bearing, both from traps this
repository paid for on 2026-08-05/06:

- The sweep (task_03) precedes the checks (task_04). Landing
  `SC-FINDING-LIFECYCLE` first would turn the repository gate red on 63
  unstamped findings inside the check task's own Verification — the task
  that repairs the state arriving after the gate that demands it, the Spec
  0075 precondition trap.
- Permission precedes obligation. task_02 carries the *permissive*
  secondbrain carve-out (sessions may create under the brain's `inbox/**`)
  so the pilot (task_05) can exercise a door the local guide no longer
  prohibits; the *obligating* clauses bind only in task_06, after the pilot
  proved the cycle — the inert-first ordering Spec 0073 proved.

The parallel waves are file-disjoint by declared Context: wave 3 pairs the
sweep (existing files under the findings directory) with the pilot (new
dated files and the spec-local report); wave 4 pairs Go work in the
consistency checker with the skills-and-guides tooling surface.

task_02 and task_06 are the two authorized tooling Tasks, bounded exactly by
`docs/workflow/authorizations/2026-08-06-fleet-knowledge-door.md`. task_05
carries the one external gate: the brain-side inbox existing under that
repository's own contract — when absent, the task settles failed naming the
dependency rather than faking the proof.

The QA gate is authored as task_07 and depends on task_04 and task_06, the
graph's only non-QA leaves.
