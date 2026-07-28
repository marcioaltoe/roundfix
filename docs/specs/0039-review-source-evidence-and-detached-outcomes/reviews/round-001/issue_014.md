---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/watch/watch_test.go
line: 831
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIp5,comment:PRRC_kwDOS0qyts7aUVEN
review_hash: a8eb8773fe097069b2597188679a1aba4d111ee2319d763054817cf9b364a154
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 014: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**`stops.calls != 5` asserts an internal poll count.**

The magic 5 encodes exactly how many times the orchestrator happens to observe stop requests across the pre-fetch wait, quiet period, fetch, and retry sleep. Any benign reordering of those observation points breaks this test without any behavior change. Assert the observable contract instead — that a stop was observed after the retry sleep and before the next evidence call — and keep the existing `len(source.requests) != 2` bound, which already pins the meaningful behavior.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/watch/watch_test.go` around lines 830 - 831, Update the assertion in
the stop/retry test to remove the exact stops.calls == 5 requirement, since it
depends on internal polling order. Assert only that a stop was observed after
the retry sleep and before the next evidence call, while preserving the existing
len(source.requests) == 2 check.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:9432742b368135f58be5509b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The retry-stop test now asserts the observable Stop Request and exactly two evidence attempts instead of coupling to an internal stop-poll count. `go test ./internal/watch` passed.
