---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/specs/0084-an-update-that-can-run/_prd.md
line: 63
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiAoq,comment:PRRC_kwDOS0qyts7fC8Q8
review_hash: 581f5d095870af06c22342fdf882be65032266b14e4fe6c91c2e06b6674d26bd
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:48Z"
---

# Issue 036: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Align the protected tooling count with the bounded paths.**

Both Project Constraints sections say that two repo-owned authorial skills are mutated, but each lists three skill files: `write-tasks`, `qa-gate`, and `roundfix`. This makes the authorization scope ambiguous.

As per coding guidelines, protected tooling work must name an exact, consistent bounded file set.
- `docs/specs/0084-an-update-that-can-run/_prd.md#L45-L63`: change the count to three, or remove the unauthorized path.
- `docs/specs/0084-an-update-that-can-run/_techspec.md#L43-L58`: apply the same correction.

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/0084-an-update-that-can-run/_prd.md#L45-L63` (this comment)
- `docs/specs/0084-an-update-that-can-run/_techspec.md#L43-L58`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0084-an-update-that-can-run/_prd.md` around lines 45 - 63, Align
the protected tooling count with the listed bounded skill paths in the Project
Constraints sections: update the count from two to three in
docs/specs/0084-an-update-that-can-run/_prd.md lines 45-63 and apply the same
correction in docs/specs/0084-an-update-that-can-run/_techspec.md lines 43-58,
or remove the unauthorized skill path consistently from both lists.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/0084-an-update-that-can-run/_prd.md</file>
<line_range>45-63</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0084-an-update-that-can-run/_techspec.md</file>
<line_range>43-58</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0c949217755a50b0c9320966 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Updated the protected tooling count from "two" to "three" in both `_prd.md` line 45 and `_techspec.md` line 43. The bounded files list correctly names three repo-owned authorial skills (`write-tasks`, `qa-gate`, and `roundfix`). The count now matches the listed files consistently.
