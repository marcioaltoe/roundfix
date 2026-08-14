---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/reviewsource/coderabbit/coderabbit.go
line: 288
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpn,comment:PRRC_kwDOS0qyts7aUVD5
review_hash: b0cd71ab0682dbca2c594d62eefbd1de26891a48321e76c839d7ffc2577360ed
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 007: _ Performance & Scalability_ _ Trivial_ _ Poor tradeoff_

## Review Comment

_🚀 Performance & Scalability_ | _🔵 Trivial_ | _⚖️ Poor tradeoff_

**Every evidence poll now makes four GitHub round trips.**

`Evidence` unconditionally fetches check runs, commit statuses, PR reviews, and review threads, and it is called on each poll interval by both the wait loop and the Merge-Ready confirmation. On a long review wait this multiplies `gh` invocations and API quota consumption several-fold versus the previous status-only path; the new 429/5xx transient handling makes rate-limit pressure a realistic outcome rather than a theoretical one.

Consider short-circuiting the reviews/threads fetches when a current-head check run is already pending or reviewing (thread counts only affect verified-vs-reviewed and detail text), or caching thread counts across polls within one wait phase.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 253 - 288,
Reduce redundant GitHub calls in Evidence by short-circuiting PullRequestReviews
and ReviewThreads when classifyEvidence inputs already show a current-head check
run is pending or reviewing; retain the existing check-run and commit-status
fetches and preserve full review/thread retrieval when needed for
verified-vs-reviewed classification or details.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:265a6d4418a4c50fe5b161ee -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Current-head pending CodeRabbit check/status evidence now returns `reviewing` before fetching pull-request reviews and threads, while structured skip evidence retains precedence. A call-count regression proves the unnecessary API calls are skipped; focused CodeRabbit evidence tests passed.
