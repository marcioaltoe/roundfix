---
source: coderabbit
pr: "125"
round: 2
round_created_at: "2026-08-05T20:14:17Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: a89da452f019b880472c798f58529ea8aebefb1b
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md
line: 30
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WxxsT,comment:PRRC_kwDOS0qyts7d80ay
review_hash: 193a8f149fff079ab7ea829ee6210be4160da69e91c14a5c87b7471656c76306
duplicate_of: ""
source_review_id: "4868376070"
source_review_submitted_at: "2026-08-05T20:14:08Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Use the complete `.coderabbit.yaml` path for `auto_review`.**

The requirements use `auto_review.*`, but the same contract places the settings under `reviews.auto_review.*`. This ambiguity can make an implementation read root-level keys, fall back to defaults, and make the wrong Preflight decision.

- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md#L26-L30`: use `reviews.auto_review.enabled`, `reviews.auto_review.auto_incremental_review`, and `reviews.auto_review.auto_pause_after_reviewed_commits`.
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/_techspec.md#L132-L144`: use the same fully qualified paths in the API contract and predicate.

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md#L26-L30` (this comment)
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/_techspec.md#L132-L144`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md` around
lines 26 - 30, The repository configuration contract must use the fully
qualified reviews.auto_review paths. Update
docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md lines 26-30
and docs/specs/_archived/0078-roundfix-asks-for-the-review/_techspec.md lines
132-144 to reference reviews.auto_review.enabled,
reviews.auto_review.auto_incremental_review, and
reviews.auto_review.auto_pause_after_reviewed_commits, preserving the existing
defaults and predicate behavior.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md</file>
<line_range>26-30</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/_techspec.md</file>
<line_range>132-144</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d72ecd415b139d927c707597 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The implementation decodes the settings beneath
  `reviews.auto_review`, while Task 03 and the TechSpec predicate omitted the
  `reviews` parent. Updated both contract locations to use the three fully
  qualified `.coderabbit.yaml` paths without changing defaults or predicate
  behavior.
- Evidence: `internal/preflight/preflight.go` confirms the nested
  `reviews.auto_review` decoder shape. The guarded
  `go test ./internal/preflight ./internal/config -count=1 -run
  'Review|Request|Coheren|Config'` focused check exited `0`, and
  `rtk git diff --check` exited `0`. The Daemon owns authoritative
  `make verify`.
