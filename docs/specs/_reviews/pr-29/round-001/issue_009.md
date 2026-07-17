---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/SKILL.md
line: 31
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYfu,comment:PRRC_kwDOS0qyts7WgS5U
review_hash: 06d7d4901b97585058f321faefac415317eb6cbdd7c98bdeac6862d4f116ad5a
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 009: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Prefix all executable examples with `rtk`.**

The new commands invoke `python3` directly at Lines 31, 78, 88, and 96, while the repository routes Python through `$(RTK)` and requires `rtk` prefixes when available. Agents following this skill will otherwise bypass the wrapper. Change each example to `rtk python3 ...`.

As per coding guidelines, shell commands must be prefixed with `rtk` when available.







Also applies to: 77-79, 87-89, 95-97

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/SKILL.md` around lines 30 - 31, Prefix every
executable Python example in the setup-context-driven skill, including the
commands near the context_setup.py invocations, with rtk so they use the
repository’s wrapper; update all four examples consistently without changing
their arguments or behavior.
```

</details>

<!-- fingerprinting:phantom:poseidon:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8dbcc41e3f3556472f8cd3d8 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The skill examples invoked `python3` directly. Updated all executable examples to use `rtk python3` without changing arguments.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
