---
task: task_05
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: medium
---

# Task 05: Enforce complete one-Run Agent Selection overrides

## Overview

Make one-Run Agent Selection overrides atomic across `resolve`, `watch`, and
`implement`. Omitting all selection flags must use category profiles; providing
any proper subset must fail before proof or mutation; a complete tuple may
replace only the Preferred Selection.

## Requirements

1. MUST treat `--agent`, `--model`, and `--reasoning-effort` as an all-or-none
   override on every Agent-starting command.
2. MUST count an explicitly empty `--reasoning-effort ""` as present and
   preserve it as model-managed intent.
3. MUST reject every non-empty proper subset with exit `2` before adapter
   proof, Session creation, Run persistence, worktree creation, or config load
   side effects.
4. MUST resolve category profiles unchanged when all three flags are absent.
5. MUST replace only the Preferred Selection when all three flags are present
   and preserve the configured Fallback Chain.
6. MUST update command help and usage errors to explain the two valid forms.
7. MUST leave Agent-free commands, including `fetch`, unchanged.

## Subtasks

- [x] Parse selection-flag presence as one atomic value.
- [x] Reject every partial flag combination before preflight side effects.
- [x] Preserve explicit empty reasoning and complete tuple values.
- [x] Apply complete overrides without replacing fallback chains.
- [x] Update CLI help and deterministic usage errors.
- [x] Cover resolve, watch, implement, detached, and Agent-free paths.

## Acceptance Criteria

- [x] Bare `--agent`, bare `--model`, bare `--reasoning-effort`, and every
      two-flag subset exit `2` with the same actionable grammar explanation.
- [x] Partial overrides create no Agent Session, Run, worktree, artifact, or
      configuration change.
- [x] Omitting all three flags selects the effective category profiles for
      Task, QA, and review work.
- [x] A complete Sol/high tuple overrides only Preferred Selection and retains
      the configured fallback sequence.
- [x] An explicitly empty reasoning value reaches selection proof as an empty
      but present value.
- [x] Help for all Agent-starting commands documents profile-led invocation and
      complete overrides; `fetch` behavior remains byte-stable.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/cli/profile_preflight.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/selection_test.go`
- interface: `internal/cli/implement_test.go`

## Verification

- `rtk go test ./internal/cli -run 'TestInvocationProfileOverride' -count=1` — expected: no-flags, all partial subsets, complete tuple, and explicit-empty reasoning cases pass.
- `rtk go test ./internal/cli -run 'Test(RunResolve|RunWatch|RunImplement).*SelectionOverride' -count=1` — expected: all three commands preserve profiles/fallbacks and reject partial overrides without side effects.
- `rtk go test ./internal/cli -run 'Test(CommandUsage|RunFetch)' -count=1` — expected: help describes the atomic grammar and Agent-free behavior is unchanged.

## References

- `_prd.md` → User Stories 5 and 6; Core Feature 7; User Experience; Success
  Metrics.
- `_techspec.md` → One-Run Override Grammar; Build Order 5; Risks and
  Considerations.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → complete tuple override requirement.

## Result

Implemented one atomic presence contract for `--agent`, `--model`, and
`--reasoning-effort`. Resolve, watch, and implement now reject a partial tuple
before configuration loading or detached execution, while a complete tuple
replaces only each category profile's Preferred Selection. Interactive Input
marks the runtime present when it supplies model or reasoning values. Fetch
does not enter the selection-validation path.

Acceptance evidence:

- `TestInvocationProfileOverrideRequiresCompleteTuple` covers all three
  one-flag and all three two-flag subsets with the same grammar error.
  `TestRunResolveSelectionOverrideRejectsPartialBeforeConfigLoad`,
  `TestRunWatchSelectionOverrideRejectsPartialBeforeConfigLoad`, and
  `TestRunImplementSelectionOverrideRejectsPartialBeforeConfigLoad` observe
  exit `2`, unchanged invalid User Config bytes, no Run Database, and no Run
  Worktree root; the resolve coverage includes `--detach`.
- `TestInvocationProfileOverrideOmittedUsesTaskQAAndReviewProfiles` proves that
  an omitted override resolves Task, QA, and review profiles.
- `TestInvocationProfileOverrideAppliesAcrossCategoriesPreservesFallbacksAndWarns`
  proves Sol/high replaces Preferred Selection while every configured Fallback
  Chain remains unchanged.
- `TestInvocationProfileOverrideParsingPreservesExplicitEmptyReasoning` and the
  resolve, watch, and implement model-managed persistence tests prove that an
  explicit empty reasoning value remains present and reaches selection proof.
- `TestCommandUsageDocumentsProfileLedAndCompleteSelectionOverrides` proves the
  two valid invocation forms for every Agent-starting command and excludes
  Agent Selection flags from fetch help. Existing `TestRunFetch*` coverage
  remains green.

Verification:

- `rtk go test ./internal/cli -run 'TestInvocationProfileOverride' -count=1`
  — passed, 14 tests.
- `rtk go test ./internal/cli -run 'Test(RunResolve|RunWatch|RunImplement).*SelectionOverride' -count=1`
  — passed, 12 tests.
- `rtk go test ./internal/cli -run 'Test(CommandUsage|RunFetch)' -count=1`
  — passed, 12 tests.
- `rtk go test ./internal/cli -count=1` — passed, 574 tests.
- `rtk make verify` — passed: 1,653 Go tests in 20 packages, 79 Python tests,
  Roundfix skill check, and CLI build.

Follow-ups: none for this Task slice.
