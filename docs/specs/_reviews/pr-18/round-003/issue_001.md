---
source: coderabbit
pr: "18"
round: 3
round_created_at: "2026-07-07T14:28:04Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: dc811027414af88955c449b2688e1d839388ebed
file: internal/cli/cli.go
line: 2719
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6O8iGZ,comment:PRRC_kwDOS0qyts7S1iD7
review_hash: 8bc49e8cb2c4fea84c7c078ea70618ee27c77d591401430805a233429c991589
duplicate_of: ""
source_review_id: "4646089284"
source_review_submitted_at: "2026-07-07T14:27:12Z"
---

# Issue 001: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Inspect the notifier implementation for an execution timeout around the shell/native command
fd -t f 'notify.go$|native_.*\.go$|shell_.*\.go$' internal/notify --exec cat -n {}
rg -nP 'WithTimeout|WithDeadline|CommandContext|context\.' internal/notify
```

Repository: marcioaltoe/roundfix

Length of output: 9871

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

# Inspect the call site and the context wrapper used before notification.
sed -n '2690,2785p' internal/cli/cli.go
printf '\n---\n'
rg -n 'withoutCancelOrBackground|notifyTerminalOutcome|completeStoppedRunRecord' internal/cli/cli.go
printf '\n---\n'
rg -n 'func withoutCancelOrBackground|WithoutCancel|Background\(\)' internal/cli
```

Repository: marcioaltoe/roundfix

Length of output: 30950

---



</details>

**Bound the native notifier call.** `notify.command` already has its own 30s timeout, but `desktopNotifier.Notify` still receives a deadline-free context here, so `osascript`/`notify-send` can block the CLI if they hang. Keep a timeout on the native path too.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2711 - 2719, The native notification path
in notifyTerminalOutcome currently strips cancellation with
withoutCancelOrBackground(ctx), so desktopNotifier.Notify can run without a
deadline and hang the CLI. Update notifyTerminalOutcome to pass a bounded
context to notifier.Notify, reusing the existing timeout behavior or applying a
short local timeout for the native path, while keeping the existing error
handling via reportOutcomeNotificationFailure intact.
```

</details>

<!-- fingerprinting:phantom:medusa:capybara -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8b7e441edc5547544414dc3d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `notifyTerminalOutcome` detached cancellation with
  `withoutCancelOrBackground(ctx)` but did not add a replacement deadline, so
  native notification runners could receive an unbounded context. Fixed by
  wrapping the detached context with the terminal notification timeout before
  calling `Notifier.Notify`.

## Resolution

- Changed `internal/cli/cli.go` so terminal outcome notifications run with a
  bounded detached context.
- Added
  `TestNotifyTerminalOutcomeDetachesCanceledParentWithBoundedDeadline` in
  `internal/cli/cli_test.go`.

## Verification

- `rtk go test ./internal/notify` — passed.
- `rtk go test ./internal/cli` — passed.
- `rtk make verify` — passed after approved rerun outside the sandbox for Go
  build-cache access.
