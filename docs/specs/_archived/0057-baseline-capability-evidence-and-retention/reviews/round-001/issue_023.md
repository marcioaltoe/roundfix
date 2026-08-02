---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/cli/baseline_human_test.go
line: 907
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ym2,comment:PRRC_kwDOS0qyts7cjgFW
review_hash: a78b01f67c79196c11a0894423a839a6312f797c98db55fc03f65ba7ecaa9ef0
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 023: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Assert the mixed-severity precondition explicitly.**

The failure message says "want mixed unsatisfied divergences", but the check only requires two or more divergences. Two blocking divergences and zero advisory divergences would satisfy it. The test then silently stops covering the mixed case that its name and message describe.

Assert that at least one blocking and at least one non-blocking divergence exist.





<details>
<summary>♻️ Proposed change</summary>

```diff
-	if alignment.Ready || len(alignment.Divergences) < 2 {
+	blocking := 0
+	advisory := 0
+	for _, divergence := range alignment.Divergences {
+		if divergence.Blocking {
+			blocking++
+			continue
+		}
+		advisory++
+	}
+	if alignment.Ready || blocking == 0 || advisory == 0 {
 		t.Fatalf("divergence fixture alignment = %+v, want mixed unsatisfied divergences", alignment)
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	blocking := 0
	advisory := 0
	for _, divergence := range alignment.Divergences {
		if divergence.Blocking {
			blocking++
			continue
		}
		advisory++
	}
	if alignment.Ready || blocking == 0 || advisory == 0 {
		t.Fatalf("divergence fixture alignment = %+v, want mixed unsatisfied divergences", alignment)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_human_test.go` around lines 905 - 907, Update the
precondition around alignment in the divergence fixture test to explicitly
require at least one blocking and one non-blocking divergence, rather than only
checking the total count. Preserve the existing Ready check and failure context
while using the divergence severity fields to validate the mixed unsatisfied
case.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:5e6236b0871e8c804b2c0a33 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The remediation fixture now requires at least one blocking and one advisory divergence before exercising the prompt. `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
