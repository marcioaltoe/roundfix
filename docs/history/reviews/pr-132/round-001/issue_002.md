---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: .agents/skills/roundfix/SKILL.md
line: 10
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kaG,comment:PRRC_kwDOS0qyts7eJyk9
review_hash: abeded13b48f4b09086b0daaee763f2a9a827c0c415b33880f05e6fd32b52c71
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:09Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Keep the owned-skill version contract unambiguous and enforce it.**

The seven skill files add duplicate top-level `version` keys, while `Makefile` validates only that at least one key exists. Keep one declaration in every skill and require exactly one match in the validator.

- .agents/skills/roundfix/SKILL.md#L10-L10: remove the duplicate `version` entry.
- .agents/skills/qa-gate/SKILL.md#L10-L10: remove the duplicate `version` entry.
- .agents/skills/setup-context-driven/SKILL.md#L11-L11: remove the duplicate `version` entry.
- .agents/skills/write-idea/SKILL.md#L11-L11: remove the duplicate `version` entry.
- .agents/skills/write-prd/SKILL.md#L11-L11: remove the duplicate `version` entry.
- .agents/skills/write-tasks/SKILL.md#L11-L11: remove the duplicate `version` entry.
- .agents/skills/write-techspec/SKILL.md#L11-L11: remove the duplicate `version` entry.
- Makefile#L181-L190: count top-level `version` matches and fail unless the count is exactly one.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.5.1)</summary>

[warning] 25: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[error] 1614: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[error] 1614: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[warning] 1090: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>📍 Affects 8 files</summary>

- `.agents/skills/roundfix/SKILL.md#L10-L10` (this comment)
- `.agents/skills/qa-gate/SKILL.md#L10-L10`
- `.agents/skills/setup-context-driven/SKILL.md#L11-L11`
- `.agents/skills/write-idea/SKILL.md#L11-L11`
- `.agents/skills/write-prd/SKILL.md#L11-L11`
- `.agents/skills/write-tasks/SKILL.md#L11-L11`
- `.agents/skills/write-techspec/SKILL.md#L11-L11`
- `Makefile#L181-L190`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md at line 10, Remove the duplicate top-level
version declaration from .agents/skills/roundfix/SKILL.md#L10-L10,
.agents/skills/qa-gate/SKILL.md#L10-L10,
.agents/skills/setup-context-driven/SKILL.md#L11-L11,
.agents/skills/write-idea/SKILL.md#L11-L11,
.agents/skills/write-prd/SKILL.md#L11-L11,
.agents/skills/write-tasks/SKILL.md#L11-L11, and
.agents/skills/write-techspec/SKILL.md#L11-L11, preserving exactly one version
per skill. Update the validator in Makefile#L181-L190 to count top-level version
matches and fail unless the count is exactly one.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>.agents/skills/roundfix/SKILL.md</file>
<line_range>10-10</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/qa-gate/SKILL.md</file>
<line_range>10-10</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/setup-context-driven/SKILL.md</file>
<line_range>11-11</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-idea/SKILL.md</file>
<line_range>11-11</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-prd/SKILL.md</file>
<line_range>11-11</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-tasks/SKILL.md</file>
<line_range>11-11</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-techspec/SKILL.md</file>
<line_range>11-11</line_range>
</site>
<site>
<role>sibling</role>
<file>Makefile</file>
<line_range>181-190</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a665eff6f923674d181a7cb7 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The claimed duplicate declarations are at distinct paths (`metadata.version` and top-level `version`), so no Skill frontmatter was removed. The validator gap was valid: `Makefile` now counts every top-level `version:` key and every non-empty value and fails unless both counts equal one. Focused evidence: `rtk make skills-version-check` exited 0, and `rtk make skills-sync-check` exited 0 with four selected Skill contract tests passing. Authoritative `make verify` remains Daemon-owned.
