---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/releaseplan/proposal.go
line: 54
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tA,comment:PRRC_kwDOS0qyts7Wt95b
review_hash: 0bb4c6f689a9535dcd652d1b377fe0727ec3ef30163dc579d485f155f4cbf502
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-32/round-002/issue_010.md
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

- Decision: `UNREVIEWED`
- Notes:
