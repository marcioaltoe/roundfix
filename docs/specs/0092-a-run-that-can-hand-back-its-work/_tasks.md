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
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_02, task_04, task_06]
---

# Tasks — A Run that can hand back its work

| id      | title                                                    | type    | complexity | needs                     |
| ------- | -------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Record the four dispositions a Run gets wrong today      | test    | high       | —                         |
| task_02 | Move the work-started boundary to the first Agent output | backend | high       | task_01                   |
| task_03 | Let a failed Batch keep what its Agent achieved          | backend | medium     | task_01                   |
| task_04 | Derive the Run outcome from unresolved work              | backend | high       | task_03                   |
| task_05 | Give a superseded Run Branch a disposition               | backend | high       | task_01                   |
| task_06 | Hand a stopped Run's settled Tasks back                  | backend | high       | task_04, task_05          |
| task_07 | Run the final QA gate                                    | qa      | high       | task_02, task_04, task_06 |

Waves: 1 → task_01 · 2 → task_02, task_03, task_05 · 3 → task_04 · 4 → task_06 · 5 → task_07

Wave 2 is file-disjoint by construction: task_02 owns `acpx_runner.go` and
`agent_session_owner.go`, task_03 owns `agent.go`, task_05 owns the reconcile
surface. task_06 is serialised behind both task_04 and task_05 because all three
reach `internal/cli`, and a shared boundary between same-wave Tasks is what cost
Spec 0084 three Runs.

task_04 depends on task_03 and on nothing else for a reason worth stating: the
outcome derivation is only safe once settlement stops being overwritten. Applying
them in the other order produces a Run that reports `Clean` on a crashed Agent,
which is exactly the state a 2026-08-09 attempt reached before it was reverted.
