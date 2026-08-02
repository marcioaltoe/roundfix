---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: invalid
terminal_reason: "The evaluator structurally appends exactly one outcome per capability in the same loop and returns immediately on error, so position mapping cannot drift under the current contract."
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/profile_alignment.go
line: 551
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymr,comment:PRRC_kwDOS0qyts7cjgFI
review_hash: 3c889d999748b04502a093798c18e058342ddb81af67040c28a6ad2fa86ddb97
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 017: _ Data Integrity & Integration_ _ Trivial_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🔵 Trivial_ | _⚡ Quick win_

**Key the probe lookup by capability ID instead of slice position.**

Line 547 reads `capabilities[index].Probe` using the index of `outcomes`. This depends on `evaluateRepositoryCapabilities` returning exactly one outcome per input capability, in the same order. That contract is implicit and unenforced. If that function ever skips or reorders a capability, each divergence receives the probe of a different capability, and no test would fail loudly because both slices stay the same length in the common case.

Resolve the probe by `outcome.ID`.





<details>
<summary>♻️ Proposed change</summary>

```diff
 	applyUniversalCapabilityRemediation(outcomes, remediationProfileID)
+	probes := make(map[string]map[string]any, len(capabilities))
+	for _, capability := range capabilities {
+		probes[capability.ID] = capability.Probe
+	}
 	divergences := make([]ProfileDivergence, 0)
-	for index, outcome := range outcomes {
+	for _, outcome := range outcomes {
 		if outcome.Status == CapabilitySatisfied {
 			continue
 		}
 		divergences = append(divergences, ProfileDivergence{
 			Code:                 outcome.Diagnostic.Code,
 			ID:                   outcome.ID,
 			Requirement:          outcome.Requirement,
 			Blocking:             outcome.Blocking,
 			Message:              outcome.Diagnostic.Message,
 			NextAction:           outcome.Diagnostic.NextAction,
-			Probe:                capabilities[index].Probe,
+			Probe:                probes[outcome.ID],
 			Evidence:             outcome.Evidence,
 			CapabilityResolution: stackCapabilityResolution(profile, outcome, catalog),
 		})
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	probes := make(map[string]map[string]any, len(capabilities))
	for _, capability := range capabilities {
		probes[capability.ID] = capability.Probe
	}

	divergences := make([]ProfileDivergence, 0)
	for _, outcome := range outcomes {
		if outcome.Status == CapabilitySatisfied {
			continue
		}
		divergences = append(divergences, ProfileDivergence{
			Code:                 outcome.Diagnostic.Code,
			ID:                   outcome.ID,
			Requirement:          outcome.Requirement,
			Blocking:             outcome.Blocking,
			Message:              outcome.Diagnostic.Message,
			NextAction:           outcome.Diagnostic.NextAction,
			Probe:                probes[outcome.ID],
			Evidence:             outcome.Evidence,
			CapabilityResolution: stackCapabilityResolution(profile, outcome, catalog),
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

In `@internal/baseline/profile_alignment.go` around lines 535 - 551, Update the
divergence construction around ProfileDivergence to resolve the probe by
outcome.ID from the capabilities collection, rather than using
capabilities[index].Probe; preserve the existing outcome iteration and
divergence fields while ensuring each divergence receives the probe belonging to
its own capability ID.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:ea1db0282122128e369e71c2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `evaluateRepositoryCapabilities` iterates the input slice in order and appends one outcome for every successful iteration; it has no skip or reorder path. Replacing this enforced positional invariant with a lookup map would address only a hypothetical future rewrite.
