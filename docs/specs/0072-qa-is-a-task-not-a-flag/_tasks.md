---
schema: spec-tasks/v1
spec: 0072-qa-is-a-task-not-a-flag
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
      needs: [task_02]
    - id: task_06
      file: task_06.md
      needs: [task_03]
---

# Tasks — QA is a Task, not a flag

| id      | title                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Teach the graph the gate and its invalidation      | backend | high       | —       |
| task_02 | Route the gate from the graph, not the request     | backend | high       | task_01 |
| task_03 | Delete the QA parameter from the Implement Command | backend | medium     | task_02 |
| task_04 | Author the gate decision in the owned skills       | docs    | medium     | task_01 |
| task_05 | Measure and trim the gate cycle's own cost         | backend | medium     | task_02 |
| task_06 | Align the agent guides with the authored gate      | docs    | low        | task_03 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03, task_05 ·
4 → task_06

Wave 2 fans out because the Daemon routing and the skill authoring rule
touch disjoint files once the graph contract exists. task_05 is the
performance slice the maintainer requested explicitly.

**This Spec's own graph stays on the v1 contract.** The Run that implements
it is driven by the binary built before it, which rejects unknown Task
Types, so no node here uses `type: qa` and this manifest carries no `qa:`
declaration. This Spec's own gate runs at close under the current flow; the
first Spec authored with a gate node is the next one in the queue
(`0074-git-spawn-economy`).
