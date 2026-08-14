---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: .agents/skills/write-tasks/SKILL.md
line: 93
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiAoS,comment:PRRC_kwDOS0qyts7fC8Qi
review_hash: 84d554b00507dbbef750aa0d204569196afe2bc36a7df21e8ffc32d051702405
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:48Z"
---

# Issue 027: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Align outside-evidence behavior across task authoring and QA.**

When the external source is unavailable, the task-authoring skill says the row never stalls the Spec, but the QA skill makes the row `partial` without equivalent evidence, and `partial` blocks PR preparation. Clarify that no human interaction is required during decomposition while preserving the QA block, or change both contracts together.

- `.agents/skills/write-tasks/SKILL.md#L85-L93`: replace the “never stalls the Spec” claim with the precise authoring and QA behavior.
- `.agents/skills/qa-gate/SKILL.md#L129-L137`: retain the blocked-row evidence requirement and state its effect on the final verdict.

Based on the supplied task-authoring and QA contracts, these two outcomes currently conflict.

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
<summary>📍 Affects 2 files</summary>

- `.agents/skills/write-tasks/SKILL.md#L85-L93` (this comment)
- `.agents/skills/qa-gate/SKILL.md#L129-L137`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-tasks/SKILL.md around lines 85 - 93, The
external-evidence contracts conflict: update
`.agents/skills/write-tasks/SKILL.md` lines 85-93 to state that decomposition
requires no human interaction, but unavailable evidence is recorded for QA
rather than treated as never stalling the Spec; update
`.agents/skills/qa-gate/SKILL.md` lines 129-137 to retain the blocked-row
evidence requirement and explicitly state its effect on the final verdict,
including that the blocked or partial row prevents PR preparation.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>.agents/skills/write-tasks/SKILL.md</file>
<line_range>85-93</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/qa-gate/SKILL.md</file>
<line_range>129-137</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:aea156fdd82d1c547995d148 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Aligned the outside-evidence contracts across both skills. `write-tasks/SKILL.md` lines 85-93 now state that decomposition proceeds without human interaction when the source is unavailable, but the row is recorded as blocked; the blocked row carries into QA and blocks PR preparation. `qa-gate/SKILL.md` lines 129-137 now explicitly state that a blocked or partial outside-evidence row prevents PR preparation until satisfied or carried forward under ADR-0097. Both skills now present a consistent contract: decomposition does not stall, QA records the block, and the blocked row gates the PR.
