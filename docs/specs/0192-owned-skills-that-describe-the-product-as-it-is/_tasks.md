---
schema: spec-tasks/v1
spec: 0192-owned-skills-that-describe-the-product-as-it-is
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — Owned skills that describe the product as it is

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Every command the CLI lists is named in the Roundfix skill | test | medium | — |
| task_02 | The user guide states what ships, and its links resolve | test | medium | task_01 |
| task_03 | The authoring skills teach a Spec the Delivery Queue accepts | docs | high | — |
| task_04 | The gate, lifecycle and discovery skills agree with the product | docs | medium | — |
| task_05 | Run the final QA gate | qa | medium | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01, task_03, task_04 · 2 → task_02 · 3 → task_05
