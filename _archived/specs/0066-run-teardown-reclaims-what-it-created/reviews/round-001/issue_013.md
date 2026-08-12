---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/worktree/worktree_test.go
line: 1572
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9YY,comment:PRRC_kwDOS0qyts7dnSb8
review_hash: 64bf258120c7f481a689fef013ce07f414122686c35e4971a7a8312624beb58f
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 013: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Add a test for the stale-proof refusal path.**

`ApplyRunBranchCandidate` refuses when `fresh.ReleasableProofs[branch] != proof` (worktree.go line 537). No test covers that branch. It is the primary safety property of a destructive operation: it stops a branch removal after the evidence changed between inspection and apply. Add a test that classifies the set, then commits new non-QA work to the candidate branch or a newer QA report to the target, and asserts the refusal plus branch survival. The existing `TestApplyTerminalRunSupersededPreservesWhenFreshProofFails` is a usable model.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree_test.go` around lines 1490 - 1572, Add a test
alongside the existing ApplyRunBranchCandidate tests that classifies the
run-branch set, changes the candidate’s evidence before applying it (for example
by committing non-QA work on the candidate branch or adding a newer QA report on
the target), then verifies ApplyRunBranchCandidate refuses with the stale-proof
error and preserves the candidate path and branch. Use
TestApplyTerminalRunSupersededPreservesWhenFreshProofFails as the behavioral
model.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:4b9367cbd8ea911c689e5553 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Added `TestApplyRunBranchCandidatePreservesStaleSupersedingProof` to mutate the QA proof after inspection, assert the stale-proof refusal, and prove the candidate worktree and branch survive.
  - This exercises the existing apply-time reclassification guard rather than duplicating production logic in the test.
  - Focused evidence: the targeted worktree suite passed (9 tests); complete affected package suites passed (1,247 tests).
  - The Daemon owns authoritative `make verify` after this Agent turn.
