---
schema: roundfix/archive-record/v1
spec: 0256-an-archive-record-that-survives-a-squash-merge
title: An Archive Record that survives a squash merge
status: archived
created: "2026-10-12"
archived: "2026-10-12"
disposition: partial
source: docs/specs/0256-an-archive-record-that-survives-a-squash-merge
source_revision: 7d30e2d2dfdd4ca0b0c0b31046046c4008dfc0e6
qa_task: task_05
qa_report: qa-report-2026-10-12-01.md
qa_verdict: partial
unproven:
  - the authoring measurement in `docs/references/2026-10-12-squash-merged-archive-source-revisions.md` and the hermetic local-remote tests of Success Metric 2.
adrs:
  - ADR-0261
sources:
  - 2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md
  - 2026-10-12-review-reads-any-markdown-in-the-archive-root-as-an-archive-record.md
regeneration:
  - command: make skills-sync
promoted: []
---

# An Archive Record that survives a squash merge

The Archive Record may carry qualifying partial QA: all local archive, recovery, delivery, transcript, documentation and history-scope checks pass; the previous two documentation findings are resolved. Real GitHub recovery remains unproven under its existing declaration, supported by the committed independent measurement and hermetic four-reader evidence. The pre-PR journey remains environment-blocked because no Pull Request is open.
