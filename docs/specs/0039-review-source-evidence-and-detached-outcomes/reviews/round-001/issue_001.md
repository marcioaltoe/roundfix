---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/cli_test.go
line: 9822
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpN,comment:PRRC_kwDOS0qyts7aUVDJ
review_hash: fb7e0739c84730b04b9bf76b37ba6b31977defe8811397d5ecf6a6c82a095359
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:18Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**`assertRecordedOutcomes` now discards exactly the fields this PR adds.**

The helper reduces both sides to `RunID/State/Kind/Target` before comparing, so every caller silently stops asserting `Reason`, `NextAction`, `ConsoleLog`, `AttachCommand`, and `ReviewIssuesKnown`. Those are the new terminal-context fields, and only `TestOutcomeNotificationCarriesTerminalContextWithBoundedDeadline` checks them at all — a regression that drops the reason or next action from notifications would pass here.

Compare the full `Outcome` and let callers pass the expected context fields, or add an explicit assertion that non-Clean recorded outcomes carry a non-empty `Reason` and `NextAction`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 9798 - 9821, Update
assertRecordedOutcomes to compare complete roundnotify.Outcome values instead of
reducing them to RunID, State, Kind, and Target. Ensure callers can provide and
validate terminal-context fields including Reason, NextAction, ConsoleLog,
AttachCommand, and ReviewIssuesKnown, while preserving the existing count and
indexed mismatch diagnostics.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:d8ecd69e467b8917bd0eaab9 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `assertRecordedOutcomes` now compares complete `roundnotify.Outcome` values, and affected expectations include the terminal-context fields. Focused CLI outcome-notification tests passed with a task-local `GOCACHE`.
