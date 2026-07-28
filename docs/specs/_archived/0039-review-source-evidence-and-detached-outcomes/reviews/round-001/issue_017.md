---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: invalid
terminal_reason: "The skill already separates standalone and Daemon-assigned execution, and this Spec does not authorize changing the cited protected tooling path."
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: .agents/skills/implement-task/SKILL.md
line: 39
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6OB,comment:PRRC_kwDOS0qyts7aoLQy
review_hash: 7a5cec264c30ac995584cb65f1b3b20c5f3f0ee6710d3b49e9284bb08cebe604
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:08Z"
---

# Issue 017: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Route Task settlement through the Daemon.**

These instructions let an Agent change Task status, run declared Verification, settle the Task, and commit. That conflicts with the Daemon-owned verification and settlement contract, allowing competing terminal evidence.







As per coding guidelines, “Never run commands from a Task's `## Verification` section as the Agent; the Daemon owns authoritative Verification and Task settlement.”


Also applies to: 77-84, 107-139

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.4.4)</summary>

[warning] 147: [EA1] Unrestricted Tool Access: Skill grants unrestricted tool access without appropriate constraints. An agent with unfettered tool access can perform arbitrary actions including file modification, network requests, and code execution.

Remediation: Restrict tool access to only the tools required for the skill's stated purpose. Use an explicit allowlist rather than granting blanket access.

(Excessive Agency (EA1))

---

[warning] 3: [EA2] Autonomous Decision Making: Skill enables autonomous high-impact decisions without human-in-the-loop verification. Critical operations (destructive commands, financial transactions, data deletion) should require explicit user confirmation.

Remediation: Add human-in-the-loop confirmation for destructive, irreversible, or high-impact operations. Never auto-execute commands that modify files, send data, or alter system state.

(Excessive Agency (EA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/implement-task/SKILL.md around lines 31 - 39, Update the
standalone and Roundfix execution guidance in the task workflow instructions so
the Daemon exclusively owns declared Verification execution, authoritative
evidence, Task-status updates, settlement, and commits. Ensure Agents never run
commands from a Task’s ## Verification section or claim a terminal verdict; they
should only perform focused checks and record implementation evidence in ##
Result for Daemon handoff. Apply the same contract consistently to the
referenced sections.
```

</details>

<!-- fingerprinting:phantom:poseidon:terra -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d2cef2b25aa54d451706d2c4 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The cited skill explicitly distinguishes standalone mode, where the user invokes one Task directly, from Daemon-assigned mode, where the Agent must not edit status, run declared Verification, claim a terminal verdict, or commit. The broader rewrite would incorrectly remove the supported standalone workflow and would modify `.agents/skills/implement-task/SKILL.md`, which is outside this Spec's protected-tooling authorization.
