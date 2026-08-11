---
task: task_03
spec: 0088-a-third-runtime-that-can-run
status: completed
type: backend
complexity: medium
---

# Task 03: Make OpenCode reasoning effort model-managed

## Overview

Roundfix applies a non-empty reasoning effort as a separate ACP config option
after ensuring the Agent Session, and on `opencode` that call cannot succeed
before the Run's first prompt — the live agent process is still on the runtime
default and advertises no `effort` option, so the adapter answers ACP `-32602`.
This Task makes `opencode` a model-managed reasoning runtime: configuration
refuses a non-empty effort with the repair named, and runtime validation refuses
it again so no invocation override slips past.

## Requirements

1. MUST stop mapping `opencode` to a reasoning-effort config key, so runtime
   validation rejects a non-empty effort on that runtime for every entry point,
   including a `--reasoning-effort` invocation override.
2. MUST refuse a non-empty `reasoning_effort` on `runtime: opencode` during Agent
   Selection normalization, so a maintainer sees it when the configuration loads
   rather than when a Run starts.
3. MUST name the repair in the refusal text: the empty value, and why — OpenCode
   advertises reasoning effort per model and only after an Agent Session's first
   prompt, so a token-free Exact Agent Selection Proof cannot apply one.
4. MUST apply the same refusal to the legacy runtime-defaults path, so a legacy
   `runtimes.opencode.reasoning_effort` cannot produce a profile that
   normalization would reject.
5. MUST keep an `opencode` Agent Selection with an empty `reasoning_effort` fully
   valid, and MUST NOT let it reach the disposable-effort application path.
6. MUST NOT change how any Codex or Claude selection is validated, mapped, or
   applied.
7. MUST keep an empty reasoning effort provable on a runtime whose adapter
   advertises a reasoning option Roundfix declines to assign, without weakening
   the rule for a runtime Roundfix does control.
8. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [x] Remove `opencode` from the reasoning-effort config key mapping.
- [x] Add the normalization refusal with its repair text.
- [x] Cover the legacy runtime-defaults path with the same refusal.
- [x] Confirm an empty-effort `opencode` selection skips effort application on
      both the disposable and the fallback paths.
- [x] Edit the break-half characterization test that pinned today's acceptance,
      and declare the break in this Task's Result.

## Acceptance Criteria

- [x] Decoding a profile with `runtime: opencode` and a non-empty
      `reasoning_effort` fails, and the message names the empty value as the
      repair and states why.
- [x] Decoding the same profile with `reasoning_effort: ""` succeeds.
- [x] Runtime validation rejects an `opencode` runtime carrying a non-empty
      reasoning effort, independently of configuration decoding.
- [x] A legacy `runtimes.opencode.reasoning_effort` that is non-empty is refused
      with the same repair.
- [x] Codex and Claude selections with a non-empty effort still validate and
      still map to their existing config keys.
