---
spec: 0019-run-outcome-notifications
status: archived
created: 2026-07-07
surfaces: [cli, docs]
archived: "2026-07-07"
source_slug: 0019-run-outcome-notifications
---


# Run Outcome Notifications

A Detached Run that dies overnight stays dead until someone happens to look:
monitors tied to an agent session end with the session, and nothing else wakes
the user. Field evidence from dogfooding: a Run failed at the commit boundary
after verification had passed and sat dead for ten and a half hours before a
human noticed. Run Outcome Notifications make the machine speak when a Run
reaches a terminal outcome — a native desktop notification by default, and a
user-configured command for any other channel — so an unattended Run's ending
is an event someone hears about, not a state someone must poll for.

## Goals

- A Run reaching any terminal outcome produces a notification without anyone
  watching a terminal, including Detached Runs whose owning session is gone.
- Users route notifications to their own channel (push service, chat, sound)
  with one configuration value and no new dependencies.
- Notifications are best-effort: they never change a Run's outcome, report,
  or exit code, and a machine without a notification tool loses nothing but
  the notification.

## User Stories

1. As a user running Detached Runs, I want a desktop notification when a Run
   ends, so that a Run that fails overnight or finishes while I work on
   something else does not sit unnoticed.
2. As a user with my own alerting channel, I want to configure a command that
   Roundfix executes on terminal outcome with the Run's context, so that
   notifications reach me off-machine.
3. As a user on a machine without notification tooling, I want Runs to behave
   exactly as today, so that notifications never become a new failure mode.

## Core Features

1. Operational Runs (`resolve`, `watch`, `implement`) fire one notification
   when they reach a terminal outcome, carrying the Run id, the outcome, the
   command kind, and the target (the Open Pull Request or the Spec slug).
   Fetch Runs and commands that create no Run do not notify.
2. By default the notification is a native desktop notification on platforms
   with a standard mechanism; where none is available, the notification is
   skipped silently.
3. A configured notification command replaces the native notification. The
   command receives the Run context through environment variables and runs
   best-effort with a bounded duration.
4. Notification configuration supports disabling notifications entirely, with
   Project Config over User Config over built-in default (enabled, native).
5. A notification failure — missing tool, non-zero exit, timeout — is recorded
   as a warning in the Run Event Journal and on stderr, and never changes the
   Run's outcome, report, or exit code.
6. Detached Runs notify from the detached process itself, so the notification
   fires even when the launching session no longer exists.

## User Experience

- Default: a Run ends and a desktop notification names the outcome and the
  target. Nothing else about the Run's output changes.
- Configured command: the user's script decides the channel and any filtering
  by outcome — Roundfix always invokes it on terminal outcome and passes the
  context; deciding what is notification-worthy belongs to the script.
- Disabled: byte-for-byte today's behavior.

## Non-Goals / Out of Scope

- Notification on non-terminal events (Round transitions, Task completions,
  Batch commits) — terminal outcome only.
- Built-in integrations with specific services (Slack, Telegram, ntfy) — the
  command is the integration point.
- Delivery guarantees, retries, or notification history.
- Watching or re-notifying Runs that ended while Roundfix was not running.

## Success Metrics

- The overnight-dead-Run scenario produces a notification at failure time on
  a default macOS setup.
- A user routes notifications through a push service with one config value
  and no Roundfix code change.

## Decisions

- Both mechanisms: native desktop notification as the zero-config default,
  one configured command as the universal override.
- Filtering by outcome lives in the user's command, not in config — Roundfix
  notifies every terminal outcome.
- Notifications are fire-and-forget with a bounded duration; failures warn
  and never block or alter the Run.

## Open Questions

None.
