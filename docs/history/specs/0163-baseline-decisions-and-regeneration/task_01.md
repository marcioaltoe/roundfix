---
task: task_01
spec: 0163-baseline-decisions-and-regeneration
status: completed
type: backend
complexity: medium
---

# Task 01: A mode change keeps the HTTP decision

## Overview

The `http-contract` branch of `promptBaselineDecision` returns `{"mode": ...}` on "Change", which silently drops every recorded exception and the source of the stored `http.contract` decision.

## Requirements

1. MUST, when a valid current `http.contract` value exists and the maintainer chooses "Change", return a clone of that value with only `mode` replaced, validated through `baseline.ValidateDecisionValue`.
2. MUST name each kept exception scope in the review line of that prompt.
3. MUST keep the prompt as today when no current value exists.
4. MUST keep a value supplied through `--decision` or `--decision-file` replacing the whole typed decision, so a deliberate exception edit stays distinguishable from retention.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The multi-exception decision of the archived Finding, with a source, keeps every exception and the source, compared field by field, after a mode-only change.
- [ ] The review names each kept exception scope.
- [ ] An explicit decision value with one exception fewer yields exactly that value.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/baseline_human.go`
- interface: `internal/baseline/project_decisions.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestHTTPContractModeChangeRetainsExceptionsAndSource|TestHTTPContractModeChangeReviewNamesKeptExceptions|TestHTTPContractExplicitValueStillReplacesExceptions)$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestHTTPContractModeChangeRetainsExceptionsAndSource TestHTTPContractModeChangeReviewNamesKeptExceptions TestHTTPContractExplicitValueStillReplacesExceptions; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The HTTP mode change
- [2026-08-07-changing-the-http-contract-discards-its-exceptions.md](../../history/findings/2026-08-07-changing-the-http-contract-discards-its-exceptions.md)

## Result

### Implementation

- A valid stored HTTP Contract now changes through a deep clone: the selected
  mode replaces `mode`, every exception and `source` stays intact, the stored
  value is not mutated, and `baseline.ValidateDecisionValue` validates the
  result.
- The change option names every retained exception scope. With no valid stored
  value, the existing mode-only prompt remains unchanged.
- Explicit `--decision` input still enters planning as the exact supplied typed
  decision, including a deliberately shorter exception list.

### Focused checks

- Before the implementation,
  `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run '^TestHTTPContractModeChangeRetainsExceptionsAndSource$' ./internal/cli`
  failed because the result contained only `mode`.
- Before the implementation,
  `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run '^TestHTTPContractModeChangeReviewNamesKeptExceptions$' ./internal/cli`
  failed because the change line named none of the four scopes.
- After the implementation, each of the three named Task tests passed when run
  individually. `TestHumanBaselineDecisionDefaults`,
  `TestProjectDecisionPrompts`, and `TestProjectDecisionReuse` also passed as
  focused regressions.
- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/cli`
  passed in 90.251 seconds after the existing GitHub-backed tests received the
  network access they require. `git diff --check` also passed.

### Acceptance evidence

1. `TestHTTPContractModeChangeRetainsExceptionsAndSource` compares the archived
   Finding fixture field by field after a mode-only change and separately proves
   that the stored value was not mutated.
2. `TestHTTPContractModeChangeReviewNamesKeptExceptions` isolates the change
   option and proves it names `/api/auth/*`, `/health`, `/openapi.json`, and
   `/reference`.
3. `TestHTTPContractExplicitValueStillReplacesExceptions` sends an explicit
   `--decision` value with one fewer exception and compares the parsed decision
   exactly with that value.

The Daemon-owned `## Verification` command was not run.

## Carry-forward provenance

- Source Run: `run_20260925T135958Z_cb4c08f055bb883c`
- Source commit: `05c25931c2b1c54c9969b7d0900d95ef4edb5da2`
