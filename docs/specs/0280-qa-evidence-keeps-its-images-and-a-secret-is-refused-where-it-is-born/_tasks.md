---
schema: spec-tasks/v1
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
qa: task_04
requires:
  - 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
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

# Tasks — QA evidence keeps its images, and a secret is refused where it is born

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | QA evidence keeps its images, and the QA Report commit refuses any other binary               | backend | high       | —       |
| task_02 | The Secret Scan reads every Task commit, and test authors learn which fixture values pass     | backend | high       | task_01 |
| task_03 | A Delivery Retry accepts a Secret Rewrite, and `secret-found` names the rewrite that works    | backend | high       | task_02 |
| task_04 | Run the final QA gate                                                                         | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 and task_02 both change `internal/daemon/task_engine.go`, the
`implement` user guide and the Roundfix `implement` reference. All three
implementation Tasks raise the Roundfix skill version and re-record
`skills/testdata/owned-skill-versions.json`, the Behavior Surface Record and
`CONTEXT.md`, so they run in series. The gate follows the last of them. The
Spec requires Spec 0279, which changes the same `deliver` guide, Roundfix
references, version record and Behavior Surface Record and releases first.
