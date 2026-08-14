---
source: coderabbit
pr: "125"
round: 3
round_created_at: "2026-08-05T20:34:12Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: a89da452f019b880472c798f58529ea8aebefb1b
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_06.md
line: 51
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WyExi,comment:PRRC_kwDOS0qyts7d9QPn
review_hash: e6066ec6e531ef7dd8b16364f9d6ffa0a0da2297dcb6c9954313b5248c032b57
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-125/round-004/issue_003.md
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

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
