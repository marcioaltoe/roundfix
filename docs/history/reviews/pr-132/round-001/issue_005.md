---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: invalid
terminal_reason: A partial report with only fully declared unreachable rows is archive-eligible under the current Roundfix contract.
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/qa/qa-report-2026-08-06-01.md
line: 9
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kaR,comment:PRRC_kwDOS0qyts7eJylN
review_hash: a0a14efeaa4666fc377ed16086b2aa66f0e4d213cfc3f2d0cf90053a9d490143
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:10Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Keep Spec 0073 active until QA returns `pass`.**

This report is already under `docs/specs/_archived/`, but its final verdict is `partial` with one declared-blocked row. Defer or reverse the archive operation until Q-15 is resolved and a fresh QA run returns `pass`.

Based on learnings, archive only after tasks, QA, and self-containment checks pass; if a check remains blocked, leave the Spec active.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/qa/qa-report-2026-08-06-01.md`
around lines 4 - 9, Keep Spec 0073 active rather than archived while the QA
report’s verdict remains partial and rows_blocked_declared is nonzero. Defer or
reverse the archive operation until Q-15 is resolved and a fresh QA run reports
pass, with tasks, QA, and self-containment checks completed.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:bcc231ea01d39a89868b8952 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The report's single blocked row is `blocked (declared: Success Metric 7)`, with `rows_blocked_environment: 0` and `rows_blocked_finding: 0`. The PRD carries the matching pre-run Unreachable Acceptance declaration and its `satisfied-by` action was stamped into `unproven`. The current Archive Command intentionally distinguishes this declared-only case from finding- or environment-blocked partial reports, which it refuses.
