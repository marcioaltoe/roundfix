---
source: coderabbit
pr: "153"
round: 3
round_created_at: "2026-08-10T22:33:54Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: 9752569db9792fe76459a9f007a196f92c2c5357
file: internal/speccheck/verification.go
line: 40
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YDrKC,comment:PRRC_kwDOS0qyts7fv-Jd
review_hash: 5e2b75b2243cae3cd789e7bbf47c1c8d84282556c11e486a874344b97b3af9c0
duplicate_of: ""
source_review_id: "4901443178"
source_review_submitted_at: "2026-08-10T22:30:02Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Do not classify variable-dependent predicates as guaranteed successes.**

`test -z "$value"` and `[ "$left" = "$right" ]` can fail when earlier chain segments set different values. For example, `make regen` can change the second `git status --porcelain` snapshot before the equality test. The matcher reports this command as vacuous, and `internal/speccheck/citations.go` then emits a blocking error finding.

Treat argument-bearing `test -z` and variable-equality predicates as `emptyInputUnknown` unless their values are statically proven. Add a regression for a two-snapshot comparison with an intervening command and expect `vacuous: false`.

- `internal/speccheck/verification.go#L37-L40`: restrict success matching to predicates that are unconditionally successful.
- `internal/speccheck/verification_test.go#L155-L183`: add the variable-dependent snapshot-comparison regression.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/speccheck/verification.go#L37-L40` (this comment)
- `internal/speccheck/verification_test.go#L155-L183`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/verification.go` around lines 37 - 40, The success matcher
in internal/speccheck/verification.go, around the predicate patterns at lines
37-40, must stop classifying argument-bearing test -z and variable-equality
predicates as unconditionally successful; restrict matching to statically
unconditional predicates and return emptyInputUnknown for variable-dependent
values. Add the requested two-snapshot comparison regression in
internal/speccheck/verification_test.go lines 155-183, including an intervening
command, and assert vacuous: false.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/speccheck/verification.go</file>
<line_range>37-40</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/speccheck/verification_test.go</file>
<line_range>155-183</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:462a7c2cff7522c7d2baa536 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: `emptyOutputSucceedsPattern` in `internal/speccheck/verification.go`
  treated argument-bearing `test -z`/`[ -z` and variable-equality predicates
  as guaranteed empty-input successes. Those forms depend on runtime values —
  a preceding chain segment can set a different operand (e.g. a second
  `git status --porcelain` snapshot after `make regen`) — so they were not
  unconditional. The success matcher now restricts `test -z`/`[ -z` to
  command-substitution or empty literals (`test -z "$(cmd)"`, `test -z ""`)
  and drops the variable-equality alternatives entirely; everything else
  falls to `emptyInputUnknown` and is never vacuous. The existing two-snapshot
  table case in `verification_test.go` was renamed to "a two-snapshot
  comparison with an intervening command" and flipped from `vacuous: true` to
  `vacuous: false`. Focused evidence: `go test ./internal/speccheck/` passes,
  `go vet ./internal/speccheck/` is clean, and the flipped case asserts the
  variable-dependent equality predicate is not reported vacuous. Daemon owns
  authoritative `make verify`.
