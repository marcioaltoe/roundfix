---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/daemon/engine.go
line: 923
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95e,comment:PRRC_kwDOS0qyts7V5taT
review_hash: 25c059c2d01caecb4eef69ea5ec700bf96dd52c8fe9fd22f4c4732b5449b86a6
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 010: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Keep failed source propagation pending instead of counting it resolved.**

On comment or resolution failure, `action.Resolved` remains true, so `ResolvedSourceThreads` is incremented despite no resolution. The error is then swallowed while the local terminal issue is excluded from future retry selection, causing permanent source drift.

Track actual completion separately from requested actions and persist/retry failed propagation before reporting success.






Also applies to: 969-1001

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/engine.go` around lines 918 - 923, Update the action
completion accounting around action.Resolved and action.Failed so failed comment
or resolution propagation is not treated as resolved. Track actual propagation
completion separately from requested actions, persist failed propagation as
pending for retry, and increment summary.Resolved/ResolvedSourceThreads only
after successful completion; preserve the failure count and retry eligibility.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fdc6cd00383453f85b0b7fca -->

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Review Source propagation now distinguishes requested resolution from completed resolution and marks propagation failures retryable; `make verify` passed.
