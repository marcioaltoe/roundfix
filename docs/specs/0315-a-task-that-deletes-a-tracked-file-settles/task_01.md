---
task: task_01
spec: 0315-a-task-that-deletes-a-tracked-file-settles
status: completed
type: backend
complexity: medium
---

# Task 01: Tracked-file readers skip a path the working tree lacks

## Overview

Three readers list tracked files from the Git index and then read each one
from the working tree. When an Agent deletes a tracked file without staging
the deletion, which is what a Task with `deletes:` does because the Daemon
makes the commit, the file is still listed and the read fails the whole call.
This Task makes `CitationSources`, `HistoryCitations` and
`StaleHistoryCitations` skip a tracked path whose working-tree entry does not
exist, and keep failing on every other read error. The Daemon's settlement
reaches `CitationSources` through the strict Spec check; task_02 proves that
path in `internal/daemon`.

## Requirements

1. MUST answer the oraculum Inbox Entry of 2026-10-22,
   [`references/2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md`](references/2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-3 of
   `_techspec.md`.
2. MUST make `CitationSources` (`internal/spec/citation_sources.go`) continue
   past a listed path when `historyRewriteRegularFile` fails with an error
   for which `errors.Is(err, os.ErrNotExist)` holds, wherever the path is
   (API Contract 1), replacing today's `docs/history`-only condition.
3. MUST make `HistoryCitations` (`internal/spec/history_entries.go`) continue
   past a listed Markdown path when `historyRewriteRegularFile` or
   `os.ReadFile` fails with `os.ErrNotExist`, wherever the path is (API
   Contract 2), replacing today's `docs/history/`-only conditions.
4. MUST make `StaleHistoryCitations` (`internal/spec/stale_citations.go`)
   continue past a listed Markdown path when `os.ReadFile` fails with
   `os.ErrNotExist` (API Contract 3).
5. MUST keep every other error of the three readers, with today's message:
   a path behind a file (`ENOTDIR`), a symbolic link, a non-regular entry and
   a read failure. State this hazard in a comment at each skip: only a missing
   entry is skipped.
6. MUST add `internal/spec/citation_sources_working_tree_test.go` with a
   `Suite:`/`Invariant:`/`Boundary` header and the two tests of `_techspec.md`
   → Testing Approach, each top-level test and subtest calling `t.Parallel()`
   first, with subtests `citation-sources`, `history-citations` and
   `stale-history-citations`:
   - `TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks`: after removing
     the tracked `docs/gone.md` and `pkg/gone.test.ts` from disk only, each
     reader returns no error and exactly `README.md` and `notes/a.md`;
   - `TestTrackedSourceReadersStillRefuseAPathBehindAFile`: after replacing
     the tracked directory `notes` by a file, each reader returns an error
     naming `notes/a.md`.
7. MUST NOT change any other production file, existing test, golden,
   fixture, guide, skill or glossary entry, except Requirement 8.
8. MUST change how two existing tests of
   `internal/cli/baseline_update_citations_test.go` inject a detector
   failure, keeping every assertion: `TestBaselineUpdateHistoryCitationDetectorFailureIsLocal`
   and `TestBaselineUpdateHistoryStaleDetectorFailureKeepsOutcome` call
   `os.Remove` on `baselineCitationHandoff`, which this Task makes a skipped
   missing path instead of a failure. Replace both calls with one helper
   `breakBaselineCitationHandoff(t, repo)` that removes the handoff's folder
   (`filepath.Dir` of the handoff path) and writes a regular file in its
   place, so the detectors fail inspecting a tracked path behind a file,
   which Requirement 5 keeps an error. Making the file unreadable does not
   work: the tests' own working-tree snapshot reads it. (Amended 2026-10-25
   by the operator after Run run_20261025T231938Z_727a0c81b3bea09d failed
   both tests; prototyped: the `TestBaselineUpdateHistory` and `TestCancel`
   tests pass.)

## Subtasks

- [ ] Skip a missing working-tree entry in the three readers.
- [ ] Add the two reader tests with their three subtests.

## Acceptance Criteria

- [ ] A tracked file removed from disk only is neither listed nor read by the three readers, and the remaining citations are still found.
- [ ] A tracked path behind a file still fails each reader, naming the path.

## Consumers

