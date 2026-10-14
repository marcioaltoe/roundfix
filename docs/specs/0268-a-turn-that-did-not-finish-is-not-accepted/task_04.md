---
task: task_04
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: completed
type: backend
complexity: medium
---

# Task 04: Claude sessions refuse background tasks and a QA turn that leaves one is an Unfinished Turn

## Overview

A Claude QA turn started its load driver as a Claude Code background task and
ended while the driver ran; the report stayed `pending` and the Run said only
that. This Task refuses background tasks in every `claude` command Roundfix
runs, records each Agent Background Task a turn starts from the adapter's
metadata, adds the foreground rule to the QA prompt contract, and settles a QA
turn that ends with its report pending after starting one as an Unfinished
Turn that names each task. It answers the Backlog Entry
[a Claude QA turn ends while its work runs in the background](references/2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md)
of 2026-10-13. It is verifiable on its own with the fake acpx and the QA
step's scripted runner.

## Requirements

1. MUST give every acpx command the runner starts for a `claude` runtime the
   environment entry `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1`, overriding an
   inherited value, through the per-runtime session environment that every
   session command already uses, and MUST NOT give it to any other runtime
   (`_techspec.md` → Invariant 12).
2. MUST add `AgentBackgroundTask` and `ExecuteResult.BackgroundTasks` with the
   shapes of Interfaces, and record the tasks of a turn per Invariant 13: the
   id from `_meta.claudeCode.toolResponse.backgroundTaskId` or `taskId`, the
   tool from `_meta.claudeCode.toolName`, and the description from the
   `rawInput.description` of an earlier update with the same `toolCallId`.
3. MUST add the rule of Invariant 15 to the QA prompt contract in
   `internal/agent/spec_prompt.go`, as one new bullet, leaving the other
   bullets byte-identical.
4. MUST make the QA step read the turn's `BackgroundTasks` and, when the seeded
   report is still pending, publish the `turn_unfinished` event and settle the
   QA Task `failed` with the reason of Invariant 14 and API Contract 4,
   returning a pending verdict without halting. A finished report, or a
   pending report with no recorded task, keeps today's handling.
5. MUST add `internal/agent/background_task_test.go` with:
   - `TestClaudeCommandsRefuseBackgroundTasks`: a `claude` prompt through the
     fake acpx sees `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1` even when the
     runner's environment sets it to `0`, and a `codex` prompt does not see it.
   - `TestTurnRecordsItsAgentBackgroundTasks`: synthetic Claude stdout with a
     Bash call whose earlier update has `"run_in_background":true` and a
     description, whose later update carries `backgroundTaskId`, a Monitor
     call whose update carries `taskId`, and a repeated id, yields two tasks
     in start order with their tools and descriptions; a stream without them
     yields none.
   - `TestQAPromptRequiresForegroundCommands`: `BuildQAPrompt` carries the new
     rule once.
6. MUST add `internal/daemon/unfinished_turn_test.go` with
   `TestQATurnThatLeavesABackgroundTaskIsUnfinished`,
   `TestQATurnWithBackgroundTaskAndFinishedReportSettlesNormally` and
   `TestQATurnWithoutBackgroundTaskKeepsAPendingReport`, driving the QA step
   with a scripted runner result. The first asserts the event payload and the
   exact failure reason for two tasks. Each calls `t.Parallel()` first
   (ADR-0259) and waits on no elapsed time (ADR-0271).
7. MUST NOT change the Provider System Error or Lost Rollout behavior, the
   Task prompt, the review path, or any existing test other than adding cases.

## Subtasks

- [ ] Add the `claude` refusal variable to the per-runtime session environment.
- [ ] Record Agent Background Tasks in the prompt stream reader.
- [ ] Add the foreground rule to the QA prompt contract.
- [ ] Settle an Unfinished Turn in the QA step.
- [ ] Add the agent and daemon tests.

## Acceptance Criteria

- [ ] Every `claude` acpx command carries the refusal variable, and no `codex`
      command does.
- [ ] A turn's Bash and Monitor background tasks are recorded once each, with
      tool and description.
- [ ] A QA turn with a recorded task and a pending report settles `failed`
      with every task named; the other two cases are unchanged.
- [ ] Every existing `internal/agent` prompt test passes.

## Context

- creates: `internal/agent/background_task_test.go`
- creates: `internal/daemon/unfinished_turn_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/codex_spawn.go`
- interface: `internal/daemon/task_engine.go`
- instruction: `internal/agent/spec_prompt_test.go`
- instruction: `internal/agent/acpx_runner_test.go`
- instruction: `internal/daemon/lost_rollout_recovery_test.go`
- instruction: `docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md`
- instruction: `docs/adr/0057-daemon-exclusively-owns-implement-task-status.md`

## Consumers

- `ExecuteResult.BackgroundTasks` and the QA prompt contract read by the QA step — consumers: `./internal/agent`, `./internal/daemon`

## Verification