- [x] An `opencode` selection with an empty effort issues no reasoning config set
      on the disposable-session path.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/fallback.go`
- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/acpx_runner.go`
- `internal/agent/acpx_runner_test.go`
- `internal/agent/selection_assignment.go`
- `internal/agent/fallback.go`
- `internal/config/profiles.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/config/profiles_characterization_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0088-a-third-runtime-that-can-run/task_03.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent ./internal/config -count=1` — expected: exits 0.
- `go test ./internal/config -run 'OpenCodeReasoning' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `go test ./internal/agent -run 'ReasoningEffortConfigKey' -count=1 -v` — expected: exits 0 and names at least one test.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.
- `grep -q 'model-managed' internal/config/profiles.go` — expected: exits 0, proving the repair is named in the refusal text rather than only in the Spec.

## References

- `_prd.md` → Goal 4; Core Features: model-managed reasoning for OpenCode.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 3.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the nine-step sequence proving the effort cannot be applied before the first
  prompt.
- ADR-0106.

## Result

OpenCode is now a model-managed reasoning runtime, refused at configuration and
at runtime validation, and its selections prove.

**What changed.** `acpxReasoningEffortConfigKey` no longer maps `opencode`; it
returns a typed `ModelManagedReasoningError` carrying the measured reason and
naming the empty value as the repair, classified as
`reasoning_control_not_advertised`. Because `validateRuntimeSelection` consults
that mapping, a `--reasoning-effort` invocation override is refused on the same
truth rather than on a second copy of it. `normalizeSelection` refuses the same
selection earlier, at configuration load, where the maintainer can act on it;
the refusal names the YAML path, the runtime, the measured reason, and
`reasoning_effort: ""`. The legacy runtime-defaults path is covered without new
code, because `Validate` runs `validateProfiles` over the profiles that path
produces.

**Scope amendment.** `internal/agent/selection_assignment.go` was added to this
Task's bounded scope during implementation, and the TechSpec and ADR-0106 were
amended to record why. The measured OpenCode session advertises a per-model
`effort` option, and `selectionStateMatches` required `state.ReasoningOption ==
nil` for the `model_managed` encoding — so an empty-effort OpenCode selection
planned and then failed its own effective-state check. Collapsing that rule for
everyone would have let an unassigned reasoning option pass the Codex and Claude
proof, so the encodings were split instead: `model_managed` keeps the strict
rule, and the new `runtime_managed` proves against a runtime whose advertised
control Roundfix declines to assign.

**Declared breaks.**

1. `TestCharacterizationTodayAcceptsOpenCodeReasoningEffort` became
   `TestCharacterizationDeclaredBreakRefusesOpenCodeReasoningEffort`. The old
   behavior is preserved in that test's comment with its measured provenance.
2. The `opencode effort` case of `TestACPXRunAppliesSelectionBeforePrompt`
   became `opencode model-managed reasoning issues no effort set`: it now
   asserts a Run issues **no** reasoning config set between the Agent Session
   and the prompt, which is the behavior this Task delivers.
3. `TestACPXRunAppliesFullAccessSessionSetup`'s opencode case lost its
   `set effort high` command, and the test helper stopped defaulting an
   `opencode` RuntimeSpec to `high`.

**Commands and outcomes.**

- `go build -buildvcs=false ./...` — exit 0.
- `go test ./internal/agent ./internal/config -count=1` — exit 0.
- `go test ./internal/config -run 'OpenCodeReasoning' -count=1 -v` — exit 0.
- `go test ./internal/agent -run 'ReasoningEffortConfigKey' -count=1 -v` — exit 0; four subtests.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — exit 0.
- `grep -q 'model-managed' internal/config/profiles.go` — exit 0.
- `make verify` — exit 0 after `go clean -testcache`.

**Evidence per acceptance criterion.**

- Non-empty OpenCode effort refused with the repair named:
  `TestCharacterizationDeclaredBreakRefusesOpenCodeReasoningEffort` asserts the
  YAML path, `must be empty for runtime "opencode"`, and `model-managed
  reasoning` in one error.
- Empty effort still valid:
  `TestCharacterizationInvariantAcceptsAnEmptyReasoningEffort`.
- Runtime validation refuses independently of configuration decoding:
  `TestValidateRuntimeSelectionRefusesOpenCodeReasoningEffort`.
- Refusal is positional-agnostic:
  `TestOpenCodeReasoningRefusalNamesTheRepairOnEveryFallbackPosition` places the
  offending selection in a Fallback Chain.
- Codex and Claude unaffected: `TestReasoningEffortConfigKeyRefusesOpenCode`
  asserts both still map to their existing keys, and
  `TestCharacterizationInvariantAcceptsCodexAndClaudeReasoningEffort` keeps
  their configuration valid.
- No reasoning config set on an OpenCode Run: the rewritten
  `TestACPXRunAppliesSelectionBeforePrompt` case asserts the exact command
  sequence, ensure then prompt, with nothing between them.
- Requirement 7: `TestRuntimeManagedReasoningProvesAgainstAnAdvertisedEffortOption`
  proves an OpenCode selection against an advertised `effort` option and, in the
  same test, asserts a Claude selection over the identical payload still fails
  the strict `model_managed` rule.

**Follow-ups.** None. The coverage record needed no re-recording: every renamed
test was introduced by this Spec and re-recorded in Task 02's commit, and
`TestCoverageEquivalence` passes.
