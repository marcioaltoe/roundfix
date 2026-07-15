---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: skills/qa-gate/SKILL.md
line: 51
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ96E,comment:PRRC_kwDOS0qyts7V5ta_
review_hash: cf34a56844264e85b7404eb14b112dc0243319625f2149d4d20ce3dd5ea1ab51
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_017.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 017: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Require the full `make verify` gate.**

Allowing unspecified “equivalents” permits closure without the repository’s required skills checks and build validation. Require `make verify`; if it cannot run, record the gate as blocked rather than substituting partial checks.





As per coding guidelines, “Run the full `make verify` gate before claiming completion; any formatting, test, or build failure is blocking.”

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.3.11)</summary>

[warning] 74: [MP2] Context Window Stuffing: Skill attempts to fill the context window with filler content, displacing legitimate instructions and safety constraints. This can degrade agent performance or bypass safety boundaries.

Remediation: Implement context-window management that detects and rejects padding or stuffing attempts. Prioritize system instructions over user-injected content.

(Memory Poisoning (MP2))

---

[warning] 74: [MP2] Context Window Stuffing: Skill attempts to fill the context window with filler content, displacing legitimate instructions and safety constraints. This can degrade agent performance or bypass safety boundaries.

Remediation: Implement context-window management that detects and rejects padding or stuffing attempts. Prioritize system instructions over user-injected content.

(Memory Poisoning (MP2))

---

[warning] 77: [MP2] Context Window Stuffing: Skill attempts to fill the context window with filler content, displacing legitimate instructions and safety constraints. This can degrade agent performance or bypass safety boundaries.

Remediation: Implement context-window management that detects and rejects padding or stuffing attempts. Prioritize system instructions over user-injected content.

(Memory Poisoning (MP2))

---

[error] 88: [P6] Direct Prompt Extraction: Skill contains instructions that could directly expose system prompts, internal rules, or hidden instructions to users or external parties.

Remediation: Remove any instructions that reveal, print, or output system prompts or internal rules. System instructions should never be exposed to end users.

(System Prompt Leakage (P6))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/qa-gate/SKILL.md` at line 51, Update the verification instruction in
SKILL.md to require running the full `make verify` command before claiming
completion. Remove the allowance for build, lint, typecheck, or test
equivalents; if `make verify` cannot run or any formatting, test, or build check
fails, record the verification gate as blocked.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:28bd8308fae5ab69f0772ee7 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
