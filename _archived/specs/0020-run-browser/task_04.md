---
task: task_04
spec: 0020-run-browser
status: completed
type: backend
complexity: medium
---

# Task 04: Browser entry points and the attach loop

## Overview

Wire the Run Browser into its two entry points: `roundfix runs` without a
subcommand and `roundfix attach` without a run id open the browser in an
interactive terminal, selection opens the existing read-only Live Run View,
and leaving the view returns to a refreshed browser until cancel. Demoable
end to end at a terminal.

## Requirements

1. MUST open the Run Browser from bare `roundfix runs` and no-run-id
   `roundfix attach` in an interactive terminal, replacing the numbered
   Interactive Input prompt.
2. MUST open the existing attach cockpit for the selected Run with
   byte-for-byte unchanged cockpit behavior, and return to the browser with
   a fresh store query when the cockpit closes.
3. MUST exit `0` with no side effects on browser cancel.
4. MUST keep the non-interactive failures from the text-surface task exactly
   as shipped there (exit `2`, wording naming `runs list`).
5. MUST leave `attach <known-run-id>` byte-for-byte unchanged.

## Subtasks

- [x] `runs` bare-command dispatch: browser in TTY, exit 2 wording in non-TTY
- [x] attach no-run-id path: browser in TTY replacing the prompt picker
- [x] Browser → cockpit → refreshed browser loop
- [x] CLI tests through the model seam: loop selection, cancel, refresh,
      unchanged explicit-id attach

## Acceptance Criteria

- [x] Selecting a Run in the browser opens the same Live Run View bytes the
      explicit `attach <run-id>` produces, and closing it returns to the
      browser with current data.
- [x] Cancelling the browser exits `0` and mutates nothing.
- [x] Bare `runs` in a non-TTY context exits `2` naming `runs list`.
- [x] Existing explicit-id attach tests pass unchanged.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  entry-point and loop tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1, 3; Core Features 1-3, 5. `_techspec.md` → System
Architecture (flow); API Contracts: runs, attach; Build Order 4.

## Result

Wired the Run Browser into both interactive entry points with a shared
browser → cockpit → refreshed-browser loop, and deleted the numbered
Interactive Input prompt picker it replaces.

What shipped:

- `internal/tui/runbrowser.go` — `RunBrowserSession(ctx, output, browser)`
  runs the task-03 `RunBrowser` model as an alternate-screen Bubble Tea v2
  program (mirroring `RunCockpit`'s program shape, including the
  context-cancel quit goroutine) and returns its `BrowserOutcome`; context
  cancellation settles as a cancel, never a selection.
- `internal/cli/runbrowser.go` — `runRunBrowserLoop`: one unbounded
  `StatesAll` listing per pass feeds both browser filters (single data
  path), selection resolves the Run, derives concurrency through the same
  `attachRunConcurrency` call the explicit path uses, and opens the cockpit
  through `browserAttachCockpit` (the `runAttachCockpit` function itself);
  cockpit close loops back to a fresh query, cancel returns exit `0`.
  `runRunsBrowserCommand` is the bare-`runs` entry; a missing Run Database
  opens one empty-state browser pass. Seams `runBrowserSession` and
  `browserAttachCockpit` let CLI tests script the loop through the model
  seam, not terminal emulation.
- `internal/cli/runs.go` — bare `runs` opens the browser when Interactive
  Input and the live TUI are available; otherwise the task-02 exit-2 wording
  naming `runs list` is byte-unchanged.
- `internal/cli/attach.go` — no-run-id attach opens the browser loop behind
  the same gate; the task-02 non-interactive wording is byte-unchanged, and
  the prompt picker (`pickAttachRun`, `collectAttachRunSelection`,
  `renderAttachRunPicker`, ordering/cancel helpers, `attachPickerInputReader`
  seam) is deleted along with the now-unused `runListState`/`runListTarget`.
  Explicit-id attach code is untouched.
- `internal/cli/cli.go` — `runs`/`attach` help text now names the Run
  Browser entry points truthfully.

Acceptance evidence:

- Same Live Run View bytes + refreshed browser:
  `TestBrowserAttachCockpitIsTheExplicitAttachCockpit` pins the loop's
  cockpit step to the identical `runAttachCockpit` function the explicit
  path calls, and `TestAttachRunBrowserLoopOpensCockpitAndRefreshes` proves
  one cockpit pass for the selected Run with the explicit path's concurrency
  fallback, repository-scoped newest-first listings, and a second browser
  pass that lists a Run created while the cockpit was open (fresh store
  query).
- Cancel is side-effect free: `TestAttachRunBrowserCancelExitsZeroWithoutAttaching`
  — exit `0`, empty stdout, zero cockpit passes, Run count and state
  untouched; `runs` cancel covered in
  `TestRunRunsWithoutSubcommandHonorsInteractivity`.
- Non-TTY bare `runs` exits `2` naming `runs list`: both the
  non-interactive-stdin and interactive-stdin/non-TTY-stdout subtests of
  `TestRunRunsWithoutSubcommandHonorsInteractivity` assert the task-02
  wording.
- Explicit-id attach tests pass unchanged: `TestAttachReplaysCompletedRunReadOnly`,
  `TestAttachUnknownRunFailsBeforeTUIStart`, `TestAttachSpecRun*`, and the
  rest of the attach suite were not modified and pass.

Verification:

- `rtk go test ./internal/cli/` — pass (342 tests).
- `rtk make verify` — pass (fmt-check, 978 tests in 19 packages,
  `roundfix skills check`, build).
- `rtk go test -race ./internal/tui/ ./internal/cli/` — new code race-clean;
  `TestOwningCockpitPollsJournalWhileOwnProcessWrites` fails under `-race`
  on the unmodified base commit too (verified via `git archive HEAD`
  export), so it is a pre-existing race-mode-only cockpit flake, not
  introduced here.

Follow-up notes (out of this Task's slice):

- Task 05 owns README/usage-guide/SKILL.md/CONTEXT.md sync for the Run
  Browser term and behavior.
- Pre-existing: `TestOwningCockpitPollsJournalWhileOwnProcessWrites` is
  timing-sensitive under the race detector and fails consistently with
  `-race` on the base commit.
