---
schema: spec-tasks/v1
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
qa: task_05
requires: [0211-a-delivery-queue-that-finishes-without-intervention]
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
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The pre-PR review omits QA evidence and upstream-managed skills and blocks above a measured bound |
| task_02 | backend | The failed-pass QA import leaves compiled source behind |
| task_03 | backend | A Task declares a path it deletes |
| task_04 | backend | Reconcile releases the Runs a merged Spec superseded |
| task_06 | backend | Apply releases a merged Run whose target branch is gone |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_05

Every implementation Task raises an owned skill's version and re-records
`skills/testdata/owned-skill-versions.json`, so the four run in series;
task_01 and task_04 also share the Roundfix Skill's `SKILL.md`, and task_04
reads the Context kinds task_03 adds.
