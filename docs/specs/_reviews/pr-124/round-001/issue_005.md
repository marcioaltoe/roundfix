---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/reviewsource/coderabbit/coderabbit_test.go
line: 915
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wua0F,comment:PRRC_kwDOS0qyts7d35tb
review_hash: b660d780411776f6babbe14726cdfb43caad914472e0425c8685403d0d54d6ae
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**This test duplicates a case already covered in the precedence table.**

`TestEvidenceRateLimitWithoutAuthoritativeCommentStaysPending` uses the same input and asserts the same result as the table case "pull request 107 rate limited status stays pending" at Lines 729-735: a `CodeRabbit` status with state `success` and description `Review rate limited`, no comments, expecting `EvidencePending` with `EvidenceKindCommitStatus`.

Keep one of the two. The standalone test carries the clearer failure message, so removing the table case is the smaller change.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit_test.go` around lines 902 - 915,
Remove the duplicate “pull request 107 rate limited status stays pending” case
from the precedence table while retaining
TestEvidenceRateLimitWithoutAuthoritativeCommentStaysPending and its assertions.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:886d32df5cd2a51d9905c05b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Removed the duplicate “pull request 107 rate limited status stays
  pending” hierarchy row. The dedicated
  `TestEvidenceRateLimitWithoutAuthoritativeCommentStaysPending` remains the
  single clearer assertion of the same input and outcome.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test ./internal/reviewsource/coderabbit -count=1` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
