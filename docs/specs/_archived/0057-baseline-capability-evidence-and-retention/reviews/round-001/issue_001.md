---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/apply.go
line: 285
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0YmO,comment:PRRC_kwDOS0qyts7cjgEm
review_hash: 1c66c6fa96b59a44c51bf44850681608da46d83eeb1749dbd159c751911f02a1
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:29Z"
---

# Issue 001: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Include `ApprovedPostimages` in the completion claim.**

The message states "Baseline update complete" when only `SemanticRetention` and `Idempotence` are verified. `ApprovedPostimages` can still be `not run` when `document.Postimages` is empty, because `verifiedPostimagesMatch` returns `false` for an empty approved set. The result then claims completion without postimage evidence. Require all three axes for the completion language.





<details>
<summary>🐛 Proposed fix for the completion condition</summary>

```diff
-	if matrix.SemanticRetention == EvidenceStatusVerified &&
-		matrix.Idempotence == EvidenceStatusVerified {
+	if matrix.ApprovedPostimages == EvidenceStatusVerified &&
+		matrix.SemanticRetention == EvidenceStatusVerified &&
+		matrix.Idempotence == EvidenceStatusVerified {
 		message = "Baseline update complete; approved Baseline Plan is already applied"
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	if matrix.ApprovedPostimages == EvidenceStatusVerified &&
		matrix.SemanticRetention == EvidenceStatusVerified &&
		matrix.Idempotence == EvidenceStatusVerified {
		message = "Baseline update complete; approved Baseline Plan is already applied"
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/apply.go` around lines 282 - 285, Update the completion
condition in the baseline update flow to also require matrix.ApprovedPostimages
to equal EvidenceStatusVerified before using the “Baseline update complete”
message; retain the existing SemanticRetention and Idempotence requirements.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:6da426a7d14e384fd44b2049 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Completion language now requires verified approved postimages, semantic retention, and idempotence. `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/baseline -count=1` passed.
