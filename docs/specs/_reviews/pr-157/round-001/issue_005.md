---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: docs/agents/docs-layout.md
line: 175
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdL,comment:PRRC_kwDOS0qyts7gSdxp
review_hash: 21a582e9ec2f53174cc72784adeef81cd0dea84074bd3eb648e7f59fd0635121
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Replace “license” with “value”.**

`absorbed_by` stores a Rollup basename or Spec slug. It is not a license. Use precise terminology so authors do not create or search for a separate license field.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/agents/docs-layout.md` around lines 171 - 175, In the archival Findings
guidance, update the description of the required absorbed_by field to call it a
value rather than a license, while retaining its requirement to resolve to an
active Rollup basename or Spec slug and preserving the YAML example.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1d8c54fd2f02901cf316b8af -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: In `docs/agents/docs-layout.md`, the archival Findings guidance now describes the required `absorbed_by:` field as a **value** that resolves to an active Rollup basename or a Spec slug, rather than calling it a "license", so authors do not create or search for a separate license field. The YAML example is unchanged. Verified: the full `make verify` gate passes.
