---
schema: spec-tasks/v1
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
qa: task_05
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
---

# Tasks — The queue shows a skipped Linux Lane and waits for a Run slot

| id      | title                                                                         | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A Lane Environment Skip reaches the item and the Pull Request body             | backend | high       | —       |
| task_02 | Doctor checks the Linux Lane, and the references name the Lane Environment Skip | backend | medium     | task_01 |
| task_03 | The engine waits for a Run slot up to runs.slot_wait                           | backend | high       | task_02 |
| task_04 | The queue recognizes the ceiling refusal and shows the Run Slot Wait            | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                          | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series because they share files: task_01 and task_02 change
`internal/cli/deliver_linux_lane.go`; task_01 and task_03 change
`internal/delivery/engine.go`; task_01 and task_04 change
`internal/cli/deliver_workflow.go` and the `deliver` guide; task_02 and
task_04 change the Roundfix `deliver` reference, the skill version and
`CONTEXT.md`; every Task re-records the Behavior Surface Record. task_04
needs the engine types task_03 adds. The gate follows the last of them. No
Spec is listed in `requires`: the other cycle AM Specs change other functions
of the shared files (`_prd.md` → Prerequisites).