- `out="$(go test -count=1 -v -run '^(TestClaudeCommandsRefuseBackgroundTasks|TestTurnRecordsItsAgentBackgroundTasks|TestQAPromptRequiresForegroundCommands|TestBuildQAPrompt.*)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestClaudeCommandsRefuseBackgroundTasks TestTurnRecordsItsAgentBackgroundTasks TestQAPromptRequiresForegroundCommands; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so their PASS lines are missing and the command fails.
- `out="$(go test -count=1 -v -run '^(TestQATurnThatLeavesABackgroundTaskIsUnfinished|TestQATurnWithBackgroundTaskAndFinishedReportSettlesNormally|TestQATurnWithoutBackgroundTaskKeepsAPendingReport|TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQATurnThatLeavesABackgroundTaskIsUnfinished TestQATurnWithBackgroundTaskAndFinishedReportSettlesNormally TestQATurnWithoutBackgroundTaskKeepsAPendingReport TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three new tests do not exist, so their PASS lines are missing.

## References

- [_prd.md](_prd.md) — Goals 4-5; User Stories 3-4; Core Feature 3; Success Metric 4; Acceptance evidence
- [_techspec.md](_techspec.md) — Interfaces; Invariants 12-15; API Contract 4; API Contract 5; Testing Approach 5; Testing Approach 6; Build Order 4
- [2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md](references/2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md)
- ADR-0274; ADR-0057; ADR-0259; ADR-0271

## Result

Task 04 implementation is handed back for Daemon Verification. Status remains
Daemon-owned. No declared Verification command, repository gate, commit, push
or Pull Request was run.

- Claude commands receive the shared per-runtime refusal override, replacing
  inherited `0`; other runtimes receive no entry. Session listing and standalone
  capability acquisition use that policy too. Requirement 1 required the small
  additional change in `internal/agent/selection_capabilities.go`.
- The stream reader returns unique background tasks in observed start order,
  correlating sparse tool updates with earlier descriptions by `toolCallId`.
- The QA prompt adds the exact foreground bullet. A successful QA turn with a
  pending seeded report and tasks emits `turn_unfinished`, settles failed with
  every task named, and returns pending without halting.

| Acceptance criterion | Focused evidence |
| --- | --- |
| Claude refusal on every command; none for Codex | `TestClaudeCommandsRefuseBackgroundTasks` checks actual fake acpx environments for session preparation, prompt, standalone `set model`, cancel, close and list, including inherited `0` and an unchanged base environment. |
| Bash and Monitor recorded once with tool and description | `TestTurnRecordsItsAgentBackgroundTasks` checks both ID fields, sparse input correlation, repeated IDs, start order, absent descriptions, empty IDs and unrelated metadata. |
| Pending QA with tasks fails with every task named; other cases retain their handling | `TestQATurnThatLeavesABackgroundTaskIsUnfinished` asserts the exact two-task event payload and reason, failed status on disk, pending report and no Run error. The other two new tests check finished pass/fail reports and pending without tasks. |
| Existing agent prompt tests pass | The final agent check runs every test containing `Prompt`; `TestQAPromptRequiresForegroundCommands` asserts the new bullet once. |

Focused checks:

- Pre-change signal: the new agent tests initially failed to compile because
  `AgentBackgroundTask` and `ExecuteResult.BackgroundTasks` did not exist.
- The extended environment test reproduced Claude's standalone `set model`
  bypass (`absent`, expected `1`) before that command received the shared policy.
- `rtk proxy go test -count=1 -run 'BackgroundTask|Foreground|Prompt|CapabilityAcquisition|^TestACPXCommandEnv' ./internal/agent` — exit 0 after final implementation edits.
- `rtk proxy go test -count=1 -run '^Test(QATurn|QA.*(LostRollout|ProviderSystemError))' ./internal/daemon` — exit 0 after final implementation edits.
- `rtk proxy go test -count=1 -run 'LostRollout|ProviderSystemError' ./internal/daemon` — exit 0.
- `rtk proxy git diff --check` — exit 0.

The new daemon tests call `t.Parallel()` first and wait on no elapsed time.
Existing tests, Task prompt, review path, recovery behavior, other Task files
and the Task Graph are unchanged. Runtime evidence uses synthetic stdout and
fake processes; no live Claude session was used.


### Verification Feedback repair

The Daemon's attempt 1 exposed an executable shell wrapper in the new
background-task environment test. That write violated the frozen executable
inventory. Replaced it with the repository's existing
`testfixture.FixtureBinary` helper: the Go tool produces the executable, and
the compiled wrapper records the fixed policy value and waits for its fake
acpx child. The assertions and production behavior are unchanged; this repair
changes only `internal/agent/background_task_test.go` and this Result section.

Focused checks after the repair:

- `rtk proxy go test -count=1 -run '^TestClaudeCommandsRefuseBackgroundTasks$' ./internal/agent` — exit 0; all captured Claude/Codex command environments retain the required behavior.
- `rtk proxy go test -count=1 -run '^TestNoTestWritesAnExecutableOutsideTheResidue$' ./internal/testfixture` — exit 0; the contract reported by the Daemon now passes.
- `rtk proxy git diff --check` — exit 0.

Task status remains unchanged. The configured repository Verification and the
Task's declared Verification remain reserved for the Daemon's rerun.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/agent/selection_capabilities.go`
