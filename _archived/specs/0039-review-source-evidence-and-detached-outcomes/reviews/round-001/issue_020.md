---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/daemon/engine.go
line: 281
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6Oa,comment:PRRC_kwDOS0qyts7aoLRQ
review_hash: 5c6e593a474befcd43c11cb5776749b6933b578266da12788dd2d983732c3bfb
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:09Z"
---

# Issue 020: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Extract the attempt/retry identity string into one helper.**

The exact `attempt %d` / `attempt %d retry %d` formatting is now built in three places; a future change to the wording has to be made three times and the progress line can silently drift from the event summaries.

<details>
<summary>♻️ Suggested helper</summary>

```go
func (req verificationAttemptRequest) identity() string {
	if req.Retry > 0 {
		return fmt.Sprintf("attempt %d retry %d", req.Attempt, req.Retry)
	}
	return fmt.Sprintf("attempt %d", req.Attempt)
}
```

Then use `req.identity()` in `summary`, `runVerificationAttempt`, and `publishVerdict`.
</details>





Also applies to: 329-332, 396-399

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/engine.go` around lines 278 - 281, Extract the duplicated
attempt/retry formatting into an identity method on verificationAttemptRequest,
preserving the existing strings and retry condition. Replace the local
formatting logic in summary, runVerificationAttempt, and publishVerdict with
req.identity() so all event and progress messages share one implementation.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:11b038ff6300f9fa2eca23a4 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added `verificationAttemptRequest.identity()` and reused it in summaries, failure/pass progress, and verdict publication so retry identity formatting cannot drift. Focused daemon verification tests and the full daemon package tests passed.
