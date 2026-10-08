---
schema: spec-tasks/v1
spec: 0251-skills-keep-up-with-the-behavior-they-describe
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

# Tasks — Skills keep up with the behavior they describe

| id      | title                                                                                        | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The skill coverage package parses the map and the record and finds Lagging Surfaces          | backend | medium     | —       |
| task_02 | The repository keeps a Skill Coverage Map and a current Behavior Surface Record              | backend | high       | task_01 |
| task_03 | The release plan refuses a release with a Lagging Surface                                    | backend | high       | task_02 |
| task_04 | Spec authoring requires a skills Task for a changed Behavior Surface                         | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                                        | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 reads the package task_01 adds. task_03 reads the package and
re-records the record task_02 creates. task_04 shares the record, the
Roundfix skill entry and its mirror, the skill version record and `CONTEXT.md`
with task_03, so the two run in series. The gate follows the only leaf.

task_03 and task_04 declare the Behavior Surface Record under `creates:`
although task_02 creates it, because a Task Context `interface:` path must
exist when the Spec is checked; both labels declare a write.
