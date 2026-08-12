---
task: task_02
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: completed
type: backend
complexity: low
---

# Task 02: State the Pull Request fact and a collision-safe report name

## Overview

The QA Agent is asked to observe Pull Request journeys without being told
whether a Pull Request exists, so it cannot tell "no Pull Request is open" from
"I failed to find it" — and the first is an environment-blocked row, not a
finding. Pass the resolved Pull Request as a labeled fact, and narrow the report
filename instruction to the one scheme the recency comparator ranks correctly.

## Requirements

1. MUST add a `PullRequest` field to the QA prompt request, rendered as a
   labeled fact when set.
2. MUST state explicitly, when no Pull Request is open, that none exists and
   that Pull Request journeys are therefore environment-blocked.
3. MUST resolve the Pull Request for the target branch best-effort through the
   existing read-only `gh` boundary before building the prompt; any failure
   leaves the field empty and MUST NOT fail the QA plan.
4. MUST instruct the report filename as `qa-report-YYYY-MM-DD.md` for a day's
   first report and numeric `-NN` suffixes for same-day reruns, removing the
   non-numeric `-<scope-or-build>` alternative.
5. MUST update the substring-pinned prompt contract assertions in the same
   change as the wording.
6. MUST leave the recency comparator untouched.

## Subtasks

- [ ] Add and render the Pull Request fact in the QA prompt contract.
- [ ] Resolve it best-effort in the Daemon's QA plan.
- [ ] Narrow the filename instruction and move the pinned assertions with it.

## Acceptance Criteria

- [ ] A QA prompt built with an open Pull Request contains it as a labeled fact.
- [ ] A QA prompt built with none states that Pull Request journeys are
      environment-blocked.
- [ ] A `gh` failure or timeout during resolution still produces a complete QA
      prompt and does not fail the plan.
- [ ] The contract instructs only numeric same-day suffixes.

## Context

- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/spec_prompt_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/agent/ ./internal/daemon/` — expected: pass,
  including the best-effort resolution failure case.

## References

`_prd.md` → Goal 1 Story 1, Goal 4 Story 6; `_techspec.md` → Build Order 2,
Interfaces, Risks (prompt contract tests).

## Result

Added the resolved Open Pull Request to the QA prompt as a labeled fact. The
Daemon now queries the existing read-only `gh` boundary for the recorded
target branch before prompt construction, with an operation-specific timeout;
lookup, timeout, and response-parse failures leave the fact empty without
failing the QA plan. An empty fact states that no Pull Request is open and
that Pull Request journeys are environment-blocked.

The QA report contract now names `qa-report-YYYY-MM-DD.md` for the day's first
report and only `qa-report-YYYY-MM-DD-NN.md` numeric suffixes for same-day
reruns. The report recency comparator was not changed.

Focused implementation evidence:

- Before implementation, the new prompt tests failed to compile because
  `QAPromptRequest.PullRequest` did not exist, and the new Daemon tests failed
  to compile because the engine had no injectable `GH` boundary.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test -count=1 -run '^TestBuildQAPrompt' ./internal/agent`
  passed 14 tests after the final implementation edit.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test -count=1 -run '^TestTaskCycleQA' ./internal/daemon`
  passed 20 tests after the final implementation edit.
- `rtk git diff --check` passed after the final implementation edit.
- The Task's two declared Daemon Verification commands were not run.

Acceptance evidence:

- Open Pull Request fact:
  `TestBuildQAPromptStatesCheckoutFactsSeparatingRunBranchFromTarget` and
  `TestTaskCycleQAPromptStatesRunBranchAndSpecTargetBranch` passed with
  `Pull Request: #40 (owner/repo)` in the prompt. The Daemon case also asserts
  the read-only lookup uses `pr list --head ma/spec-work --state open`.
- No open Pull Request:
  `TestBuildQAPromptStatesPullRequestJourneysAreEnvironmentBlockedWhenNoneIsOpen`
  passed with the explicit `none open` and `environment-blocked` fact.
  `TestTaskCycleQAPromptStaysUsableWithoutRecordedTargetBranch` passed with
  the same fact and no attempted lookup.
- Best-effort failure behavior:
  `TestTaskCycleQAPromptSurvivesPullRequestResolutionFailure` passed for both
  a `gh` error and `context.DeadlineExceeded`; each case still ran the QA
  Agent with a complete prompt and settled its report verdict.
- Numeric suffix contract:
  `TestBuildQAPromptStatesQAGateContract` passed with the first-report name,
  the full numeric `-NN` rerun name, and an assertion excluding the removed
  `-<scope-or-build>` alternative.
