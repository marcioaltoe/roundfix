---
schema: spec-tasks/v1
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
qa: task_04
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
---

# Tasks — An older Roundfix refuses a repository a newer Roundfix wrote

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Configuration loading refuses a newer Recorded Roundfix Version and the recovery commands warn | backend | high       | —       |
| task_02 | A Baseline plan records its writer's release and the Baseline commands refuse a newer one      | backend | high       | task_01 |
| task_03 | The Roundfix baseline and setup references describe the record and the refusal                 | docs    | medium     | task_02 |
| task_04 | Run the final QA gate                                                                          | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The Tasks run in series. task_02 calls the check task_01 adds and shares
`internal/cli` test helpers and guides with it; each code Task documents its
own surface in the command guides, and task_03 the strings both print in the
Roundfix skill's references.
Each code Task runs the suites that consume its changed surfaces in its own
Verification, so no separate Integration Checkpoint is needed. No other Spec
is required; Specs 0377, 0376, 0378, 0379, 0380 and 0381 deliver first and
edit other sections of the shared files named in `_prd.md` → Prerequisites.
