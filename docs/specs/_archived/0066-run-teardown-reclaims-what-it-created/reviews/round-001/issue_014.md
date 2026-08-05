---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/worktree/worktree.go
line: 511
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9Yc,comment:PRRC_kwDOS0qyts7dnScB
review_hash: bbc0d5859eb724da465df1708cae90c2074f16298a6f96f24b3d8daee08c20db
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 014: _ Maintainability & Code Quality_ _ Major_ _ Heavy lift_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _🏗️ Heavy lift_

**Three new functions exceed the configured length and complexity limits.** The coding guidelines cap functions at 120 lines and 80 statements and cap cyclomatic complexity at 13. Each of these three functions inlines validation, grouping, per-item classification, and result assembly in one body, and each also repeats a near-identical result literal several times. The shared remediation is the same: extract the grouping loop, the per-item decision, and the result construction into named helpers.
- `internal/worktree/worktree.go#L323-L511`: split `classifyRunBranchSet` (about 189 lines); extract the candidate-collection loop at lines 388-448 and the current-selection block at lines 450-499, sharing the `release` and `preserve` closures through a small classifier struct.
- `internal/cli/reconcile.go#L402-L561`: split `inspectReconcileRunBranches` (about 160 lines, nesting depth 5); add one constructor for the seven repeated preserved `reconcileDebrisResult` literals and move the per-run decision body into its own function.
- `internal/cli/cli.go#L2076-L2222`: split `filterPendingRunWorkByTarget` (about 146 lines); extract the per-spec classification loop at lines 2139-2210 and replace the four positional return slices with a named result struct.

As per coding guidelines: "Keep cyclomatic complexity at or below 13, avoid deeply nested conditionals, limit functions to 120 lines and 80 statements, and avoid duplication exceeding the configured 100-token threshold."

<details>
<summary>📍 Affects 3 files</summary>

- `internal/worktree/worktree.go#L323-L511` (this comment)
- `internal/cli/reconcile.go#L402-L561`
- `internal/cli/cli.go#L2076-L2222`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree.go` around lines 323 - 511, Reduce complexity and
duplication across the three affected functions: in
internal/worktree/worktree.go lines 323-511, refactor classifyRunBranchSet by
introducing a classifier struct that owns release/preserve behavior, extracting
candidate collection and current-selection logic into named helpers, and moving
result assembly into a helper as needed; in internal/cli/reconcile.go lines
402-561, refactor inspectReconcileRunBranches by adding a constructor for
repeated reconcileDebrisResult values and extracting per-run decisions; in
internal/cli/cli.go lines 2076-2222, refactor filterPendingRunWorkByTarget by
extracting the per-spec classification loop and returning a named result struct
instead of four positional slices. Keep behavior unchanged while bringing each
function within the stated length, statement, nesting, complexity, and
duplication limits.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/worktree/worktree.go</file>
<line_range>323-511</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/reconcile.go</file>
<line_range>402-561</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/cli.go</file>
<line_range>2076-2222</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:8ea9a57b4ae0915a9ff9d7ef -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes:
  - The repository has no `.golangci*` configuration, and its authoritative `verify` target is `fmt-check test spec-budget skills-sync-check skills-check build spec-check`; it does not configure the numeric function-length, statement, cyclomatic, or duplication thresholds claimed by the finding.
  - `rtk rg --files -g '.golangci*'` found no configuration. `rtk rg -n '^verify:|^lint:' Makefile` found the verification gate and no lint target.
  - A broad decomposition of three working orchestration functions would therefore be an ungoverned refactor rather than a repair of a repository contract. The focused behavior changes in this Batch were made at their actual boundaries and remain covered by the 1,247 passing affected-package tests.
