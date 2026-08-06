---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/0069-review-run-targets-its-pull-request/task_03.md
line: 71
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v6U,comment:PRRC_kwDOS0qyts7eEK7p
review_hash: a2ae8b1db861727125588fae57dc350274a568255e937575254c96462878749f
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 021: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Add effect assertions to the Task Verification commands.**

Both Tasks verify broad gates or artifact presence, but neither verifies the contract that its Task must deliver. A clean baseline or a stale artifact can satisfy the current commands.

- `docs/specs/0069-review-run-targets-its-pull-request/task_03.md#L63-L71`: assert each newly required Skill statement before checking synchronization and broad gates.
- `docs/specs/0069-review-run-targets-its-pull-request/task_04.md#L44-L49`: assert the successful QA verdict, closed rows, typed blocked-row counts, and evidence for each required observed behavior.

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/0069-review-run-targets-its-pull-request/task_03.md#L63-L71` (this comment)
- `docs/specs/0069-review-run-targets-its-pull-request/task_04.md#L44-L49`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0069-review-run-targets-its-pull-request/task_03.md` around lines
63 - 71, Add effect assertions to the Verification commands in
docs/specs/0069-review-run-targets-its-pull-request/task_03.md (lines 63-71),
asserting every newly required Skill statement before synchronization and
broad-gate checks. Also update
docs/specs/0069-review-run-targets-its-pull-request/task_04.md (lines 44-49) to
assert the successful QA verdict, closed rows, typed blocked-row counts, and
evidence for each required observed behavior.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/0069-review-run-targets-its-pull-request/task_03.md</file>
<line_range>63-71</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0069-review-run-targets-its-pull-request/task_04.md</file>
<line_range>44-49</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2db602945f630e4b455b852d -->

_Sources: Coding guidelines, Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
