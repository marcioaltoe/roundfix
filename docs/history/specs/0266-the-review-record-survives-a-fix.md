---
schema: roundfix/archive-record/v1
spec: 0266-the-review-record-survives-a-fix
title: The review record survives a fix
status: archived
created: "2026-10-13"
archived: "2026-10-13"
disposition: pass
source: docs/specs/0266-the-review-record-survives-a-fix
source_revision: 028aba827f9b3b7fc29ef0c3183062a0433c46cd
qa_task: task_05
qa_report: qa-report-2026-10-13.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0272
  - ADR-0233
sources:
  - 2026-10-13-the-review-correction-proof-misses-its-record-under-a-symlinked-root.md
  - 2026-10-13-review-round-two-overwrites-the-parked-findings-record.md
regeneration:
  - command: make skills-sync
promoted: []
---

# The review record survives a fix

On 2026-10-13 the Delivery Queue parked items 0259 and 0258 as `corrective-spec-required` after their pre-PR review. The operator fixed each finding, disposed it with `roundfix review dispose --fixed-by`, and ran `roundfix deliver retry`; both retries were refused for a missing review correction record, and both Pull Requests (#512 and #514) were opened by hand (operator log entries 307, 308 and 312). Two defects compound, one per adopted Backlog Entry:
