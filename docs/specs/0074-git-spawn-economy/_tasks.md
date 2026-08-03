---
schema: spec-tasks/v1
spec: 0074-git-spawn-economy
qa: task_07
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_02, task_03, task_04, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Tasks — Git spawn economy

| id      | title                                             | type    | complexity | needs                              |
| ------- | ------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | Commit the spawn census baseline                  | docs    | medium     | —                                  |
| task_02 | Batch object reads in skills restore              | backend | high       | task_01                            |
| task_03 | Batch object reads in assets-sync provenance      | backend | medium     | task_02                            |
| task_04 | Combine repository resolution queries             | backend | medium     | task_01                            |
| task_05 | Give the agent runner its environment explicitly  | backend | high       | task_01                            |
| task_06 | Publish the before-and-after                      | docs    | low        | task_02, task_03, task_04, task_05 |
| task_07 | Run the final QA gate                             | qa      | high       | task_06                            |

Waves: 1 → task_01 · 2 → task_02, task_04, task_05 · 3 → task_03 ·
4 → task_06 · 5 → task_07

Wave 2 fans out because its three Tasks touch disjoint surfaces: the skills
restore reader, repository resolution, and the agent runner's environment.
task_03 waits for task_02 because it reuses the batch-reader shape proven
there.

The QA gate is authored as task_07 — this is the first Task Graph born under
the contract Spec 0072 shipped, so the gate exists here from decomposition,
not as a close-time amendment. It depends on task_06, which transitively
covers every other node.
