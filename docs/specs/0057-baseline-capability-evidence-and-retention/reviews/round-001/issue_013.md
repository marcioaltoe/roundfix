---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: invalid
terminal_reason: "This loader assigns both sourcePath and carrierPath from the same unique carrier.Path, so the proposed tie-breaker cannot distinguish equal primary keys."
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/preservation.go
line: 804
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymk,comment:PRRC_kwDOS0qyts7cjgFA
review_hash: 64e4e6737012c501d26d7445d5a3061f83bf238b398a9ef4b94fe83557aaac80
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 013: _ Functional Correctness_ _ Trivial_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🔵 Trivial_ | _⚡ Quick win_

**Use a stable sort for preservation sources.**

`sort.Slice` is not stable. The comparison uses only `sourcePath`. Today each stale carrier contributes one entry, so the result is deterministic. If two carriers ever resolve to the same `sourcePath`, the order becomes unspecified, and the derived Source Baseline and Plan Digest become nondeterministic. `loadRootPreservationSources` already guards this by comparing `carrierPath` as a tiebreaker.





<details>
<summary>♻️ Proposed change</summary>

```diff
 	sort.Slice(sources, func(i, j int) bool {
-		return sources[i].sourcePath < sources[j].sourcePath
+		if sources[i].sourcePath != sources[j].sourcePath {
+			return sources[i].sourcePath < sources[j].sourcePath
+		}
+		return sources[i].carrierPath < sources[j].carrierPath
 	})
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	sort.Slice(sources, func(i, j int) bool {
		if sources[i].sourcePath != sources[j].sourcePath {
			return sources[i].sourcePath < sources[j].sourcePath
		}
		return sources[i].carrierPath < sources[j].carrierPath
	})
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/preservation.go` around lines 802 - 804, Update the sorting
in loadRootPreservationSources to use a stable sort for sources, preserving
input order when sourcePath values are equal; alternatively, add the existing
carrierPath-based tiebreaker used by the loader. Keep sourcePath as the primary
ordering key so derived Source Baseline and Plan Digest remain deterministic.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:fcaab28c45b391838e202b8b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: In `loadStaleManagedCarrierSources`, both fields are set to `carrier.Path`; equal `sourcePath` values therefore have equal proposed `carrierPath` tie-breakers. The finding's nondeterministic two-carrier scenario is not representable by the inspected carrier inventory.
