---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/agent/acpx_runner_test.go
line: 2524
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA2l,comment:PRRC_kwDOS0qyts7e5EBB
review_hash: 75ac3b40198888dccb8691b440c1eb43febeb42f2d99b70ba7d840490c137697
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 007: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Bound both event waits.**

Both tests perform unconditional channel receives. If a regression prevents the expected start event, the test blocks until the global suite timeout and hides the actual failure.

- `internal/agent/acpx_runner_test.go#L2520-L2523`: wait for `promptStarted.done` with a bounded predicate wait and fail with diagnostics when it expires.
- `internal/cli/implement_test.go#L1197-L1202`: wait for each `probe.started` event with the same bounded failure behavior.

As per coding guidelines, use explicit predicate-based waits with timeouts for deterministic tests.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/agent/acpx_runner_test.go#L2520-L2523` (this comment)
- `internal/cli/implement_test.go#L1197-L1202`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner_test.go` around lines 2520 - 2523, Bound the event
waits in internal/agent/acpx_runner_test.go:2520-2523 and
internal/cli/implement_test.go:1197-1202 using explicit predicate-based waits
with timeouts; on expiry, fail each test with diagnostic context instead of
unconditionally receiving from promptStarted.done or probe.started. Update both
affected test sites consistently.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/agent/acpx_runner_test.go</file>
<line_range>2520-2523</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/implement_test.go</file>
<line_range>1197-1202</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3ed5afd629f95a74520fe8a3 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added 5-second timeout selects around channel receives in both acpx_runner_test.go:2524 and implement_test.go:1197-1202. `rtk go build ./...` passes.
