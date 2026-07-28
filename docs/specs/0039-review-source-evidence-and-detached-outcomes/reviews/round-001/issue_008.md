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
line: 671
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpq,comment:PRRC_kwDOS0qyts7aUVD8
review_hash: 5672046f65dba3687cc19be0c64429135b4564e0ecf8d207dcbfc7dad889b26f
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 008: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Extract the repeated verified-vs-reviewed selection.**

The same `state := EvidenceVerified; if unresolvedThreads > 0 { state = EvidenceReviewed }` block appears three times across the check-run, commit-status, and review-approval branches. A single helper keeps the unresolved-thread rule in one place so a future change to that rule cannot be applied to only two of the three signal kinds.



<details>
<summary>♻️ Proposed helper</summary>

```diff
+func settledState(unresolvedThreads int) reviewsource.EvidenceState {
+	if unresolvedThreads > 0 {
+		return reviewsource.EvidenceReviewed
+	}
+	return reviewsource.EvidenceVerified
+}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 644 - 671,
Extract the repeated state-selection logic into a helper near the existing
evidence helpers, returning EvidenceReviewed when unresolvedThreads is greater
than zero and EvidenceVerified otherwise. Replace the duplicate blocks in the
check-run, commit-status, and review-approval branches with calls to this
helper, preserving each branch’s existing evidence construction.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:65d8cdbdd9f85f3bafbf2f74 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Extracted `settledEvidenceState` and reused it for successful check, status, and approval evidence, removing repeated verified/reviewed branching. Focused CodeRabbit classification tests passed.
