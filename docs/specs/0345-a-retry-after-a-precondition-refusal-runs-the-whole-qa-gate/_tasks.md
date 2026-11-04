---
schema: spec-tasks/v1
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
qa: task_03
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
---

# Tasks — A retry after a QA precondition refusal runs the whole QA gate

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The QA Report Import refuses a Precondition Refusal, and its report's format path is pinned    | backend | medium     | —       |
| task_02 | The implement reference, the configuration guide and the glossary state the import rules       | docs    | medium     | task_01 |
| task_03 | Run the final QA gate                                                                         | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 needs task_01 because it documents the refusal reason and the rule
task_01 makes true; the two share no file. task_01 runs the existing import,
carry and format suites that consume the changed contract in its own
Verification, and task_02 runs the skill and documentation contract suites,
so no separate Integration Checkpoint is needed. No other Spec is required.
