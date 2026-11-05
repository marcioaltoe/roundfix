---
schema: spec-tasks/v1
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
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

# Tasks — Only a .git entry that Git accepts marks a Git root

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Git root walk accepts only a .git entry that Git accepts                                   | backend | medium     | —       |
| task_02 | The skills project root uses the same walk, and the commands under a rejected .git say so      | backend | medium     | task_01 |
| task_03 | The Roundfix skill, the runs guide and ADR-0357 state that only a .git Git accepts marks a Git root | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                          | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 changes the walk and every fixture together, because the repository
Verification at its settlement runs the old fixtures under the new walk.
task_02 needs `config.FindGitRoot` and `gittest`'s fixtures from task_01.
task_03 documents what tasks 01 and 02 make true. Each implementation Task
runs the suites that consume its changed contract in its own Verification, so
no separate Integration Checkpoint is needed. No other Spec is required.
