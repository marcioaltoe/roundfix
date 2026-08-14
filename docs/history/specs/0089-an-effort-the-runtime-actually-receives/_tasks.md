---
schema: spec-tasks/v1
spec: 0089-an-effort-the-runtime-actually-receives
qa: task_08
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
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_04]
    - id: task_08
      file: task_08.md
      needs: [task_05, task_06, task_07]
---

# Tasks — An effort the runtime actually receives

| id      | title                                                     | type    | complexity | needs                     |
| ------- | --------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Record the characterization corpus before anything moves  | test    | medium     | —                         |
| task_02 | Plan a deferred effort and fail closed on an unadvertised one | backend | high    | task_01                   |
| task_03 | Stop refusing an OpenCode reasoning effort                | backend | medium     | task_02                   |
| task_04 | Warm the session, apply the effort, and publish the receipt | backend | high     | task_03                   |
| task_05 | Route this repository to deepseek-v4-pro at xhigh         | infra   | low        | task_04                   |
| task_06 | Ship the roundfix skill with the removed refusal          | docs    | low        | task_04                   |
| task_07 | Record that the runtime hands back the floor              | docs    | low        | task_04                   |
| task_08 | Run the final QA gate                                     | qa      | high       | task_05, task_06, task_07 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05, task_06, task_07 · 6 → task_08

The chain from task_01 through task_04 is serial by file boundary: tasks 02, 03,
and 04 each edit `internal/agent/selection_assignment.go`,
`internal/agent/acpx_runner.go`, or the corpus task_01 writes, and a shared
boundary between parallel Tasks is how Spec 0084 lost three Runs. Tasks 05, 06,
and 07 have disjoint boundaries and are the graph's only parallel wave.
