---
schema: spec-tasks/v1
spec: 0287-reports-that-count-what-they-say
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

# Tasks — Reports that count what they say

| id      | title                                                              | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The Linux Lane's pass line counts the packages that ran            | backend | medium     | —       |
| task_02 | The storage report measures each Run under the root it wrote to    | backend | medium     | task_01 |
| task_03 | An Archive Record summary carries what its colon introduces        | backend | medium     | task_02 |
| task_04 | Run the final QA gate                                              | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

One Task per defect. The three slices touch different packages, but each
raises the Roundfix skill version and re-records
`skills/testdata/owned-skill-versions.json`, the skill mirror and the Behavior
Surface Record, so they run in series. The gate follows the last of them. No
other Spec is required.
