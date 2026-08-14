---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/scripts/context_setup.py
line: 2019
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYfr,comment:PRRC_kwDOS0qyts7WgS5O
review_hash: 837f871a351a98655a5499a3b80f12eb693480d38c98b5cf6b79708536b9a5f9
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 008: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Validate `skills` as a list before iterating it.**

`"skills": null` raises `TypeError`, while strings or objects are iterated into invalid pseudo-entries before the later type check. Return a setup-snapshot validation finding immediately when the field is not a list.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/scripts/context_setup.py` around lines 1985 -
2019, Validate source_doc["skills"] as a list before the raw_skill loop in the
setup snapshot normalization flow. When the field is missing or not a list,
immediately append the existing setup_snapshot_invalid_finding for the current
source_path and setup_id and return the appropriate invalid result, preventing
iteration over null, strings, or objects; preserve the existing per-skill
normalization for valid lists.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:33039d95ccc30e4cbb913e1c -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Setup snapshot normalization iterated `skills` before validating it. Added early non-empty list validation that returns a controlled finding.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
