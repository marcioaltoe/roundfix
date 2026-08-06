---
source: coderabbit
pr: "125"
round: 1
round_created_at: "2026-08-05T19:38:23Z"
status: invalid
terminal_reason: "ADR-0080 permits a pass with environment-blocked rows when each row has equivalent observed or supervised evidence, which this report records."
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 1384d928c1d0af4d0bca06e506eb8b2953f9f341
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md
line: 10
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WxMJQ,comment:PRRC_kwDOS0qyts7d790e
review_hash: 6f46369ddf26b74fe638bb391629f344ae44b84e1257f0795d7a30f2b6731766
duplicate_of: ""
source_review_id: "4868102567"
source_review_submitted_at: "2026-08-05T19:37:35Z"
---

# Issue 005: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Set the QA verdict to `partial`.**

Lines 121-125 record 11 unresolved environment-blocked rows. Line 133 declares `pass`. A `pass` verdict is not valid while permitted blocks remain unresolved. Update the frontmatter and final verdict text to `partial`.

As per coding guidelines, “partial” applies to unresolved permitted blocks, while “pass” requires all runnable rows and evidence requirements to be satisfied.






Also applies to: 121-135

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`
around lines 3 - 10, The QA report currently declares pass despite 11 unresolved
environment-blocked rows. Update the frontmatter status/verdict and the final
verdict text in the report to partial, preserving the recorded blocked-row
counts and other metadata.
```

</details>

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:bf8d08c7c139f0a6518f2c0e -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The repository's accepted QA policy distinguishes environment-blocked rows from finding-blocked, failed, or unevidenced rows. Environment-only blockers do not cap the verdict when equivalent observed or supervised evidence is recorded; the report declares 11 environment-blocked rows, zero finding/declared blockers, and provides the required equivalents.
- Evidence: `docs/adr/0080-qa-verdicts-distinguish-environment-blocked-rows.md` defines this verdict rule, and the QA report's frontmatter and matrix conform to it.
