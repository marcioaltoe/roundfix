---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/classification_test.go
line: 368
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0YmU,comment:PRRC_kwDOS0qyts7cjgEu
review_hash: d93e6d6b76c6ee40ff4246840ee43a8c4eb3e78fb1ce3b1a40f320ee16cec2eb
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:29Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Add `t.Parallel()` to the new tests and subtests.**

`TestCarrierClassification`, its two subtests, and `TestUnclassifiableCarrierStillWarns` each create their own repository through `newInspectionRepository` or `newBaselinePlanCharacterizationRepository`. They share no global state, so they are independently runnable. The second subtest runs a full plan-apply-replan cycle and is slow. Mark the tests and subtests parallel.






As per coding guidelines: "Independent tests SHOULD use `t.Parallel()` when possible."


Also applies to: 494-495, 528-529

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/classification_test.go` around lines 367 - 368, Add
t.Parallel() to TestCarrierClassification and both of its subtests, plus
TestUnclassifiableCarrierStillWarns. Place the calls at the start of each test
function so these independent repository-backed tests, including the slow
plan-apply-replan subtest, can run concurrently.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:de3bb21ee070f826513f27b6 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The independent parent tests and repository-backed subtests now call `t.Parallel()`. The full Baseline package test passed.
