---
schema: spec-tasks/v1
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
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

# Tasks — The Roundfix skill references agree with the CLI

| id      | title                                                                                                   | type | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The review, events and deliver documents and the glossary state what the CLI does                       | docs | low        | —       |
| task_02 | Documentation contract tests read the conventions and the event categories from the code               | test | medium     | task_01 |
| task_03 | Run the final QA gate                                                                                   | qa   | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 corrects the documents; task_02 adds the guards that require the
corrected documents, so on the tree before task_01 they fail and task_02
follows it. No production code changes, so no Shared Contract has a consumer
outside `internal/docscontract` and no Integration Checkpoint is needed. The
gate follows both. No other Spec is required.
