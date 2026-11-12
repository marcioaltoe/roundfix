---
schema: spec-tasks/v1
spec: 0386-a-message-that-asks-for-a-completion-notice-receives-one
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

# Tasks — A message that asks for a completion notice receives one

| id      | title                                                                                         | type | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The Baseline sends one completion notice to a message that asks for it                        | docs | medium     | —       |
| task_02 | The follow-up mode and the release reference carry a Notice Request to its completion notice | docs | medium     | task_01 |
| task_03 | Run the final QA gate                                                                         | qa   | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The Tasks run in series. task_02's follow-up step 4 and release section
restate the receiving clause task_01 rewrites, and both Tasks run the
documentation contracts, which read the glossary and the managed guide. Each
implementation Task runs the suites that consume its changed contract in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required; Specs 0383, 0384 and 0385 deliver first and edit other
sections of the shared files named in `_prd.md` → Prerequisites.
