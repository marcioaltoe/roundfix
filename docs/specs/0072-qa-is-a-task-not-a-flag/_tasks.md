---
schema: spec-tasks/v1
spec: 0072-qa-is-a-task-not-a-flag
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
      needs: [task_02]
    - id: task_06
      file: task_06.md
      needs: [task_03]
    - id: task_07
      file: task_07.md
      needs: [task_04, task_05, task_06]
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
| task_07 | Run the final QA gate                              | qa      | high       | task_04, task_05, task_06 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03, task_05 ·
4 → task_06 · 5 → task_07

Wave 2 fans out because the Daemon routing and the skill authoring rule
touch disjoint files once the graph contract exists. task_05 is the
performance slice the maintainer requested explicitly.

**The gate was authored into this graph the moment the contract went
live.** The decomposition originally kept this manifest on the v1 shape,
because the binary that would execute the implementation predated the
`qa` type it introduces — and the six implementation Tasks did run that
way. At close, the rebuilt binary carried the new contract and refused
`--qa` with the remediation task_03 wrote for it, so the premise expired:
task_07 is the authored terminal gate, depending on every leaf, and this
Spec closes under the contract it built. Its first gate attempt failed and
was remediated (see `qa/qa-report-2026-08-03.md` and the 2026-08-03
authorization addendum); the gate re-runs as the graph's terminal node.
