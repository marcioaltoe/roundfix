---
task: task_05
spec: 0007-command-lifecycle
status: completed
type: backend
complexity: medium
---

# Task 05: Push spec Runs at Clean via Project Config

## Overview

Implement ADR-0021: a Project Config key (default off) makes a spec Run push
its branch after a Clean outcome — and only then — through the existing push
machinery; pull requests are never created. Verifiable through the
outcome-matrix CLI tests over the fake pusher.

## Requirements

1. MUST add `implement.auto_push` (bool, default false) to Project Config
   with validation mirroring the existing push-key conventions; the key is
   honored from Project Config (per-repository decision), with User Config
   allowed to set a default the project can override.
2. MUST push only on a Clean outcome — including the QA-pass case — using
   the existing pusher against the branch's detected upstream; Unresolved,
   Failed, and Stopped Runs never invoke the pusher; a missing upstream is
   one stderr note and never a failure.
3. MUST journal the push through the existing push event kind and append
   one stdout line `pushed <remote>/<branch>` only when a push happened;
   push failure marks the Run Failed exactly as the review path's Final
   Push failure semantics do.
4. MUST leave review-path push behavior byte-identical, and spec Runs with
   the key off byte-identical to today.

## Subtasks

- [x] Config key, defaults, validation
- [x] Clean-only push wiring after the implement outcome
- [x] Journal event and stdout line
- [x] Outcome-matrix tests over the fake pusher

## Acceptance Criteria

- [x] Matrix tests: Clean+key pushes (invocation captured, stdout line
      present); Clean without key, Unresolved, Stopped, QA-fail all record
      zero pusher calls; missing upstream notes and exits per the Clean
      path.
- [x] Push failure ends the Run Failed with exit 1, journaled.
- [x] Config validation rejects non-boolean values with a named-key error.
- [x] Full suite passes; review-path push tests unchanged.

## Verification

- `rtk go test ./internal/config/ ./internal/cli/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 5. `_techspec.md` → Push at Clean,
Build Order 5. ADR-0013 (superseded half), ADR-0021. Round-1 dogfood
finding 25.

## Result

- Added `implement.auto_push` to config with default `false`, User Config
  precedence and Project Config override coverage, generated config output,
  and a named-key parse error for non-boolean values.
- Implemented Clean-only spec Run push through the existing `Pusher` and
  `daemon.push` event kind. `TestRunImplementAutoPushOutcomeMatrix` covers
  Clean+key with QA pass pushing `origin/ma/widget-flow`, Clean without key,
  Unresolved, Stopped, and QA-fail with zero pusher calls.
- Added missing-upstream coverage proving one stderr note, Clean exit
  semantics, no pusher call, and a journaled skipped push decision.
- Added push-failure coverage proving exit 1, Run state `Failed`, lock
  release, one pusher attempt, and a journaled failed push decision.
- Verification evidence: `rtk go test ./internal/config/ ./internal/cli/`
  passed with 244 tests; `rtk go test ./...` passed with 656 tests.
