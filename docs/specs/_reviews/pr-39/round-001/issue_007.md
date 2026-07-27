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
line: 857
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMy,comment:PRRC_kwDOS0qyts7aNkLZ
review_hash: b299f5390d70d2d3dd97d892568148a7657cf284b1f8532eb8f3714436df1e4d
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 007: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**`TerminalRunLookup` cannot receive the prune context.**

The lookup performs a Run Store read but takes only `runID`, so implementations must capture a context from the enclosing scope (or fall back to `context.Background()`), and the cancellation/deadline of the `ctx` passed to `PruneTerminalReport` is not propagated into the database call. Adding `ctx` as the first parameter keeps context propagation explicit through this boundary.

As per coding guidelines: "Declare `ctx context.Context` as the first parameter of functions that accept a context" and "Use context-aware APIs for downstream operations".




<details>
<summary>♻️ Proposed signature change</summary>

```diff
-type TerminalRunLookup func(runID string) (store.Run, bool, error)
+type TerminalRunLookup func(ctx context.Context, runID string) (store.Run, bool, error)
```

```diff
-		run, found, err := loadTerminalRun(runID)
+		run, found, err := loadTerminalRun(ctx, runID)
```
</details>


Also applies to: 899-903

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree.go` around lines 849 - 857, Update the
TerminalRunLookup function type to accept context.Context as its first
parameter, then pass the PruneTerminal ctx through every lookup invocation
within PruneTerminal and related reconciliation paths. Update all
implementations and callers to use the context-aware Run Store read, preserving
existing lookup behavior while propagating cancellation and deadlines.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:69df6b4ed7989ea308da9e35 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `TerminalRunLookup` forced callers to capture an outer context, so
  the prune operation could not express context propagation at the function
  boundary. The lookup now accepts `context.Context` first,
  `PruneTerminalReport` passes its context through, and store-backed callers
  use the received context.

## Verification

- `rtk proxy env GOCACHE=/private/tmp/roundfix-pr39-batch001-gocache go test ./internal/worktree ./internal/cli ./internal/store -run 'Test(CountRetainedTerminalRunsBatchesGitInspectionByRepository|InspectTerminalRunUnsafePath|RunRunsList.*Retained.*Worktree|RunRunsListReportsRetainedWorktreeInspectionFailure|PruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|ReconcileIntegration)' -count=1`
  — passed.
