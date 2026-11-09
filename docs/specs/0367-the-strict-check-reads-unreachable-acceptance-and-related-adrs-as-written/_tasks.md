---
schema: spec-tasks/v1
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
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

# Tasks — The strict check reads Unreachable Acceptance and related ADRs as written

| id      | title                                                                                                           | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The strict check names Unreachable Acceptance text the QA gate cannot read                                      | backend | medium     | —       |
| task_02 | Only an ADR the Spec applies draws in the ADRs that cite it                                                     | backend | medium     | task_01 |
| task_03 | The guides and the authoring skills say how to mark an ADR not applicable and what the new finding means        | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                                           | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 and task_02 both change `CONTEXT.md` and the `internal/speccheck`
package, so they run in series; task_03 documents the delivered behavior of
both and regenerates the skill version record and the Behavior Surface Record
on its own base. Each Task runs the suites that consume its contract in its
own Verification, so no separate Integration Checkpoint is needed. No other
Spec is required; the queue delivers 0365 and 0366 first.
