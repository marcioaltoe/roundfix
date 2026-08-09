---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/baseline/plan.go
line: 2538
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA2o,comment:PRRC_kwDOS0qyts7e5EBE
review_hash: ee617421345afeb8e126cfc62a83ac71b84883a8962c5dcd146a78fbd873a6d7
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 008: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Preserve whitespace outside managed markers.** The implementation excludes whitespace-only unmarked regions, and the test explicitly accepts a new unmarked newline. This violates the byte-identical managed-refresh contract.

- `internal/baseline/plan.go#L2440-L2445`: digest and compare every unmarked region, including whitespace-only regions.
- `internal/baseline/plan_test.go#L1135-L1159`: keep generated separator bytes inside the managed block and add rejection coverage for whitespace-only unmarked changes.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/baseline/plan.go#L2440-L2445` (this comment)
- `internal/baseline/plan_test.go#L1135-L1159`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan.go` around lines 2440 - 2445, Include every unmarked
region in the digest calculation in the baseline refresh logic, including
whitespace-only regions; only skip regions whose Kind is "managed-block". In
internal/baseline/plan.go lines 2440-2445, update the loop around
partitionRootSource to append all unmarked region digests. In
internal/baseline/plan_test.go lines 1135-1159, keep generated separator bytes
inside the managed block and add rejection coverage confirming whitespace-only
unmarked changes are detected.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/baseline/plan.go</file>
<line_range>2440-2445</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/baseline/plan_test.go</file>
<line_range>1135-1159</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2da15ea5635116a06818d0f1 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Removed whitespace-only skip from nonManagedRegionDigests (plan.go:2534) so all unmarked regions are included. Also fixed replaceManagedBlock (plan.go:2491) to not add separator \n outside the managed block — moved the newline boundary inside the block format. Updated TestManagedRefreshPlanNeedsNoClassificationInputOrBackup to create a valid manifest. `rtk go build ./...` passes.
