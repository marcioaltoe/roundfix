---
schema: spec-tasks/v1
spec: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
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

# Tasks — One blocking message per sender goes first, and every direct question is answered

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Follow-up Command gives the declared precedence to one message per sender                      | backend | medium     | —       |
| task_02 | The Baseline's receiving clause answers every direct question and bounds the declared precedence   | docs    | medium     | task_01 |
| task_03 | The follow-up mode, the vocabulary and the decisions follow the bounded precedence and reply rule  | docs    | medium     | task_02 |
| task_04 | Run the final QA gate                                                                              | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The Tasks run in series. task_01 and task_03 both raise the Roundfix skill
version and change its `SKILL.md` pair and version record; task_03's step 4
restates the receiving clause task_02 rewrites; task_02 follows task_01 so one
Task at a time writes the shared generated records. Each implementation Task
runs the suites that consume its changed contract in its own Verification, so
no separate Integration Checkpoint is needed. No other Spec is required;
Specs 0377, 0376, 0378 and 0379 deliver first and edit other sections of the
shared files named in `_prd.md` → Prerequisites.
