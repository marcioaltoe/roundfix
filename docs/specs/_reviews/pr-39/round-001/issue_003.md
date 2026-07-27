---
source: coderabbit
pr: "39"
round: 1
round_created_at: "2026-07-27T21:22:30Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-run-worktree-reconciliation
head_sha: 44fa422cea404a2d8c47e4b8011f065c4c0481ba
file: internal/store/store.go
line: 540
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMf,comment:PRRC_kwDOS0qyts7aNkLD
review_hash: 2e81ee1680c458414f68d809f15448dbfba10523fa79e90fea4eabcd8b3946bf
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Trivial_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🔵 Trivial_ | _⚡ Quick win_

**Run Branch naming is duplicated across packages.**

`"roundfix/run-" + req.RunID` re-implements `worktree.BranchName` (which uses `runBranchPrefix`). Since `internal/worktree` imports `internal/store`, the constant cannot be imported back, so a change to the prefix in the worktree package would make every reconciliation request fail this check with no compile-time signal. Extract the prefix into a small shared package (or a `store`-owned exported constant that `worktree.BranchName` consumes) so both sides derive the name from one source.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/store.go` around lines 532 - 540, The Run Branch prefix is
duplicated in the reconciliation check and worktree naming logic. Extract the
prefix into a dependency-safe shared symbol, or make it a store-owned exported
constant consumed by worktree.BranchName, then update the reconciliation
validation around expectedRunBranch to derive the expected name through that
single source instead of concatenating "roundfix/run-" directly.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:d3a51e7e868a24c683d2c104 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The Run Branch prefix was independently encoded in `internal/store`
  and `internal/worktree`, so either package could drift without a compile-time
  signal. `store.RunBranchPrefix` is now the single source consumed by both
  reconciliation validation and worktree branch naming.

## Verification

- `rtk proxy env GOCACHE=/private/tmp/roundfix-pr39-batch001-gocache go test ./internal/worktree ./internal/cli ./internal/store -run 'Test(CountRetainedTerminalRunsBatchesGitInspectionByRepository|InspectTerminalRunUnsafePath|RunRunsList.*Retained.*Worktree|RunRunsListReportsRetainedWorktreeInspectionFailure|PruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|ReconcileIntegration)' -count=1`
  — passed.
