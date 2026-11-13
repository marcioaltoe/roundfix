---
schema: spec-tasks/v1
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
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

# Tasks — A nested carrier Git does not list raises no conflict

| id      | title                                                                        | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Baseline inventory keeps only the nested carriers Git lists               | backend | medium     | —       |
| task_02 | The baseline guide and reference state the rule; the amendment notice lands  | docs    | low        | task_01 |
| task_03 | Run the final QA gate                                                        | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The Tasks run in series: task_02 documents the behavior task_01 delivers and
shares no file with it. task_01 runs the suites that consume the inventory
(`./internal/baseline` and the `Baseline` tests of `./internal/cli`) in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required; Spec 0389 delivers first and touches only the shared
version fields and records named in `_prd.md` → Prerequisites.
