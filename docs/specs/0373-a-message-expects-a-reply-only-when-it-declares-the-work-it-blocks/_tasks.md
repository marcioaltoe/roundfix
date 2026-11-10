---
schema: spec-tasks/v1
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
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

# Tasks — A message expects a reply only when it declares the work it blocks

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Follow-up Command ranks a message that declares blocked work first, without the judge     | backend | medium     | —       |
| task_02 | The Baseline's message clauses ask for a reply only when the sender's work is blocked         | docs    | medium     | task_01 |
| task_03 | The follow-up mode replies only to declared blocked work, and the vocabulary and decisions say so | docs | medium     | task_02 |
| task_04 | Run the final QA gate                                                                         | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The Tasks run in series. task_01 and task_03 both raise the Roundfix skill
version and change its `SKILL.md` pair and version record; task_03's step 8
cites the cross-project continuity rule task_02 rewrites; task_02 follows
task_01 so one Task at a time writes the shared generated records. Each
implementation Task runs the suites that consume its changed contract in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required; Specs 0372 and 0371 deliver first and edit other sections of
the shared files named in `_prd.md` → Prerequisites.
