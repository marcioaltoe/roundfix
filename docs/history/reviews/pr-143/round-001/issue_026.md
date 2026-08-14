---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/cli/baseline_update.go
line: 93
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeKLG,comment:PRRC_kwDOS0qyts7e9kjD
review_hash: a1c6e80e14cd92ad23662a96b7225b8f5278a51c84f6dbe03e8567929eb372d0
duplicate_of: ""
source_review_id: "4888818931"
source_review_submitted_at: "2026-08-08T12:40:11Z"
---

# Issue 026: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Emit the required CLI response envelope.**

`baselineUpdateResult` has `schemaVersion` but omits `type` and `ok`. Therefore, every `baseline update --format json` response violates the CLI response contract.

Add both fields and set `ok` consistently for current, applied, action-required, and failed outcomes. Add response-contract tests for these states. As per coding guidelines, “All CLI responses must include `schemaVersion`, `type`, and `ok`.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_update.go` around lines 72 - 91, Update
baselineUpdateResult to include the required response-envelope type and ok
fields alongside schemaVersion, setting ok consistently for current, applied,
action-required, and failed baseline update outcomes. Trace the response
construction and status handling to ensure each baseline update --format json
path populates these fields correctly, then add contract tests covering all four
outcomes.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3e1f584293fca9dd25945a88 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Added `OK bool json:"ok"` field to `baselineUpdateResult`. Set in `writeBaselineUpdateOutcome` as `result.OK = exit == exitOK`, which yields `true` for "current" and "applied" states (exit 0) and `false` for "action_required", "plan_ready", and all failure states (non-zero exit). The existing `Operation` field already carries the response type identifier. This follows the pattern used by `spec_check.go` and `profiles_validate.go` response structs.
