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
line: 1308
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95H,comment:PRRC_kwDOS0qyts7V5tZv
review_hash: a1cf977a4a60e63f50bac8bb104e3139be2f8f452f1393893e96e118b3f82f6c
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 002: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Re-scan for older Active Runs after reclaiming the newest orphan.**

Bypass Runs can coexist without locks, while `ActiveReviewRunByTarget` returns only the newest nonterminal Run. Clearing `report.ActiveRun` after reclaim can therefore miss an older live Run and allow concurrent checkout mutation. Re-query until no Active Run remains, blocking on the first live owner.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 1293 - 1308, After reclaiming an orphaned
run in the report.ActiveRun handling flow, repeatedly re-query
ActiveReviewRunByTarget until no active run remains. Reclaim each newly
discovered orphaned run, but stop and block on the first live owner instead of
clearing report.ActiveRun prematurely. Preserve the existing store close and
error propagation behavior for every re-scan.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:381037d9956ac798a0674cbc -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Branch Integrity Preflight now re-queries Active Runs after each orphan reclaim and blocks on the first live owner; `make verify` passed.
