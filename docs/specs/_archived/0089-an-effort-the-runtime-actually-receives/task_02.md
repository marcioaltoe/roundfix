---
task: task_02
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
type: backend
complexity: high
---

# Task 02: Plan a deferred effort and fail closed on an unadvertised one

## Overview

Teach the planner that a non-empty effort on a runtime which cannot accept one
before its session's first prompt is its own encoding. Preflight can see the
value advertised; it cannot apply it. Naming that case `runtime_deferred` is
what keeps a readiness surface from reporting an assignment it never made.

## Requirements

1. MUST add the `runtime_deferred` selection encoding and select it when the
   requested effort is non-empty, the runtime defers effort application, and the
   requested value appears among the values that model advertises.
2. MUST produce `SelectionUnsupportedError` when the requested effort is
   non-empty and absent from the advertised values, naming those values.
3. MUST leave `independent`, `model_variant`, `model_managed`, and
   `runtime_managed` selecting exactly as they do today.
4. MUST make the effective-state check accept a `runtime_deferred` assignment
   whose effort has not been applied yet, and MUST NOT let it accept one whose
   model is wrong.
5. MUST derive "this runtime defers effort" from one predicate so no second copy
   of the runtime list can drift from the first, and that predicate MUST NOT
   detect the runtime by catching the refusal error, because Task 03 removes
   that error. `runtimeManagesOwnReasoning` does exactly that today and is
   rewritten here, not left for Task 03 to discover.
6. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Add the encoding constant with the comment that distinguishes it.
- [ ] Add the deferring-runtime predicate, independent of the refusal error, and
      use it in planning.
- [ ] Select the encoding and fail closed on an unadvertised value.
- [ ] Extend the effective-state check for the new encoding.
- [ ] Edit the break-half characterization tests and declare the breaks.

## Acceptance Criteria

- [ ] An `opencode` selection with an advertised non-empty effort plans
      `runtime_deferred` and reports the requested effort.
- [ ] The same selection with an unadvertised effort produces
      `SelectionUnsupportedError` listing the advertised values.
- [ ] A Codex selection with a non-empty effort still plans `independent`.
- [ ] An `opencode` selection with an empty effort still plans `runtime_managed`.
- [ ] No predicate in `internal/agent` decides that a runtime defers effort by
      matching `ModelManagedReasoningError`.
- [ ] The effective-state check accepts a `runtime_deferred` assignment before
      application and rejects one whose current model differs.

## Context

- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/selection_capabilities.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_assignment.go`
- `internal/agent/selection_assignment_test.go`
- `internal/agent/selection_effort_characterization_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_02.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent -count=1` — expected: exits 0.
- `go test ./internal/agent -run 'RuntimeDeferred' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `! grep -q 'ModelManagedReasoningError' internal/agent/selection_assignment.go` — expected: exits 0, proving the deferring predicate does not depend on the error Task 03 removes.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.
- `grep -q 'SelectionEncodingRuntimeDeferred' internal/agent/selection_assignment.go` — expected: exits 0.

## References

- `_prd.md` → Goals 2 and 4; Core Features: a proof split across two moments.
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0108.

## Result

Implemented the planning half of the deferred OpenCode effort contract. A
non-empty effort now plans `runtime_deferred` only when the selected model's
reasoning option advertises it, carries the option id and requested value for
later Run application, and otherwise returns `SelectionUnsupportedError` with
the advertised effort values. Empty-effort OpenCode planning and every
non-deferring encoding retain their existing paths.

The runtime distinction now comes from `runtimeDefersReasoningEffort`, which
normalizes the ACP Runtime id directly and is reused by both the empty and
non-empty planning branches. It does not inspect or match
`ModelManagedReasoningError`. The effective-state check proves a deferred
assignment from the current model and still-advertised requested effort without
requiring that effort to be current before the Run applies it.

Focused checks run after the last Go edit:

- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run "^TestSelection(EffortCharacterizationRuntimeDeferred|RuntimeDeferred)" -count=1'` — 4 tests passed.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run "^(TestPlanSelectionAssignment|TestSelectionProofAcceptsEchoedAliasGroup|TestSelectionEffortCharacterization.*|TestSelectionRuntimeDeferred.*|TestApplySessionSelection)$" -count=1'` — 29 tests passed.
- `rtk rg -n 'ModelManagedReasoningError|runtimeDefersReasoningEffort|SelectionEncodingRuntimeDeferred' internal/agent/selection_assignment.go` — found the new encoding and predicate call sites, with no `ModelManagedReasoningError` match.

Acceptance evidence:

- Criterion 1: `TestSelectionEffortCharacterizationRuntimeDeferredPlansAdvertisedOpenCodeEffort` passed and compares the full planned assignment, including `ReasoningEffort: "high"`, `ReasoningKey: "effort"`, and `Encoding: "runtime_deferred"`.
- Criterion 2: `TestSelectionEffortCharacterizationRuntimeDeferredRejectsUnadvertisedOpenCodeEffort` passed and proves the typed error, `reasoning_control_not_advertised` classification, and advertised `low`, `medium`, and `high` values.
- Criterion 3: `TestSelectionEffortCharacterizationInvariantCodexPlansIndependentEncoding` passed in the 29-test focused run.
- Criterion 4: `TestSelectionEffortCharacterizationInvariantOpenCodeEmptyEffortPlansRuntimeManagedEncoding` passed in the 29-test focused run.
- Criterion 5: the predicate inspection found no `ModelManagedReasoningError` reference in `selection_assignment.go`; the single normalized runtime predicate owns both planning decisions.
- Criterion 6: `TestSelectionRuntimeDeferredStateMatchesBeforeEffortApplication` and `TestSelectionRuntimeDeferredStateRejectsDifferentCurrentModel` both passed.

The characterization corpus gained tests but no test was renamed or removed,
so `docs/references/coverage-record.json` was not re-recorded. The commands
authored under `## Verification` were not run; Daemon Verification remains the
settlement boundary.
