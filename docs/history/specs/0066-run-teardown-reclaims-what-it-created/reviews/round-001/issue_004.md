---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/cli/cli_test.go
line: 777
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9Xy,comment:PRRC_kwDOS0qyts7dnSbR
review_hash: 81a942661cb30d28b25e4d198a188ee03e982847e18e7faf705af315641d2260
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Add coverage for the unproven-termination outcome.**

`reconcileProcessControllerStub.TerminateTreeAndWait` always returns `Proven: true`, so no test exercises the refusal at `internal/cli/reconcile.go` lines 649-657. That branch implements the stated guarantee that unproven termination is not reported as successful. Add a stub mode that returns an outcome with `Proven: false` and a reason, then assert that the candidate action becomes `"preserve"`, `DebrisSummary.ProcessesApplied` stays 0, and the command exits with the operational-failure code.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 677 - 777, Extend
reconcileProcessControllerStub with a mode that makes TerminateTreeAndWait
return Proven: false and a non-empty reason, then add coverage in
TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches (or a focused
neighboring test). Assert the process candidate action is "preserve",
DebrisSummary.ProcessesApplied remains 0, and the reconcile command returns the
operational-failure exit code.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:a7fe81a3ee0bf2ce530ae746 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Added `TestRunReconcileApplyPreservesProcessWhenTerminationIsUnproven` for the missing `--apply` refusal boundary.
  - The regression proves an unproven termination outcome keeps the candidate at `preserve`, leaves `ProcessesApplied` at zero, increments `OperationalFailures`, returns exit 1, and retains the process tree.
  - Focused evidence: `rtk go test ./internal/cli -run 'TestRunReconcileApplyPreservesProcessWhenTerminationIsUnproven'` passed; the complete affected package suites passed (1,247 tests).
  - The Daemon owns authoritative `make verify` after this Agent turn.
