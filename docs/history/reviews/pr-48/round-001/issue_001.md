---
source: coderabbit
pr: "48"
round: 1
round_created_at: "2026-07-29T21:58:48Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/repository-derived-skill-requirements
head_sha: 3ef6a563f8be4a4e72a2a063463d904fd0e0a9a1
file: .agents/skills/roundfix/SKILL.md
line: 71
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6U6ury,comment:PRRC_kwDOS0qyts7bPmQC
review_hash: 147530dee50eaaa3dc2b85cdf6c468dc53aa9c8b358eb34d939302d77628c82a
duplicate_of: ""
source_review_id: "4813239038"
source_review_submitted_at: "2026-07-29T21:57:47Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Document the generic external fallback remediation.**

This section documents only the skill-scoped `bunx skills add marcioaltoe/skills@<skill>` remediation for *named* external gaps. It omits `externalSkillsNextAction` (`bunx skills experimental_install && bunx skills update -p -y`), which Doctor prints instead whenever a failure doesn't identify specific skill names (e.g. an unreadable/symlinked `skills-lock.json` or an unclassified checker error — see `doctorExternalSkillNextActions`'s fallback in `internal/cli/doctor.go`). A user seeing that command with no explanation in the docs may be confused about when/why it appears versus the per-skill command.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.4.4)</summary>

[warning] 24: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[warning] 912: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md around lines 66 - 71, Update the
remediation documentation in SKILL.md to include the generic external fallback
command from externalSkillsNextAction, `bunx skills experimental_install && bunx
skills update -p -y`. Explain that Doctor prints this fallback when external
skill failures do not identify specific skill names, while named gaps use
skill-scoped commands; preserve the existing owned and mixed-failure remediation
guidance.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:79ea10b9cceeb0c6a374205e -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - The Roundfix Skill documented named external gaps but omitted Doctor's existing generic `externalSkillsNextAction` branch for external failures without skill names.
  - Updated the canonical and distributed Roundfix Skill to distinguish named skill-scoped remediation from `bunx skills experimental_install && bunx skills update -p -y`, while preserving owned and mixed named failure guidance.
  - `rtk make baseline-digests`: passed; regenerated the deterministic Baseline pins derived from the authorized Roundfix-owned Skill edit.
  - `rtk make skills-sync-check`: passed.
  - `rtk proxy cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`: passed.
  - Daemon Verification `make verify` was not run by this Agent; the Daemon owns authoritative Verification after this turn.
