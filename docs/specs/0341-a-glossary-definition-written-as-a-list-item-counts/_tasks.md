---
schema: spec-tasks/v1
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
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

# Tasks — A glossary definition written as a list item counts

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The glossary reader accepts list-item and em dash definitions and the findings name the forms | backend | medium     | —       |
| task_02 | The authoring guides state the forms and ask for the line-anchored glossary Verification      | docs    | low        | task_01 |
| task_03 | Run the final QA gate                                                                         | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 needs task_01 because its test reads `GlossaryDefinitionForms` and
`GlossaryVerificationCommand`, which task_01 adds, and because the spec
command guide task_01 edits is a source of the Behavior Surface the Roundfix
`spec` reference covers. Each Task runs the suites that consume its changed
contract in its own Verification, so no separate Integration Checkpoint is
needed. No other Spec is required.
