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
line: 498
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYfd,comment:PRRC_kwDOS0qyts7WgS49
review_hash: bd7a4125d70091737e262a3c06eb995af1be66da8661b5f202ab54a42ae6eb08
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Audit installed skill content, not only its directory name.**

Setup snapshots store `contentDigest`, but a required skill passes solely when its name exists. A stale or substituted `SKILL.md` therefore satisfies the audit while exposing different capabilities than the canonical setup.

Compare installed content with the snapshot digest and report drift.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/scripts/context_setup.py` around lines 476 - 498,
Update the required-skill validation in context setup to compare each installed
skill’s SKILL.md content against the setup snapshot’s contentDigest, rather than
accepting presence based only on skill name. Use the existing installed-skill
metadata and snapshot fields, and report a finding for digest mismatches while
preserving the current missing-skill finding for absent required skills.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5df9b5056d77395424fb541f -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Required skill audit only checked names. Added installed `SKILL.md` digest checks against selected setup snapshot digests when the canonical local skill file is available, and refreshed snapshot digests.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
