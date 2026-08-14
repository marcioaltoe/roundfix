---
source: coderabbit
pr: "125"
round: 4
round_created_at: "2026-08-05T20:44:41Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 535c9dd97cb583f418deeca1bc639b5030e5e728
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_06.md
line: 51
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WyExi,comment:PRRC_kwDOS0qyts7d9QPn
review_hash: 9cccaf1f75a8e73bcd67cc973e0094fa05c16f4f128e19230c65f8fac6e0f380
duplicate_of: ""
source_review_id: "4868508392"
source_review_submitted_at: "2026-08-05T20:33:32Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Validate the required QA observations, not only report headers.**

The commands accept any report with the correct `spec`, `verdict`, and blocked-count fields. They do not verify the required one-request-per-Round observation, both independent Preflight refusals, the `fetch` exemption, or the no-retry result. A stale or incomplete report can pass. Add checks for the required evidence rows and confirm that `task_05.md` settled `completed` before accepting the report.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0078-roundfix-asks-for-the-review/task_06.md` around
lines 47 - 51, Strengthen the QA validation commands in the report checks to
verify required evidence rows, not just headers: confirm one request per Round,
both independent Preflight refusals, the fetch exemption, and the no-retry
result, and require task_05.md to have settled as completed. Keep the existing
spec, verdict, and typed blocked-row validations while rejecting stale or
incomplete reports.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2021e78555948c7e390fbcf5 -->

_Source: Coding guidelines_

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Task 06 checked only report frontmatter. Its Verification contract now
  also requires Task 05 status `completed` and validates terminal Results rows
  R05-R07 and R10-R12 for both Preflight refusals, one request per Round, the
  artifact-only descendant, the `fetch` exemption, and the no-retry result.
- Focused evidence: the strengthened Ruby validator exited 0 against
  `qa-report-2026-08-05.md` and the completed Task 05 frontmatter. The Daemon
  owns `make verify`.
