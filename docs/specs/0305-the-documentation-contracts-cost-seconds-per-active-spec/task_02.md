---
task: task_02
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
status: pending
type: test
complexity: medium
---

# Task 02: The corpus tests share one Check per active Spec

## Overview

`TestCheckCorpusGolden`, `TestCheckCorpusBudget`,
`TestCheckActiveCorpusHasNoErrors` and `TestStageScopeDefaultSweepIsUnchanged`
each run `speccheck.Check` on every active Spec, and the last runs it twice,
so `internal/docscontract` pays five Checks per active Spec. This Task gives
them one memoized sweep per test process, keeps the one fresh
`CheckStage(StageAll)` call whose equality with `Check` the stage test
asserts, and proves the memo on a fixture corpus.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-the-documentation-contract-suite-costs-about-a-minute-per-active-spec-and-runs-twice-in-one-ci-job.md`](references/2026-10-22-the-documentation-contract-suite-costs-about-a-minute-per-active-spec-and-runs-twice-in-one-ci-job.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 2 and 3 of
   `_techspec.md` in `internal/docscontract/corpus_test.go`, an existing
   `//verify:always` file; no new file is added to the package.
2. MUST add `sharedCorpusSweep(t *testing.T, specsRoot, repoRoot string)
   *corpusSweep` with the shape `_techspec.md` → Interfaces gives: one sweep
   per `(specsRoot, repoRoot)` pair per test process, guarded for concurrent
   callers, computed the first time the pair is requested by calling
   `checkCorpusSpec` once per Spec directory that `sweepCorpus` selects
   today, in the same order, and returned unchanged to every later request.
   A read or Check error and an uncharacterized finding code are stored and
   reported with today's messages by every test that requests the sweep.
3. MUST make the four corpus tests read the shared sweep exactly as API
   Contract 3 states: the golden counts and the budget's `work` and elapsed
   time come from it; `TestCheckActiveCorpusHasNoErrors` keeps iterating
   `spec.ListActive` and calls `speccheck.Check` only for a slug the sweep did
   not check; `assertDefaultStageSweepUnchanged` keeps one subtest per Spec,
   takes the `Check` result from the sweep and compares it with a fresh
   `speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageAll)`.
   Every failure message, the corpus golden file and
   `maxCorpusCheckOperationsPerSpec` keep their bytes. Remove `sweepCorpus`
   only if nothing calls it any more.
4. MUST add `TestTheCorpusTestsCheckEachSpecOncePerProcess` to
   `internal/docscontract/corpus_test.go`, calling `t.Parallel()` first: in
   `t.TempDir()` it writes two carrier Specs under `docs/specs/` (the
   `loop-order` carrier `writeLoopOrderCarrier` writes and a second one with
   the same shape and another slug), requests `sharedCorpusSweep` for that
   root three times, and asserts that the three requests return the same
   pointer, that the sweep's `work` counts 2 Specs and 2 Check operations,
   and that its results hold both slugs. A sweep recomputed per request fails
   on the pointer; a second Check per Spec fails on the operation count.
5. MUST NOT change any production code, `testdata/corpus-golden.json`, the
   Contract Relevance of any test, the Makefile, the CI workflows or any
   other file of the package.

## Subtasks

- [ ] Add the memoized corpus sweep.
- [ ] Make the four corpus tests read it, keeping the fresh `CheckStage` comparison.
- [ ] Add the fixture test of the memo.

## Acceptance Criteria

- [ ] In one test process the four corpus tests run one `Check` and one `CheckStage` per active Spec.
- [ ] The corpus golden, the no-error rule and the stage equality still fail on the inputs they fail on today.
- [ ] Three requests for one fixture corpus return one sweep that ran exactly one Check per Spec.

## Consumers

None. The change is test code inside `internal/docscontract`; no production code or other package reads it.

## Context

- interface: `internal/docscontract/corpus_test.go`

## Verification

- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestTheCorpusTestsCheckEachSpecOncePerProcess|TestCheckCorpusGolden|TestCheckCorpusBudget|TestCheckActiveCorpusHasNoErrors|TestStageScopeDefaultSweepIsUnchanged|TestCheckLoopOrderRepositoryAgrees|TestCheckLoopOrderDivergent)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheCorpusTestsCheckEachSpecOncePerProcess TestCheckCorpusGolden TestCheckCorpusBudget TestCheckActiveCorpusHasNoErrors TestStageScopeDefaultSweepIsUnchanged TestCheckLoopOrderRepositoryAgrees TestCheckLoopOrderDivergent; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task `TestTheCorpusTestsCheckEachSpecOncePerProcess` does not exist, so the command fails at its first pass check. Afterwards the memo test and every corpus and loop-order test pass, and the package's tests compile for Windows.

## References

- `_prd.md` → Goals; Core Feature 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 2-3; Invariants 3-4; Testing Approach; Build Order 2
