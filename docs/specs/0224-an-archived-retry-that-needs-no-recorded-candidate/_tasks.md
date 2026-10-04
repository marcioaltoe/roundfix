---
schema: spec-tasks/v1
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_03]
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03, task_05]
    - id: task_05
      file: task_05.md
      needs: [task_01]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A Delivery Retry resumes an operator-archived item from its Run start, whatever its park |
| task_02 | backend | An environment-only QA partial parks for the operator even when only Pull Request rows block it |
| task_03 | docs | The Roundfix Skill and the deliver guide describe both rules |
| task_04 | qa | Run the final QA gate |
| task_05 | backend | An override archive without a candidate or History refuses with the archived-head text |

Waves: 1 → task_01, task_03 · 2 → task_02, task_05 · 3 → task_04

The first three implementation Tasks share no declared file. task_01 changes the
retry in the delivery engine and task_03 the `deliver` guide and the Roundfix
Skill, written from the TechSpec, so they run together. task_02 changes the
park classification in the CLI workflow, whose guide task_03 writes, so it
follows task_03. task_05, the corrective Task for finding F1 of the
2026-10-04 QA Report, changes the retry in the delivery engine that task_01
changed, so it follows task_01. The gate follows all four.
