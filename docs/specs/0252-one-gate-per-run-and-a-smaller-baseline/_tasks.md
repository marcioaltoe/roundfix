---
schema: spec-tasks/v1
spec: 0252-one-gate-per-run-and-a-smaller-baseline
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

# Tasks — One gate per Run and a smaller Baseline

| id      | title                                                                                        | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A Daemon-assigned Task runs focused tests and leaves the full gate to the Daemon            | backend | medium     | —       |
| task_02 | The recorded Verification and runtime values are the ones CI and the Project Config use     | backend | medium     | task_01 |
| task_03 | Each review, tracker and local-research rule is stated once                                  | backend | high       | task_02 |
| task_04 | The root block says when a session needs the docs-layout and Secondbrain guides              | backend | medium     | task_03 |
| task_05 | Run the final QA gate                                                                        | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

Every implementation Task changes Baseline source and shares
`internal/baseline/module-versions.json`, the derived catalog files, the
profile digest and `docs/agents/setup-context.json`, so the four run in series
(Wave collision). task_01 goes first, so that task_02 to task_04 start from
guides that already send a Daemon-assigned Task to focused tests. task_01 and
task_02 both change `docs/agents/agent-instructions.md`; task_02 and task_03
both change the `autonomous-work` module; task_03 and task_04 both change the
`secondbrain` module. The gate follows the only leaf.
