---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/acpx_runner_test.go
line: 760
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWiS,comment:PRRC_kwDOS0qyts7fswNs
review_hash: 4b75cf28a1772cd9eba26147770fecdd40c864d8c724bc3f059b1f89d2d851d4
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:13Z"
---

# Issue 003: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Three tests index the invocation slice without a length guard.** The new override-free catalogue ensure shifted every invocation index, and these tests now read higher indices than before. If preflight emits fewer calls, each site panics with an index-out-of-range rather than reporting the expectation. All tests in this file use `t.Parallel()`, so a panic aborts the test binary and discards unrelated parallel results. `TestACPXProbeValidatesSelectionWithDisposableSession` (line 667) shows the guarded pattern.
- `internal/agent/acpx_runner_test.go#L748-L760`: check `len(invocations) != 7` before reading `invocations[1]` through `invocations[6]`.
- `internal/agent/acpx_runner_test.go#L848-L853`: check the length covers index 4 before reading `invocations[1]` through `invocations[4]`.
- `internal/agent/acpx_runner_test.go#L1289-L1295`: check the length covers index 3 before reading `disposableInvocations[0]` through `disposableInvocations[3]`.

<details>
<summary>📍 Affects 1 file</summary>

- `internal/agent/acpx_runner_test.go#L748-L760` (this comment)
- `internal/agent/acpx_runner_test.go#L848-L853`
- `internal/agent/acpx_runner_test.go#L1289-L1295`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner_test.go` around lines 748 - 760, Add
invocation-count guards before indexing in all three sites:
internal/agent/acpx_runner_test.go lines 748-760 must require len(invocations)
== 7, lines 848-853 must verify the slice covers index 4, and lines 1289-1295
must verify disposableInvocations covers index 3. Follow the guarded pattern in
TestACPXProbeValidatesSelectionWithDisposableSession so failures report through
the test instead of panicking.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/agent/acpx_runner_test.go</file>
<line_range>748-760</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/agent/acpx_runner_test.go</file>
<line_range>848-853</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/agent/acpx_runner_test.go</file>
<line_range>1289-1295</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0e593885cd0b2a34ab1fd227 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Added invocation-count guards before slot indexing in all three sites in internal/agent/acpx_runner_test.go. `TestProfileProofIsHonestAboutTheSelection` now requires `len(invocations) == 7` before reading [1..6]; the missing-session on selection-failure test now guards `len(invocations) < 5` before reading [1..4]; the disposable/live test now guards `len(disposableInvocations) < 4` before reading [0..3]. Mirrors the guarded pattern so a short invocation list reports through the test instead of panicking under `t.Parallel()`. Focused: `go test ./internal/agent -run 'TestProfileProof|TestProof|TestMissing'` ok.
