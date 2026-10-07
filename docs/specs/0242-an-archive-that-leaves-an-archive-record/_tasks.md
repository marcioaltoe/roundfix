---
schema: spec-tasks/v1
spec: 0242-an-archive-that-leaves-an-archive-record
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_06, task_07]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_04]
---

# Tasks — An archive that leaves an Archive Record

| id      | title                                                                                                   | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | An archive writes the Archive Record and removes the Spec folder, and the Delivery Queue accepts it     | backend | high       | —                |
| task_02 | Reconcile, the pre-PR review and the queue's other readers work from the Archive Record                 | backend | high       | task_01          |
| task_03 | The Spec check, the audit, the suite guard and the repository tests need no archived Spec folder        | backend | high       | task_01          |
| task_04 | Archive Advice and promotion, and the skills, guides, Baseline clauses and glossary describe the record | docs    | medium     | task_02, task_03 |
| task_05 | Run the final QA gate                                                                                   | qa      | high       | task_06, task_07 |
| task_06 | The archive confirmation names promoted files, and a keyless plan says no advice                        | backend | low        | task_04          |
| task_07 | Absorbed history resolves without the Spec folders, and the archive guidance describes only the record  | backend | medium     | task_04          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_06, task_07 · 5 → task_05

task_01 ships the record, the cut and the Delivery Queue's archive stage
together. The queue's and supersede's existing tests run the real Archive
Command, so neither can follow it. task_02 and task_03 read the resolver
task_01 adds. They share no declared file, so they run in parallel. task_04
extends the Archive Command and `internal/cli/archive.go`, which task_01
changes. It also extends `internal/judge/questions_test.go`, which task_03
changes, and it documents what task_02 and task_03 settle, so it follows
both.

task_06 and task_07 are the corrective Tasks for the first QA gate's
findings. task_06 fixes the promotion suffix and the keyless plan line
(F-01, F-02) in `internal/cli/archive.go`. task_07 makes absorbed history
resolve through the record or Git and rewrites the stale archive guidance
(F-03, F-04). They share no declared file, so they run in parallel, and the
gate now follows both.
