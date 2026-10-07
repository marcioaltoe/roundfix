---
schema: roundfix/archive-record/v1
spec: 0027-review-loop-integrity
title: Review Loop Integrity
status: archived
created: "2026-07-14"
archived: "2026-07-15"
disposition: pass
source: docs/history/specs/0027-review-loop-integrity
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-14.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0042
  - ADR-0043
sources: []
regeneration: []
promoted: []
---

# Review Loop Integrity

Field runs against real pull requests (PR #4 and PR #17 sessions, 2026-07-14) reproduced three integrity gaps in the review loop. First, a review Run can start from a HEAD that omits work stranded on Run Branches in kept worktrees — it then reviews stale code, and its Final Push can silently publish a HEAD missing completed work. Second, when the Review Source check never appears for the pushed head, watch declares Clean with only a stderr warning, so a script caller cannot tell a verified Merge-Ready outcome from an unverified one. Third, Review Issue outcomes are invisible or misleading on GitHub: invalid issues are resolved silently, failed and duplicated threads stay open with no explanation, and the final report mixes per-Run and cumulative counts. This spec makes the review loop anchored, verifiable, and auditable end to end.
