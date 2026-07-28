---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/reviewsource/reviewsource_test.go
line: 43
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpw,comment:PRRC_kwDOS0qyts7aUVEC
review_hash: 386f0e0c24a7b2e20059b1ebc6f543fa348d7283d8365e71f4feeb2cdacb8df6
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 010: _ Functional Correctness_ _ Trivial_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🔵 Trivial_ | _⚡ Quick win_

**The rune-boundary backtrack in `BoundEvidenceDetail` is not covered.**

`strings.Repeat("é", MaxEvidenceDetailLength)` is 4096 bytes of 2-byte runes, so byte offset 2048 already lands on a rune start and the `!utf8.RuneStart` loop never executes. Add a case with a 3-byte rune (for example `strings.Repeat("→", ...)`) so the cut actually falls mid-rune and the result is asserted to be valid UTF-8.



<details>
<summary>💚 Suggested addition</summary>

```go
func TestBoundEvidenceDetailCutsOnRuneBoundary(t *testing.T) {
	bounded := BoundEvidenceDetail(strings.Repeat("→", MaxEvidenceDetailLength))
	if !utf8.ValidString(bounded) {
		t.Fatalf("bounded detail is not valid UTF-8: %q", bounded)
	}
	if len(bounded) > MaxEvidenceDetailLength+len("…") {
		t.Fatalf("bounded detail length = %d", len(bounded))
	}
}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/reviewsource_test.go` around lines 34 - 43, The tests
for BoundEvidenceDetail do not exercise truncation in the middle of a multi-byte
rune. Add a focused test using repeated 3-byte characters such as “→”, assert
the bounded result is valid UTF-8, and retain the maximum byte-length assertion
including the ellipsis.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:89f85b15c952db483a672855 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added a three-byte rune boundary case that verifies `BoundEvidenceDetail` returns valid UTF-8 and respects the byte limit. Focused reviewsource rune-boundary tests passed.
