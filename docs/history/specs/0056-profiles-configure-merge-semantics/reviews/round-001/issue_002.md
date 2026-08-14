---
source: coderabbit
pr: "67"
round: 1
round_created_at: "2026-08-02T11:30:39Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/profiles-configure-merge-semantics
head_sha: ffcc15ebed0a055d329cb3215ae0878b90931948
file: internal/cli/cli_test.go
line: 1563
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vwh2G,comment:PRRC_kwDOS0qyts7ceBLn
review_hash: 13ae5b5fdbf5a57ce7be7f140374bc33c3d65f55616a42a1076f39ed6aac57e5
duplicate_of: ""
source_review_id: "4838273774"
source_review_submitted_at: "2026-08-02T11:29:41Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Replace the direct `os.Stdin` swap with a seam.**

This subtest reassigns the `os.Stdin` package variable and restores it in `t.Cleanup`. The restore is correct, but the mutation makes the subtest unable to run with `t.Parallel()` and couples the test to process-global state. The package already exposes `profilesConfigureInput` as an injectable reader seam; extend the confirmation path to read through a comparable seam, then inject an empty reader here.

If the confirmation reader cannot be injected today, keep the current approach and add a short comment stating why the global swap is required.

As per coding guidelines: "Keep tests isolated ... and avoid shared global state and module-level singletons."

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 1525 - 1563, The non-interactive
confirmation test directly mutates the process-global os.Stdin. Extend the
confirmation input path used by the profiles configure flow with an injectable
reader seam comparable to profilesConfigureInput, then inject an empty reader in
the “non-interactive confirmation EOF” subtest instead of swapping os.Stdin;
preserve the existing EOF refusal assertions and cleanup behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:2aef3c3ce27a5e10bf275355 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `defaultConfirmProfilesConfigure` read `os.Stdin` directly even though the command already owns the `profilesConfigureInput` reader seam. Confirmation now reads through that seam, and the EOF subtest injects an empty reader without mutating process-global stdin.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache rtk go test ./internal/cli -run '^(TestProfilesConfigureExitCodes|TestProfilesConfigureChangeSummary|TestProfilesConfigureProofScope|TestPrintProfilesConfigureRefusalWritesJSONBeforeDiagnostic)$' -count=1` passed 17 tests, including `TestProfilesConfigureExitCodes/non-interactive_confirmation_EOF`.
