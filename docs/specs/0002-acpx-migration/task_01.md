---
task: task_01
spec: 0002-acpx-migration
status: completed
type: backend
complexity: high
---

# Task 01: Build the acpx invocation core

## Overview

Create the acpx-backed runner core inside the agent package: constructing acpx command lines, streaming `--json-strict` NDJSON into the existing Run Event conversion, writing the agent log, and mapping acpx exit codes to Roundfix outcomes. This task builds and tests the core in isolation behind a scripted fake acpx; nothing is wired into Runs yet and the existing SDK runner keeps working untouched. Verifiable alone through the helper-process test rig.

## Requirements

1. MUST construct prompt invocations exactly per the TechSpec's acpx invocation mapping: `acpx <agent> prompt -s <session> --cwd <workdir> --format json --json-strict --approve-all [--model <id>] -f -` with the built prompt delivered on stdin, and the runtime always named explicitly (never acpx's implicit default).
2. MUST map the command override escape hatch to the global `--agent "<command>"` in place of the adapter name, preserving the existing stdio escape hatch semantics.
3. MUST parse each stdout line as a raw ACP JSON-RPC message: `session/update` notifications flow into the existing stream→Run Event conversion with the raw line as the journaled payload (ADR-0008 unchanged); the `session/prompt` response line yields the stop reason for the execute result.
4. MUST append every stdout line to the existing agent log path convention, and capture stderr only into error context (never parsed, never journaled as events).
5. MUST implement the exit-code mapping: 0 → success; 1 → Batch failure (agent/protocol error); 3 → Batch failure with timeout reason; 5 → Batch failure journaled loudly (unreachable under `--approve-all`); 2 and 4 → infrastructure errors that fail the Run; 130 → Stop Request semantics.
6. MUST build the helper-process test rig with stdlib only: the test binary re-execs itself as a scripted fake acpx selected by environment (canned NDJSON on stdout, chosen exit code, invocation args captured to a file for assertions).
7. MUST NOT change the `Runner` interface, any CLI/daemon wiring, or the existing SDK runner in this task.

## Subtasks

- [x] Command-line construction for prompt invocations and the command override
- [x] NDJSON parsing into the existing stream→Run Event conversion with raw payloads
- [x] Agent log writing and stderr capture
- [x] Exit-code mapping table
- [x] Helper-process fake-acpx test rig

## Acceptance Criteria

- [x] Command-construction tests assert every flag and its order for: default adapter, model set, command override, prompt via stdin.
- [x] A scripted NDJSON stream produces Run Events whose payloads are byte-identical to the emitted lines, and the stop reason lands in the execute result.
- [x] The agent log contains every stdout line in order; a stream mixing valid updates and the final response journals only the updates.
- [x] Each exit code in the mapping table has a test asserting the classified outcome (Batch failure vs infrastructure error vs stop), including the timeout reason for exit 3.
- [x] No production file outside the agent package changes; the full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/agent/` — expected: all tests pass, including the helper-process rig.
- `rtk go build ./...` — expected: builds cleanly.
- `rtk go test ./...` — expected: full suite passes unchanged.

## References

`_prd.md` → User Stories 5, 6; Core Features 3, 5, 6. `_techspec.md` → Interfaces, acpx invocation mapping, Stream and journaling, Exit-code mapping, Testing Approach, Build Order 1. ADR-0008, ADR-0010, ADR-0017.

## Result

Implemented an isolated `internal/agent` acpx invocation core. The core builds explicit acpx prompt commands, maps command overrides through global `--agent`, sends the prompt on stdin, parses stdout as JSON-RPC NDJSON, journals only `session/update` lines as byte-exact Run Event payloads, appends every stdout line to the Agent log, stores the final `session/prompt` stop reason in `ExecuteResult.StopReason`, and classifies documented acpx exit codes into Batch failure, infrastructure error, or Stop Request outcomes.

Fresh verification:

- `rtk go test ./internal/agent/` — passed, 52 tests in 1 package.
- `rtk go build ./...` — passed.
- `rtk go test ./...` — passed, 421 tests in 16 packages.
- `rtk make verify` — passed: full suite, `roundfix skills check`, and build.

Acceptance evidence:

- Command construction: `TestACPXPromptArgsMatchTechSpecOrder` covers default adapter, model flag placement, and command override; `TestACPXRunPromptSendsPromptOnStdin` proves prompt delivery via stdin and captured argv order.
- Stream and journaling: `TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason` proves emitted update lines are stored byte-identically in Run Events, the final response is not journaled as an event, and `end_turn` lands in `ExecuteResult.StopReason`.
- Agent log and stderr handling: the same stream test proves every stdout line is written to the log in order; `TestACPXExitCodeMapping/usage_infrastructure_error` proves stderr appears only in error context, not the Agent log or Run Events.
- Exit mapping: `TestACPXExitCodeMapping` covers exits 0, 1, 2, 3, 4, 5, and 130, including timeout reason for exit 3 and the loud permission-denied status event for exit 5.
- Scope: `git status --short` after implementation shows production changes only under `internal/agent`; no CLI, daemon, `Runner` interface, SDK runner wiring, Task Graph manifest, or other task files were edited.

Follow-up notes:

- Wiring `ACPXRunner` into real Runs, Agent Session lifecycle, cancellation commands, and preflight/version checks remain deferred to later tasks in this Spec.
