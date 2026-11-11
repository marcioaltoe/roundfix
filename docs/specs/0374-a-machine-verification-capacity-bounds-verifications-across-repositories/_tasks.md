---
schema: spec-tasks/v1
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
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

# Tasks — A Machine Verification Capacity bounds Verifications across repositories

| id      | title                                                                                                  | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The Machine Verification Capacity has its User Config key and its entry store                          | backend | high       | —       |
| task_02 | Every Daemon Verification takes the Machine Verification Capacity after the Run's own                  | backend | high       | task_01 |
| task_03 | The repository gate takes the capacity, and runs list, deliver status and doctor name its holders     | backend | high       | task_02 |
| task_04 | The glossary defines the Machine Verification Capacity, and ADR-0056 names its amendment              | docs    | low        | task_03 |
| task_05 | Run the final QA gate                                                                                  | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series. task_02 and task_03 call the store task_01 adds;
tasks 01 to 03 each re-record the Behavior Surface Record, and tasks 02 and
03 each raise the Roundfix skill version; task_04 defines the vocabulary the
guides of tasks 02 and 03 use. Each
implementation Task runs the suites that consume its changed contract in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required; Spec 0373 delivers first and edits other sections of the
shared files named in `_prd.md` → Prerequisites.
