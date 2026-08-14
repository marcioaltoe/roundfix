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
line: 1112
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95k,comment:PRRC_kwDOS0qyts7V5taZ
review_hash: 6f646fc13bed19919f51255abe0909901028217acc07045db31497de6e156810
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_011.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 011: _ Security & Privacy_ _ Major_ _ Heavy lift_

## Review Comment

_🔒 Security & Privacy_ | _🟠 Major_ | _🏗️ Heavy lift_

**Do not publish internal paths or untrusted reasons verbatim.**

`DuplicateOf` contains a local artifact path, while `TerminalReason` can contain diagnostics paths or agent-generated Markdown. Posting these directly to the remote Review Source leaks workstation details and permits unwanted mentions/content.

Resolve the canonical issue to a public `SourceRef`, and emit a bounded, sanitized public reason while retaining full diagnostics locally.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/engine.go` around lines 1093 - 1112, Update
outcomeCommentBody to avoid publishing raw issue.DuplicateOf or
issue.TerminalReason values. Resolve duplicated issues to their public SourceRef
before constructing the comment, and sanitize and length-bound terminal reasons
for invalid and failed outcomes while retaining full diagnostics only in local
logs or state. Keep the existing fallback messages and outcome-specific next
actions intact.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:83fe72b808f61ee0fc15a860 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
