---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/cli/reconcile.go
line: 672
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9YB,comment:PRRC_kwDOS0qyts7dnSbh
review_hash: 123a1d9cd9cf0537ffc3a7306efc1cc2e0a61697abc41041bbd36edd3328ac75
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 009: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Terminate process trees before releasing Run Branches.**

`applyReconcileReport` releases Run Worktrees and Run Branches first, then terminates owned process trees. A process that still holds a handle inside a Run Worktree can make `git worktree remove` fail even when the worktree is clean, because `ApplyRunBranchCandidate` only detects dirty content, not open handles. The failure is recorded as an operational failure and needs a second `--apply` run. Move the process-termination loop above the run-branch loop so the reclaimed process tree releases its handles first.



<details>
<summary>🔧 Proposed reordering</summary>

```diff
 	applyReconcileWorktrees(ctx, homeDir, opts, report)
+	applyReconcileProcesses(ctx, report)
 	for index := range report.RunBranchCandidates {
```

Move lines 641-672 into `applyReconcileProcesses` and call it before the run-branch loop.
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/reconcile.go` around lines 626 - 672, In applyReconcileReport,
terminate owned process trees before applying Run Branch candidates so their
worktree handles are released first. Move the existing process-candidate loop
into applyReconcileProcesses, preserving its outcome validation, action/refusal
updates, counters, and already-absent handling, then invoke that helper before
the run-branch loop.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:35d8c33464a1dc71f8049ecd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Confirmed `applyReconcileReport` reclaimed worktrees and Run Branches before terminating owned process trees.
  - Reordered apply so process termination and absence proof complete before any worktree or branch cleanup.
  - The end-to-end apply test now asserts both superseded worktrees and branches still exist at the instant `TerminateTreeAndWait` runs.
  - Focused evidence: the reconcile regression and complete affected package suites passed (1,247 tests).
  - The Daemon owns authoritative `make verify` after this Agent turn.
