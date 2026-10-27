---
schema: spec-tasks/v1
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A profile change keeps the recorded decisions and names what it wrote

| id      | title                                                                                                     | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The guides, skills and glossary describe the profile change, the Root Backup and the Claude Skill Links | docs    | medium     | —                |
| task_02 | A profile change carries the recorded decisions and refuses once                                          | backend | high       | task_01          |
| task_03 | The apply names each Root Backup and reconcile says what it compared                                      | backend | medium     | task_02          |
| task_04 | A restore keeps the repository's Claude Skill Links                                                       | backend | high       | task_01          |
| task_05 | Run the final QA gate                                                                                     | qa      | high       | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05

task_01 writes the guides, skills and glossary first, so each code Task names
its command's guide through that dependency; it records the owned-skill
versions and the Behavior Surface fingerprint last. task_03 follows task_02
because both edit `internal/cli/baseline_profile.go`. task_04 touches only
the restore and the transaction and runs beside task_02. Each code Task runs
the plan, apply, reconcile, restore and transaction suites that consume its
contract, so no separate Integration Checkpoint is needed. No other Spec is
required.
