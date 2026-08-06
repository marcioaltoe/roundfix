---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: invalid
terminal_reason: Spec 0073 satisfies the Archive Command's declared-only partial eligibility contract and must remain archived.
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_tasks.md
line: 4
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kaL,comment:PRRC_kwDOS0qyts7eJylG
review_hash: e3a8b5e121b6aa2087d6ee95d5f630aec1728dbe57af5684dc325f5d080eaaf6
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:10Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Major_ _ Heavy lift_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _🏗️ Heavy lift_

**Keep Spec 0073 active until QA passes.**

The PR objective reports a `partial` QA verdict. The archive policy requires completed tasks, passing QA, and self-containment checks before archival. Move this Spec out of `docs/specs/_archived/` until a new QA run reports `pass`. Preserve the existing QA evidence as historical evidence.

- `docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_tasks.md#L1-L4`: restore the task manifest to the active Spec location.
- `docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_techspec.md#L1-L5`: restore the technical specification to the active Spec location.

Based on learnings: “Archive only after tasks, QA, and self-containment checks pass; if any check fails, stop, report the offending item, and leave the spec active.”

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_tasks.md#L1-L4` (this comment)
- `docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_techspec.md#L1-L5`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_tasks.md`
around lines 1 - 4, Move Spec 0073 back to the active specs location because its
QA verdict is partial; restore both
docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_tasks.md
(lines 1-4) and _techspec.md (lines 1-5) to the corresponding active Spec
directory, preserving their existing QA evidence and content. Do not archive the
spec again until tasks, QA, and self-containment checks all pass.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_tasks.md</file>
<line_range>1-4</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_techspec.md</file>
<line_range>1-5</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:764b90f9c25df4e9c3fbfbc1 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `internal/spec/archive.go`, `internal/cli/archive_test.go`, and the canonical Roundfix Skill allow a partial report whose only unmet rows are declared unreachable. The Spec has all Tasks completed, zero environment- or finding-blocked rows, one declared-blocked row, and the declaration's satisfying action recorded under `unproven`. Restoring `_tasks.md` or `_techspec.md` to the active tree would reverse a valid archive.
