---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: invalid
terminal_reason: Roundfix permits archival of a declared-only partial QA report and preserves each satisfying action under unproven.
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_prd.md
line: 9
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kaI,comment:PRRC_kwDOS0qyts7eJylC
review_hash: cb6f3eb496fe7e37050f1a677d2acc8391c7d09f800ea5e9f5c182f6103bf206
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:09Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Major_ _ Heavy lift_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _🏗️ Heavy lift_

**Do not archive the Spec while an acceptance criterion remains unproven.**

`status: archived` is set while `unproven` records Success Metric 7 as incomplete. The supplied PR summary also reports a `partial` QA verdict with one declared-unreachable row. Leave the Spec active until the missing authorization, bounded Task, QA, and self-containment checks are complete.

Based on learnings, archive only after all tasks, QA, and self-containment checks pass; if a check fails, leave the Spec active.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0073-skill-versions-decoupled-from-the-binary/_prd.md`
around lines 3 - 9, The Spec must remain active while the unproven acceptance
criterion and partial QA verdict remain. Update the front matter in _prd.md to
remove the archived status/date and preserve the unresolved requirement until
its authorization, bounded Task, QA, and self-containment checks pass; archive
it only after all checks succeed.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ed7abaeb0844927e0599fa67 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The current Archive Command contract explicitly accepts `partial` when every unmet row is covered by the Spec's pre-run `## Unreachable Acceptance` declaration. Spec 0073's newest report records `rows_blocked_environment: 0`, `rows_blocked_finding: 0`, and `rows_blocked_declared: 1`; the PRD declares Success Metric 7 unreachable and the archived frontmatter preserves its `satisfied-by` action under `unproven`. Moving it active would contradict the implemented and tested archive-eligibility policy.
