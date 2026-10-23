---
schema: spec-tasks/v1
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
qa: task_03
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
---

# Tasks — The review drops a Result Status Note, and a stale fingerprint names its command

| id      | title                                                                                     | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The pre-PR review dismisses a Result Status Note of a completed Task without the validator | backend | high       | —                |
| task_02 | A stale Behavior Surface fingerprint names the command that re-records it                  | test    | low        | —                |
| task_03 | Run the final QA gate                                                                     | qa      | high       | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03

task_01 changes the review conventions, the convention validator, their
tests and the documents the documentation contract holds equal to the
conventions, so code and text land in one Task; it proves itself in
`internal/cli` and `internal/docscontract`, the packages that consume every
contract it changes, so it needs no separate Integration Checkpoint. task_02
changes only `internal/docscontract/skill_coverage_test.go`, which task_01
does not declare, so the two share no file and run in the same Wave. The gate
follows both. No other Spec is required.
