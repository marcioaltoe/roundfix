---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: .agents/skills/write-tasks/SKILL.md
line: 26
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdF,comment:PRRC_kwDOS0qyts7gSdxi
review_hash: 77102866510b75b9f10789ec0572b5f7b261cc8e007c99fb99b814fe41547039
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Resolve archive paths from the configured Spec Root across all authoring contracts.**

The runtime contract maps the default `docs/specs` root to `_archived/specs`, but maps non-default roots to `<spec-root>/_archived`. These documents use only the default path, which can misclassify external archives and reuse Spec number prefixes.

- `.agents/skills/write-tasks/SKILL.md#L26-L26`: classify Specs using the resolved archive directory.
- `.agents/skills/write-techspec/SKILL.md#L23-L23`: scan active and archived directories under the effective Spec Root.
- `docs/agents/docs-layout.md#L190-L190`: document both default and non-default archive locations.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.5.1)</summary>

[warning] 20: [EA2] Autonomous Decision Making: Skill enables autonomous high-impact decisions without human-in-the-loop verification. Critical operations (destructive commands, financial transactions, data deletion) should require explicit user confirmation.

Remediation: Add human-in-the-loop confirmation for destructive, irreversible, or high-impact operations. Never auto-execute commands that modify files, send data, or alter system state.

(Excessive Agency (EA2))

</details>

</details>

<details>
<summary>📍 Affects 3 files</summary>

- `.agents/skills/write-tasks/SKILL.md#L26-L26` (this comment)
- `.agents/skills/write-techspec/SKILL.md#L23-L23`
- `docs/agents/docs-layout.md#L190-L190`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-tasks/SKILL.md at line 26, Update
.agents/skills/write-tasks/SKILL.md at line 26 to classify Specs using the
archive directory resolved from the configured Spec Root. Update
.agents/skills/write-techspec/SKILL.md at line 23 to scan active and archived
directories under the effective Spec Root. Update docs/agents/docs-layout.md at
line 190 to document both the default docs/specs-to-_archived/specs mapping and
the non-default <spec-root>/_archived location.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>.agents/skills/write-tasks/SKILL.md</file>
<line_range>26-26</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-techspec/SKILL.md</file>
<line_range>23-23</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/agents/docs-layout.md</file>
<line_range>190-190</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3c1ce5451a00ec99a2a528d3 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Updated all three affected files so the archive directory is resolved from the configured Spec Root rather than hard-coded to the default:
  - `.agents/skills/write-tasks/SKILL.md` (and `skills/` mirror): Project Constraint preflight step 1 classifies legacy Specs against the resolved archive directory (`_archived/specs` for the built-in `docs/specs` root, `<spec-root>/_archived/` otherwise).
  - `.agents/skills/write-techspec/SKILL.md` (and `skills/` mirror): the refactor-mint numbering rule now scans the configured Spec Root and its resolved archive directory.
  - `docs/agents/docs-layout.md`: documented both the default `docs/specs` → `_archived/specs/` mapping and the non-default `<spec-root>/_archived/` location.
  - Verified: `make skills-sync-check`, `roundfix skills check`, and the full `make verify` gate pass.
