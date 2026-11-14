---
schema: spec-tasks/v1
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — A feature is judged for the Idea stage before its PRD

| id      | title                                                                              | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The judge answers the Idea-stage criteria for one item or one request of the session | backend | high       | —       |
| task_02 | The route clauses judge a feature for the Idea stage before its PRD                 | docs    | medium     | task_01 |
| task_03 | The authoring skills and the follow-up mode run and record the Idea-stage judgment  | docs    | medium     | task_01, task_02 |
| task_04 | Run the final QA gate                                                              | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The Tasks run in series. task_02's first clause names the command task_01
adds, and task_03's skills restate the rule task_02 fixes; task_01 and task_03
both raise the Roundfix skill version, and every Task runs the documentation
contracts, which read the guide, the glossary and the skills. Each
implementation Task runs the suites that consume its changed contract in its
own Verification, so no separate Integration Checkpoint is needed. The Spec is
one context in three implementation Tasks that each earn their verification
cycle: the command with its guide, the Baseline with its derived files, and
the skills with the glossary; the amendment notice in ADR-0201 joins task_03
because alone it would not earn one. No other Spec is required; Specs 0391 and
0392 deliver before this one and edit other sentences of the shared files
named in `_prd.md` → Prerequisites.
