---
source: coderabbit
pr: "128"
round: 1
round_created_at: "2026-08-06T03:35:45Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0069-review-run-targets-its-pull-request
head_sha: 62cd2ea6f84aa181570ef18f0e05225c6e4ebb88
file: internal/preflight/preflight.go
line: 313
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W1Guf,comment:PRRC_kwDOS0qyts7eBurl
review_hash: 3b35dca7559b08d4a0f3e6627919b4227f1e0d30fe580335307e9407360b2ed3
duplicate_of: ""
source_review_id: "4869925235"
source_review_submitted_at: "2026-08-06T00:16:28Z"
---

# Issue 006: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Validate the PR head revision before accepting the branch.**

Lines 311-313 return success when only the branch matches. They do not compare `pullRequest.HeadSHA` with `gitState.HEAD`. A stale checkout on the correct branch can therefore pass preflight and run against a revision that is not the PR head.

Return `TargetMismatch` when both revisions are known and differ. Update the matching fixture to use the checkout SHA, and add a same-branch, different-SHA case.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/preflight/preflight.go` around lines 311 - 313, Update the preflight
branch-acceptance logic around expectedBranch and foundBranch to also compare
pullRequest.HeadSHA with gitState.HEAD when both revisions are known, returning
TargetMismatch on a difference. Adjust the matching fixture to use the checkout
SHA and add coverage for a same-branch, different-SHA scenario.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:b305b816fcfae3884a93f4a5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The branch-equality early return skipped the already-resolved
  `PullRequest.HeadSHA`. The guard now accepts matching branches only when the
  known checkout and PR revisions also match, and returns `TargetMismatch`
  otherwise. The canonical table now uses `abc123` for the matching fixture
  and covers same-branch `abc123` versus `remote-head`. Before the production
  fix, that focused case failed with `expected TargetMismatch, got <nil>`;
  afterward, `go test ./internal/preflight -count=1` passed 51 tests and the
  focused CLI no-side-effects suite passed 4 tests. Authoritative `make
  verify` remains Daemon-owned.
