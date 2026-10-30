---
schema: spec-tasks/v1
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
qa: task_03
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
---

# Tasks — The follow-up command takes every documented message type, an absent Spec Root counts as empty, and the follow-up guides name the JSON fields

| id      | title                                                                                            | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | An absent Spec Root inside the repository reads as an empty one                                  | backend | high       | —       |
| task_02 | The follow-up command takes every documented message type, and its guides name the JSON it emits | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                            | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 follows task_01 because both edit `internal/cli/follow_up_inventory.go`,
the Roundfix `SKILL.md` pair, `skills/testdata/owned-skill-versions.json`
and `docs/references/behavior-surfaces.json`, and its follow-up guides
describe the `specs` source task_01 changes. Each Task carries its own guide
and skill text and records the owned-skill version and the Behavior Surface
fingerprint last. Each code Task runs the suites that consume its contract,
so no separate Integration Checkpoint is needed. No other Spec is required.
