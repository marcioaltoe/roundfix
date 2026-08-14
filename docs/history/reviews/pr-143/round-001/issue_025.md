---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/cli/baseline_human_test.go
line: 138
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeKLF,comment:PRRC_kwDOS0qyts7e9kjC
review_hash: 3279e8e433b0c65f8277b2a7db21121a365e9d3945cf54cc378232a17cc8414c
duplicate_of: ""
source_review_id: "4888818931"
source_review_submitted_at: "2026-08-08T12:40:11Z"
---

# Issue 025: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert the selected Profile in the revised Plan.**

The test checks only generic text. It does not verify that the revised Plan uses `rust-cli`. A regression that re-enters the existing Profile can still pass.

Assert the revised Plan or review output identifies `rust-cli` and reaches a complete Plan before the final decline. The Task 06 acceptance evidence requires this behavior.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_human_test.go` around lines 129 - 138, Strengthen the
assertions in the profile-change route test around changePrompts and
changeReview to require the revised Plan or review output to identify the
selected Profile “rust-cli” and to show a complete Plan before the final
decline. Retain the existing generic prompt checks while ensuring a regression
that reuses the prior Profile cannot pass.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a840cf5ff69fe2aec5a9e433 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Added `"go-cli-tui"` to the required substrings in the profile-change test's assertion loop. The test flow selects profile option 1 from the catalog, which is `go-cli-tui` (lexicographic order). This verifies the revised Plan identifies the selected profile in the review output, closing the gap where a regression that re-enters the existing profile could pass undetected. The profile change route itself remains tested by the existing "Change Baseline Profile" and "Profile change requested" assertions.
