---
schema: spec-tasks/v1
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | No clause is rendered twice, and every clause keeps its force |
| task_02 | backend | The loop clause follows the Delivery Queue |
| task_03 | backend | Authorization and evidence clauses state today's refusals |
| task_04 | backend | Lifecycle wording is consistent and cites no repository record |
| task_05 | qa | Run the final QA gate |
| task_06 | backend | A removed Source Baseline clause is replaced, not lost |
