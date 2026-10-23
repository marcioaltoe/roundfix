---
task: task_01
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
status: pending
type: backend
complexity: medium
---

# Task 01: The removed-history citation pass looks each token up in a set

## Overview

`spec.HistoryCitationsInSource` compares every path-like token of a Citation
Source with every Removed History Path and builds `target + "/"` for each
comparison, which makes one `speccheck.Check` of this repository cost about
15 s of CPU. This Task replaces the loop with a set lookup of the token and
its parent directories, keeps every reported History Citation, and proves
with an allocation count that the pass no longer grows with the number of
Removed History Paths.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-the-documentation-contract-suite-costs-about-a-minute-per-active-spec-and-runs-twice-in-one-ci-job.md`](references/2026-10-22-the-documentation-contract-suite-costs-about-a-minute-per-active-spec-and-runs-twice-in-one-ci-job.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`.
2. MUST change only the `Locate` closure of `HistoryCitationsInSource` in
   `internal/spec/citation_sources.go` and what it needs: build a set of
   `targets` once per call, before the line loop, and for each candidate,
   after `citationProseTarget`, look up the candidate and then each parent
   directory obtained by cutting at the last `/`, recording a citation for
   every target found. The closure still returns
   `GitLocator{Path: candidate}` and whether any target matched. The Markdown
   pass `archiveContentCitations`, `CitationRewrite` and every other function
   keep their bytes.
3. MUST keep the result byte for byte (Invariant 1): a candidate equal to a
   target, a candidate below a target, nested targets that both contain the
   candidate, an explicit `./` or `../` prose path, and a candidate that only
   shares a string prefix with a target (`docs/history/adr/0001-old.md.bak`,
   `docs/history/backlog-old/x.md`) give the same citations as today.
4. MUST add the new file `internal/spec/history_citation_cost_test.go` with:
   - `TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths`, a
     Sequential Test whose doc comment is `// Sequential: testing.AllocsPerRun
     sets runtime.GOMAXPROCS to 1 and counts process-wide allocations.`
     (`AllocsPerRun` panics during a parallel test). It builds 2,000 Removed
     History Paths of the form
     `docs/history/backlog/2026-01-01-a-removed-backlog-entry-NNNN.md` (each
     longer than 32 bytes, so the old concatenation allocates), a one-line
     content for the path `Makefile` with 200 tokens
     `internal/spec/source-file-NNN.go` followed by the last removed path,
     measures `testing.AllocsPerRun(5, ...)` of `HistoryCitationsInSource`
     with only the last removed path and with all 2,000, requires exactly one
     citation per call, logs both counts, and fails when the second exceeds
     the first by more than 64.
   - `TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants`,
     calling `t.Parallel()` first, a table over the targets
     `docs/history/adr/0001-old.md`, `docs/history/specs/0001-old`,
     `docs/history/specs/0001-old/references/source.md` and
     `docs/history/backlog`, each case a two-line content whose second line
     holds the citation, asserting path, line 2 and the exact ordered targets:
     exact path; descendant (`docs/history/specs/0001-old/_prd.md` gives
     `docs/history/specs/0001-old`); nested removed paths (the reference file
     gives both Spec-folder and file targets); relative prose path
     (`../history/adr/0001-old.md` from `docs/agents/guide.txt`); string
     prefix only (no citation).
   Both compile for Windows and use no Git, temporary home or real time.
5. MUST keep the `internal/spec` Sequential Test count within its ceiling of
   2 in `internal/testfixture/parallel_tests_test.go` without changing that
   file.
6. MUST NOT change `internal/speccheck`, any finding text, the corpus golden,
   any existing test, golden or fixture, the Makefile or the CI workflows.

## Subtasks

- [ ] Replace the per-target loop with the set lookup of the candidate and its parent directories.
- [ ] Add the allocation proof and the descendant table.
- [ ] Run the existing History Citation and `SC-HISTORY-RETAINED` tests unchanged.

## Acceptance Criteria

- [ ] With 2,000 Removed History Paths the pass allocates at most 64 more per call than with one.
- [ ] Every descendant-table case and every existing History Citation, history rewrite and `SC-HISTORY-RETAINED` test passes unchanged.
- [ ] `internal/spec` still holds at most two Sequential Tests and compiles for Windows.

## Consumers

- `spec.HistoryCitationsInSource` results — consumers: `./internal/speccheck`, `./internal/cli`

## Context

- instruction: `internal/spec/git_locator.go`
- instruction: `internal/spec/live_citations.go`
- instruction: `internal/spec/citation_sources_test.go`
- instruction: `internal/spec/history_entries_test.go`
- instruction: `internal/speccheck/history_retained.go`
- instruction: `internal/speccheck/history_retained_test.go`
- instruction: `internal/cli/baseline_update_citations_test.go`
- instruction: `internal/testfixture/parallel_tests_test.go`
- creates: `internal/spec/history_citation_cost_test.go`
- interface: `internal/spec/citation_sources.go`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths|TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants|TestHistoryRewriteReachesEveryCitationSource|TestHistoryRewriteResolvesARelativeProsePath|TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown|TestHistoryCitationsNameMarkdownThatCitesARemovedPath|TestHistoryRetained.*|TestEveryTestInAParallelTestPackageRunsInParallel)$' ./internal/spec ./internal/speccheck ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants TestHistoryRewriteReachesEveryCitationSource TestHistoryCitationsNameMarkdownThatCitesARemovedPath TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator TestEveryTestInAParallelTestPackageRunsInParallel; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 15m -run '^TestBaselineUpdate.*Citation.*$' ./internal/cli && GOOS=windows go vet ./internal/spec` — expected: exit 0. Before this Task the two new tests do not exist, so the command fails at its first pass check; with the per-target loop restored the allocation test fails (1,603,622 against 2,423 allocations measured on 2026-10-22). Afterwards both new tests, the existing History Citation, history rewrite, `SC-HISTORY-RETAINED`, Parallel Test Package and `baseline update` citation tests pass, and `internal/spec`'s tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Success Metric 1; Success Metric 2; Success Metric 5
- `_techspec.md` → Current behavior; API Contract 1; Invariants 1-2; Testing Approach; Build Order 1
