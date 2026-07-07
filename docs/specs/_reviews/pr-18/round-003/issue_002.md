---
source: coderabbit
pr: "18"
round: 3
round_created_at: "2026-07-07T14:28:04Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: dc811027414af88955c449b2688e1d839388ebed
file: internal/notify/notify.go
line: 186
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6O8iGj,comment:PRRC_kwDOS0qyts7S1iEH
review_hash: 2241d51f26df6082f3732abd8204e7dcbc7f16b28a03324bffa7e0521216b81d
duplicate_of: ""
source_review_id: "4646089284"
source_review_submitted_at: "2026-07-07T14:27:12Z"
---

# Issue 002: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Missing tool lookup failure is silently swallowed with no diagnostic.**

When `lookPath` fails (native tool not installed/misconfigured), `Notify` returns `nil` with zero signal to the user. golangci-lint's `nilerr` check flags this at Line 170-171. Per path instructions, `internal/**/*.go` diagnostics/warnings must go to stderr — currently there's no way for a user to learn why native notifications never fire.

As per path instructions, "Stdout must contain only requested command output; diagnostics, progress, and warnings must go to stderr."





<details>
<summary>💡 Suggested diagnostic on missing native tool</summary>

```diff
 	if _, err := lookPath(notifier.tool); err != nil {
-		return nil
+		fmt.Fprintf(os.Stderr, "native notification %q unavailable: %v\n", notifier.tool, err)
+		return nil
 	}
```
</details>

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 golangci-lint (2.12.2)</summary>

[error] 171-171: error is not nil (line 170) but it returns nil

(nilerr)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/notify/notify.go` around lines 159 - 186, The native tool lookup
failure in desktopNotifier.Notify is being swallowed, so missing/misconfigured
notification binaries produce no diagnostic. Update the lookPath failure path to
emit a warning/diagnostic to stderr before returning, using the existing Notify
method context and notifier.tool so users can see why notifications are not
firing. Keep stdout untouched, and preserve the current runner.Run error
handling while ensuring the missing-tool branch no longer returns nil silently.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:587265673741408988fc898f -->

_Sources: Path instructions, Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The requested stderr diagnostic conflicts with the shipped native
  notification contract. The archived notification spec requires native
  notifiers to be a silent no-op when the platform tool is missing, and
  `skills/roundfix/SKILL.md` documents the same behavior. The existing
  `TestDesktopNotifierMissingToolIsNoop` also asserts that a missing native
  tool returns no error and skips the runner. Emitting a warning would change
  default `notify.enabled: true` behavior on machines without `osascript` or
  `notify-send`, so no code change is appropriate for this finding.

## Verification

- `rtk make verify` — passed after approved rerun outside the sandbox for Go
  build-cache access.
