---
task: task_06
spec: 0037-terminal-outcome-integrity
status: completed
type: docs
complexity: medium
---

# Task 06: Align terminal-outcome operator guidance

## Overview

Publish the completed Force Stop, graceful Stop Request, terminal replay, and
cleanup contracts through supported operator documentation and command help.
This documentation slice leaves every protected Skill path unchanged so the
isolated tooling Task can publish the final Agent-facing wording.

## Requirements

1. MUST document that Force Stop reports Stopped only after owner exit proof.
2. MUST document the failed-closed result, retained Active Run lock, and
   actionable inspection or retry guidance.
3. MUST document graceful Stop Request observation during Review Source waits.
4. MUST document idempotent same-outcome replay and conflicting-outcome
   rejection without exposing internal storage details as user API.
5. MUST describe registered Agent Session cleanup and primary-before-secondary
   diagnostics.
6. MUST preserve canonical glossary terms and finding/ADR traceability.
7. MUST leave all protected tooling files unchanged.

## Subtasks

- [x] Update Stop Command and operational usage guidance.
- [x] Align command help with proof-before-completion behavior.
- [x] Document graceful watch interruption and terminal replay.
- [x] Document cleanup eligibility and diagnostic ordering.
- [x] Resolve glossary, ADR, Spec, and finding links.
- [x] Check supported examples and terminology.

## Acceptance Criteria

- [x] A reader can predict successful and failed Force Stop state and lock
      behavior.
- [x] Guidance states when a graceful watch stop is observed and what work will
      not run afterward.
- [x] Registered-session cleanup and secondary-warning behavior are explicit.
- [x] User docs do not promise terminal overwrite or warning-only owner
      reclamation.
- [x] Every cited ADR, finding, and command reference resolves.
- [x] No protected tooling path changes in this Task.

## Context

- instruction: `docs/agents/cli.md`
- instruction: `docs/agents/domain.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `CONTEXT.md`
- interface: `README.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/findings/2026-07-16-vortex-pr87-detached-watch-notification.md`

## Verification

- `rtk grep -n 'stop --force\|Force Stop\|Stop Request' docs/user-guide/commands.md docs/user-guide/usage.md`
  — expected: supported guidance covers force and graceful stop contracts.
- `rtk go test ./internal/cli -run 'Test.*Stop.*(Help|Usage|Report)' -count=1`
  — expected: command help and public reports match documented behavior.
- `rtk git diff --check`
  — expected: documentation contains no whitespace errors.

## References

- `_prd.md` → User Stories 1–5; User Experience; Non-Goals; Decisions.
- `_techspec.md` → API Contracts; Risks & Considerations; Build Order 6.
- `../../adr/0052-run-completion-is-compare-and-set.md` → public terminal
  integrity contract.

## Result

Aligned the Stop Command reference, operational guide, public help, and Force
Stop success report with proof-before-completion behavior. The guidance now
covers graceful Review Source wait interruption, failed-closed Force Stop
recovery, stable terminal replay, registered Agent Session cleanup, and
primary-before-secondary diagnostics without exposing internal persistence as
a user API.

### Acceptance criterion evidence

- Successful Force Stop guidance states that Roundfix proves the recorded owner
  exited before reporting Stopped or releasing the Active Run lock. Failure
  guidance states that the Run remains Active, the lock stays retained, and the
  operator can inspect Active Runs and retry the exact Force Stop command.
- Graceful watch guidance names the next configured poll boundary and states
  that no later fetch, check, commit, push, or Review Source mutation runs after
  the Stop Request is observed.
- Cleanup guidance limits cancellation to registered active Agent Sessions,
  treats an already-absent registered session as idempotent, and places other
  cleanup failures after the primary failure as secondary warnings.
- Replay guidance reports an already Stopped outcome idempotently and rejects a
  different terminal outcome unchanged. The obsolete immediate-completion and
  warning-only reclamation promises are absent.
- The glossary, ADR-0052, Spec 0037, Stop Command, and detached-watch finding
  targets exist; the cited finding heading resolves to its documented anchor.
- `rtk git diff --name-only -- .agents/skills skills` returned no paths, proving
  that protected tooling remains unchanged.

### Verification evidence

- `rtk grep -n 'stop --force\|Force Stop\|Stop Request' docs/user-guide/commands.md docs/user-guide/usage.md`
  — passed; both supported guides contain the graceful and Force Stop
  contracts.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test ./internal/cli -run 'Test.*Stop.*(Help|Usage|Report)' -count=1`
  — passed. The task-local cache was required because the managed sandbox
  cannot write the default Go build cache.
- `rtk git diff --check` — passed after the final implementation and test edit.
- `rtk ls CONTEXT.md docs/adr/0052-run-completion-is-compare-and-set.md docs/specs/0037-terminal-outcome-integrity/_prd.md docs/findings/2026-07-16-vortex-pr87-detached-watch-notification.md docs/user-guide/commands.md`
  — passed; all cited files resolve.
- `rtk rg -n '^## 4\\. Cleanup noise appeared before the actionable failure$' docs/findings/2026-07-16-vortex-pr87-detached-watch-notification.md`
  — passed; the cited finding anchor source heading exists.

The Daemon owns the authoritative execution of this Task's declared
Verification after the Agent turn.

### Verification retry — attempt 1

The Daemon exposed an escaping defect in the declared grep pattern: doubled
backslashes inside the single-quoted basic regular expression searched for
literal separators instead of alternation. The Task now declares single
backslashes, and the corrected command exits `0` against both supported user
guides. No product or operator-guidance change was required for this retry.
