---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/notify/notify.go
line: 54
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpj,comment:PRRC_kwDOS0qyts7aUVD0
review_hash: 837b4665d6e5a17ec547cfb32ec03bf17a3f06d229683db7b3d214419bc7886f
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 006: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Consider defined types for `Route` and `Status`.**

Both are untyped string constants, so nothing stops a caller from passing a route where a status is expected — and `internal/cli` already builds journal keys with `"outcome_notification_" + receipt.Status`. Defining `type Route string` and `type Status string` makes those mix-ups compile errors while this API still has a single consumer.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/notify/notify.go` around lines 17 - 54, The Route and Status
constants are untyped strings, allowing routes and statuses to be mixed
accidentally. Define named Route and Status string types, update their constants
and the corresponding NotificationReceipt fields, and adjust the internal/cli
journal-key construction to explicitly convert Status to string where
concatenation requires it.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:a973484423e56d007b3edb21 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Introduced named `notify.Route` and `notify.Status` types and updated receipt producers/consumers. `go test ./internal/notify` and focused CLI receipt tests passed.
