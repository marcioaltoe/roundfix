---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: internal/cli/detach_test.go
line: 45
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95P,comment:PRRC_kwDOS0qyts7V5tZ-
review_hash: 8dcbda6933c8856ffc324d11a57a563eb2bae884f6fa811dd4831a1645cce757
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_006.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 006: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Give subprocess startup enough time before testing phase transitions.**

Lines 45 and 64 can expire before a loaded CI worker starts the child, producing the wrong timeout diagnostic. Increase the liveness budget while keeping Run creation delayed beyond it.

<details>
<summary>Proposed fix</summary>

```diff
-	stdout, stderr, code := runDetachParentForTest(t, detachTestChildRunCreated, 20*time.Millisecond, time.Second)
+	stdout, stderr, code := runDetachParentForTest(t, detachTestChildRunCreated, 500*time.Millisecond, 2*time.Second)
...
-	stdout, stderr, code := runDetachParentForTest(t, detachTestChildSlowRunCreation, 50*time.Millisecond, 50*time.Millisecond)
+	stdout, stderr, code := runDetachParentForTest(t, detachTestChildSlowRunCreation, 500*time.Millisecond, 50*time.Millisecond)
...
-		time.Sleep(75 * time.Millisecond)
+		time.Sleep(750 * time.Millisecond)
```
</details>






Also applies to: 63-64, 154-168

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/detach_test.go` around lines 44 - 45, Update the timing
arguments in TestRunDetachedCommandAllowsRunCreationPastLivenessDeadline and the
related tests around runDetachParentForTest so subprocess startup has a larger
liveness budget on slow CI workers. Keep the child run creation delay beyond
that budget, preserving the intended phase-transition assertion and timeout
diagnostic.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:43a0c18dc032d596fa4a383c -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
