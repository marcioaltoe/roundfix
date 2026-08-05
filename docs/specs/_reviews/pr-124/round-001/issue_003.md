---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/cli/cli.go
line: 3233
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wuaz0,comment:PRRC_kwDOS0qyts7d35tI
review_hash: ab0dc98cc6026410f5bbfdada0047115c973c788d14fcdcfc5fdc87ff616aa0a
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**The pending-evidence predicate is duplicated across packages.**

`watch.unrecognisedEvidenceReason` decides when `TerminalReason` carries the unrecognised-signal diagnostic. It uses exactly the condition `State == EvidencePending && Kind != EvidenceKindNone`. Line 3229 repeats that condition to decide whether to print the reason. The two copies must stay in sync, and nothing enforces that.

Expose the decision from the `watch` package instead. One option is a boolean field on `watch.Result`, for example `TerminalReasonIsDiagnostic`. Another option is an exported helper in `watch` that the CLI calls. Either removes the need for the CLI to re-derive evidence semantics.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 3228 - 3233, Expose the diagnostic decision
from the watch package, using either a Result field such as
TerminalReasonIsDiagnostic or an exported helper tied to
watch.unrecognisedEvidenceReason. Update the CLI timeout handling around
result.Outcome and result.TerminalReason to use that exposed decision instead of
duplicating the EvidencePending/EvidenceKindNone predicate, while preserving the
existing output behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:96a76f7a3f30513e365397bc -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `watch.resultForTimedOut` now owns the unrecognised-Evidence decision
  and exposes it as `Result.TerminalReasonIsDiagnostic`. The CLI prints the
  diagnostic from that result flag instead of re-deriving the
  `pending`/non-`none` predicate. Watch tests cover both an unrecognised
  diagnostic and an ordinary timeout without one.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test ./internal/watch -count=1` passed, as did the CLI regression via `rtk
  env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli
  -count=1 -run
  '^TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush$'`.
- Daemon Verification: `make verify` not run; Daemon-owned.
