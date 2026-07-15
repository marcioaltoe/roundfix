---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/rounds/rounds.go
line: 486
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95z,comment:PRRC_kwDOS0qyts7V5taq
review_hash: 9055d57b6fe193492fd94a1ce0a9b78382b9437480ca156f3144b67a059888f2
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 014: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Clear stale reasons when the terminal status changes.**

A blank reason currently preserves metadata from the previous status. For example, `failed → duplicated` retains the old failure reason. Preserve an empty argument only when updating the same status.

<details>
<summary>Proposed fix</summary>

```diff
-	nextReason := terminalReason
-	if strings.TrimSpace(nextReason) == "" {
+	nextReason := strings.TrimSpace(terminalReason)
+	if nextReason == "" && status == frontmatter.Status {
 		nextReason = frontmatter.TerminalReason
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	nextReason := strings.TrimSpace(terminalReason)
	if nextReason == "" && status == frontmatter.Status {
		nextReason = frontmatter.TerminalReason
	}
	if status == StatusResolved {
		nextReason = ""
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/rounds/rounds.go` around lines 480 - 486, Update the nextReason
handling in the status transition logic around terminalReason and
frontmatter.TerminalReason so a blank terminalReason clears stale metadata when
status changes, including failed → duplicated. Only reuse
frontmatter.TerminalReason when the incoming status matches the existing
terminal status; continue clearing the reason for StatusResolved.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:bedcfdf5ed895be1fe01e03c -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `SetIssueStatus` now preserves a blank terminal reason only when the status is unchanged, clearing stale metadata on status transitions; `make verify` passed.
