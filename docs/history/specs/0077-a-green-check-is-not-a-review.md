---
schema: roundfix/archive-record/v1
spec: 0077-a-green-check-is-not-a-review
title: A green check is not a review
status: archived
created: "2026-08-05"
archived: "2026-08-05"
disposition: pass
source: docs/history/specs/0077-a-green-check-is-not-a-review
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_05
qa_report: qa-report-2026-08-05-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "124"
delivery_commit: 830326d31ea53bb8395a6548599e5e216f8d4824
---

# A green check is not a review

When CodeRabbit declines to review because the account hit its rate limit, it publishes a check named `Review rate limited` whose conclusion is **success**, deliberately, so the block does not prevent a merge. The authoritative signal is the comment it posts, not the green check.
