---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/config/config_test.go
line: 3467
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiApW,comment:PRRC_kwDOS0qyts7fC8Rq
review_hash: 3596ec39d2f9f5979d9c142d139624e0966ff23bb64281257e9c3c9ede60325b
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:49Z"
---

# Issue 048: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Move these tests to `internal/config/profiles_test.go`.**

`ConfiguredWorkCategories` is defined in `internal/config/profiles.go`, but these new tests are in `internal/config/config_test.go`. Keep the source-to-test mapping clear.

As per coding guidelines, test files “MUST be named after the source file under test.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/config_test.go` around lines 3397 - 3467, Move
TestConfiguredWorkCategoriesCoversDefinedOptionalCategories and
TestConfiguredWorkCategoriesOmitsInheritedCategories from config_test.go into
profiles_test.go, keeping their implementations and assertions unchanged. This
preserves the required source-to-test mapping for ConfiguredWorkCategories,
defined in profiles.go.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:a007e3231c11f1a003e84013 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
