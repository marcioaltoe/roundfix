---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/cli/baseline_update.go
line: 399
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA23,comment:PRRC_kwDOS0qyts7e5EBZ
review_hash: 9c55d523eb5480fc14f1f6c8bba8f2ba7d9936a79a9124ec07747b48342b7c05
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 015: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Report installed skills from the install result rather than the owned-name list.**

`dependencies.install` returns a `roundskills.InstallResult`, and the code discards it. `result.Installed` and `result.InstalledCount` then describe every owned skill name, not what the install stage actually wrote. The JSON field name `installedCount` implies observed writes. Use `InstallResult.Targets` to populate the reported names and count.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_update.go` around lines 380 - 389, Capture the
InstallResult returned by dependencies.install in the baseline update flow, then
populate result.Installed from its Targets field instead of
dependencies.ownedNames(). Keep sorting the reported targets and derive
InstalledCount from their length.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:0637bd2f559382e35a342b48 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Captured the install call return, but InstallResult.Targets only exposes file counts (3 per roundfix skill), not skill name count. Reverted to len(result.Installed) for InstalledCount. The InstallResult interface doesn't expose installed skill names; a separate API change would be needed.
