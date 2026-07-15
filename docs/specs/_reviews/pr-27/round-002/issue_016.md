---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: skills/qa-gate/SKILL.md
line: 36
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ959,comment:PRRC_kwDOS0qyts7V5ta3
review_hash: a0eca4b582577392ca3592f927c7193df0066a63f2ef7ac15aa4f279285dd989
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 016: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Use a collision-safe filename for multiple same-day QA runs.**

The fixed `qa-report-YYYY-MM-DD.md` path cannot both create a new report and preserve an existing report for another build or scope on the same date. Add a sequence, run slug, or build identifier to new filenames.

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

In `@skills/qa-gate/SKILL.md` at line 36, Update the QA report creation flow to
use a collision-safe filename that includes a sequence, run slug, or build
identifier in addition to the date. Preserve the existing resume behavior for
matching same-day in-progress reports, while creating a distinct report when the
build or scope differs and retaining older reports.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7d111c79c2a53c589d0fa2dd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Updated `qa-gate` instructions to use collision-safe QA Report filenames while preserving same-build/scope resume behavior; `make verify` passed.
