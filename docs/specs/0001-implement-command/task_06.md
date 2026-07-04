---
task: task_06
spec: 0001-implement-command
status: pending
type: backend
complexity: high
---

# Task 06: Ship the Implement Command end to end

## Overview

Wire the Implement Command into the CLI: flags, Preflight Validation, Run creation, cockpit, deterministic output, exit codes, and resume. After this task a developer can execute a Spec non-interactively from the command line — the user-visible spine of the feature, verifiable through buffer-captured CLI tests over the existing collaborator seams.

## Requirements

1. MUST accept `--spec <slug>` plus the shared agent flags (`--agent`, `--model`, `--agent-command`, `--agent-full-access`, `--interactive`/`--no-input`) with the existing config-then-flag precedence; `--qa` arrives in a later task.
2. MUST run Preflight Validation before creating a Run, each failure exiting 2 with one actionable message naming the check and the fix: Spec exists and is active; Task Graph valid and acyclic; every task file parses with at least one Verification command; working tree clean; current branch is not the detected repository default; no Active Run for this work target and none in this working tree (naming the blocking run id and `roundfix stop`); runtime probe passes.
3. MUST create the Run with the `implement` Kind and the spec slug, start the Live Run View exactly as the existing operational commands do, execute the cycle, and complete the Run with the terminal outcome: every Task completed → `Clean`, any non-completed Task → `Unresolved`, Stop Request → `Stopped`, infrastructure error → `Failed`.
4. MUST print to stdout only: one line per Task in graph order — `task_NN <status> — <title>` with status `completed | failed | skipped | pending` — and one final outcome line; diagnostics, progress, Run id, and the agent log path go to stderr.
5. MUST map outcomes to the existing exit codes: 0 Clean, 1 Unresolved/Failed, 2 Preflight Validation, 130 Stop Request.
6. MUST resume by construction: a new Run picks up every non-completed Task, including a stale `in_progress` left by a dead Run; when every Task is already completed, report that and exit 0 without creating a Run.
7. MUST update the top-level usage and command help text truthfully; the dirty-worktree message MUST say how a failed attempt's leftovers are cleared (commit, stash, or discard).

## Subtasks

- [ ] Flag parsing and config precedence following the operational-command shape
- [ ] Preflight Validation wiring with per-check actionable messages
- [ ] Run creation, cockpit start, cycle execution, terminal outcome and completion
- [ ] Deterministic stdout contract and stderr diagnostics
- [ ] Resume and all-completed no-op behavior
- [ ] Usage and help text

## Acceptance Criteria

- [ ] A buffer-captured end-to-end run over a fake Agent completes a multi-Task Spec: correct stdout lines in graph order, outcome line, exit 0, one commit per Task, journaled Run reaching `Clean`.
- [ ] Each preflight failure from Requirement 2 has a test asserting exit 2 and its named check and fix; nothing is written to the Run Database in these cases.
- [ ] An induced Task failure ends the Run `Unresolved` with exit 1; a second invocation executes only the non-completed Tasks and finishes the graph (PRD resume metric), including the stale-`in_progress` case.
- [ ] A canceled run ends `Stopped` with exit 130; the lock is released on every terminal path.
- [ ] Implement help text lists exactly the implemented flags; the full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix implement --help` — expected: help showing `--spec` and the shared agent flags, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 3, 4; Core Features 1, 2, 8; User Experience. `_techspec.md` → API Contracts, Build Order 8, Risks (dirty worktree after failure). ADR-0012, ADR-0013.
