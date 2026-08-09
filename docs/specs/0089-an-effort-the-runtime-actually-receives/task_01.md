---
task: task_01
spec: 0089-an-effort-the-runtime-actually-receives
status: pending
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
- `git diff --name-only -- internal | grep -v '_characterization_test\.go$'` — expected: prints nothing, proving no existing source or test file moved.

## References

- `_prd.md` → Goals 2, 4, 5.
- `_techspec.md` → Testing Approach; Build Order 1.
- `references/2026-08-09-the-opencode-runtime-hands-back-the-floor-of-every-range.md`
  → the per-model effort vocabularies the fixtures reproduce.
- ADR-0108.
