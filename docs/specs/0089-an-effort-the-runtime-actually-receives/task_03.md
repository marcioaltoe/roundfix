---
task: task_03
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
type: backend
complexity: medium
---

# Task 03: Stop refusing an OpenCode reasoning effort

## Overview

Spec 0088 refused a non-empty `reasoning_effort` on `opencode` at two gates:
configuration normalization and runtime validation. Both refusals exist because
the effort could not be applied. Task 02 made it plannable; this Task removes
the refusals so a maintainer can configure one at all.

## Requirements

1. MUST remove the configuration refusal so a non-empty `reasoning_effort` on
   `runtime: opencode` loads.
2. MUST restore `opencode` to the generic reasoning-effort config key so runtime
   validation accepts the selection.
3. MUST keep an empty `reasoning_effort` on `opencode` valid, still planning
   `runtime_managed`.
4. MUST NOT change how Codex or Claude map to their reasoning keys.
5. MUST also edit Spec 0088's configuration corpus, which pins the refusal this
   Task removes. A characterization corpus outlives the Spec that wrote it, so a
   later declared break lands in an earlier Spec's file; that is a declared
   break like any other and is recorded in this Task's Result.
6. MUST leave no unreachable remnant of the refusal — the error type and the
   runtime list it consulted go with it if nothing else uses them. Task 02
   already rewrote the deferring predicate so it does not catch that error; if
   any reference survives, remove it here rather than keeping the type alive.
7. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Remove the normalization refusal and its runtime list.
- [ ] Restore the `opencode` reasoning-effort key mapping.
- [ ] Remove the now-unused error type if nothing else references it.
- [ ] Edit the break-half characterization tests and declare the breaks.

## Acceptance Criteria

- [ ] Configuration with `runtime: opencode` and a non-empty `reasoning_effort`
      loads and resolves with that effort intact.
- [ ] Configuration with `reasoning_effort: ""` on `opencode` still loads.
- [ ] Runtime validation accepts an `opencode` runtime carrying a non-empty
      effort.
- [ ] Codex and Claude still map to their existing reasoning keys.
- [ ] `grep -rn "must be empty for runtime" internal/config` finds nothing.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/agent/acpx_runner.go`

## Bounded scope

This Task may create or modify only:

- `internal/config/profiles.go`
- `internal/config/config_test.go`
- `internal/config/opencode_effort_characterization_test.go`
- `internal/config/profiles_characterization_test.go`
- `internal/agent/acpx_runner.go`
- `internal/agent/acpx_runner_test.go`
- `internal/agent/acpx_session_effort_characterization_test.go`
- `internal/agent/selection_assignment.go`
- `docs/references/coverage-record.json`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_03.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/config ./internal/agent -count=1` — expected: exits 0.
- `go test ./internal/config -run 'OpenCodeEffortAccepted' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `! grep -rq 'must be empty for runtime' internal/config` — expected: exits 0, proving the refusal text is gone rather than reworded. The leading `!` is required: a Verification command passes only by exiting 0.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.

## References

- `_prd.md` → Goal 1; Core Features: a configuration that stops lying.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 3.
- ADR-0108.

## Result

Removed both OpenCode refusal gates. Agent Selection normalization now accepts
and preserves a non-empty `reasoning_effort`, and the acpx runtime mapping uses
the generic `effort` key for both `opencode` and `opencode-custom`. Codex keeps
its `model_reasoning_effort` key, Claude keeps the generic key, and an empty
OpenCode effort still bypasses runtime assignment and plans `runtime_managed`.

The obsolete `runtimesManagingReasoning` list,
`validateModelManagedReasoning`, `ModelManagedReasoningError`, and its measured
reason constant were removed with their last uses. Spec 0088's surviving
configuration corpus now declares the later Spec 0089 break: Preferred and
Fallback OpenCode selections both load with their configured efforts intact.
The OpenCode config and Agent Session characterization tests were renamed from
the old `Today...Refuses...` contract to the accepted key/loading contract.

Pre-change signal:

- `rtk proxy env GOCACHE=<worktree>/.gocache go test ./internal/config -run '^TestOpenCodeEffortAcceptedConfigurationLoadsAndResolvesNonEmptyReasoningEffort$' -count=1 -v` — failed because normalization returned `reasoning_effort must be empty for runtime "opencode"`.
- `rtk proxy env GOCACHE=<worktree>/.gocache go test ./internal/agent -run '^(TestACPXSessionEffortCharacterizationDeclaredBreakReasoningKeyMapsOpenCode|TestReasoningEffortConfigKeyMapsSupportedRuntimes|TestValidateRuntimeSelectionAcceptsOpenCodeReasoningEffort)$' -count=1 -v` — failed because the OpenCode key mapping returned `ModelManagedReasoningError`.

Focused checks run after the last Go edit:

- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/config -run "^(TestOpenCodeEffortAcceptedConfigurationLoadsAndResolvesNonEmptyReasoningEffort|TestCharacterizationDeclaredBreakAcceptsOpenCodeReasoningEffort|TestCharacterizationDeclaredBreakAcceptsOpenCodeReasoningEffortInFallbackChain|TestCharacterizationInvariantAcceptsAnEmptyReasoningEffort)$" -count=1 -v'` — 4 tests passed.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run "^(TestACPXSessionEffortCharacterizationDeclaredBreakReasoningKeyMapsOpenCode|TestReasoningEffortConfigKeyMapsSupportedRuntimes|TestValidateRuntimeSelectionAcceptsOpenCodeReasoningEffort|TestSelectionEffortCharacterizationInvariantOpenCodeEmptyEffortPlansRuntimeManagedEncoding)$" -count=1 -v'` — 8 tests/subtests passed.
- `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record` — exited 0 and re-recorded `docs/references/coverage-record.json` after the test renames; the record was generated, not hand-edited.
- `rtk rg -n "must be empty for runtime|runtimesManagingReasoning|validateModelManagedReasoning|ModelManagedReasoningError|openCodeModelManagedReasoning" internal/config internal/agent` — produced no matches.

Acceptance evidence:

- Criterion 1: `TestOpenCodeEffortAcceptedConfigurationLoadsAndResolvesNonEmptyReasoningEffort` passed after parsing and resolving an OpenCode Preferred Selection at `high`, then asserted that `high` remained intact. The two Spec 0088 declared-break tests passed for Preferred and Fallback positions at `max`.
- Criterion 2: `TestCharacterizationInvariantAcceptsAnEmptyReasoningEffort` passed for configuration loading, and `TestSelectionEffortCharacterizationInvariantOpenCodeEmptyEffortPlansRuntimeManagedEncoding` passed with `runtime_managed`.
- Criterion 3: `TestValidateRuntimeSelectionAcceptsOpenCodeReasoningEffort` passed for both non-empty `max` and empty OpenCode efforts.
- Criterion 4: all four subtests of `TestReasoningEffortConfigKeyMapsSupportedRuntimes` passed: Codex maps to `model_reasoning_effort`, while Claude, OpenCode, and the OpenCode command override map to `effort`.
- Criterion 5: the refusal/remnant search produced no match in `internal/config` or `internal/agent`; the narrower Daemon-owned negative `grep` remains in authored Verification.

No follow-up work was found inside this Task's slice. The commands authored
under `## Verification` were not run; Daemon Verification remains the
settlement boundary.
