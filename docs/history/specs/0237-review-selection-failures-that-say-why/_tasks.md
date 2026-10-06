---
schema: spec-tasks/v1
spec: 0237-review-selection-failures-that-say-why
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — Review selection failures that say why

| id      | title                                                                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The review reference and guide describe the failing step, the adapter's message and the retry      | docs    | low        | —                |
| task_02 | The ACPX Runner keeps the failing protocol step, the adapter's message and whether the prompt was sent | backend | medium     | —                |
| task_03 | The review names the failing step and retries a failure before the prompt once                     | backend | medium     | task_01, task_02 |
| task_04 | Run the final QA gate                                                                              | qa      | high       | task_03          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

task_01 writes the `review` guide and skill reference the CLI Task changes,
and task_02 adds the runner detail task_03 reads, so task_03 follows both.
task_01 and task_02 share no declared file. The gate follows the only leaf.
