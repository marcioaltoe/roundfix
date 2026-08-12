---
source: coderabbit
pr: "113"
round: 1
round_created_at: "2026-08-05T02:12:07Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0068-implementation
head_sha: c9af2617f988bd63e1bd8f22c6178758a8e5fd40
file: internal/specaudit/audit.go
line: 663
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WhZqR,comment:PRRC_kwDOS0qyts7dk969
review_hash: fa7311ef903e6a7512208f8d64381aac7c431daa45634ea9bf37f931bdd7888d
duplicate_of: ""
source_review_id: "4860420451"
source_review_submitted_at: "2026-08-05T02:11:27Z"
---

# Issue 005: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**`git branch -d` cannot delete content-integrated residue.** `classifyBranch` chooses the delete flag without consulting the integration proof from `classifyGitRef`. A squash-merged branch is classified `residue` through content equality while remaining unreachable from the default branch, and `git branch -d` refuses to delete it. Two tests assert the failing string instead of executing it, so the defect passes CI.
- `internal/specaudit/audit.go#L655-L663`: select the delete flag from the proof. Emit `-d` only when the ref is reachable from the default branch, and `-D` when only the content is represented.
- `internal/specaudit/audit_test.go#L95-L101`: expect `-D` for the squash-merged fixture, and execute the emitted command as `TestAuditScratchWorktreeReclaimCommandRuns` does.
- `internal/cli/spec_check_test.go#L281-L289`: update the expected `reclaim:` line to match the corrected flag.

<details>
<summary>📍 Affects 3 files</summary>

- `internal/specaudit/audit.go#L655-L663` (this comment)
- `internal/specaudit/audit_test.go#L95-L101`
- `internal/cli/spec_check_test.go#L281-L289`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/specaudit/audit.go` around lines 655 - 663, Update classifyBranch in
internal/specaudit/audit.go:655-663 to choose the reclaim delete flag from
classifyGitRef’s integration proof, using -d only for refs reachable from the
default branch and -D for content-only integration. In
internal/specaudit/audit_test.go:95-101, expect -D for the squash-merged fixture
and execute the emitted reclaim command as
TestAuditScratchWorktreeReclaimCommandRuns does. Update the expected reclaim
line in internal/cli/spec_check_test.go:281-289 to match.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/specaudit/audit.go</file>
<line_range>655-663</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/specaudit/audit_test.go</file>
<line_range>95-101</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/spec_check_test.go</file>
<line_range>281-289</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ca1441034f6688166e0c99a5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: A squash-integrated local branch is residue by content but not ancestry, so `git branch -d` refuses the audit's emitted reclaim command.

## Resolution

- Propagated the integration proof from `classifyGitRef`; local branches use `-d` for reachability proof and `-D` for content-only proof.
- Updated both Spec Audit expectations and executed the squash-integrated branch's emitted reclaim command in the real Git fixture.
- Focused evidence: the regression first failed because the command was `git branch -d`; `rtk env GOCACHE=/private/tmp/roundfix-review0068-audit-cache go test ./internal/specaudit -run '^(TestAuditClassifiesResidueBranch|TestAuditScratchWorktreeReclaimCommandRuns|TestAuditPreservesActiveRunSurvivors)$' -count=1` exited 0 after the fix, and the full affected-package command exited 0.
- Daemon Verification: `make verify` was not run; the Daemon owns that command.
