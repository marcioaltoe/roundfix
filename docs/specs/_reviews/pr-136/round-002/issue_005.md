---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/findings/2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md
line: 61
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5YW,comment:PRRC_kwDOS0qyts7eY0i_
review_hash: d3dd50fd3f3134eda9242175204c914345f7a245c86c53a44772e01595baf356
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Keep the Finding evidence-only.**

Lines 51-61 prescribe the fix and assign it to Spec 0080. That content is intent, not immutable evidence. Move the recommendation to a Backlog Entry or the Spec workflow. Keep this Finding limited to the observed corpus regression and its root cause.

Based on learnings, Findings record immutable evidence and never commitments; Backlog Entries record intended work.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/findings/2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md`
around lines 51 - 61, Remove the prescriptive “The shape of a fix” section,
including the proposed detector scoping and assignment to Spec 0080, from this
Finding. Retain only the observed corpus regression and its root cause; move any
recommended work to the appropriate Backlog Entry or Spec workflow.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ec12509a4ae4226453f35ccd -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The repository's Findings Operational Contract explicitly gives each
  finding an `Action / suggestion` field and permits routing to a Spec or
  direct fix. This section records a suggested detector scope and points to the
  already relevant Spec; it does not change the observed corpus regression or
  declare implementation complete. Deleting it would also conflict with the
  rule that Finding history is immutable and later routing evidence is
  appended rather than removed.
- Focused evidence: independently compared the section with
  `docs/agents/docs-layout.md`'s Findings template, routing convention, and
  immutable-history clauses. No product change is required.
- Daemon Verification: `make verify` not run; Daemon-owned.
