---
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
prd: _prd.md
created: 2026-10-22
---

# The documentation contracts cost seconds per active Spec — Technical Spec

## Executive Summary

Each corpus test of `internal/docscontract` calls `speccheck.Check` once per
active Spec, and `TestStageScopeDefaultSweepIsUnchanged` calls it twice, so
the package runs five Checks per active Spec. Each Check of this repository
costs about 9 s of wall time and 15 s of CPU, and almost all of it is
`spec.HistoryCitationsInSource`: for every path-like token of every Citation
Source it compares the token with every Removed History Path, and each
comparison builds the string `target + "/"`. This Spec makes that lookup a set
lookup of the token and its parent directories, which keeps the reported
History Citations byte for byte, and makes the four corpus tests read one
memoized sweep per test process. Measured on 2026-10-22 the lookup alone cut
one Check from 15.7 s to 1.3 s of user CPU; the TechSpec's Measurements table
holds the package numbers.

The trade-off accepted: the corpus tests share a sweep computed once per test
process, so a `-count=N` run repeats the comparisons and not the Checks, and
an equality between two separate sweeps is no longer exercised by the
package; `TestStageScopeDefaultSweepIsUnchanged` keeps the one equality that
was ever asserted, `Check` against `CheckStage(StageAll)`, with a fresh
`CheckStage` call.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; tests use in-memory content, temporary
  fixture corpora and this repository's documents. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — the Spec records no decision and
  changes no behavior; no accepted ADR constrains how a Check computes its
  findings or how a test package shares its own results, as
  `_prd.md` → Project Constraints states. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the Makefile and the CI workflows stay unchanged. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0305-the-documentation-contracts-cost-seconds-per-active-spec/_authorization.md`;
  bounded files: none (`paths: []`).

## Current behavior

Read and characterized on 2026-10-22 at `42f523c2` (v0.109.0) in scratch
clones, with Specs 0299 and 0300 active. Line numbers move when other Specs
merge; the function names do not.

1. `internal/docscontract/corpus_test.go`: `TestCheckCorpusGolden` and
   `TestCheckCorpusBudget` each call `sweepCorpus`, which calls
   `checkCorpusSpec` (one `speccheck.Check`) per Spec directory with a
   `_prd.md`. `TestCheckActiveCorpusHasNoErrors` calls `speccheck.Check` per
   `spec.ListActive` entry. `assertDefaultStageSweepUnchanged`, under
   `TestStageScopeDefaultSweepIsUnchanged`, calls `speccheck.Check` and
   `speccheck.CheckStage(..., speccheck.StageAll)` per Spec and compares them
   with `reflect.DeepEqual`; `CheckStage` returns `Check` for `StageAll`
   (`internal/speccheck/coherence.go`). Five Checks per active Spec.
   `TestCheckLoopOrderRepositoryAgrees` adds one Check of the repository
   against a one-Spec carrier corpus, whatever the number of active Specs.
2. `go test -count=1 -timeout 20m -tags docscontract -v ./internal/docscontract`
   took 104.9 s with two active Specs: `TestCheckCorpusGolden` 20.8 s,
   `TestCheckCorpusBudget` 18.7 s, `TestCheckActiveCorpusHasNoErrors` 18.6 s,
   `TestStageScopeDefaultSweepIsUnchanged` 37.2 s,
   `TestCheckLoopOrderRepositoryAgrees` 9.0 s, everything else under 0.2 s.
   So a Check costs about 9 s whether its Spec is real or a carrier: the cost
   is repository-wide, not per Spec.
3. A CPU profile of that run: of 16.8 s of samples under `speccheck.Check`,
   16.2 s were in `detectHistoryRetained`
   (`internal/speccheck/history_retained.go`), 11.1 s of them in
   `spec.HistoryCitationsInSource` (`internal/spec/citation_sources.go`),
   dominated by `runtime.concatstring2` under its `Locate` closure; garbage
   collection took most of the remaining process CPU (176 % CPU for one
   sequential Check).
4. A logging `git` wrapper on `PATH` counted 31 Git processes in one Check,
   0.16 s of wall time in all. Process spawns are not the cost.
5. `HistoryCitationsInSource` runs `CitationRewrite.ApplyText` once per line
   with a `Locate` closure that, for each path-like token, loops over every
   target and tests `candidate == target || strings.HasPrefix(candidate,
   target+"/")`. `detectHistoryRetained` passes the 160 Removed History Paths
   of this repository and calls it once per Citation Source, every tracked or
   untracked unignored text file outside `docs/history/`. `SourceHistoryCitations`,
   the `baseline update` history plan's and `cancel`'s path, calls it the same
   way.
6. Prototype A, the set lookup of API Contract 1 in a scratch clone: one Check
   (`TestCheckLoopOrderRepositoryAgrees` run from a compiled test binary)
   went from 15.7 s to 22.5 s of user CPU in three runs under machine loads
   of 2.5 to 50 to 1.3 s to 1.5 s, with identical results. Also replacing
   the `strings.TrimLeft` copy of the rest of the content that
   `CitationRewrite.paths` makes at each token by a byte-slice form changed
   nothing measurable (1.27 s against 1.42 s), so it is left alone.
7. Prototype B, the shared sweep of API Contract 2: the package passed with
   the same corpus golden counts; the Measurements table below holds its
   timings.

### Measurements

Run with `go test -count=1 -timeout 20m -tags docscontract -v
./internal/docscontract` on the maintainer's 10-core Mac, one run at a time,
on 2026-10-22; the machine also ran Delivery Queues, so the load average
column matters and user CPU is the steadier figure. "0" removes Specs 0299
and 0300 with `git rm`; "5" adds the Spec folders of 0296, 0297 and 0298
restored from their Archive Records' `source_revision`.

| active Specs | code | wall | user CPU | load (1 min) |
| --- | --- | --- | --- | --- |
| 0 | today | 10.0 s | 17.5 s | 4.5 |
| 2 | today | 104.9 s | not timed (180.9 s of profile samples) | about 2.5 |
| 2 | today | 254.7 s | 260.0 s | 9.2, rising to about 20 |
| 5 | today | 316.5 s | 504.4 s | 24.0 |
| 5 | today | 631.1 s | 635.1 s | 4.6 at the end, about 24 during the run |
| 5 | set lookup only | 39.1 s | 37.5 s | 6.8 |
| 0 | set lookup and shared sweep | 2.3 s | 2.2 s | 5.2 |
| 2 | set lookup and shared sweep | 7.7 s | 7.3 s | 5.2 |
| 5 | set lookup and shared sweep | 17.5 s | 16.9 s | 5.0 |

With five active Specs, the four corpus tests took 84.1 s, 64.2 s, 55.1 s
and 102.7 s today (the 316.5 s run); with the set lookup only, 7.2 s, 8.1 s,
7.9 s and 14.0 s; with both changes, 7.8 s, 0.0 s, 0.0 s and 7.7 s. The
wall time of each row includes the package's compilation and its other tests
(about 1 s plus `TestCheckLoopOrderRepositoryAgrees`).

## System Architecture

- `internal/spec/citation_sources.go`: `HistoryCitationsInSource` gains a set
  of its targets, built once per call, and its `Locate` closure looks up the
  candidate and each parent directory in it.
- `internal/docscontract/corpus_test.go`: one memoized corpus sweep per
  `(specsRoot, repoRoot)` per test process, read by the four corpus tests; a
  new contract test of the memo.
- `internal/spec/history_citation_cost_test.go` (new): the allocation proof
  and the descendant table.

No package, import edge, exported symbol, schema or document changes.

## Implementation Design

### Interfaces

```go
// internal/docscontract/corpus_test.go (test code only)
type corpusSweep struct {
	slugs   []string                    // Spec directories in ReadDir order
	results map[string]speccheck.Result // by slug
	counts  map[string]int              // by finding code, all codes present
	work    corpusSweepWork
	elapsed time.Duration
	err     error
}

