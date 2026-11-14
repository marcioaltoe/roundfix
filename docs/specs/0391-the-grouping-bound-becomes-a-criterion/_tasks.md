---
schema: spec-tasks/v1
spec: 0391-the-grouping-bound-becomes-a-criterion
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

# Tasks — The grouping bound becomes a criterion

| id      | title                                                                                   | type | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The Baseline sizes a Spec by one context and by Tasks that earn their verification cycle | docs | medium     | —       |
| task_02 | The authoring skills, the guide and the glossary state the grouping bound as a criterion | docs | medium     | task_01 |
| task_03 | Run the final QA gate                                                                   | qa   | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The Tasks run in series. task_02's skills restate the clause task_01 rewrites,
and both Tasks run the documentation contracts, which read the glossary and
the managed guide. Each implementation Task runs the suites that consume its
changed contract in its own Verification, so no separate Integration
Checkpoint is needed. Applied to this Spec, the criterion it writes gives two
Tasks, one per carrier family: a third Task for the glossary and the ADR
notice alone would not earn its verification cycle. No other Spec is
required; Specs 0392 and 0393 deliver after this one and edit other sections
of the shared files named in `_prd.md` → Prerequisites.
