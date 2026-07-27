---
source: coderabbit
pr: "39"
round: 1
round_created_at: "2026-07-27T21:22:30Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-run-worktree-reconciliation
head_sha: 44fa422cea404a2d8c47e4b8011f065c4c0481ba
file: internal/cli/runs.go
line: 114
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMO,comment:PRRC_kwDOS0qyts7aNkKv
review_hash: 4a561e8878c0559f6f5c6fddd6c5d60f4150c85146d4071a7a1708f4d40df054
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 001: _ Performance & Scalability_ _ Major_ _ Quick win_

## Review Comment

_🚀 Performance & Scalability_ | _🟠 Major_ | _⚡ Quick win_

**`runs list` now forks Git per terminal Run on every invocation.** `countRetainedTerminalRunWorktrees` inspects every terminal implement Run returned by the unbounded all-states query, even when the user asked for `--state active` or a small `--limit`. Each `InspectTerminalRun` call executes several Git subprocesses, so listing cost grows linearly with terminal Run history and turns a cheap DB read into a multi-second operation on repositories with a long Run history.

Consider bounding this: skip inspection when the recorded `WorkDir` is empty, cap the number of inspected Runs, or compute the note only for states the listing is actually reporting.

Also note the swallowed error at Line 176 — an inspection failure is silently counted as retained with no diagnostic, which will mask Git problems.






Also applies to: 169-181

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/runs.go` at line 114, Update countRetainedTerminalRunWorktrees
and its InspectTerminalRun loop to avoid inspecting every terminal Run from the
unbounded query: skip empty WorkDir values and bound inspections to Runs
relevant to the requested state/limit, or otherwise enforce a safe cap. In the
inspection error path around the retained-worktree count, report the failure
with diagnostic context instead of silently treating it as retained.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fe34e99afa238b1ee8514b24 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Runs List performed a full reconciliation inspection, including
  multiple Git subprocesses, for every terminal implement Run in the
  unbounded history query. It also converted inspection errors into retained
  counts without a diagnostic. The fix batches Git root validation and Run
  Branch listing once per repository, checks recorded paths without Git, and
  emits contextual warnings for repository or path inspection failures.

## Verification

- `rtk proxy env GOCACHE=/private/tmp/roundfix-pr39-batch001-gocache go test ./internal/worktree ./internal/cli ./internal/store -run 'Test(CountRetainedTerminalRunsBatchesGitInspectionByRepository|InspectTerminalRunUnsafePath|RunRunsList.*Retained.*Worktree|RunRunsListReportsRetainedWorktreeInspectionFailure|PruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|ReconcileIntegration)' -count=1`
  — passed.
