---
task: task_04
spec: 0002-acpx-migration
status: pending
type: backend
complexity: medium
---

# Task 04: Wire the Agent Session through both Run paths

## Overview

Connect the acpx runner to real Runs: the CLI derives the Agent Session name from the run id, both engine paths (review resolve cycle and spec TaskCycle) carry it into every Agent invocation, and every terminal Run outcome — Clean, Unresolved, Failed, Stopped, on both paths — closes the session. After this task the acpx runner is the live default; the SDK code still exists but nothing invokes it (deleted in task_05). Verifiable through the existing CLI and daemon suites over runner fakes.

## Requirements

1. MUST derive one Agent Session per Run, named `roundfix-<run-id>`, created lazily at the Run's first Agent work and shared by every Batch and Task in that Run (ADR-0018); sessions are never reused across Runs.
2. MUST thread the session reference through both engine plans into `ExecuteRequest` without changing any engine decision logic, prompts, statuses, commit behavior, or events.
3. MUST close the session (`EndSession`) on every terminal path of both engines — including preflight-passed-but-failed Runs and Stop Requests — with best-effort semantics that never mask the Run's real outcome; journal the session start and close through existing event kinds (no new Run Event vocabulary).
4. MUST switch the default runner wiring to the acpx runner (probe included) for fetch-with-agent-free paths untouched: fetch never starts an Agent and gains nothing; resolve, watch, and implement all run through acpx sessions.
5. MUST keep the full behavioral contract byte-identical: stdout lines, exit codes, Run states, journal event kinds, and both Batch contracts — the existing CLI and daemon suites pass with only mechanical fake extensions for the session surface.

## Subtasks

- [ ] Session name derivation and lazy creation at first Agent work
- [ ] Session reference through the resolve cycle and TaskCycle plans
- [ ] `EndSession` on every terminal path of both paths, including stop
- [ ] Default runner switch to acpx with probe at the existing call sites
- [ ] Session lifecycle journaled via existing event kinds

## Acceptance Criteria

- [ ] A multi-Batch review Run and a multi-Task spec Run each show exactly one ensure per Run and one close at the end (asserted via runner fakes recording session calls).
- [ ] Terminal-path matrix tests: Clean, Unresolved, Failed, and Stopped each close the session on both paths; a failing close never changes the Run outcome or exit code.
- [ ] Watch Runs reuse one session across their Rounds' resolve cycles within the same Run and close it at the watch outcome.
- [ ] The full existing suite passes with unchanged assertions on stdout, exit codes, states, and events.

## Verification

- `rtk go test ./internal/cli/ ./internal/daemon/` — expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 2; Core Features 2, 8; Success Metrics (one spawn per Run). `_techspec.md` → System Architecture (edges), Coverage Map, Build Order 4. ADR-0018.
