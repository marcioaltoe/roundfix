---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/profile_alignment.go
line: 986
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymt,comment:PRRC_kwDOS0qyts7cjgFL
review_hash: 57ee6fe2944805d1b32aa93f43752b34b84302a76baf81f84ed99ea1c7a61eff
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 018: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Report hop-limit exhaustion with a reason distinct from `link-cycle`.**

The `visited` set at lines 969-973 already detects every true cycle before the hop counter can matter. Therefore the guard at line 983 fires only for an acyclic symlink chain longer than `maxExecutableLinkHops`. That case is reported as `link-cycle`.

The reason string reaches the operator. `collectExecutableEvidence` copies it into `CapabilityEvidence.Detail`, and `rejectedExecutableNextAction` renders it as "Repair the inspected candidate X (link-cycle)". An operator then looks for a cycle that does not exist.

Add a distinct reason for the hop limit.





<details>
<summary>🐛 Proposed fix</summary>

```diff
 	executableProbeReasonNotFound      = "not-found"
 	executableProbeReasonBrokenLink    = "broken-link"
 	executableProbeReasonLinkCycle     = "link-cycle"
+	executableProbeReasonLinkDepth     = "link-depth-exceeded"
 	executableProbeReasonNotExecutable = "not-executable"
```

```diff
 		if result.HopCount >= maxExecutableLinkHops {
-			result.Reason = executableProbeReasonLinkCycle
+			result.Reason = executableProbeReasonLinkDepth
 			return result, true
 		}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/profile_alignment.go` around lines 966 - 986, Introduce a
distinct executable-probe reason for hop-limit exhaustion and use it in the
maxExecutableLinkHops guard within the symlink-resolution loop, while retaining
executableProbeReasonLinkCycle for visited-set cycle detection. Ensure the new
reason propagates through the existing evidence and operator messaging paths.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:88f4b89cf5af9a44cc43125e -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Hop-limit exhaustion now reports `link-depth-exceeded`, while visited-path detection retains `link-cycle`. A regression with an acyclic chain beyond the bound and the full Baseline package test passed.
