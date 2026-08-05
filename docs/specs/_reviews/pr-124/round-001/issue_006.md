---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/reviewsource/coderabbit/coderabbit_test.go
line: 1249
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wua0K,comment:PRRC_kwDOS0qyts7d35tg
review_hash: 06cde0b4ffa84c98c2851ac4f4a698f343cd4d688b690fd7d8a18d01ed0c6f1e
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 006: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**The fixture body contains literal `\n` text, not newlines.**

Line 1223 is inside a raw string literal delimited by backticks, so `\\n` is a backslash followed by `n`. JSON decoding turns `\\` into one backslash and keeps `n`, so `comments[0].Body` holds the two characters `\` and `n` where a line break is intended.

The current assertions still pass because they only use `strings.Contains`. The defect becomes visible as soon as this recorded body is fed to `rateLimitCommentReason`, which splits on `"\n"` and would find no `## ` heading. Use a single backslash so the JSON escape produces a real newline.



<details>
<summary>💚 Proposed fixture fix</summary>

```diff
 	fixture := `[{
 		"id": 5182301262,
-		"body": "<!-- This is an auto-generated comment: rate limited by coderabbit.ai -->\\n\\n> ## Review limit reached",
+		"body": "<!-- This is an auto-generated comment: rate limited by coderabbit.ai -->\n\n> ## Review limit reached",
 		"user": {"login": "coderabbitai[bot]"}
 	}]`
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func TestIssueCommentsMapGitHubRateLimitCommentJSON(t *testing.T) {
	fixture := `[{
		"id": 5182301262,
		"body": "<!-- This is an auto-generated comment: rate limited by coderabbit.ai -->\n\n> ## Review limit reached",
		"user": {"login": "coderabbitai[bot]"}
	}]`
	var calls [][]string
	withRunGH(t, func(_ context.Context, args ...string) ([]byte, error) {
		calls = append(calls, append([]string(nil), args...))
		return []byte(fixture), nil
	})

	comments, err := (GHClient{}).IssueComments(context.Background(), "marcioaltoe/roundfix", "107")
	if err != nil {
		t.Fatalf("map issue comments: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("issue comment count = %d, want 1", len(comments))
	}
	if comments[0].DatabaseID != 5182301262 || comments[0].Author != coderabbitBotLogin ||
		!strings.Contains(comments[0].Body, "Review limit reached") {
		t.Fatalf("issue comment = %#v, want recorded CodeRabbit refusal", comments[0])
	}
	if len(calls) != 1 {
		t.Fatalf("gh calls = %#v, want one issue-comment request", calls)
	}
	assertStringSlicesEqual(t, calls[0], []string{
		"api", "--paginate", "repos/marcioaltoe/roundfix/issues/107/comments",
	})
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit_test.go` around lines 1220 -
1249, Update the raw JSON fixture in
TestIssueCommentsMapGitHubRateLimitCommentJSON to use a single backslash in the
\n escape sequence, so JSON decoding produces actual newlines in
comments[0].Body. Leave the existing assertions and request setup unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:038f776857c56fe6d92a2cc5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The raw JSON fixture now contains JSON `\n` escapes rather than
  escaped backslashes, and the test asserts that decoding produces actual
  newlines before the refusal heading.
- Focused evidence: the new decoded-newline assertion failed against the old
  fixture, then `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/reviewsource/coderabbit -count=1` passed after the fixture fix.
- Daemon Verification: `make verify` not run; Daemon-owned.
