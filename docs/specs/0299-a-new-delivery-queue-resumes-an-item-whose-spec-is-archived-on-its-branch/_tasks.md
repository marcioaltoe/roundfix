---
schema: spec-tasks/v1
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
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

# Tasks — A new Delivery Queue resumes an item whose Spec is already archived on its branch

| id      | title                                                                                                  | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The deliver guide, the skill reference and the glossary describe the resume of an archived item branch | docs    | low        | —       |
| task_02 | The Delivery Engine resumes a proven archived item at review and parks an unproven one                 | backend | medium     | task_01 |
| task_03 | The queue proves an archive on a continued item branch and delivers it through deliver start           | backend | medium     | task_02 |
| task_04 | Run the final QA gate                                                                                  | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 writes the `deliver` guide, the Roundfix `deliver` reference and the
glossary for the surface the later Tasks change, so they follow it in one
dependency chain. task_02 adds the engine dependency that task_03 implements
in `internal/cli`; task_03 is its Integration Checkpoint, because its tests
drive the engine through `roundfix deliver start` with the real workflow. No
two Tasks declare the same file. No other Spec is required.
