---
task: task_01
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
type: test
complexity: medium
---

# Task 01: Record the characterization corpus before anything moves

## Overview

Pin today's behavior of the OpenCode reasoning contract before any of it moves.
The corpus has two halves named apart on purpose: invariants that must survive
this Spec, and current behavior it intends to break. Each later Task edits only
the second half, so every break is a visible, deliberate edit.

## Requirements

1. MUST add characterization tests in three new files, one per package, and MUST
   NOT modify any existing test file.
2. MUST pin these invariants: a Codex selection with a non-empty effort plans
   the `independent` encoding; a Claude selection does the same; an empty effort
   on `opencode` plans `runtime_managed`; a model absent from the advertised
   values produces `SelectionUnsupportedError`.
3. MUST pin this current behavior, with `Today` in each test name: configuration
   refuses a non-empty `reasoning_effort` on `runtime: opencode`;
   `acpxReasoningEffortConfigKey` returns a `ModelManagedReasoningError` for
   that runtime; and a Run's recorded acpx command sequence for an `opencode`
   session contains no reasoning config set between the session and the prompt.
4. MUST derive fixtures from the shape the adopted measurement recorded — an
   `effort` option whose advertised values differ per model — rather than from
   invented vocabularies.
5. MUST NOT change any non-test file and MUST NOT change production behavior.

## Subtasks

- [ ] Add the planning and encoding characterization file.
- [ ] Add the configuration characterization file.
- [ ] Add the session-lifecycle characterization file covering the command order.
- [ ] Name every break-half test with `Today`.
- [ ] Confirm the whole suite is green with the corpus added.

## Acceptance Criteria

- [ ] Three new test files exist and no pre-existing test file differs.
- [ ] Every invariant in Requirement 2 has a test that fails if the behavior moves.
- [ ] Every current behavior in Requirement 3 has a test whose name contains `Today`.
- [ ] The recorded command sequence test asserts the absence of a reasoning set,
      not merely the presence of the prompt.
- [ ] No file outside the three new test files changed.

## Context

- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/config/profiles.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_effort_characterization_test.go`
- `internal/agent/acpx_session_effort_characterization_test.go`
- `internal/config/opencode_effort_characterization_test.go`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_01.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent ./internal/config -run 'EffortCharacterization' -count=1 -v` — expected: exits 0 and names at least one test in each package; `no tests to run` fails this Task.
- `go test ./internal/agent ./internal/config -run 'EffortCharacterizationToday' -count=1 -v` — expected: exits 0 and names the break-half tests, proving Requirement 3 is pinned rather than described.
- `test -z "$(git diff --name-only -- internal | grep -v '_characterization_test\.go$')"` — expected: exits 0, proving no existing source or test file moved. `test -z` is what makes an empty result a pass; a bare `grep -v` exits 1 when it matches nothing, which the Daemon reads as failure.

## References

- `_prd.md` → Goals 2, 4, 5.
- `_techspec.md` → Testing Approach; Build Order 1.
- `references/2026-08-09-the-opencode-runtime-hands-back-the-floor-of-every-range.md`
  → the per-model effort vocabularies the fixtures reproduce.
- ADR-0108.

## Result

Implemented the characterization corpus without changing production behavior.
The planning tests preserve Codex and Claude `independent` encoding, OpenCode
empty-effort `runtime_managed` encoding, and the typed refusal for an Agent
Model absent from advertised values. Their capability fixtures reproduce the
measured per-model `effort` shapes: Grok 4.5 advertises `low`, `medium`, and
`high`, while DeepSeek V4 Pro advertises `high` and `xhigh`.

The declared-break tests all contain `Today` in their names. They record the
configuration refusal, the `ModelManagedReasoningError` key-mapping refusal,
and the current Run command order. The lifecycle assertion locates `sessions
ensure` and `prompt` in the recorded acpx invocations, then rejects either
`set effort` or `set reasoning_effort` anywhere between them.

Focused checks run after the last Go edit:

- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run "^TestSelectionEffortCharacterization" -count=1'` — 4 tests passed.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run "^TestACPXSessionEffortCharacterization" -count=1'` — 2 tests passed.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/config -run "^TestOpenCodeEffortCharacterizationTodayConfigurationRefusesNonEmptyReasoningEffort$" -count=1'` — 1 test passed.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/agent ./internal/config -count=1'` — 471 tests passed across both affected packages.
- `rtk git -c core.fsmonitor=false diff --check` — exited 0.

Acceptance evidence:

- Criterion 1: `rtk git -c core.fsmonitor=false status --short --untracked-files=all`
  reports exactly the three new characterization test files plus this Task
  file; no pre-existing test file differs.
- Criterion 2: the four `TestSelectionEffortCharacterizationInvariant...`
  tests exercise every Requirement 2 invariant and passed in the focused run.
- Criterion 3: the three declared breaks are named
  `TestOpenCodeEffortCharacterizationToday...` and
  `TestACPXSessionEffortCharacterizationToday...`; all passed in the focused
  runs.
- Criterion 4: the Run lifecycle test asserts both command order and the
  absence of the two reasoning-set command keys between session creation and
  the work prompt.
- Criterion 5: the changed-path inspection shows no implementation or
  pre-existing test change. The Task file differs because the Daemon set its
  status to `in_progress` and this Result records Agent evidence.

The commands authored under `## Verification` were not run; Daemon
Verification remains the settlement boundary.
