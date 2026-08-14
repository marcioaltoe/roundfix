---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: internal/reviewsource/coderabbit/coderabbit.go
line: 205
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95s,comment:PRRC_kwDOS0qyts7V5tak
review_hash: 02b69d28f98c1670f5f8efcfd53b7c9eb678af68cab573bf52d3180ca4a37b66
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_013.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 013: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Use the standalone marker detector for idempotency.**

`strings.Contains` treats quoted or embedded marker text as an existing outcome, incorrectly suppressing the reply. Use the helper whose tested contract requires a standalone marker line.

<details>
<summary>Proposed fix</summary>

```diff
 		for _, comment := range thread.Comments {
-			if strings.Contains(comment.Body, marker) {
+			if HasRoundfixCommentMarker(marker, comment.Body) {
 				return reviewsource.IssueCommentResult{Skipped: true}, nil
 			}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		for _, comment := range thread.Comments {
			if HasRoundfixCommentMarker(marker, comment.Body) {
				return reviewsource.IssueCommentResult{Skipped: true}, nil
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 203 - 205,
Replace the strings.Contains check in the thread.Comments loop with the existing
standalone marker detector, preserving the Skipped result when that detector
confirms a standalone marker line and allowing replies when the marker is only
quoted or embedded.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c2f08a6a29775c74ab96a756 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
