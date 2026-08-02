---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/cli/baseline_human.go
line: 1003
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ym6,comment:PRRC_kwDOS0qyts7cjgFb
review_hash: 8b5e38548a21e5490970fadf33154f0380a93ac6c17635cfd5164ce928b2bdad
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 025: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Distinguish blocking divergences from advisory divergences in the remediation list.**

`renderBaselineProfileRemediation` prints every divergence under the heading "Repository remediation before Baseline re-run:". Advisory divergences such as `capability.firecrawl` and `capability.rtk` appear in the same flat list as blocking ones.

This removes the signal that the rest of this change adds. `RenderProfileDivergences` groups divergences into blocking, advisory, and informational sections, and marks advisory items with "This advisory does not block readiness or apply." The remediation list drops that distinction, so the operator cannot tell which entries must be fixed before the re-run succeeds.

Compare `uniqueProfileDivergenceActions` at line 1244, which the non-removable path calls with blocking divergences only.

Mark or separate the advisory entries.





<details>
<summary>🔧 Proposed change</summary>

```diff
 	fmt.Fprintln(output, "\nRepository remediation before Baseline re-run:")
 	for _, divergence := range divergences {
 		action := strings.TrimSpace(divergence.NextAction)
 		if action == "" {
 			action = "review the reported evidence and remediate this divergence"
 		}
-		fmt.Fprintf(output, "- %s: %s\n", divergence.ID, action)
+		if divergence.Blocking {
+			fmt.Fprintf(output, "- [blocking] %s: %s\n", divergence.ID, action)
+			continue
+		}
+		fmt.Fprintf(output, "- [advisory] %s: %s\n", divergence.ID, action)
 	}
```

The test at `internal/cli/baseline_human_test.go` lines 933-938 asserts the exact `"- " + divergence.ID + ": " + divergence.NextAction` prefix and needs a matching update.
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	fmt.Fprintln(output, "\nRepository remediation before Baseline re-run:")
	for _, divergence := range divergences {
		action := strings.TrimSpace(divergence.NextAction)
		if action == "" {
			action = "review the reported evidence and remediate this divergence"
		}
		if divergence.Blocking {
			fmt.Fprintf(output, "- [blocking] %s: %s\n", divergence.ID, action)
			continue
		}
		fmt.Fprintf(output, "- [advisory] %s: %s\n", divergence.ID, action)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_human.go` around lines 996 - 1003, Update
renderBaselineProfileRemediation to distinguish advisory divergences from
blocking divergences in its remediation output, using each divergence’s existing
classification and preserving the blocking-only behavior used by
uniqueProfileDivergenceActions. Mark or separate advisory entries so they are
clearly non-blocking, and update the exact-prefix assertion in the related test
to match the new output.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3f11d58108dbd2151f2c250a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Repository remediation entries now retain their blocking, advisory, or informational classification, with a fallback derived from requirement strength for unclassified values. The mixed-divergence prompt regression asserts the exact labels. `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
