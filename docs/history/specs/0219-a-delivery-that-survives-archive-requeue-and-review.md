---
schema: roundfix/archive-record/v1
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
title: A delivery that survives archive, requeue and review
status: archived
created: "2026-10-03"
archived: "2026-10-03"
disposition: qa-override
source: docs/history/specs/0219-a-delivery-that-survives-archive-requeue-and-review
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-03-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked row is the Pull Request row (R10)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 6d7bb6f1722526228a978b6133c72af9b61d6a1a
adrs:
  - ADR-0223
  - ADR-0226
sources:
  - 2026-10-01-archiving-or-requeueing-a-spec-loses-or-breaks-its-delivery.md
  - 2026-10-02-the-review-flags-an-authorized-qa-override-archive.md
  - 2026-09-30-authored-verification-runs-without-a-provenance-check.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "364"
delivery_commit: 2ef4556401dc0ff35e9761fe09aac0719a74d813
---

# A delivery that survives archive, requeue and review

Between 2026-10-01 and 2026-10-03 the operator finished five Delivery Queue items by hand although their work was correct, and each stop came after the Spec's own work was done. Spec 0205's test read its own active TechSpec and failed the repository gate once the archive commit moved it. The Pull Request then failed the corpus check, because Spec 0218's Task Context still named Spec 0205's active path. The Delivery Retry after the operator's correction refused because the archived item head had moved, so the operator opened the Pull Request by hand.
