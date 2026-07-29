---
source: coderabbit
pr: "42"
round: 1
round_created_at: "2026-07-29T02:33:36Z"
status: invalid
terminal_reason: "runtimeDisplayName only receives supported adapter runtime IDs from validated Agent selections or fixed lineage contracts, so non-ASCII identifiers cannot reach the byte slice."
head_repository: marcioaltoe/roundfix
head_branch: ma/claude-adapter-standardization
head_sha: 7155ba4d2ef353257a1bacf697027202d4750492
file: internal/agent/acpx_runner.go
line: 855
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UmtXj,comment:PRRC_kwDOS0qyts7ayCz7
review_hash: 25d230e0da7c026d81a37e68cff91c3209905bff5b47404c4c7eccadb7c7796b
duplicate_of: ""
source_review_id: "4803488138"
source_review_submitted_at: "2026-07-29T02:32:46Z"
---

# Issue 002: _ Functional Correctness_ _ Trivial_ _ Low value_

## Review Comment

_🎯 Functional Correctness_ | _🔵 Trivial_ | _💤 Low value_

**`runtimeDisplayName` slices bytes, not runes.**

`runtimeID[:1]` splits the first byte, so a non-ASCII runtime identifier would produce mangled output in every adapter error message. Use rune-aware capitalization.





<details>
<summary>♻️ Proposed refactor</summary>

```diff
 func runtimeDisplayName(runtimeID string) string {
 	if runtimeID == "" {
 		return "Adapter"
 	}
-	return strings.ToUpper(runtimeID[:1]) + runtimeID[1:]
+	first, size := utf8.DecodeRuneInString(runtimeID)
+	return string(unicode.ToUpper(first)) + runtimeID[size:]
 }
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func runtimeDisplayName(runtimeID string) string {
	if runtimeID == "" {
		return "Adapter"
	}
	first, size := utf8.DecodeRuneInString(runtimeID)
	return string(unicode.ToUpper(first)) + runtimeID[size:]
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 850 - 855, Update
runtimeDisplayName to capitalize the first Unicode rune rather than slicing
runtimeID by bytes, while preserving the empty-string fallback and the remainder
of the identifier unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:5e3a44e8659a3738c89a43ac -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes:
  - Agent Selection validation accepts only `codex`, `claude`, and `opencode`; adapter lineage errors use those fixed runtime contract IDs.
  - No non-ASCII runtime identifier can reach `runtimeDisplayName`, so rune-aware capitalization would add imports and code for an impossible product input.
