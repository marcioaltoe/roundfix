---
source: coderabbit
pr: "32"
round: 2
round_created_at: "2026-07-17T13:23:47Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: d7ab1933ac9fdcf0c94d73e2f417d99d38e43fe7
file: internal/store/agent_selection.go
line: 118
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tL,comment:PRRC_kwDOS0qyts7Wt95p
review_hash: 01d39a800183c1b74710d33ca52f4906adeee669300acb1d0842a950b049db59
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 013: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Allow lifecycle transitions on the same selection attempt.**

After attempt 1 is stored as `attempting`, `ensureNextAgentSelectionAttempt` requires attempt 2, so `active`, `failed`, or `closed` cannot subsequently be persisted for attempt 1. Treat the same attempt number as a status update with immutable selection fields verified, while reserving `previous+1` for a new candidate. Add transition tests for `attempting → active → closed` and `attempting → failed`.







Also applies to: 338-357

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/agent_selection.go` around lines 105 - 118, The
ensureNextAgentSelectionAttempt flow must allow lifecycle status updates for an
existing attempt number while preserving immutable selection fields. Update
ensureNextAgentSelectionAttempt and the surrounding append logic so the same
attempt number validates unchanged selection data for attempting → active →
closed or attempting → failed transitions, while requiring previous+1 only for a
new candidate; add tests covering both transition sequences.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:94e3afb9e0a14324dfc38c27 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Fixed Agent Selection persistence to treat same-attempt lifecycle changes as immutable-field-checked status updates while reserving `previous+1` for new candidates. Evidence: `GOCACHE=/private/tmp/roundfix-go-build rtk go test ./internal/agent ./internal/cli ./internal/config ./internal/daemon ./internal/releaseplan ./internal/spec ./internal/store ./internal/tui` passed.
