---
task: task_10
spec: 0064-spec-artifact-consistency-gate
status: completed
type: test
complexity: medium
---

# Task 10: Prove the sweep budget under conditions that mean something

## Overview

QA finding F-001: `TestCheckCorpusGoldenAndBudget` asserts the corpus sweep
finishes under one second, but it measures wall clock inside
`go test -parallel 16 ./...`. The same sweep measured 0.64 s alone and
3.0–4.0 s under the suite, so `make verify` exits 2 for a product that meets
its budget. The assertion measures the runner's load, not the check.

This slice separates the two things that test conflates. The golden per-code
counts are deterministic and belong in the ordinary suite. The duration proof
needs to own the machine while it runs, so it becomes its own test that takes
its execution conditions explicitly rather than inheriting whatever the suite
happens to be doing.

The one-second limit does not move. Raising it would hide the conflation
instead of removing it, and the QA report says so: do not weaken the acceptance
limit or hide the failing exit.

## Requirements

1. MUST split the existing test into a deterministic golden-count test that
   runs in the ordinary suite and a separate sweep-budget test.
2. MUST keep the sub-second limit exactly as the Spec states it. Widening the
   limit, removing the assertion, or downgrading it to a log line is out of
   scope and fails this Task.
3. MUST make the budget test take its execution conditions explicitly, per
   ADR-0089, so it asserts only where the measurement means what it claims and
   is inert inside the contended package sweep.
4. MUST make the budget test able to fail: under its declared conditions a
   sweep exceeding the limit exits non-zero. A test that passes because it
   selected no work or skipped silently does not satisfy this.
5. MUST leave the golden counts and every detector behavior unchanged — this
   Task changes when and how the budget is measured, never what the check
   finds.
6. MUST NOT change `Makefile`; wiring the serial gate step is task_11's
   bounded tooling scope.
7. SHOULD record, in the test or beside it, why the budget cannot be proven
   inside the parallel suite, so the next reader does not re-merge them.

## Subtasks

- [x] Split the golden-count assertion from the duration assertion.
- [x] Give the budget test its explicit execution conditions.
- [x] Prove the budget test fails when the limit is exceeded.
- [x] Confirm the ordinary suite no longer carries a wall-clock assertion.

## Acceptance Criteria

- [ ] The golden-count test passes inside `go test -parallel 16 ./...` and its
      counts are byte-identical to the checked-in golden.
- [ ] No wall-clock assertion remains in any test that runs under the ordinary
      parallel package sweep, proven by a repository-wide check rather than by
      inspecting one file.
- [ ] Under its declared conditions the budget test measures the sweep and
      passes at the unchanged sub-second limit.
- [ ] The budget test fails when the limit is exceeded, proven by a temporary
      limit reduction reverted within the same check.
- [ ] `go test -parallel 16 ./internal/speccheck` passes with no duration
      failure.
- [ ] `Makefile` is unchanged by this Task.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/speccheck/constraints_characterization_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -parallel 16 2>&1 | grep -q "want under 1s" && exit 1 || exit 0`
  — expected: exit 0; the contended sweep no longer carries the wall-clock
  assertion that F-001 reported.
- `go test ./internal/speccheck -count=1 -run 'Golden' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the golden-count test ran and passed.
- `go test ./internal/speccheck -count=1 -run 'Budget' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the budget test ran and passed under its conditions.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0; the full parallel sweep is green.
- `if git diff --name-only HEAD~1 2>/dev/null | grep -q "^Makefile$"; then exit 1; fi`
  — expected: exit 0; this Task did not touch the bounded tooling file.

## References

- `_prd.md` → Core Feature 1; Success Metric 2.
- `_techspec.md` → Testing Approach (budget execution conditions).
- `qa/qa-report-2026-08-03.md` → F-001; rows Q-11 and Q-15.
- ADR-0089.

## Result

### Implementation

- Split `TestCheckCorpusGoldenAndBudget` into an ordinary
  `TestCheckCorpusGolden` count comparison and an opt-in
  `TestCheckCorpusBudget` duration proof.
- Kept the limit at `time.Second`. The budget test accepts only the dedicated
  `Budget` and `^TestCheckCorpusBudget$` selectors; every broader package run
  skips before materializing or timing the corpus and explains the required
  invocation.
- Kept the detector calls, corpus roots, archived-corpus materialization, and
  golden comparison unchanged. Added a focused timing log while retaining the
  non-zero assertion on an exceeded limit.

### Focused checks

- `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false rtk proxy go test ./internal/speccheck -count=1 -run '^TestCheckCorpusBudget$' -v`
  — exit 0 after the final edit; measured the active and archived sweep at
  `309.29275ms` against the unchanged `1s` budget.
- Temporary mutation check: changed only `corpusBudget` from `time.Second` to
  `time.Nanosecond`, ran the same dedicated command, and observed exit 1 with
  `full Spec corpus sweep took 321.8035ms, want under 1ns`; restored
  `time.Second` immediately, then reran the preceding command successfully.
- `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false rtk proxy go test ./internal/speccheck -count=1 -parallel=16 -run '^(TestCheckCorpusGolden|TestCheckCorpusBudget|TestCheckActiveCorpusHasNoErrors)$' -v`
  — exit 0; the golden and active-corpus tests passed, while the budget test
  explicitly skipped before timing because this was not a dedicated selector.
- `rtk zsh -c 'if rtk rg -n "TestCheckCorpusGoldenAndBudget|want under 1s" --glob "*_test.go" .; then exit 1; fi'`
  — exit 0; the repository-wide test-file check found neither the conflated
  test nor its ordinary-suite wall-clock diagnostic.
- `rtk git -c core.fsmonitor=false diff --quiet HEAD -- Makefile` — exit 0;
  `Makefile` is unchanged.

### Acceptance-criterion evidence

1. The golden-count test passed in the focused `-parallel=16` package subset
   with the checked-in golden unchanged. The Daemon-owned full parallel sweep
   remains the authoritative criterion check.
2. The repository-wide absence check exited 0, and the focused ordinary-run
   output proves the remaining budget test skips before any wall-clock work.
3. The dedicated budget test swept both corpus roots and passed at
   `309.29275ms < 1s` after the final edit.
4. The temporary `1ns` limit made that same dedicated sweep exit 1; restoring
   the unchanged `1s` limit made it pass again.
5. The relevant `internal/speccheck` subset passed with `-parallel=16`. The
   complete package command is declared Verification and was not run here.
6. The focused Git check proves `Makefile` has no diff.

### Follow-up

- task_11 owns wiring the dedicated `^TestCheckCorpusBudget$` invocation into
  the serial gate step; this Task intentionally did not edit `Makefile`.
