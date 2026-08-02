---
source: coderabbit
pr: "67"
round: 1
round_created_at: "2026-08-02T11:30:39Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/profiles-configure-merge-semantics
head_sha: ffcc15ebed0a055d329cb3215ae0878b90931948
file: internal/config/config_test.go
line: 895
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vwh2K,comment:PRRC_kwDOS0qyts7ceBLs
review_hash: 9d3571f3ec074143a053c39b310918ca8d9f9d06fab0eabd78f22fd1090dc8b0
duplicate_of: ""
source_review_id: "4838273774"
source_review_submitted_at: "2026-08-02T11:29:42Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Write the fixture at its target indentation instead of rewriting it.**

`strings.ReplaceAll(original, "    ", "  ")` rewrites every run of four spaces in the literal, not only the leading indentation. It produces the intended result today only because all indentation in the literal is a multiple of four. If a maintainer adds a scalar value that contains four consecutive spaces, the fixture is corrupted silently and the failure is hard to trace.

Declare the literal with 2-space indentation directly and delete the rewrite.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/config_test.go` around lines 894 - 895, Update the fixture
literal near the original assignment to use 2-space indentation directly, then
remove the strings.ReplaceAll transformation. Keep the existing strings.Replace
call that changes “backend-primary” to “backend-updated” unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:97dc90b679372a41150f1692 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The fixture-wide four-space replacement could alter scalar content unrelated to indentation. The literal now declares its intended two-space indentation directly; the targeted model replacement remains unchanged.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache rtk go test ./internal/config -run '^(TestProfilesConfigureRemovalPreservesSpacing|TestProfilesConfigureMergePreservesOtherCategories|TestEffectiveChangeSet)$' -count=1` passed 21 tests.
