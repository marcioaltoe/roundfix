---
schema: spec-tasks/v1
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
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

# Tasks — The Spec skills announce the Spec Events the user defines

| id      | title                                                                                     | type | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The glossary defines the Spec Events and the PRD writers announce new and deadline events | docs | medium     | —       |
| task_02 | The graph and archive skills announce authorized, completed and cancelled                 | docs | low        | task_01 |
| task_03 | The Supervisor announces the Run and queue events and the guide lists every owner         | docs | medium     | task_02 |
| task_04 | A skill-contract test holds the step generic and every owned skill free of personal names | test | medium     | task_03 |
| task_05 | Run the final QA gate                                                                     | qa   | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 to task_03 run in series because each edits owned skills and
re-records `skills/testdata/owned-skill-versions.json`. task_04 needs
task_03 because its test reads every section the three Tasks write. Each Task
runs the suites that consume its changed text in its own Verification, so no
separate Integration Checkpoint is needed. No other Spec is required.
