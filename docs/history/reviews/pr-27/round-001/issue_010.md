---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: internal/daemon/engine.go
line: 923
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95e,comment:PRRC_kwDOS0qyts7V5taT
review_hash: 7f35df0388c94cb1620c8506fdf030295f4b4cfeaa5120aeaca833e6bf0ba447
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_010.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
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

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
