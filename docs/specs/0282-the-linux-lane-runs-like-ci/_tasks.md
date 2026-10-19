---
schema: spec-tasks/v1
spec: 0282-the-linux-lane-runs-like-ci
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

# Tasks — The Linux Lane runs like CI

| id      | title                                                                                   | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A test that needs an External Test Tool skips naming it                                  | test    | medium     | —       |
| task_02 | The Linux Lane runs under an init on the configured image                                | backend | medium     | task_01 |
| task_03 | A Lane Container Failure is told apart from a Lane Change Failure                        | backend | high       | task_02 |
| task_04 | `roundfix deliver lane` runs the Linux Lane on demand and exits with its outcome          | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                                    | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

Every implementation Task writes a term into `CONTEXT.md`. task_02, task_03
and task_04 also change `internal/cli/deliver_linux_lane.go`, the `deliver`
user guide and the Roundfix `deliver` reference, raise the Roundfix skill
version and re-record `skills/testdata/owned-skill-versions.json` and the
Behavior Surface Record, so all four run in series. The gate follows the last
of them. No other Spec is required.
