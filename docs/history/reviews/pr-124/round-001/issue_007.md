---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/reviewsource/coderabbit/coderabbit.go
line: 296
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wua0R,comment:PRRC_kwDOS0qyts7d35tr
review_hash: dea266faf5c63c77ca3128798a148786f601f15f1e7d94e8df406e70672ee00c
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 007: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Refusal classification runs twice on the `Evidence` path.**

`Evidence` calls `refusalEvidence` at Line 294 with `unresolvedThreads = 0`, and `classifyEvidence` calls it again at Line 682 with the real count. The early call always short-circuits first, so the second call is only reachable when `classifyEvidence` is invoked directly. The behavior is correct, but the duplicated precedence rule is easy to break later: a future change to refusal precedence must be applied in two places.

Consider moving the refusal check into a single place. One option is to keep only the `classifyEvidence` check and let the early path stop at `reviewingEvidence`, at the cost of one extra reviews/threads fetch on refusal. Another option is to pass the already-fetched comments into `classifyEvidence` and drop the early call, keeping a dedicated helper for the fetch-skipping decision.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 287 - 296, The
refusal classification is duplicated between the Evidence flow and
classifyEvidence, creating two precedence rules to maintain. Consolidate refusal
handling into a single classification path by updating the surrounding Evidence
logic and classifyEvidence to reuse the already-fetched comments where
applicable, while retaining a separate helper only for deciding whether comment
fetching is needed.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:7b8c4a1be38d66f7548142e1 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Refusal classification now runs only in `Client.Evidence` before the
  review/thread fetch boundary. `classifyEvidence` no longer repeats the
  precedence rule or accepts comments. The recorded commit-status corpus now
  exercises `Client.Evidence`, including the already-fetched authoritative
  comments path.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test ./internal/reviewsource/coderabbit -count=1` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
