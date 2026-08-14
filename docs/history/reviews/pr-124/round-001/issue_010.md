---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/reviewsource/coderabbit/coderabbit.go
line: 934
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wua0n,comment:PRRC_kwDOS0qyts7d35uH
review_hash: 960268226a20302668f586dd3323a8dcd4ce40f445fe551fa20627386acb90cb
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:40Z"
---

# Issue 010: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**`rateLimitCommentReason` returns the oldest matching comment and strips only one blockquote level.**

Two edge cases exist in the reason extraction:

1. The GitHub issue-comments endpoint returns comments oldest first. The loop returns the first matching comment. If CodeRabbit posted a rate-limit comment for an earlier head and another for the current head, the reason text comes from the earlier comment. The classification stays `skipped`, so only the reported reason can be stale. Iterating in reverse selects the most recent refusal.

2. Line 923 removes one `>` prefix. CodeRabbit renders the heading inside a GitHub alert block, which the recorded fixture shows as `> ## Review limit reached`. A nested quote (`> > ## ...`) would not match, and the function would return no reason. That makes `refusalReason` report `false`, and the green rate-limit status would resolve to `pending` instead of `skipped`. Trimming `>` repeatedly removes that dependency on the exact nesting depth.



<details>
<summary>🐛 Proposed fix for both edge cases</summary>

```diff
 func rateLimitCommentReason(comments []IssueComment) (string, bool) {
 	const marker = "auto-generated comment: rate limited by coderabbit.ai"
-	for _, comment := range comments {
+	for index := len(comments) - 1; index >= 0; index-- {
+		comment := comments[index]
 		if !isCodeRabbitAuthor(comment.Author) || !strings.Contains(normalized(comment.Body), marker) {
 			continue
 		}
 		for _, line := range strings.Split(comment.Body, "\n") {
 			line = strings.TrimSpace(line)
-			line = strings.TrimSpace(strings.TrimPrefix(line, ">"))
+			for strings.HasPrefix(line, ">") {
+				line = strings.TrimSpace(strings.TrimPrefix(line, ">"))
+			}
 			if !strings.HasPrefix(line, "## ") {
 				continue
 			}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func rateLimitCommentReason(comments []IssueComment) (string, bool) {
	const marker = "auto-generated comment: rate limited by coderabbit.ai"
	for index := len(comments) - 1; index >= 0; index-- {
		comment := comments[index]
		if !isCodeRabbitAuthor(comment.Author) || !strings.Contains(normalized(comment.Body), marker) {
			continue
		}
		for _, line := range strings.Split(comment.Body, "\n") {
			line = strings.TrimSpace(line)
			for strings.HasPrefix(line, ">") {
				line = strings.TrimSpace(strings.TrimPrefix(line, ">"))
			}
			if !strings.HasPrefix(line, "## ") {
				continue
			}
			reason := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			if reason != "" {
				return reviewsource.BoundEvidenceDetail(reason), true
			}
		}
	}
	return "", false
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 915 - 934,
Update rateLimitCommentReason to iterate matching comments from newest to oldest
so it returns the latest refusal reason. In the heading extraction loop,
repeatedly trim leading blockquote markers and surrounding whitespace before
checking for the “## ” heading, allowing nested quote levels to match while
preserving the existing non-empty reason and fallback behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8d3bc5573e580db6802ae993 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `rateLimitCommentReason` now scans issue comments newest-first and
  removes every leading blockquote marker before matching the heading. Focused
  cases prove that a newer refusal reason wins and that a nested blockquote
  still resolves `skipped` with its reason.
- Focused evidence: both new cases failed against the former first-match,
  single-prefix implementation, then `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/reviewsource/coderabbit -count=1` passed after the fix.
- Daemon Verification: `make verify` not run; Daemon-owned.
