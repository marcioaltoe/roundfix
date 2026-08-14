---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/preservation.go
line: 556
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymi,comment:PRRC_kwDOS0qyts7cjgE-
review_hash: 445fe35f22d43915321d8ab98edd47dcc1072b03ef3d4c1996c853182833cb59
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 012: _ Performance & Scalability_ _ Trivial_ _ Quick win_

## Review Comment

_🚀 Performance & Scalability_ | _🔵 Trivial_ | _⚡ Quick win_

**Carrier classification repeats the same filesystem work several times per plan.**

`classifyCarriers` opens the repository root, reads every bounded carrier, hashes its bytes, and parses the Setup Manifest. `buildPlanWithCatalog` in `internal/baseline/plan.go` already calls `classifyCarriers` at lines 325-328 and again at lines 437-440, and it calls `planRootPreservationWithCatalog` twice (lines 373 and 447). With `classifyCarriers: true`, the same classification therefore runs up to four times for one `baseline plan` invocation.

Compute the classification once per snapshot and pass the result into `planRootPreservationWithCatalog` through `RootPreservationRequest`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/preservation.go` around lines 546 - 556, Compute carrier
classifications once in buildPlanWithCatalog for each snapshot when
classifyCarriers is enabled, then pass the cached result through
RootPreservationRequest to both planRootPreservationWithCatalog calls. Update
the request and preservation flow to reuse that result, and use it for warnings
instead of invoking classifyCarriers again in preservation.go.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:3a91f4512664f19a554feff1 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `buildPlanWithCatalog` now computes carrier classifications once for the initial snapshot and once for the refreshed snapshot, then passes each cached result through `RootPreservationRequest`. The full Baseline package test passed.