// sharedCorpusSweep returns the one sweep of specsRoot for this test process.
func sharedCorpusSweep(t *testing.T, specsRoot, repoRoot string) *corpusSweep
```

### Data Models

None. No persisted format changes.

### API Contracts

1. API Contract: `spec.HistoryCitationsInSource(content, path, targets)`
   keeps its signature, its error, and its result: the same History
   Citations, deduplicated and sorted by line then target, for every input.
   A candidate token matches a target when it equals the target or when one
   of its parent directories, found by cutting at each `/` from the right,
   equals the target; every matching target is reported, so nested targets
   each yield a citation. The set is built once per call from `targets`. The
   Markdown pass (`archiveContentCitations`) is unchanged.
2. API Contract: `sharedCorpusSweep` runs `checkCorpusSpec` once per Spec
   directory carrying a `_prd.md`, in `os.ReadDir` order, skipping names that
   start with `_`, exactly as `sweepCorpus` does today, the first time a
   `(specsRoot, repoRoot)` pair is requested in the test process; every later
   request returns the same `*corpusSweep`. A read or Check error is stored
   and reported by every test that requests the sweep, with today's
   messages. An uncharacterized finding code is reported the way `sweepCorpus`
   reports it.
3. API Contract: `TestCheckCorpusGolden` compares the shared sweep's counts
   with `testdata/corpus-golden.json`; `TestCheckCorpusBudget` logs and
   bounds the shared sweep's `work` and elapsed time;
   `TestCheckActiveCorpusHasNoErrors` reads each `spec.ListActive` slug's
   result from the shared sweep and calls `speccheck.Check` only for a slug
   the sweep did not check; `assertDefaultStageSweepUnchanged` takes each
   Spec's `Check` result from the shared sweep and compares it with a fresh
   `speccheck.CheckStage(..., speccheck.StageAll)`. Every failure message and
   the corpus golden keep their bytes.

### Surface Transcripts

None. No command surface changes.

### Invariants

1. For every content, path and target list, `HistoryCitationsInSource`
   returns the same value before and after this Spec.
2. With a fixed content, its allocations per call do not grow with the
   number of targets beyond a constant.
3. In one test process, the four corpus tests together call
   `speccheck.Check` at most once per active Spec, and
   `speccheck.CheckStage` once per active Spec.
4. A corpus golden mismatch, an error-severity finding in an active Spec,
   and a difference between `Check` and `CheckStage(StageAll)` still fail
   the test that asserts them today.

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1 and 2 → API Contract 1,
  Invariants 1 and 2, task_01.
- Goal 2 → Invariant 2, `TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths`, task_01.
- Goal 3, Core Feature 2, Success Metrics 3 and 4 → API Contracts 2 and 3,
  Invariant 3, task_02.
- Goal 4, Success Metric 5 → Invariants 1 and 4, the existing `internal/spec`,
  `internal/speccheck` and `internal/docscontract` tests, task_01, task_02,
  task_03.
- Success Metric 6 → Unreachable Acceptance, operator after merge.

Deviations from the PRD: None.

## Absence Evidence

1. The set lookup of Removed History Paths in `HistoryCitationsInSource` —
   command: `grep -n "targetSet\|for _, target := range targets" internal/spec/citation_sources.go`;
   result: only the loop `for _, target := range targets` inside the
   `Locate` closure; no set of targets exists.
2. The shared corpus sweep — command:
   `grep -n "sync.Once\|sharedCorpusSweep\|sync\"" internal/docscontract/corpus_test.go`;
   result: no match; every corpus test calls `sweepCorpus` or
   `speccheck.Check` itself.
3. An allocation-based cost test in the repository — command:
   `grep -rln "AllocsPerRun" internal cmd`; result: no match.

## ADR Alternatives

None. The Active ADR obligations row lists no ADR: the Spec records no
decision and changes no behavior.

## Integration Points

None. No external system, network or provider is reached.

## Testing Approach

- `internal/spec/history_citation_cost_test.go` (new):
  - `TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths`, a
    Sequential Test (`// Sequential: testing.AllocsPerRun sets
    runtime.GOMAXPROCS to 1 and counts process-wide allocations.`), because
    `AllocsPerRun` panics during a parallel test: 2,000 Removed History Paths
    of more than 32 bytes, a one-line source `Makefile` with 200 path-like
    tokens and one citation of the last removed path; it asserts one citation
    per call and at most 64 more allocations per call than with only that
    path. Before the fix it measured 1,603,622 against 2,423; after, 1,633
    against 1,625. `internal/spec` then holds two Sequential Tests, its
    ceiling.
  - `TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants`,
    parallel, a table over a non-Markdown source with the cases exact path,
    descendant of a removed directory, nested removed paths (both reported),
    relative prose path from `docs/agents/`, and string prefix only (none
    reported). It passes before and after; exact-only and first-match-only
    sabotages fail it (checked on 2026-10-22).
- `internal/docscontract/corpus_test.go` (existing `//verify:always` file;
  no new file in the package): `TestTheCorpusTestsCheckEachSpecOncePerProcess`
  builds a fixture corpus of two carrier Specs in a temporary directory with
  `writeLoopOrderCarrier`'s shape, requests `sharedCorpusSweep` three times,
  and asserts one pointer, two Specs and two Check operations. A recompute
  per request fails on the pointer, a second Check per Spec on the count.
