---
schema: spec-tasks/v1
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
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
      needs: [task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The guides and the Roundfix Skill describe the credit floor, the named refusal and the relay |
| task_02 | backend | The Jev Router gate refuses below the account credit floor |
| task_03 | backend | A loopback relay records the routed model on the Judge Log |
| task_04 | backend | An OpenRouter credit refusal has its own name |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 changes the gate whose guide and skill reference task_01 writes.
task_03 and task_02 share the daemon's session owner and gate tests, and
task_04 builds on the relay's refusal note and shares the runner, the ledger
and those tests with task_03, so the backend Tasks run in series. The gate
follows all four.
