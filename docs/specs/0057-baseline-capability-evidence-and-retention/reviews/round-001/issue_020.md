---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: invalid
terminal_reason: "The private sorter only receives classified divergences whose Group already matches Requirement, while the exported renderer independently uses Group-first fallback."
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/profile_alignment.go
line: 1815
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymw,comment:PRRC_kwDOS0qyts7cjgFO
review_hash: c82a487a559e174f16e1036703fa64abe90a5b1c62a53d3ecd0980d4b7bada5e
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 020: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Sort by the classified `Group` and fall back to `Requirement`.**

`sortProfileDivergences` derives the group from `Requirement`, but `RenderProfileDivergences` prefers the already-assigned `Group` field and only falls back to `Requirement`. The two functions therefore use different sources of truth for the same concept. Inside this package `classifyProfileDivergences` always keeps them in agreement, so no defect exists today. `RenderProfileDivergences` is exported, so an external caller can construct divergences whose `Group` and `Requirement` disagree; the rendered sections would then be out of order.

Use the same resolution rule in both functions.





<details>
<summary>♻️ Proposed change</summary>

```diff
+func resolvedDivergenceGroup(divergence ProfileDivergence) ProfileDivergenceGroup {
+	if divergence.Group != "" {
+		return divergence.Group
+	}
+	return profileDivergenceGroup(divergence.Requirement)
+}
+
 func sortProfileDivergences(divergences []ProfileDivergence) {
 	sort.Slice(divergences, func(i, j int) bool {
-		leftGroup := profileDivergenceGroup(divergences[i].Requirement)
-		rightGroup := profileDivergenceGroup(divergences[j].Requirement)
+		leftGroup := resolvedDivergenceGroup(divergences[i])
+		rightGroup := resolvedDivergenceGroup(divergences[j])
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func resolvedDivergenceGroup(divergence ProfileDivergence) ProfileDivergenceGroup {
	if divergence.Group != "" {
		return divergence.Group
	}
	return profileDivergenceGroup(divergence.Requirement)
}

func sortProfileDivergences(divergences []ProfileDivergence) {
	sort.Slice(divergences, func(i, j int) bool {
		leftGroup := resolvedDivergenceGroup(divergences[i])
		rightGroup := resolvedDivergenceGroup(divergences[j])
		if leftGroup != rightGroup {
			return profileDivergenceGroupRank(leftGroup) < profileDivergenceGroupRank(rightGroup)
		}
		if divergences[i].ID != divergences[j].ID {
			return divergences[i].ID < divergences[j].ID
		}
		return divergences[i].Code < divergences[j].Code
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

In `@internal/baseline/profile_alignment.go` around lines 1803 - 1815, Update
sortProfileDivergences to resolve each divergence’s group using the same
Group-first, Requirement-fallback rule as RenderProfileDivergences. Sort by the
resolved group rank, then retain the existing ID and Code tie-breakers so
externally supplied Group values determine rendered section order consistently.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:335f6dd984f8cf02060730d2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `classifyProfileDivergences` sets every internal Group immediately before sorting, and `RenderProfileDivergences` resolves Group first and Requirement second for exported callers. No external caller can invoke the private sorter, so the admitted mismatch cannot affect current output.
