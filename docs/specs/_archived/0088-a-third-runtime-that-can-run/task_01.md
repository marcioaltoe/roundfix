---
task: task_01
spec: 0088-a-third-runtime-that-can-run
status: completed
type: test
complexity: medium
---

# Task 01: Record the characterization corpus before anything moves

## Overview

Pin today's behavior of the capability projection, the OpenCode reasoning-effort
mapping, and profile readiness scope as executable tests, before any of the three
changes. The corpus has two halves that are named apart on purpose: invariants
that must survive this Spec, and current behavior this Spec intends to break.
Each later Task edits only the second half, so every break is a visible,
deliberate edit rather than a silent one.

## Requirements

1. MUST add characterization tests in three new files, one per package, and MUST
   NOT modify any existing test file.
2. MUST pin these invariants, which no later Task may change: a `select` option
   at or below the retained bound keeps every advertised value in advertised
   order; a bracketed variant identifier still parses into its canonical model
   and effort; a requested model absent from the advertised list produces
   `SelectionUnsupportedError` rather than invalid capability evidence.
3. MUST pin this current behavior, which later Tasks will break, with `Today` in
   each test name so the break is legible: a `select` option above the retained
   bound fails with `too_many_option_values` and cascades to
   `missing_model_state` and `contradictory_response`; a profile naming
   `runtime: opencode` with a non-empty `reasoning_effort` is accepted by
   configuration decoding; `roundfix doctor` reports profile readiness over
   exactly the five required Agent Work Categories and ignores a configured
   optional-category profile.
4. MUST derive the oversized-option fixture from the shape the adopted
   measurement recorded — a `model` option whose advertised values carry a
   provider-prefixed identifier — at a size above the bound rather than at 417.
5. MUST NOT change any non-test file, and MUST NOT change any production
   behavior.

## Subtasks

- [x] Add the capability-projection characterization file with both halves.
- [x] Add the configuration characterization file covering the accepted
      `opencode` reasoning effort.
- [x] Add the Doctor Command characterization file covering readiness scope.
- [x] Name every break-half test with `Today` so a later edit is self-declaring.
- [x] Confirm the whole suite is green with the corpus added.

## Acceptance Criteria

- [x] Three new test files exist, one in each of the agent, config, and cli
      packages, and no pre-existing test file differs.
- [x] Every invariant in Requirement 2 has a test that fails if the behavior
      moves.
- [x] Every current behavior in Requirement 3 has a test whose name contains
      `Today`.
- [x] The oversized-option fixture exceeds the retained bound and uses
      provider-prefixed model identifiers.
- [x] No file outside the three new test files changed.

## Context

- interface: `internal/agent/selection_capabilities.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/config/profiles.go`
- interface: `internal/cli/doctor.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_capabilities_characterization_test.go`
- `internal/config/profiles_characterization_test.go`
- `internal/cli/doctor_characterization_test.go`
- `docs/specs/0088-a-third-runtime-that-can-run/task_01.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent ./internal/config ./internal/cli -run 'Characterization' -count=1 -v` — expected: exits 0 and names at least one test in each of the three packages; a package reporting `no tests to run` fails this Task.
- `go test ./internal/agent ./internal/config ./internal/cli -run 'CharacterizationToday' -count=1 -v` — expected: exits 0 and names the break-half tests, proving Requirement 3 is pinned and not merely described.
- `git diff --name-only -- internal | grep -v '_characterization_test\.go$'` — expected: prints nothing, proving no existing source or test file moved.

## References

- `_prd.md` → Goals 1, 3, 4; Core Features.
- `_techspec.md` → Testing Approach; Build Order 1.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the advertised payload shape the fixture reproduces.
- ADR-0105, ADR-0106, ADR-0107.

## Result

Three characterization files were added and no existing file moved.

**What changed.** The capability projection, Agent Selection normalization, and
Doctor Command readiness scope now each have executable descriptions of their
current behavior, split into an invariant half the Spec must preserve and a
`Today` half the Spec intends to break. The oversized-option fixture is built
from the retained-bound constant rather than a literal, so it follows the bound
when Task 02 moves it, and it uses provider-prefixed identifiers matching the
shape the adopted measurement recorded.

**Commands and outcomes.**

- `go build -buildvcs=false ./...` — exit 0.
- `go test ./internal/agent ./internal/config ./internal/cli -run 'CharacterizationToday|CharacterizationInvariant' -count=1 -v` — exit 0; eight tests, four in agent, three in config, one in cli.
- `git status --porcelain -- internal` — three untracked characterization files and nothing else.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — exit 0 on a clean test cache.
- `make verify` — exit 0 after `go clean -testcache`.

**Evidence per acceptance criterion.**

- Three new test files, no pre-existing file changed: the `git status` output
  above lists exactly `internal/agent/selection_capabilities_characterization_test.go`,
  `internal/config/profiles_characterization_test.go`, and
  `internal/cli/doctor_characterization_test.go`.
- Invariants pinned: `TestCharacterizationInvariantRetainsEveryValueAtOrBelowTheBound`,
  `TestCharacterizationInvariantParsesBracketedVariantIntoCanonicalAndEffort`,
  `TestCharacterizationInvariantUnadvertisedModelIsUnsupportedNotInvalidEvidence`,
  `TestCharacterizationInvariantAcceptsCodexAndClaudeReasoningEffort`,
  `TestCharacterizationInvariantAcceptsAnEmptyReasoningEffort`.
- Break-half pinned with `Today` in the name:
  `TestCharacterizationTodayRefusesAnOversizedAdvertisedOption`,
  `TestCharacterizationTodayAcceptsOpenCodeReasoningEffort`,
  `TestCharacterizationTodayDoctorIgnoresConfiguredOptionalCategory`.
- Oversized fixture exceeds the bound and uses provider-prefixed identifiers:
  `oversizedModelCapabilityFixture` fills until `len(values) > maxCapabilityValues`
  with `openrouter/vendor-NNN/model` values plus the requested
  `opencode-go/kimi-k3`.
- The oversized refusal reproduces the measured cascade: the test asserts
  `too_many_option_values`, `missing_model_state`, and `contradictory_response`
  together, which is the exact triple `roundfix profiles validate --category data`
  returned on 2026-08-08.

**Follow-ups.** Tasks 02, 03, and 04 each edit exactly one `Today` test; each
must declare that edit as a break in its own Result.
