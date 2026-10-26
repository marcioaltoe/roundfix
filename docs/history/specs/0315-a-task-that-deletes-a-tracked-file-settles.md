---
schema: roundfix/archive-record/v1
spec: 0315-a-task-that-deletes-a-tracked-file-settles
title: A Task that deletes a tracked file settles
status: archived
created: "2026-10-24"
archived: "2026-10-25"
disposition: partial
source: docs/specs/0315-a-task-that-deletes-a-tracked-file-settles
source_revision: ddcae12c554d678b0a98c0d2b4060abfaf2a4e3e
qa_task: task_03
qa_report: qa-report-2026-10-25.md
qa_verdict: partial
unproven:
  - oraculum's next Run of a Task that declares `deletes:` with the released binary, reported back through the Secondbrain inbox.
adrs: []
sources:
  - 2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md
regeneration: []
promoted: []
pull_request: "678"
---

# A Task that deletes a tracked file settles

The local gate verifies disk-only and staged deletion settlement, surviving link and authorization checks, all three readers and their error boundaries, and scope preservation. The verdict is partial solely because the post-release oraculum observation is declared unreachable and no Pull Request exists; the archive plan result is recorded in row 12. The operator must obtain the after-release oraculum report for Success Metric 7. No repository Verification was repeated, no Task state changed and no commit, push or Pull Request was created.
