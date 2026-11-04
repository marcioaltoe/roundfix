---
schema: spec-tasks/v1
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
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

# Tasks — The baseline update leaves owned skill lock entries unpinned

| id      | title                                                                                           | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The bundle install removes the ref of each owned lock entry it rewrites and reports it          | backend | medium     | —       |
| task_02 | Doctor warns with DR-SKILL-LOCK-STALE for an owned lock entry that pins another copy            | backend | medium     | task_01 |
| task_03 | The guides, the Roundfix skill and the glossary state the unpinned owned lock entry             | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                           | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 needs task_01 because both change
`internal/baseline/owned_skills_lock.go` and re-record
`docs/references/behavior-surfaces.json`, and its test calls the writer
task_01 changes. task_03 needs task_02 because it describes both behaviors
and re-records the same record. Each Task runs the suites that consume its
changed contract in its own Verification, so no separate Integration
Checkpoint is needed. No other Spec is required.
