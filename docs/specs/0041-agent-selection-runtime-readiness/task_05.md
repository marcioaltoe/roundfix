---
task: task_05
spec: 0041-agent-selection-runtime-readiness
status: pending
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

- [ ] Parse selection-flag presence as one atomic value.
- [ ] Reject every partial flag combination before preflight side effects.
- [ ] Preserve explicit empty reasoning and complete tuple values.
- [ ] Apply complete overrides without replacing fallback chains.
- [ ] Update CLI help and deterministic usage errors.
- [ ] Cover resolve, watch, implement, detached, and Agent-free paths.

## Acceptance Criteria

- [ ] Bare `--agent`, bare `--model`, bare `--reasoning-effort`, and every
      two-flag subset exit `2` with the same actionable grammar explanation.
- [ ] Partial overrides create no Agent Session, Run, worktree, artifact, or
      configuration change.
- [ ] Omitting all three flags selects the effective category profiles for
      Task, QA, and review work.
- [ ] A complete Sol/high tuple overrides only Preferred Selection and retains
      the configured fallback sequence.
- [ ] An explicitly empty reasoning value reaches selection proof as an empty
      but present value.
- [ ] Help for all Agent-starting commands documents profile-led invocation and
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

