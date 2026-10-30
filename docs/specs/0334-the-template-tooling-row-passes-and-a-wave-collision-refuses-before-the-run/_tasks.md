---
schema: spec-tasks/v1
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
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

# Tasks — The template's tooling row passes the Spec Check, and a Wave collision is refused before the Run

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The templates' no-mutation Tooling authority row passes the Spec Check when its record grants no path   | backend | medium     | —       |
| task_02 | A completed Task never collides, and implement refuses a Wave collision before the Run                  | backend | high       | task_01 |
| task_03 | Run the final QA gate                                                                                   | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 follows task_01 because both raise owned-skill versions recorded in
`skills/testdata/owned-skill-versions.json` (task_01 for `write-prd` and
`write-techspec`, task_02 for `roundfix`). Each code Task runs the suites
that consume its contract in its own Verification, so no separate
Integration Checkpoint is needed. No other Spec is required; Spec 0333
delivers first (see `_techspec.md` → Build Order).
