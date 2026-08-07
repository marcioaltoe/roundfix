---
source: coderabbit
pr: "137"
round: 1
round_created_at: "2026-08-07T03:22:12Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/promoted-backlog-entries-have-a-home
head_sha: ea93c68b70d066c1ee7f322e40ac1d547420e8be
file: internal/speccheck/constraints.go
line: 102
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XKdm6,comment:PRRC_kwDOS0qyts7eg8Kb
review_hash: fae98cced73d1056014a4ff28aab817d098bb989e861ddbe633c98ae7458eb29
duplicate_of: ""
source_review_id: "4879615443"
source_review_submitted_at: "2026-08-07T03:12:07Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Wrap the detector error with context.**

Return `fmt.Errorf("detect backlog promotion: %w", err)` here. The current return loses the operation name when the caller receives a filesystem or YAML parsing error.

As per coding guidelines, wrap propagated internal errors with context using `fmt.Errorf("{context}: %w", err)`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/constraints.go` around lines 100 - 102, Update the error
return after detectBacklogPromotion in the surrounding constraints flow to wrap
err with the operation context using fmt.Errorf and %w, preserving the
underlying error while identifying backlog promotion detection; ensure fmt is
imported if needed.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1bdca05442f75ba61e8359ee -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `Check` returned `detectBacklogPromotion` errors directly, unlike the
  repository rule that each propagated internal error name the failed
  operation and preserve the chain with `%w`.

## Resolution

- Wrapped the detector failure as `fmt.Errorf("detect backlog promotion: %w",
  err)`.
- Added a regression case with invalid Backlog Entry YAML that asserts the
  outer detector context, the existing inner `parse Backlog Entry` context,
  and a preserved wrapped cause.
- Reproduction evidence before the production fix: the focused test received
  `parse Backlog Entry ...` without the `detect backlog promotion` prefix.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/speccheck -run TestCheckBacklogUnmoved -count=1` passed after the
  fix; `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/speccheck -count=1` also passed.
- Authoritative Verification `make verify` was not run; the Daemon owns it
  after this Agent turn.
