---
source: coderabbit
pr: "32"
round: 2
round_created_at: "2026-07-17T13:23:47Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: d7ab1933ac9fdcf0c94d73e2f417d99d38e43fe7
file: internal/releaseplan/proposal.go
line: 54
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tA,comment:PRRC_kwDOS0qyts7Wt95b
review_hash: 0bb4c6f689a9535dcd652d1b377fe0727ec3ef30163dc579d485f155f4cbf502
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 010: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Do not emit an empty `Approval.Increment` enum.**

Manual-required, no-release, and patch proposals leave `Approval` zero-valued. Since JSON always emits `approval.increment`, consumers receive `""`, which is not one of the declared `IncrementKind` values. Initialize it to `IncrementNone` and update the test expecting an empty string.

   


Also applies to: 66-80

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/releaseplan/proposal.go` around lines 50 - 54, Initialize
Approval.Increment to IncrementNone for manual-classification, no-release, and
patch proposal paths in the proposal construction logic, including the branches
around the referenced manual-required case, so JSON never emits an empty enum
value. Update the corresponding test assertions to expect IncrementNone instead
of an empty string.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:50a68cce594af50ad7f333fb -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Initialized non-approval proposal `Approval.Increment` to `IncrementNone` and updated approval tests so JSON never emits an empty increment enum. Evidence: `GOCACHE=/private/tmp/roundfix-go-build rtk go test ./internal/agent ./internal/cli ./internal/config ./internal/daemon ./internal/releaseplan ./internal/spec ./internal/store ./internal/tui` passed.
