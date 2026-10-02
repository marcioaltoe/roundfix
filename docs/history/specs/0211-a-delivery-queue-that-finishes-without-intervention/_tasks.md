---
schema: spec-tasks/v1
spec: 0211-a-delivery-queue-that-finishes-without-intervention
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The Roundfix Skill and the command guides describe the queue that finishes |
| task_02 | backend | An archived Delivery Retry finds the start head of a Run the queue started |
| task_03 | backend | The queue merges only on GitHub's mergeable state and resumes a delivery-error retry |
| task_04 | backend | The agent environment drops a NODE_OPTIONS preload whose file is missing |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01, task_04 · 2 → task_02, task_03 · 3 → task_05

task_02 and task_03 change the `deliver` command, whose guide task_01 writes,
so they follow it; they share no declared file with each other. task_04
changes no command file and shares no file with any other Task.