- `spec.CitationSources` result for a tracked path missing on disk — consumers: `./internal/speccheck`
- `spec.HistoryCitations` and `spec.StaleHistoryCitations` results for a tracked Markdown file missing on disk — consumers: `./internal/cli`

## Context

- instruction: `internal/speccheck/history_retained.go`
- instruction: `internal/spec/citation_sources_test.go`
- instruction: `internal/spec/stale_citations_test.go`
- instruction: `internal/spec/live_citations.go`
- instruction: `internal/speccheck/history_retained_test.go`
- instruction: `internal/cli/cancel_test.go`
- interface: `internal/cli/baseline_update_citations_test.go`
- instruction: `internal/cli/cancel.go`
- instruction: `internal/cli/baseline_update_history.go`
- interface: `internal/spec/citation_sources.go`
- interface: `internal/spec/history_entries.go`
- interface: `internal/spec/stale_citations.go`
- creates: `internal/spec/citation_sources_working_tree_test.go`

## Verification

- `out="$(go test -timeout 10m -count=1 -v -run '^(TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks|TestTrackedSourceReadersStillRefuseAPathBehindAFile|TestStaleHistoryCitationsReportsFilesTheTagHolds|TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles|TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator|TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator|TestCancelTranscriptWarnsAboutAParkedItemAndACitation|TestBaselineUpdateReportsStaleHistoryCitationsTranscript)$' ./internal/spec ./internal/speccheck ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks/citation-sources TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks/history-citations TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks/stale-history-citations TestTrackedSourceReadersStillRefuseAPathBehindAFile/citation-sources TestTrackedSourceReadersStillRefuseAPathBehindAFile/history-citations TestTrackedSourceReadersStillRefuseAPathBehindAFile/stale-history-citations TestStaleHistoryCitationsReportsFilesTheTagHolds TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator TestCancelTranscriptWarnsAboutAParkedItemAndACitation TestBaselineUpdateReportsStaleHistoryCitationsTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/spec` — expected: exit 0. Before this Task the two new tests do not exist, so the command fails at its first pass check, and with the tests added but the readers unchanged every `TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks` subtest fails. Afterwards the six subtests pass with the unchanged Stale History Citations tests and the unchanged consumer tests of the strict check's retired-history pass, `roundfix cancel` and `roundfix baseline update`, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-2; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior; API Contracts 1-3; Invariants 1-2; Testing Approach; Build Order 1

## Result

Implemented API Contracts 1–3: the three tracked-source readers skip only
`errors.Is(err, os.ErrNotExist)` at the specified inspection/read boundaries,
regardless of directory. Each skip documents that a path behind a file and
other errors still fail; existing error messages and symlink/binary handling
remain unchanged.

Added the two parallel reader regressions with all three named parallel
subtests. Each uses a committed temporary Git repository, modifies only the
working tree and runs with the existing private-HOME test helper. Updated
only the two authorized baseline detector-failure injections to share
`breakBaselineCitationHandoff(t, repo)`, replacing the handoff directory with
a regular file; every assertion is retained.

Acceptance evidence:

- Missing tracked sources: the initial focused run failed in all three
  missing-path subtests on `docs/gone.md`. After the reader changes,
  `go test ./internal/spec -run '^TestTrackedSourceReaders' -count=1 -v -timeout 2m`
  exited 0 with all six named subtests passing. The missing-path cases assert
  exactly `README.md` and `notes/a.md`; citation cases also assert the expected
  target and line, proving surviving citations remain visible.
- Other read errors: the same focused run passed all three path-behind-a-file
  cases, each requiring an error naming `notes/a.md`.
- Consumer failure handling:
  `go test ./internal/cli -run '^TestBaselineUpdateHistory(CitationDetectorFailureIsLocal|StaleDetectorFailureKeepsOutcome)$' -count=1 -timeout 3m`
  exited 0 with both amended tests passing and their assertions unchanged.
- `git diff --check` exited 0. Changed-path inspection found only the three
  authorized production files, the new reader test, the authorized baseline
  test file and this Task file. Its pre-existing `in_progress` status was
  preserved.

Environment: the executor refused the symlinked writable root; shell access
used the resolved checkout with approved elevated execution. No repository
configuration was changed.

The Task's declared Verification and repository Verification were not run;
those and Task settlement remain Daemon-owned. No follow-up work identified.
