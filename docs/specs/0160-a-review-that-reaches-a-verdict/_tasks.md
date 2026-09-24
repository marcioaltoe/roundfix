---
schema: spec-tasks/v1
spec: 0160-a-review-that-reaches-a-verdict
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
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06, task_07, task_08, task_09, task_10]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
    - id: task_09
      file: task_09.md
      needs: [task_08]
    - id: task_10
      file: task_10.md
      needs: [task_09]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Classify the verdict by substance and keep the answer |
| task_02 | backend | Run the claude provider |
| task_03 | backend | Review a Spec's delivery against its decisions |
| task_04 | docs | Keep the shipped skill and the guide true |
| task_06 | backend | No findings header escapes, no answer is invented |
| task_07 | backend | Spec context that never blocks and stays bounded |
| task_08 | docs | Describe the Spec-context contract the correctives delivered |
| task_09 | backend | A no-findings verdict must account for the whole answer |
| task_10 | docs | Describe the whole-answer verdict contract |
| task_05 | qa | Run the final QA gate |
