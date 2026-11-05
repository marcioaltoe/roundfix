---
schema: spec-tasks/v1
spec: 0354-an-external-skill-is-checked-only-for-presence
qa: task_04
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
---

# Tasks — An external skill is checked only for presence

| id      | title                                                                                                     | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Doctor and the managed refresh stop comparing external skills with the Setup Snapshot                     | backend | medium     | —       |
| task_02 | A required external skill is checked only for presence                                                    | backend | medium     | task_01 |
| task_03 | The guides, the README, the Roundfix skill, the glossary and three ADRs state that an external skill is checked only for presence | docs    | medium     | task_02 |
| task_04 | Run the final QA gate                                                                                     | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 removes the snapshot comparison first, because it deletes
`snapshotComparableSkills`, the last reader of `OutdatedExternal` outside the
readiness renderers, and creates `internal/cli/external_skill_presence_test.go`.
task_02 removes `OutdatedExternal`, changes the same CLI files and adds the
lock rows to that test file. task_03 documents what both make true. Each
implementation Task runs the suites that consume its changed contract in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required.
