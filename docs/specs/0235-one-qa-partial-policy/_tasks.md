---
schema: spec-tasks/v1
spec: 0235-one-qa-partial-policy
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

# Tasks — One QA partial policy

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A partial whose only unmet rows nobody in the Run can reach qualifies under the one eligibility policy   | backend | medium     | —       |
| task_02 | Settlement, archive, settle and the queue agree on every partial, and settle names the report it refused | backend | medium     | task_01 |
| task_03 | The skills and user guides state the one partial policy and the network-denied marker                   | docs    | medium     | task_02 |
| task_04 | The Spec routing and autonomous-work guides exempt a network-denied outside-evidence row                 | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                                   | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 depends on the policy task_01 puts in `internal/spec`. task_03
documents both and shares `docs/user-guide/commands/settle.md` with task_02.
task_04 follows task_03 so that only one Task at a time edits guidance and
runs the derived-file regeneration. The gate follows the only leaf.
