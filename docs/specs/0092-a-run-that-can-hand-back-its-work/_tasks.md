---
schema: spec-tasks/v1
spec: 0092-a-run-that-can-hand-back-its-work
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
      needs: [task_01]
    - id: task_08
      file: task_08.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03, task_08]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_04, task_05]
    - id: task_09
      file: task_09.md
      needs: [task_06]
    - id: task_10
      file: task_10.md
      needs: [task_05, task_06]
    - id: task_11
      file: task_11.md
      needs: [task_10]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_04, task_06, task_09, task_10, task_11]
---

# Tasks — A Run that can hand back its work

| id      | title                                                    | type    | complexity | needs                     |
| ------- | -------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Record the four dispositions a Run gets wrong today      | test    | high       | —                         |
| task_02 | Move the work-started boundary to the first Agent output | backend | high       | task_01                   |
| task_03 | Let a failed Batch keep what its Agent achieved          | backend | medium     | task_01                   |
| task_08 | Name the seventh break the enumeration missed            | test    | medium     | task_02                   |
| task_04 | Derive the Run outcome from unresolved work              | backend | high       | task_03, task_08          |
| task_05 | Give a superseded Run Branch a disposition               | backend | high       | task_01                   |
| task_06 | Hand a stopped Run's settled Tasks back                  | backend | high       | task_04, task_05          |
| task_09 | Let the reconcile JSON contract know about carry-forward | test    | low        | task_06                   |
| task_10 | Let the two new acts be discoverable and named           | docs    | low        | task_05, task_06          |
| task_11 | Make the assembled tree pass its own gate                | test    | low        | task_10                   |
| task_07 | Run the final QA gate                                    | qa      | high       | task_02, task_04, task_06, task_09, task_10, task_11 |

Waves: 1 → task_01 · 2 → task_02, task_03, task_05 · 3 → task_08 · 4 → task_04 · 5 → task_06 · 6 → task_09, task_10 · 7 → task_11 · 8 → task_07

Four of these eleven Tasks were minted after implementation began, and all four
have the same cause: a Task delivered a public change whose surrounding
contract — the test pinning a payload, the help copy naming a flag, the
glossary entry naming a term — sat outside every bounded scope, so no Task
could complete it. task_05's own Result recorded its instance while it happened.
Enumerating a boundary is not the same as enumerating what the change obliges.

task_08 was minted on 2026-08-11, after task_04's whole-package gate found a
seventh test broken by task_02's boundary where task_01 had enumerated six. It
sits between them because task_04 proves both packages green, and it cannot do
that while a break nobody declared is still failing.

Wave 2 is file-disjoint by construction: task_02 owns `acpx_runner.go` and
`agent_session_owner.go`, task_03 owns `agent.go`, task_05 owns the reconcile
surface. task_06 is serialised behind both task_04 and task_05 because all three
reach `internal/cli`, and a shared boundary between same-wave Tasks is what cost
Spec 0084 three Runs.

task_04 depends on task_03 and on nothing else for a reason worth stating: the
outcome derivation is only safe once settlement stops being overwritten. Applying
them in the other order produces a Run that reports `Clean` on a crashed Agent,
which is exactly the state a 2026-08-09 attempt reached before it was reverted.
