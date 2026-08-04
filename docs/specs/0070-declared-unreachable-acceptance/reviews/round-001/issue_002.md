---
source: coderabbit
pr: "110"
round: 1
round_created_at: "2026-08-04T22:55:35Z"
status: invalid
terminal_reason: "ADR-0080 explicitly permits pass with environment-blocked rows backed by equivalent evidence."
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0070-implementation
head_sha: a588c6ca3ab9d977284ba1f9e80a89b0e6336786
file: internal/spec/archive.go
line: 81
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WeYX2,comment:PRRC_kwDOS0qyts7dggqW
review_hash: 94c6fc888e3c7df3b66e473943350be0b25c0e150b04233fb97d2ad08208f517
duplicate_of: ""
source_review_id: "4859094834"
source_review_submitted_at: "2026-08-04T21:23:48Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Reject blocked-row counts in `pass` reports.**

At Line 80, `VerdictPass` returns before the blocked-row checks at Lines 86-93. A `pass` report with nonzero `rows_blocked_*` values can archive. `internal/cli/archive_test.go`, Line 24, currently supplies `rows_blocked_environment: 3` and expects that result.

Require all blocked-row counters to be zero before accepting `VerdictPass`. Update that test to expect Preflight Validation failure.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/archive.go` around lines 80 - 81, Update the VerdictPass
handling in the archive validation flow to verify every rows_blocked_* counter
is zero before returning success; otherwise return the existing Preflight
Validation failure. Adjust the relevant archive test to expect failure when
rows_blocked_environment is nonzero.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2366bb17b519225d75754988 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: ADR-0080 and the canonical `qa-gate` verdict rules explicitly allow a nonzero `rows_blocked_environment` count on `pass` when equivalent observed or supervised evidence exists. The archived Spec 0053 QA contract also names environment-blocked `pass` as readable and archivable. Requiring every blocked counter to be zero would break that accepted behavior and the Spec 0070 non-regression requirement.
- Evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli -run '^TestRunArchive'` passed, including the compatibility case cited by the finding.
