---
schema: roundfix/archive-record/v1
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
title: A turn that did not finish is not accepted
status: archived
created: "2026-10-13"
archived: "2026-10-14"
disposition: pass
source: docs/specs/0268-a-turn-that-did-not-finish-is-not-accepted
source_revision: ff37d99daf0715d9c5c83582fedf80363c9bfe26
qa_task: task_05
qa_report: qa-report-2026-10-14.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0274
sources:
  - 2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md
  - 2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md
regeneration:
  - command: make skills-sync
promoted: []
---

# A turn that did not finish is not accepted

The Archive Record may carry this QA pass and the report path. It carries no QA override or unreachable acceptance claim. Archiving, QA Task settlement, commits, pushes and Pull Request operations remain with the parent Daemon/delivery flow.
