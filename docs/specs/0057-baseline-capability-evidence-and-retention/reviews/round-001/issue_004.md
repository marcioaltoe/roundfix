---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/classification_test.go
line: 394
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0YmX,comment:PRRC_kwDOS0qyts7cjgEy
review_hash: 7a3c4799a325e4328bdbec70103946b70f9fa97728543a933b54ec66a81797f6
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:29Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Assert that the two selected guides use distinct paths.**

The test takes the first two artifacts with `Kind == "guide"` and then uses `guides[0].Path` and `guides[1].Path` as separate map keys in `want`. If two guide artifacts ever resolve to the same path, the `want` map keeps only one entry, and the test asserts `carrierStaleManaged` for a file that also holds the current block. The failure would be confusing and unrelated to the behavior under test.





<details>
<summary>♻️ Proposed change</summary>

```diff
 		if len(guides) != 2 {
 			t.Fatalf("active guide artifacts = %d, want at least two", len(guides))
 		}
+		if guides[0].Path == guides[1].Path {
+			t.Fatalf("selected guides share path %q, want distinct carriers", guides[0].Path)
+		}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		guides := make([]plannedArtifact, 0, 2)
		for _, artifact := range artifacts {
			if artifact.Kind == "guide" {
				guides = append(guides, artifact)
				if len(guides) == 2 {
					break
				}
			}
		}
		if len(guides) != 2 {
			t.Fatalf("active guide artifacts = %d, want at least two", len(guides))
		}
		if guides[0].Path == guides[1].Path {
			t.Fatalf("selected guides share path %q, want distinct carriers", guides[0].Path)
		}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/classification_test.go` around lines 383 - 394, Update the
guide-selection assertions in the test to verify that guides[0].Path and
guides[1].Path are distinct before using them as separate want-map keys. Fail
the test with a clear message if the paths match, while preserving the existing
requirement that two guide artifacts are selected.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:b26cf7501625548ed07c3057 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The characterization now fails with a direct diagnostic if the two selected guide artifacts share one carrier path. The full Baseline package test passed.
