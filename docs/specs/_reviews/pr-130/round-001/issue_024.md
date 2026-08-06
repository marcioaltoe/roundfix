---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: internal/speccheck/citations_test.go
line: 123
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v6b,comment:PRRC_kwDOS0qyts7eEK7z
review_hash: fecd857e8765d6f9cb8a56ea70e46d4eca5e42a6aed2f65a85ff1d841a745340
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:29Z"
---

# Issue 024: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Return a valid archived Spec input.**

Line 119 returns `"_archived/<name>"`. `speccheck.Check` rejects that value because it accepts only a basename slug. If no active Spec exists, this test fails before it checks loop-order consistency.

Return the selected Spec root separately and pass the archived root with `entry.Name()`, or remove the archived fallback.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/citations_test.go` around lines 99 - 123, The
anyCheckableSpecSlug helper currently returns an archived path that
speccheck.Check rejects; update its archived fallback to return only
entry.Name() as the slug, while still locating the Spec under _archived.
Preserve the active Spec behavior and fallback selection order.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2462cd2bcdaf943aeda85fdb -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
