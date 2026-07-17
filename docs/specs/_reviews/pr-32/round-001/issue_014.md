---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/store/agent_selection.go
line: 153
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tO,comment:PRRC_kwDOS0qyts7Wt95u
review_hash: 6e43c6e29f127594cf3ffa00c25e0662f05bd724ecf1f1e52f84c4b7a1e1604f
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 014: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Reject exhaustion unless the latest attempt failed.**

Any nonempty history currently permits an exhausted event, including a latest `active` or `closed` attempt. Require the latest persisted status to be `failed`, and derive the category from that history—or reject a mismatched request category.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/agent_selection.go` around lines 146 - 153, Update the
exhaustion flow around selectAgentSelectionAttemptsForScope and
agentSelectionExhaustedEvent to inspect the latest persisted attempt, requiring
its status to be failed before creating an exhausted event. Derive the event
category from that failed attempt’s history, or reject the request when its
category does not match; continue rejecting empty histories.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8a267c7cf56443647105bf55 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
