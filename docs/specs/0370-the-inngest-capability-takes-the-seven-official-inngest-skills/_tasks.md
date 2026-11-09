---
schema: spec-tasks/v1
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
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

# Tasks — The Inngest capability takes the seven official Inngest skills

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The backend module takes the seven Inngest skills, inngest is retired and the snapshots follow the catalog | backend | medium     | —       |
| task_02 | ADR-0206, the glossary and the guides say that inngest was replaced and retired                    | docs    | low        | task_01 |
| task_03 | Run the final QA gate                                                                              | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 shares no file with task_01 but documents its delivered behavior, so
it runs after it. task_01 runs the suites that consume its module, snapshots
and Retired Skill set in its own Verification, so no separate Integration
Checkpoint is needed. No other Spec is required; the queue delivers 0366 to
0369 first.
