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
line: 1112
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95k,comment:PRRC_kwDOS0qyts7V5taZ
review_hash: 83e8eb2cde2c168809d4887faf4eb9938ac7a0ece6e94b90f75d17d786f1cc36
duplicate_of: ""
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

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Outcome Comments now use public Source References for duplicates and sanitized, bounded public reasons; `make verify` passed.
