---
schema: spec-tasks/v1
spec: 0223-adjustments-the-adopters-asked-for
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The qa-gate skill, the Roundfix Skill and the guides describe the adjustments |
| task_02 | backend | The branch prefix becomes an optional Baseline decision |
| task_03 | backend | roundfix release plan reports the skills and baseline checks |
| task_04 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04

task_02 and task_03 change the behavior whose wording task_01 writes, so they
follow it. They share no declared file: task_02 changes the Baseline catalog,
its renderer, its derived artifacts, this repository's Setup Manifest and one
characterization case; task_03 changes the release plan command and the
Doctor's skills helper.
