---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: invalid
terminal_reason: "DERIVED_DIGEST_PATHS is an authoritative Makefile constant, and ADR-0081 requires computed fallout rather than copied per-Spec path lists."
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/0069-review-run-targets-its-pull-request/_techspec.md
line: 71
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v6L,comment:PRRC_kwDOS0qyts7eEK7g
review_hash: c8718d91a026b69c0a5c04cbb2c5f9aabd8a0fc09fd80ae6e8bf00a0bba35791
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 018: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Replace `DERIVED_DIGEST_PATHS` with an exact authorized file list.**

The tooling scope permits digest fallout through an undefined symbolic path set. A fresh agent cannot determine which repository files are authorized to change, so the boundary is not enforceable.

- `docs/specs/0069-review-run-targets-its-pull-request/_techspec.md#L61-L71`: list every authorized digest-fallout path explicitly.
- `docs/specs/0069-review-run-targets-its-pull-request/task_03.md#L38-L40`: copy the same exact repository-relative allowlist into the Task scope and acceptance criteria.

As per coding guidelines: tooling authorization must name “exact bounded repository-relative files.”

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/0069-review-run-targets-its-pull-request/_techspec.md#L61-L71` (this comment)
- `docs/specs/0069-review-run-targets-its-pull-request/task_03.md#L38-L40`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0069-review-run-targets-its-pull-request/_techspec.md` around
lines 61 - 71, Replace DERIVED_DIGEST_PATHS in
docs/specs/0069-review-run-targets-its-pull-request/_techspec.md lines 61-71
with an explicit list of every authorized digest-fallout repository-relative
file. Copy the identical allowlist into
docs/specs/0069-review-run-targets-its-pull-request/task_03.md lines 38-40,
covering only the bounded Roundfix Skill paths permitted by the authorization
and preserving the CLI-synchronisation scope.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/0069-review-run-targets-its-pull-request/_techspec.md</file>
<line_range>61-71</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0069-review-run-targets-its-pull-request/task_03.md</file>
<line_range>38-40</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e6c4776fd3718fdd51e67aa0 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `Makefile` defines `DERIVED_DIGEST_PATHS`, and
  `internal/baseline/derived_ownership_test.go` reads and validates that
  declaration mechanically. ADR-0081 makes regenerated pins computable fallout
  of the authorized Roundfix Skill edit and explicitly rejects widening the
  authorization boundary to those outputs. Copying a static list into the Spec
  and Task would duplicate an authoritative constant and drift when it changes.
- Daemon Verification: `make verify` not run; Daemon-owned.
