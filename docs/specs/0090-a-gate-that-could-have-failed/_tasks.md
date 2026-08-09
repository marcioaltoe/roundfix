---
schema: spec-tasks/v1
spec: 0090-a-gate-that-could-have-failed
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
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_07
      file: task_07.md
      needs: [task_04, task_05, task_06]
---

# Tasks — A gate that could have failed

| id      | title                                                    | type    | complexity | needs                     |
| ------- | -------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Record what a vacuous and an unobserved gate do today    | test    | medium     | —                         |
| task_02 | Separate a command's verdict from the runner's sight     | backend | medium     | task_01                   |
| task_03 | Refuse a Task whose gate already passes                  | backend | high       | task_02                   |
| task_04 | Publish the probe's finding where a reader can see it    | backend | medium     | task_03                   |
| task_05 | Carry the negative control a Task declares               | backend | medium     | task_01                   |
| task_06 | Source every wait budget from its shared constant        | test    | low        | —                         |
| task_07 | Run the final QA gate                                    | qa      | high       | task_04, task_05, task_06 |

Waves: 1 → task_01, task_06 · 2 → task_02, task_05 · 3 → task_03 · 4 → task_04 · 5 → task_07

The two parallel waves are file-disjoint by construction. Wave 1 splits
`internal/daemon` (task_01) from `internal/agent` and `internal/store`
(task_06); wave 2 splits `internal/daemon` (task_02) from `internal/spec`
(task_05). Spec 0084 lost three Runs to a shared boundary between same-wave
Tasks, so each Task below enumerates its files rather than describing them.

The serial chain task_02 → task_03 → task_04 exists because all three edit
`internal/daemon/engine.go` or `internal/daemon/task_engine.go`, and because
each needs the one before it to be meaningful: the probe cannot classify a
command it cannot observe, and the event cannot report a finding the probe does
not yet produce.

From task_03 onward this Spec is subject to its own mechanism — task_04's
Verification is probed by the code task_03 introduces. That is intended, and it
is the cheapest possible acceptance evidence for the feature.
