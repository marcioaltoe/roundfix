---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: invalid
terminal_reason: "The active Task paths were correct when Verification ran; archive relocation does not make the historical declarations live again."
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md
line: 76
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v53,comment:PRRC_kwDOS0qyts7eEK7G
review_hash: b5d46441c64237858c79d2c12194e9d292dcc3ea4945b61b19b8b4c93b05d789
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 011: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Use exact `DERIVED_DIGEST_PATHS` and archived Task paths.**

Both allowlists permit broad Baseline directories and name the active Spec path instead of the archived Task path.
- `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md#L75-L76`: use the exact derived-file list and `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md`.
- `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md#L76-L77`: use the exact derived-file list and `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md`.

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md#L75-L76` (this comment)
- `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md#L76-L77`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md`
around lines 75 - 76, Update the changed-file allowlist checks in
docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md lines
75-76 and task_05.md lines 76-77 to use the exact DERIVED_DIGEST_PATHS entries
rather than broad internal/baseline directories, and reference each file’s
corresponding archived Task path instead of the active docs/specs path.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md</file>
<line_range>75-76</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md</file>
<line_range>76-77</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8a5f2b7d06d04cd298d6a9b0 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The declared paths named the active Spec location that existed when
  each Task ran. Archiving later relocated the historical files; the Tasks are
  no longer executable and must not be rewritten to pretend they ran from the
  archived path. ADR-0081 separately owns the computed digest fallout.
- Daemon Verification: `make verify` not run; Daemon-owned.
