---
source: coderabbit
pr: "39"
round: 1
round_created_at: "2026-07-27T21:22:30Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-run-worktree-reconciliation
head_sha: 44fa422cea404a2d8c47e4b8011f065c4c0481ba
file: internal/worktree/worktree.go
line: 222
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMu,comment:PRRC_kwDOS0qyts7aNkLV
review_hash: 86612da12c2d2aca59fd975bc0997a6c95952b8c6e38c583739dad0349528181
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 006: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Misleading reason for an unregistered worktree directory.**

`recordedWorktree` returns `unsafe=true` both for genuinely unsafe paths (symlink, unclean, non-absolute) and for a directory that exists but is simply not registered in the Git root (line 1116). In the latter case operators see "recorded Run Worktree path is unsafe" instead of the dedicated `reconciliationReasonWorktreeUnregistered`. Consider distinguishing the "exists but unregistered" outcome so the surfaced reason matches the actual evidence.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree.go` around lines 218 - 222, Update the worktree
reconciliation flow around recordedWorktree so an existing directory that is not
registered in the Git root produces reconciliationReasonWorktreeUnregistered,
while genuinely unsafe paths continue using reconciliationReasonWorktreePath.
Distinguish these outcomes using the recordedWorktree result and preserve the
existing return behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:891a974c07bd82eec84d395d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `recordedWorktree` collapsed an existing unregistered directory into
  the same boolean used for unsafe paths. It now returns an explicit state, so
  existing unregistered directories produce
  `reconciliationReasonWorktreeUnregistered` while malformed, symlinked, or
  non-directory paths retain the unsafe-path reason.

## Verification

- `rtk proxy env GOCACHE=/private/tmp/roundfix-pr39-batch001-gocache go test ./internal/worktree ./internal/cli ./internal/store -run 'Test(CountRetainedTerminalRunsBatchesGitInspectionByRepository|InspectTerminalRunUnsafePath|RunRunsList.*Retained.*Worktree|RunRunsListReportsRetainedWorktreeInspectionFailure|PruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|ReconcileIntegration)' -count=1`
  — passed.
