---
source: coderabbit
pr: "39"
round: 1
round_created_at: "2026-07-27T21:22:30Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-run-worktree-reconciliation
head_sha: 44fa422cea404a2d8c47e4b8011f065c4c0481ba
file: internal/worktree/worktree_test.go
line: 1562
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMo,comment:PRRC_kwDOS0qyts7aNkLP
review_hash: 7319fd8371fbdb3856b92a8f7a32a49d15265f5e7966a5c0d8839958cdbd22d7
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Use `slices.Equal` instead of a hand-rolled comparison.**

`slicesEqual` duplicates the standard library helper.



<details>
<summary>♻️ Proposed replacement</summary>

```diff
-func slicesEqual(left, right []string) bool {
-	if len(left) != len(right) {
-		return false
-	}
-	for index := range left {
-		if left[index] != right[index] {
-			return false
-		}
-	}
-	return true
-}
```

Import `slices` and call `slices.Equal(...)` at the three comparison sites (lines 1293, 1347-1348, 1333).
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree_test.go` around lines 1552 - 1562, Remove the
hand-rolled slicesEqual helper and import the standard library slices package.
Replace all three slicesEqual call sites in the affected tests with
slices.Equal, preserving the existing comparison arguments and behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:5da1abc0035ea3d6f15e79f1 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The test-only `slicesEqual` helper duplicated the Go standard
  library. Its three call sites now use `slices.Equal`, and the redundant
  helper was removed without changing assertions.

## Verification

- `rtk proxy env GOCACHE=/private/tmp/roundfix-pr39-batch001-gocache go test ./internal/worktree ./internal/cli ./internal/store -run 'Test(CountRetainedTerminalRunsBatchesGitInspectionByRepository|InspectTerminalRunUnsafePath|RunRunsList.*Retained.*Worktree|RunRunsListReportsRetainedWorktreeInspectionFailure|PruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|ReconcileIntegration)' -count=1`
  — passed.
