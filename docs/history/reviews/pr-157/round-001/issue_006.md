---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md
line: 88
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdO,comment:PRRC_kwDOS0qyts7gSdxr
review_hash: 9332683d43b76da784c7c0c8c5597f176534d0de1da5ecd82292154e049ff859
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 006: _ Security & Privacy_ _ Major_ _ Heavy lift_

## Review Comment

_🔒 Security & Privacy_ | _🟠 Major_ | _🏗️ Heavy lift_

**Bound the generated Skill mirrors explicitly.**

This record names six `.agents/skills` files but authorizes their `skills/` mirrors without listing the mirror paths. An active authorization must let postflight validation identify every allowed path. Enumerate each corresponding mirror path or enforce an exact one-to-one derivation before this record authorizes the change.

Based on learnings: active authorization records must list exact bounded repository-relative file paths; directory globs are not sufficient.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`
around lines 84 - 88, Update the authorization record’s generated Skill mirror
list to explicitly enumerate the six corresponding skills/ paths, or enforce an
exact one-to-one path derivation before authorizing the change. Ensure
postflight validation can identify every allowed repository-relative path
without relying on an unbounded directory reference.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:752ad47e5f01b4bda53c0545 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Updated `docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md` to enumerate each of the six `.agents/skills/` sources together with its matching `skills/` mirror path (six one-to-one pairs) instead of a directory-glob-style reference. Postflight validation can now identify every allowed repository-relative path: `skills/{archive-spec,roundfix,write-idea,write-prd,write-tasks,write-techspec}/SKILL.md` alongside the corresponding `.agents/skills/` files. Each pair is still bounded to the archive-layout subject. Verified: `make skills-sync-check` passes (mirrors match canonical sources).
