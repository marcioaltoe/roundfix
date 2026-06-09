---
name: go-tui
description: Use when designing, building, modifying, or reviewing terminal UI in Go for this repository, including task sidebars, issue lists, streaming agent consoles, keyboard focus, panes, status bars, and Bubble Tea or similar TUI libraries.
---

# Go TUI

Build a terminal UI that helps users watch work progress without hiding errors.

## Product shape

The default Roundfix TUI should favor this layout:

- Left pane: pull request, round, task batch, issue status, and selected issue.
- Right pane: streaming runtime console and tool output for the selected job.
- Bottom bar: compact key hints and run state.

Keep it operational. Do not turn the TUI into a landing page or long in-app
manual.

## Workflow

1. Start from the CLI/run state model. Do not invent separate TUI state that can
   drift from execution state.
2. Keep `Init`, `Update`, and `View` responsibilities separate.
3. Keep `View` pure: no IO, no process control, no database writes.
4. Put process, database, and stream effects behind commands or injected
   interfaces.
5. Run `rtk gofmt -w <changed-go-files>` and `rtk go test ./...`.

## Layout rules

- Account for borders and fixed bars before sizing panes.
- Use proportional sizing rather than hardcoded widths where possible.
- Truncate long issue titles, file paths, branches, and log lines explicitly.
- Keep text inside its pane. Avoid relying on terminal auto-wrap.
- Recalculate layout on terminal resize.
- Keep stable dimensions for sidebars, status bars, and counters so updates do
  not shift the whole screen.

## Interaction rules

- Keyboard focus must be visible and predictable.
- Use tab or arrow navigation for pane focus and issue selection.
- Provide a clear quit/detach path that does not accidentally cancel running
  work.
- Never hide preflight errors behind an empty screen. Render the error and the
  next action before any remote fetch or agent startup.
- If a selected ACP runtime fails to start, show that runtime error. Do not
  fallback to another runtime automatically.

## Streaming rules

- Store console output in a bounded buffer.
- Preserve enough recent output to diagnose failures.
- Avoid unbounded memory growth during long watch runs.
- Tie stream readers to `context.Context` and close them on cancellation.
- Treat stdout and stderr as distinct streams when the runtime exposes them.

## Dependency rule

Do not add a TUI dependency until implementation needs it. When needed, add it
with `rtk go get` and isolate framework-specific code behind package-local
types so the run engine remains testable without a terminal.
