---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: internal/cli/detach.go
line: 307
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95W,comment:PRRC_kwDOS0qyts7V5taK
review_hash: 042d9be0581d50a793feb53eec0ad8649958e1669b11bd9a8c7f1676073205e4
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_008.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 008: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Do not treat an unfinished child wait as exit status 0.**

`waitDetachedProcess` returns `(nil, false)` when the child has not exited, but the boolean is discarded. The command then claims the child was killed with `exit status 0`, while the `cmd.Wait` goroutine may remain blocked. Retain one owned wait channel and confirm shutdown after killing before reporting the result.





As per coding guidelines, “Every goroutine must have an owner, cancellation, and a clear shutdown path; no fire-and-forget goroutines.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/detach.go` around lines 298 - 307, Update
handleDetachedHandshakeTimeout to retain and use the completion signal returned
by waitDetachedProcess after killDetachedProcess. Confirm the child has actually
exited before reporting its result; if the wait remains unfinished, preserve the
existing timeout/failure outcome without treating a nil wait error as exit
status 0, and ensure the cmd.Wait goroutine has a clear shutdown path.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:cbfce190880da32e10f4206d -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
