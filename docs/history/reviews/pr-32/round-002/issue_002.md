---
source: coderabbit
pr: "32"
round: 2
round_created_at: "2026-07-17T13:23:47Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: d7ab1933ac9fdcf0c94d73e2f417d99d38e43fe7
file: internal/agent/sessions.go
line: 35
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5sP,comment:PRRC_kwDOS0qyts7Wt94g
review_hash: 74a8cd112a23c3f918a1bd6db9c1475e5b3045e9a3704b862f04e77e0b12b8d6
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:30Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Update session-name parsing for the new QA and review suffixes.**

`ParseRoundfixSessionName` treats these suffixes as part of the run ID, so generated names do not round-trip correctly. Parse `-qa`, `-review`, and `-review-NNN` explicitly and add helper-to-parser round-trip tests.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/sessions.go` around lines 19 - 35, Update
ParseRoundfixSessionName to recognize and remove the generated -qa, -review, and
zero-padded -review-NNN suffixes before extracting the run ID and batch number,
while preserving existing parsing behavior for other names. Add round-trip tests
covering SessionRefForQA and SessionRefForReview with and without a positive
batch number.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3d91edaae6edfee4ced6f4d5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Fixed `ParseRoundfixSessionName` to strip generated `-qa`, `-review`, and `-review-NNN` suffixes, with helper round-trip tests. Evidence: `GOCACHE=/private/tmp/roundfix-go-build rtk go test ./internal/agent ./internal/cli ./internal/config ./internal/daemon ./internal/releaseplan ./internal/spec ./internal/store ./internal/tui` passed.
