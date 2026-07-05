---
task: task_06
spec: 0001-implement-command
status: completed
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

- [x] Flag parsing and config precedence following the operational-command shape
- [x] Preflight Validation wiring with per-check actionable messages
- [x] Run creation, cockpit start, cycle execution, terminal outcome and completion
- [x] Deterministic stdout contract and stderr diagnostics
- [x] Resume and all-completed no-op behavior
- [x] Usage and help text

## Acceptance Criteria

- [x] A buffer-captured end-to-end run over a fake Agent completes a multi-Task Spec: correct stdout lines in graph order, outcome line, exit 0, one commit per Task, journaled Run reaching `Clean`.
- [x] Each preflight failure from Requirement 2 has a test asserting exit 2 and its named check and fix; nothing is written to the Run Database in these cases.
- [x] An induced Task failure ends the Run `Unresolved` with exit 1; a second invocation executes only the non-completed Tasks and finishes the graph (PRD resume metric), including the stale-`in_progress` case.
- [x] A canceled run ends `Stopped` with exit 130; the lock is released on every terminal path.
- [x] Implement help text lists exactly the implemented flags; the full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix implement --help` — expected: help showing `--spec` and the shared agent flags, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 3, 4; Core Features 1, 2, 8; User Experience. `_techspec.md` → API Contracts, Build Order 8, Risks (dirty worktree after failure). ADR-0012, ADR-0013.

## Result

`roundfix implement` now runs a Spec end to end. New `internal/cli/implement.go` follows the `runOperationalCommand` shape: flags parse over config defaults (`--spec`, `--agent`, `--model`, `--agent-command`, `--agent-full-access`, `--interactive`, `--no-input`; no `--qa` yet), Preflight Validation runs in the documented order (Spec active and graph valid via `spec.Load` typed errors, clean working tree with the commit/stash/discard guidance, default-branch veto naming branch, detected default, and detection source, `ActiveRunInGitRoot` naming the blocking run id and `roundfix stop <run-id>`, runtime probe), and nothing touches the Run Database until every check passes. The Run is created with Kind `implement`, the cockpit starts through `startRunUI` exactly like resolve, `Engine.TaskCycle` executes the graph with `QA: false`, and terminal outcomes map to `Clean`/`Unresolved`/`Stopped`/`Failed` with exits 0/1/2/130 through the existing interrupt mapping. stdout carries only the per-Task lines in graph order (`task_NN <status> — <title>`, statuses `completed|failed|skipped|pending`) plus one outcome line; Run id, agent log path, and diagnostics go to stderr. All-completed Specs report and exit 0 without creating a Run; resume executes only non-completed Tasks, including stale `in_progress`. The engine wires `Pusher`/`Source` only to satisfy `NewEngine` construction — the Task cycle never invokes them (asserted in the e2e test).

Commands run: `rtk go test ./internal/cli/` → 115 passed; `rtk go run ./cmd/roundfix implement --help` → help with exactly the implemented flags, exit 0; `rtk go test ./...` → 377 passed in 16 packages; `make verify` → fmt-check, tests, `skills check`, and build all green.

Evidence per acceptance criterion: e2e → `TestRunImplementExecutesSpecEndToEnd` (deterministic stdout equality, one commit per Task with frontmatter-derived messages and trailers, Run row reaches Clean, lock released); preflight → `TestRunImplementValidationFailures` (8 cases) plus `TestRunImplementPreflightFailures` (spec not found, inactive, missing Verification, cycle, dirty tree, default-branch veto), `TestRunImplementPreflightRejectsActiveRunInWorkingTree`, and `TestRunImplementPreflightProbeFailureCreatesNoRun`, each asserting exit 2, the named check and fix, and zero Run rows; resume → `TestRunImplementFailedTaskEndsUnresolvedAndResumeFinishesGraph` (exit 1, second invocation runs only task_02/task_03) and `TestRunImplementResumesStaleInProgressTask`; stop → `TestRunImplementStopRequestEndsStoppedWithInterruptMapping` (exit 130 via `exitForInterrupt`, Stopped in the store, lock released); help → `TestRunImplementHelpListsExactlyImplementedFlags` and `TestRunHelpListsImplementCommand`; full existing suite unchanged.

Follow-ups for later tasks: task_07 adds the `--qa` flag in `parseImplementCommand` and flips `TaskPlan.QA` in `executeImplementCycle` (plus the QA verdict line in the stdout report); task_08 hooks Interactive Input between parse and validate (mirror `maybeCollectInteractiveInput`, replace the two "not available yet" validation messages and the help-text caveat, and consider `rememberInteractiveDefaults`, which implement intentionally skips today); task_09/task_10 can read the Run row fields already populated here (`Kind=implement`, `SpecSlug`, `GitRoot`, `LocalBranch`, `HeadSHA`; PR fields empty) and replace the placeholder `implementLiveRunView` header, which currently shows the slug in the Repository slot and no Task Work Items. The `.agents/skills/roundfix` update for implement ships with Build Order item 11 before the feature PR opens, per the repo hard rule.
