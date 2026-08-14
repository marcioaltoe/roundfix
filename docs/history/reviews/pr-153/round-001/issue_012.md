---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/speccheck/verification.go
line: 23
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWjN,comment:PRRC_kwDOS0qyts7fswO3
review_hash: e5726446f0f2fd095d7a6263b9f872b9dab22410eba1288c88e34288735eb77f
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 012: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Do not classify nonempty-diff assertions as checks that pass on an unchanged tree.**

`git diff --name-status HEAD | grep -q .` writes no input to `grep` on an unchanged tree. `grep -q .` then exits with status 1. The command does not satisfy the stated vacuity condition, but `workingTreeStatePattern` reports it because it only inspects the Git flags. This creates a false `SC-VERIFY-VACUOUS-COMMAND` finding.

- `internal/speccheck/verification.go#L21-L23`: inspect the complete command predicate, or conservatively exclude piped commands unless the complete form is known to exit zero for an unchanged tree.
- `internal/speccheck/citations.go#L1288-L1298`: emit the unchanged-tree diagnostic only for commands that the corrected predicate proves pass on an unchanged tree.
- `internal/speccheck/verification_test.go#L107-L109`: mark the `grep -q .` command as non-vacuous and add coverage for an explicitly inverted empty-output check.

<details>
<summary>📍 Affects 3 files</summary>

- `internal/speccheck/verification.go#L21-L23` (this comment)
- `internal/speccheck/citations.go#L1288-L1298`
- `internal/speccheck/verification_test.go#L107-L109`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/verification.go` around lines 21 - 23, Update
workingTreeStatePattern in internal/speccheck/verification.go to evaluate the
complete command predicate, or conservatively exclude piped commands unless
unchanged-tree success is proven; update the diagnostic emission in
internal/speccheck/citations.go:1288-1298 to rely only on that corrected
predicate. In internal/speccheck/verification_test.go:107-109, classify the grep
-q . command as non-vacuous and add coverage for an explicitly inverted
empty-output check.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/speccheck/verification.go</file>
<line_range>21-23</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/speccheck/citations.go</file>
<line_range>1288-1298</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/speccheck/verification_test.go</file>
<line_range>107-109</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3a2dc7c86706478733f5c6a5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Reworked the vacuity predicate in internal/speccheck/verification.go to classify the complete command, not git flags. `workingTreeCleanlinessCheck` now evaluates the terminal segment after the last shell separator; `git diff --name-status HEAD | grep -q .` (and any `| grep -q .` terminal, including `--exit-code` variants) is no longer classified vacuous because its terminal predicate fails on empty output. `citations.go` already emits via `VacuousVerificationCommands`, so the corrected predicate is the sole gate. Added regression coverage in internal/speccheck/verification_test.go including an inverted empty-output check. Focused: `go test ./internal/speccheck` ok.
