---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: invalid
terminal_reason: The files have one required top-level version and one distinct nested metadata.version field, not duplicate YAML keys.
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: .agents/skills/archive-spec/SKILL.md
line: 11
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kaC,comment:PRRC_kwDOS0qyts7eJyk3
review_hash: 064a5a50123010ba36b7797fc7c04ed2763a51ac65fb356dd1a7e793ce56b9e3
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:09Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Remove duplicate skill frontmatter version keys.**

Each affected skill contains an added `version: 0.0.2` entry that duplicates an existing key. Keep exactly one canonical `version` key in each file to avoid parser rejection or parser-dependent metadata resolution.

- `.agents/skills/archive-spec/SKILL.md#L11-L11`: remove the duplicate `version` entry.
- `.agents/skills/implement-spec/SKILL.md#L12-L12`: remove the duplicate `version` entry.
- `.agents/skills/implement-task/SKILL.md#L10-L10`: remove the duplicate `version` entry.

<details>
<summary>📍 Affects 3 files</summary>

- `.agents/skills/archive-spec/SKILL.md#L11-L11` (this comment)
- `.agents/skills/implement-spec/SKILL.md#L12-L12`
- `.agents/skills/implement-task/SKILL.md#L10-L10`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md at line 11, Remove the duplicate
version entry, retaining exactly one canonical version key, in
.agents/skills/archive-spec/SKILL.md lines 11-11,
.agents/skills/implement-spec/SKILL.md lines 12-12, and
.agents/skills/implement-task/SKILL.md lines 10-10.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>.agents/skills/archive-spec/SKILL.md</file>
<line_range>11-11</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/implement-spec/SKILL.md</file>
<line_range>12-12</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/implement-task/SKILL.md</file>
<line_range>10-10</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e0b8f91d14f9027eb3113050 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/implement-spec/SKILL.md`, and `.agents/skills/implement-task/SKILL.md` each contain exactly one `version` at the YAML document root. The other occurrence is `metadata.version`, a distinct nested field used by skill metadata. Removing either field would weaken one of the two contracts rather than fix a duplicate key.
