---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/cli/baseline_update_test.go
line: 1339
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA20,comment:PRRC_kwDOS0qyts7e5EBT
review_hash: 5fb753edc9c5026d6a5ebc66e7c34cfd93c3cbe27995a3a74dd4377ea4cc8818
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 014: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Share one Setup Manifest test helper across package `cli`.** Both files are in package `cli` and now carry independent read, decode, modify, and write logic for `docs/agents/setup-context.json`. The shared root cause is a missing single helper pair, which lets the two copies diverge in path, indentation, trailing newline, and permissions.
- `internal/cli/baseline_update_test.go#L688-L713`: replace the hardcoded `filepath.Join(repository, "docs", "agents", "setup-context.json")` in `readBaselineUpdateManifest` and `writeBaselineUpdateManifest` with the existing `baselineSetupManifestPath` constant, and export these two helpers as the single read/write pair for the package.
- `internal/cli/baseline_human_test.go#L1614-L1636`: rewrite `removeHumanBaselineManifestDecision` to call that shared read/write pair instead of repeating the read, `json.Unmarshal`, `json.MarshalIndent`, newline append, and `os.WriteFile` sequence.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/cli/baseline_update_test.go#L688-L713` (this comment)
- `internal/cli/baseline_human_test.go#L1614-L1636`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_update_test.go` around lines 688 - 713, In
internal/cli/baseline_update_test.go:688-713, use baselineSetupManifestPath in
readBaselineUpdateManifest and writeBaselineUpdateManifest, then export them as
the package-wide shared read/write helpers. In
internal/cli/baseline_human_test.go:1614-1636, update
removeHumanBaselineManifestDecision to call those shared helpers and remove its
duplicated path, JSON, newline, and file-writing logic.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/cli/baseline_update_test.go</file>
<line_range>688-713</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/baseline_human_test.go</file>
<line_range>1614-1636</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:bf471b45f1ea81633dc6a3b1 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Replaced hardcoded manifest path in read/write helpers with baselineSetupManifestPath. Exported as ReadBaselineSetupManifest/WriteBaselineSetupManifest. Updated removeHumanBaselineManifestDecision to use shared helpers. `rtk go build ./...` passes.
