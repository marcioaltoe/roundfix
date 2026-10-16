---
schema: spec-tasks/v1
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Authoring guards hold the plan, the glossary and the PRD to their rules

| id      | title                                                                              | type    | complexity | needs            |
| ------- | ---------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Spec check reports a guide Task that does not record its surface fingerprint       | backend | medium     | —                |
| task_02 | The QA secret refusal is documented as the Daemon prints it                        | docs    | low        | task_01          |
| task_03 | A contract test holds the glossary to the rule it applies to Specs                 | test    | high       | —                |
| task_04 | A PRD carries a Clarification Record and spec check reports one without it         | backend | high       | task_02, task_03 |
| task_05 | Run the final QA gate                                                              | qa      | medium     | task_04          |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05

task_01, task_02 and task_04 each raise the Roundfix skill version and
re-record `skills/testdata/owned-skill-versions.json` and the Behavior Surface
Record, and task_01 and task_04 both edit the `spec` guide and reference, the
corpus pins and the detector wiring, so they run in series. task_03 touches
only its new test and `CONTEXT.md`, so it runs beside task_01; task_04 adds a
glossary entry that must pass task_03's contract, so it follows both. The gate
follows the only leaf.
