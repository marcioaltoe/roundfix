---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/cli/cli.go
line: 2015
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95I,comment:PRRC_kwDOS0qyts7V5tZ0
review_hash: 5c6b596708db72acd9410018f97ee258488b9ba459a9efed4c595601db8a0262
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 003: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Do not report per-Run fallback data as cumulative totals.**

When cumulative artifact loading fails, the current Run’s issues are silently printed as “Pull Request cumulative.” Report cumulative data as unavailable and emit the load diagnostic instead of publishing incorrect totals.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2009 - 2015, Update reviewIssueReportData
so a loadReviewIssueReportIssues failure does not substitute refreshedRunIssues
for cumulativeIssues. Mark cumulative data as unavailable using the existing
reviewIssueReport representation and emit the load error diagnostic, while
preserving refreshedRunIssues as the per-Run issue data.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:588c51a5be542ac54995ecac -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Cumulative Review Issue report load failures now print unavailable cumulative data and emit a stderr diagnostic; `make verify` passed.
