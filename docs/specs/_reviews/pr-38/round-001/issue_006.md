---
source: coderabbit
pr: "38"
round: 1
round_created_at: "2026-07-27T15:34:32Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 9ed57622bb92f138aa3e23d4d59e260ebbff0116
file: internal/store/store.go
line: 425
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-PF,comment:PRRC_kwDOS0qyts7aENCf
review_hash: bd89437ed3c245af355b6846cef53f6df813ef12b2049665adb904008c5c511c
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260727T152947Z_936cd84aa803ba5d/verification/batch-001-attempt-2.log'
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:14Z"
---


# Issue 006: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Terminal-state set is enumerated in three places in `internal/store/store.go`.** Both compare-and-set guards inline the same ten states already listed by `IsTerminalState` (Line 1045); adding a future terminal state to only one location silently makes that outcome overwritable, defeating the immutability invariant this PR establishes.
- `internal/store/store.go#L376-L395`: build the `state NOT IN (...)` clause and its arguments from a shared `terminalStates` slice instead of the inline literals.
- `internal/store/store.go#L752-L789`: reuse the same generated clause and arguments here, and have `IsTerminalState` range over that slice.

<details>
<summary>📍 Affects 1 file</summary>

- `internal/store/store.go#L376-L395` (this comment)
- `internal/store/store.go#L752-L789`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/store.go` around lines 376 - 395, The terminal-state list is
duplicated across compare-and-set updates and IsTerminalState, allowing future
additions to become inconsistent. In internal/store/store.go lines 376-395,
create shared terminalStates-derived placeholders and arguments for the state
NOT IN clause; in lines 752-789, reuse that generated clause and arguments and
update IsTerminalState to range over terminalStates.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/store/store.go</file>
<line_range>376-395</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/store/store.go</file>
<line_range>752-789</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:4af78a2dd9c764213810a8cd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `CompleteRun`, `UpdateRunState`, and `IsTerminalState` maintained separate terminal-state enumerations. They now share `terminalStates`; both SQL compare-and-set guards derive their placeholders and arguments from that set. Focused evidence: `rtk proxy env GOCACHE=/tmp/roundfix-run-936cd84aa803ba5d-gocache go test ./internal/store -run '^(TestTerminalOutcomeEveryStoredTerminalStateIsImmutable|TestTerminalOutcomeRejectsIntermediateStateUpdate|TestCompleteRunWinnerAndIdenticalReplay)$' -count=1` passed.
