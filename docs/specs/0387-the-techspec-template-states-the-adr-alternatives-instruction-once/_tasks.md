---
schema: spec-tasks/v1
spec: 0387-the-techspec-template-states-the-adr-alternatives-instruction-once
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — The TechSpec template states the ADR Alternatives instruction once

| id      | title                                                                                                                         | type | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | The TechSpec template states the ADR Alternatives instruction once, and owned skills refuse a phrase repeated across a line break | docs | low        | —       |
| task_02 | Run the final QA gate                                                                                                         | qa   | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 changes one template line and its mirror, raises and records the
skill's version, and adds the guard; it runs the whole `./skills` package that
consumes the shipped skill and the version record in its own Verification, so
no separate Integration Checkpoint is needed. The gate follows it. No other
Spec is required; Specs 0383, 0385, 0384 and 0386 deliver first and append
only `roundfix` versions to the shared version record.
