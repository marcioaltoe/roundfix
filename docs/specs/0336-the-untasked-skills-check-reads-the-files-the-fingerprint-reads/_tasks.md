---
schema: spec-tasks/v1
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — The untasked-skills check reads the files the fingerprint reads

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The untasked-skills check and the Behavior Surface fingerprint share one file predicate         | backend | medium     | —       |
| task_02 | Run the final QA gate                                                                         | qa      | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 runs the suites that consume the changed lookup and fingerprint in
its own Verification, so no separate Integration Checkpoint is needed. No
other Spec is required; Specs 0332 to 0335 deliver first (see
`_techspec.md` → Build Order).
