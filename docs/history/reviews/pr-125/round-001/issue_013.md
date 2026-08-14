---
source: coderabbit
pr: "125"
round: 1
round_created_at: "2026-08-05T19:38:23Z"
status: invalid
terminal_reason: "Normal review Runs are serialized by the Active Run work-target lock and Branch Integrity also rejects lockless bypass Runs, so concurrent requesters cannot race this lookup/post sequence."
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 1384d928c1d0af4d0bca06e506eb8b2953f9f341
file: internal/reviewsource/coderabbit/coderabbit.go
line: 157
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WxMJ3,comment:PRRC_kwDOS0qyts7d791U
review_hash: 7e04a9dc0bdff589510bb14c9d91bf6fe51aabd78b602264c4e2496a0aead842
duplicate_of: ""
source_review_id: "4868102567"
source_review_submitted_at: "2026-08-05T19:37:36Z"
---

# Issue 013: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/usr/bin/env bash
set -euo pipefail

rg -n -C 4 '\bRequestReview\s*\(|\.RequestReview\s*\(|\bLock\s*\(|flock|singleflight' internal
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/usr/bin/env bash
set -euo pipefail

printf '%s\n' '--- coderabbit outline ---'
ast-grep outline internal/reviewsource/coderabbit/coderabbit.go --view expanded
printf '%s\n' '--- coderabbit implementation ---'
sed -n '1,240p' internal/reviewsource/coderabbit/coderabbit.go
printf '%s\n' '--- coderabbit tests ---'
sed -n '1,280p' internal/reviewsource/coderabbit/coderabbit_test.go
printf '%s\n' '--- reviewsource types ---'
sed -n '1,150p' internal/reviewsource/reviewsource.go
printf '%s\n' '--- related GitHub methods ---'
rg -n -C 5 'func .*IssueComments|func .*CommentOnPullRequest|type .*GitHub|IssueComments\(|CommentOnPullRequest\(' internal/reviewsource internal | head -n 240
printf '%s\n' '--- state and locking candidates ---'
rg -n -i -C 3 'review.?request|coderabbit|idempot|dedup|marker|process.?lock|repository.?lock|flock|lockfile|exclusive' internal --glob '*.go' --glob '!**/*_test.go' | head -n 320
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/usr/bin/env bash
set -euo pipefail

printf '%s\n' '--- focused call graph ---'
rg -n -C 8 'RequestReview\(|publishReviewRequestEvent|RoundfixCommentBody|isRoundfixReviewRequestComment|IssueComments\(|CommentOnPullRequest\(' internal/reviewsource internal/watch internal/cli
printf '%s\n' '--- client construction and process lifetime ---'
rg -n -C 5 'coderabbit\.Client|reviewsource\.ReviewRequester|New.*Client|type Client struct' internal --glob '*.go' --glob '!**/*_test.go'
printf '%s\n' '--- test synchronization helpers ---'
rg -n -C 5 'blocking|barrier|entered|release|concurrent|WaitGroup|go func' internal/reviewsource/coderabbit internal/reviewsource --glob '*_test.go'
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---



</details>

**Serialize per-head marker creation.** `RequestReview` lists comments and then posts without an atomic boundary, so separate Roundfix processes can both publish the same head request. Use shared idempotency state or a cross-process lock keyed by repository, pull request, and head that spans both operations. Add a regression test that blocks both lookups before either publication.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 141 - 157,
Update RequestReview to serialize the comment lookup and publication for each
repository, pull request, and headSHA key, using shared idempotency state or a
cross-process lock that spans both IssueComments and CommentOnPullRequest
operations. Preserve the existing outcome and error behavior, and add a
regression test that blocks concurrent lookups until both requests are waiting,
then verifies only one publication occurs.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e8505c3efe2e1789f60330a2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The requester primitive is deliberately idempotent within Roundfix's serialized Run lifecycle, not a cross-process transactional API. The Run Database rejects a second Active review Run for the same `(Head Repository, PR Head Branch)`, and Branch Integrity separately discovers and rejects Runs created through the explicit lock bypass before another normal Run starts in the checkout.
- Evidence: ADR-0012 and `internal/store/store.go` define the work-target lock; `TestBranchIntegrityRejectsLocklessBypassRun` covers the bypass case; the focused `internal/cli` and CodeRabbit suites passed. Adding external locking or a new persistence schema would expand the reviewed contract without a reachable production race.
