---
schema: spec-tasks/v1
spec: 0264-baseline-group-usage-errors-name-their-own-help
requires: [0260-help-on-every-command-path-exits-zero]
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

# Tasks — Baseline group usage errors name their own help

| id      | title                                                        | type    | complexity | needs   |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The roundfix skill describes the bare-group usage error      | docs    | low        | —       |
| task_02 | Baseline usage errors name the group once and their own help | backend | low        | task_01 |
| task_03 | Run the final QA gate                                        | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 owns the skill text and its version record, task_02 owns the
dispatcher and its tests; they share no declared file. task_02 depends on
task_01 because the Task that changes the `baseline` command surface needs
the Task that documents it in its ancestry. The whole graph requires Spec
0260 merged: task_01's sentence follows the one Spec 0260 adds to the same
reference, and task_02's pointer names a help Spec 0260 makes exit 0 and its
class guard reuses Spec 0260's `helpCommandPaths`.
