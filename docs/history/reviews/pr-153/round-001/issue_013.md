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
line: 83
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YCprG,comment:PRRC_kwDOS0qyts7fukY8
review_hash: a3053d035a5c23663f29ed5b1560e6f0d55d8d3977876042e1f9df8f8d7aef09
duplicate_of: ""
source_review_id: "4901105231"
source_review_submitted_at: "2026-08-10T21:29:18Z"
---

# Issue 013: _ Functional Correctness_ _ Major_ _ Heavy lift_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _🏗️ Heavy lift_

**Classify the final shell predicate, not Git flags inside it.**

`git diff --exit-code HEAD | grep -q .` exits `1` on an unchanged tree because `grep` receives no input. Line 83 still reports it as vacuous because `--exit-code` matches. Conversely, `git diff --name-only HEAD | cat` exits `0` on an unchanged tree but does not match the success pattern.

This causes both false `SC-VERIFY-VACUOUS-COMMAND` findings and missed vacuous commands. Add regression cases first, then restrict the matcher to compound forms whose terminal predicate is proven to succeed on empty output.

- `internal/speccheck/verification.go#L71-L83`: determine unchanged-tree success from the complete pipeline or command chain.
- `internal/speccheck/verification_test.go#L107-L117`: add `git diff --exit-code HEAD | grep -q .` as non-vacuous and `git diff --name-only HEAD | cat` as vacuous.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/speccheck/verification.go#L71-L83` (this comment)
- `internal/speccheck/verification_test.go#L107-L117`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/verification.go` around lines 71 - 83, Update
workingTreeCleanlinessCheck to classify unchanged-tree success from the complete
pipeline or command chain, not Git flags such as --exit-code; only accept
compound commands when their terminal predicate is proven to succeed on empty
output. Add regression cases in internal/speccheck/verification_test.go: treat
`git diff --exit-code HEAD | grep -q .` as non-vacuous and `git diff --name-only
HEAD | cat` as vacuous.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/speccheck/verification.go</file>
<line_range>71-83</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/speccheck/verification_test.go</file>
<line_range>107-117</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ac39f926f01c36f9087f4ae8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Reworked `workingTreeCleanlinessCheck` in internal/speccheck/verification.go to classify from the complete pipeline/chain rather than git flags. It now extracts the terminal segment after the last shell separator (`terminalSegment`) and matches `emptyOutputSucceedsPattern` against that terminal predicate only. `git diff --exit-code HEAD | grep -q .` is now classified non-vacuous (terminal `grep -q .` fails on empty output despite `--exit-code`), and `git diff --name-only HEAD | cat` is now classified vacuous (terminal `cat` succeeds on empty output). Added the two regression cases plus an inverted empty-output check in internal/speccheck/verification_test.go. Focused: `go test ./internal/speccheck` ok.
