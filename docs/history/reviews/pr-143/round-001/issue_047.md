---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/baseline/plan_test.go
line: 2921
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiApS,comment:PRRC_kwDOS0qyts7fC8Rn
review_hash: 62c4a4d0c8f31d04911056818f5a464723b7ffbfe5ed44303d1bd225926ad1bc
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:49Z"
---

# Issue 047: _ Maintainability & Code Quality_ _ Trivial_ _ Poor tradeoff_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚖️ Poor tradeoff_

**The test reads unexported catalog internals.**

Line 2866 walks `catalog.modules[moduleID]["rules"]` directly. The assertion depends on the internal catalog map shape rather than a resolved public contract. A later catalog refactor breaks this test without any behavior change.

Prefer an existing accessor that returns resolved clauses. Keep the guidance and enforcement assertions unchanged.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan_test.go` around lines 2818 - 2921, Update
TestStandardTypeScriptStructuralClauseRetention to obtain clauses through the
existing catalog accessor that returns resolved clauses, replacing the direct
catalog.modules[moduleID]["rules"] traversal. Preserve the currentClauses
construction and leave the guidance and enforcement assertions unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:7029754f0cc48b2b0511a3f2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
