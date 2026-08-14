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
severity: review
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95k,comment:PRRC_kwDOS0qyts7V6VrK
review_hash: 1438372997fe1969b14f4fb7dcf0bc72eb392e734a0c240c1a070a3e0bf879e8
duplicate_of: ""
source_review_id: "4705912088"
source_review_submitted_at: "2026-07-15T16:00:46Z"
---

# Issue 024: @marcioaltoe, this particular outcome looks safe: its reason is generic and c...

## Review Comment

`@marcioaltoe`, this particular outcome looks safe: its reason is generic and contains no local path or injected Markdown.

However, the finding still applies generally: `outcomeCommentBody` publishes `issue.TerminalReason` and `issue.DuplicateOf` verbatim, and those fields are not guaranteed to have this safe shape. A failed verification reason can include a diagnostics path, and `DuplicateOf` is currently a local artifact reference. The values should still be sanitized/bounded, with duplicate references mapped to public `SourceRef` values before posting.

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `outcomeCommentBody` now publishes bounded public reasons through `publicOutcomeReason` and maps duplicate references through `publicDuplicateReference`; `make verify` passed.