- The existing tests prove Invariants 1 and 4: `internal/spec`
  `TestHistoryRewriteReachesEveryCitationSource`,
  `TestHistoryRewriteResolvesARelativeProsePath`,
  `TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown`,
  `TestHistoryCitationsNameMarkdownThatCitesARemovedPath`;
  `internal/speccheck` `TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator`
  and the other `TestHistoryRetained` tests; the corpus golden.

## Build Order

1. The set lookup in `HistoryCitationsInSource` with its two tests (task_01).
2. The shared corpus sweep and its contract test (task_02) (depends on: 1,
   so that its measurement and the corpus golden run on the cheaper Check).
3. The QA gate (task_03) (depends on: 1, 2).

No other Spec is required.

## Risks & Considerations

- The `internal/spec` Sequential Test ceiling is 2 and this Spec reaches it.
  A plan that adds another Sequential Test there must lower the count or
  raise the ceiling in its own change.
- A memo shared by tests can hide a test that mutates the repository between
  them. No corpus test writes to the repository; `-record-skill-coverage`
  writes `docs/references/behavior-surfaces.json` only when that flag is set
  and `-run` selects its own test.
- The Markdown pass still scales with lines × Removed History Paths; it is
  about a fifth of a Check after this Spec and is named under Non-Goals.
- The CI job still runs the package twice; the follow-up Backlog Entry keeps
  it.

## Glossary

None.

## Decisions

- Fix the cost where the profile puts it, the per-token loop over every
  Removed History Path, and share the corpus sweep across tests; no Makefile
  or workflow change.
- Prove the cost by an allocation count, which machine load does not move,
  and not by elapsed time.
- No ADR: the change alters no contract a later reader must know.
