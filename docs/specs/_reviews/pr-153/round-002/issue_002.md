---
source: coderabbit
pr: "153"
round: 2
round_created_at: "2026-08-10T22:09:34Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/speccheck/verification.go
line: 102
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YDU5o,comment:PRRC_kwDOS0qyts7fvgDK
review_hash: 74cb062070eb211dc718d81be8bcac80f8b590f64f16876ef9582102792c4d43
duplicate_of: ""
source_review_id: "4901334183"
source_review_submitted_at: "2026-08-10T22:08:33Z"
---

# Issue 002: _ Potential issue_ _ Major_

## Review Comment

_⚠️ Potential issue_ | _🟠 Major_

**Evaluate the complete shell predicate before reporting vacuity.**

`terminalSegment` returns the text after the last separator, but the last segment may not execute. On an unchanged tree, `git diff --name-only HEAD | grep -q . && cat` returns 1 because `grep` finds no input and `cat` is skipped. The current code sees `cat` and returns true. Parse `&&`/`||` control flow or conservatively reject chains that are not proven safe, then add this case as a non-vacuous regression. This repeats the final-predicate issue from the previous review, but this path remains unhandled.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/verification.go` around lines 88 - 102, The vacuity check
in terminalSegment incorrectly trusts the final command even when short-circuit
control flow may skip it. Update terminalSegment and its caller to evaluate
&&/|| chains conservatively, returning vacuous only when execution of the
relevant consuming command is guaranteed; otherwise reject the chain as
non-vacuous. Add a regression covering `git diff --name-only HEAD | grep -q . &&
cat`, which must not be reported as vacuous.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c480f6ffc389dd0595cf1544 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: `terminalSegment` returned the text after the last separator of any
  kind, so in `... | grep -q . && cat` the final `cat` was trusted even though
  the failing `grep -q .` short-circuits the `&&` and the tail never runs.
  `internal/speccheck/verification.go` now evaluates the `&&`/`||`/`;` chain
  with control flow: `chainPassesOnEmptyOutput` only reports vacuous when the
  chain provably reaches a passing predicate — an `&&` tail is skipped when the
  left side failed, an `||` tail runs when the left failed, and an unknown
  intermediate never lets an unexecuted tail decide the result. The old
  `commandSeparatorPattern`/`terminalSegment` are removed. Focused evidence:
  `go test ./internal/speccheck/` passes; new regressions pin `git diff
  --name-only HEAD | grep -q . && cat` as non-vacuous and `|| cat` as vacuous.
  Daemon owns authoritative `make verify`.
