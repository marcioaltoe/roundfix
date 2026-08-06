---
source: coderabbit
pr: "129"
round: 1
round_created_at: "2026-08-06T03:33:20Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0075-typed-docs-backlog
head_sha: 04b156c2a36969a67c06958bbc366fc47a6db816
file: internal/baseline/preservation_test.go
line: 262
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W1gkM,comment:PRRC_kwDOS0qyts7eCVqL
review_hash: 64ed238d7dc6f57ab60165ed96ebaa9340c1da2bfeed0fc5afbc4c024590a77b
duplicate_of: ""
source_review_id: "4870101613"
source_review_submitted_at: "2026-08-06T00:57:05Z"
---

# Issue 002: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Remove the fixed count checks from the invariant test.**

`maintainedSourceBaselineEntries` and `maintainedSourceBaselineAccounting` still require exactly 106 entries and 51 accounting records. Any legitimate corpus expansion fails even when `Identity.EntryCount` matches the entries slice.

This conflicts with Line 321 through Line 326, which state that the protected invariant holds for any corpus size. Keep structural invariants, or document this as an exact snapshot test instead.





Also applies to: 321-329

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/preservation_test.go` around lines 255 - 262, Remove the
fixed-value assertions using maintainedSourceBaselineEntries and
maintainedSourceBaselineAccounting from the invariant test in
preservation_test.go. Keep the structural checks, including consistency between
Identity.EntryCount and the entries slice, so the test remains valid for corpus
changes of any size.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:26275717498085f7caf488b5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The suggested removal was rejected because this maintained fixture is
  intentionally an exact-shape characterization gate: the named expectations
  force legitimate corpus changes through the sanctioned re-recording
  workflow. The finding correctly identified that the adjacent comment instead
  described a corpus-size-independent invariant. The comment now states both
  protected properties without weakening the exact entry and accounting count
  assertions.
- Focused evidence:
  `rtk env GOCACHE=/Users/marcio/dev/roundfix-c/.gocache go test ./internal/baseline -run '^TestReadoptionCompatibilityMaintainedFixture$' -count=1`
  exited 0 (`ok roundfix/internal/baseline`). `rtk git diff --check` for the
  tracked test diff exited 0, and an `rtk awk` trailing-whitespace check over
  the changed test plus both assigned issue files exited 0. The Daemon still
  owns authoritative `make verify`.
