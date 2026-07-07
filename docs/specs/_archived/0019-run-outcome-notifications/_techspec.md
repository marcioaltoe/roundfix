---
spec: 0019-run-outcome-notifications
prd: _prd.md
created: 2026-07-07
---

# Run Outcome Notifications — Technical Spec

## Executive Summary

A small notifier sits behind the existing Run completion boundaries: when an
operational Run reaches a terminal outcome, the owning process fires one
best-effort notification — a platform-native desktop notification by default,
or the user's configured command. The primary trade-off is per-platform code
for the native default: shelling out to `osascript` (macOS) and `notify-send`
(Linux) keeps the stdlib-only rule and needs no daemon or service, at the cost
of two thin platform paths and a silent skip where neither exists (Windows,
headless boxes). The configured command escapes that limit entirely — any
channel the user can script. Notifications run in the Run-owning process, so
Detached Runs notify from the detached child, which is precisely the
unattended case the feature exists for.

## System Architecture

One new package, two touched modules:

- `internal/notify` (new) — the notifier: outcome payload, native desktop
  path per platform (build-tagged, mirroring the `detach_unix.go` /
  `detach_windows.go` pattern), configured-command execution with a bounded
  context, and best-effort error reporting. New package earns its place: the
  notifier is used by three commands and owned by none of them, and
  `internal/cli` is already the repository's largest package.
- `internal/config` — new `notify` section (`enabled`, `command`) with the
  standard precedence and strict-key validation.
- `internal/cli` — the terminal-outcome call sites for `resolve`, `watch`,
  and `implement` invoke the notifier after `CompleteRun` succeeds. Fetch
  completion, settle, and archive do not notify (they create no long Run or
  no Run at all). Stop-forced completion notifies through the same boundary
  as any other terminal outcome reached by the owning process.

Flow: `CompleteRun` returns the terminal Run → build payload (run id,
outcome, kind, target) → notifier: configured command if set, else native
desktop, else silent skip → failure becomes one stderr warning plus one
Daemon-source Run Event; success is silent.

## Implementation Design

### Interfaces

```go
// Package notify fires one best-effort notification per terminal Run outcome.

// Outcome carries the Run context a notification names.
type Outcome struct {
    RunID   string
    State   string // terminal Run state, e.g. Clean, Failed, Unresolved
    Kind    string // resolve | watch | implement
    Target  string // "pr:<number>" or "spec:<slug>"
}

// Notifier sends one notification; implementations are best-effort.
type Notifier interface {
    Notify(ctx context.Context, outcome Outcome) error
}

// New picks the implementation: command when configured, native otherwise,
// no-op when disabled or no native mechanism exists.
func New(cfg Config) Notifier
```

Command environment (the configured command's contract):

```text
ROUNDFIX_RUN_ID   ROUNDFIX_OUTCOME   ROUNDFIX_KIND   ROUNDFIX_TARGET
```

### Data Models

No schema change. Config gains:

```yaml
notify:
  # false disables notifications entirely.
  enabled: true
  # Shell command run on terminal outcome; empty uses the native desktop
  # notification. Receives ROUNDFIX_* environment variables.
  command: ""
```

### API Contracts

- Native default: macOS uses `osascript -e 'display notification ...'`;
  Linux uses `notify-send` when present on PATH. Anything else — including
  Windows and headless Linux — skips silently. The notification text names
  the outcome and target (`Clean — spec:0017-run-discovery` shape) with a
  Roundfix title.
- Configured command: run through the shell with the `ROUNDFIX_*`
  environment, bounded by a 30s timeout, stdout/stderr discarded except into
  the warning on failure.
- Failure reporting: one stderr line shaped like
  `roundfix: outcome notification failed: <reason>` and one Daemon-source Run
  Event; the Run's report, outcome, and exit code are untouched. Exit code
  contracts of all commands are unchanged.
- `notify.enabled: false` restores byte-for-byte today's behavior.

## Coverage Map

- Story 1 (desktop on terminal outcome, detached included) → `notify.New`
  native path + CLI completion call sites (detached child owns completion)
- Story 2 (own channel via command) → command Notifier + env contract
- Story 3 (no tooling, no new failure mode) → silent-skip no-op path +
  best-effort error handling
- Core Feature 5 (journaled warning) → CLI boundary warning + Run Event
- Core Feature 4 (config precedence) → `internal/config` notify section

## Integration Points

`osascript` and `notify-send` as optional local binaries, invoked
best-effort; absence is a silent skip, never an error.

## Testing Approach

Existing seams; the Notifier interface is the one new seam (needed so CLI
tests assert firing without desktop side effects):

- `internal/notify` tests: command notifier env and timeout via a fake
  command (shell script writing its env to a file); native selection logic
  behind a lookup fake; no-op paths.
- `internal/config` table tests for the notify section.
- `internal/cli` tests: a captured fake Notifier asserts one notification
  per terminal outcome for resolve, watch, and implement, none for fetch,
  and the warning path on notifier failure leaves report and exit code
  unchanged.

## Build Order

1. `notify` config section (enabled, command) with validation and tests.
2. `internal/notify` package: Outcome, command notifier, native notifiers
   (build-tagged), selection, tests (depends on: 1).
3. CLI wiring at the terminal-outcome boundaries for resolve, watch, and
   implement, with warning + Run Event on failure and CLI tests
   (depends on: 2).
4. Docs and skill sync: README Config and Command Boundaries, usage guide,
   roundfix SKILL.md, `make skills-sync` (depends on: 3).

## Risks & Considerations

- Notification runs after `CompleteRun`; a crash between completion and
  notification loses the notification, not the Run — acceptable for
  best-effort.
- The configured command runs with the Run's environment; document that it
  executes with user privileges and should be short-lived (the 30s bound
  kills stragglers).
- macOS `osascript` notifications require no permission prompt for the
  calling terminal in the default case; when the OS suppresses them, the
  command exits zero and Roundfix cannot tell — out of scope.

## Decisions

- One new `internal/notify` package over growing `internal/cli` — three
  consumers, one owner.
- Native mechanism per platform via build tags, mirroring the detach
  platform-file pattern.
- 30s command bound, discard output on success — notifications must never
  become the slow or noisy part of a Run.
