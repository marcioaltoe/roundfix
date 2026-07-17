---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/daemon/engine.go
line: 520
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5sx,comment:PRRC_kwDOS0qyts7Wt95L
review_hash: 8ae4144fbd0c53ec0a247709498d1587d2c28d529b9fcff62c1876d059b81ddf
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 007: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Do not discard owned-session close failures.**

A failed `owner.Close` is currently invisible, so the batch can report success while its session remains active. Join the close error into the return path or publish a terminal cleanup failure, and cover that path with a test.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/engine.go` around lines 517 - 520, Update the defer cleanup
around owner.Close in the batch execution flow to preserve any close error
instead of discarding it: join it into the function’s returned error or publish
a terminal cleanup failure. Add a test covering an owner.Close failure and
verifying the batch reports the cleanup failure.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:6738e3dc6260c5ea58ce8453 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
