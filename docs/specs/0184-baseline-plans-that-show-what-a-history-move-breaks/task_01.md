---
task: task_01
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: pending
type: backend
complexity: high
---

# Task 01: Find the Relocation Citations a set of History Relocations breaks

## Overview

Baseline planning names each History Relocation but not the citations it
breaks. This Task adds the scan, `relocationCitationFindings`, in a new
`internal/baseline/history_citations.go`. Given the repository root and the
ordered History Relocation ledger, it returns one Baseline warning per tracked
file whose citations resolve before the relocations and not after them. It adds
at most one omitted summary and one unscanned summary. The scan reads only
paths in the Git index of the adopter's repository, and it reads them from the
working tree. Its findings reach whoever reviews the plan: the maintainer in
text, and an Agent in JSON. The scan therefore must never open an untracked or
ignored file, follow a symbolic link, leave the repository root, or quote file
content other than path text. This Task does not wire the scan into planning.

## Requirements

1. MUST add, in a new `internal/baseline/history_citations.go`,
   `relocationCitationFindings(ctx context.Context, root string, moves []HistoryMove) ([]Finding, error)`.
   It returns nil with no Git call and no filesystem access when `moves` is
   empty.
2. MUST list the index with `listTrackedPaths(ctx, root)`. That function runs
   `git -C <root> -c core.fsmonitor=false ls-files -z --cached --full-name`
   with `GIT_OPTIONAL_LOCKS=0` and `GIT_TERMINAL_PROMPT=0`. It splits the
   output on NUL without trimming, drops empty and duplicate entries, keeps
   only paths for which `repositoryPathIsSafe` holds, and sorts them. A Git
   failure MUST be returned as an error.
3. MUST implement the resolution model in `_techspec.md` → Implementation
   Design → The resolution model:
   - files before and files after;
   - unit directories strictly below a move's `historyMoveSourceRoot`, with
     family and legacy roots never targets;
   - the citing file's after-location;
   - repository-path tokens in every scanned file;
   - Markdown link destinations (inline links, images and reference
     definitions, outside fenced code blocks) in `.md` and `.markdown` files;
   - one citation per `(line, before-resolution)`.

   A citation MUST count only when its before-resolution exists before and its
   after-resolution does not exist after.
4. MUST read a tracked path only after `os.Lstat` shows that neither the path
   nor any of its parent directories is a symbolic link, and that the path is
   a regular file. It MUST skip a file whose first 8,000 bytes contain a NUL
   byte. A regular file larger than 4 MiB, or one that fails to open or read,
   MUST go into the unscanned summary instead. It MUST NOT stat or open any
   path that is not in the index.
5. MUST return the findings in `_techspec.md` → Data Models:
   - one `baseline.history.citation` finding per citing file, whose path is the
     citing file and whose message follows the exact per-citation phrasing,
     capped at three citations plus `; and <N> more`;
   - at most 200 citing files, then one `baseline.history.citation.omitted`
     finding with path `.`;
   - one `baseline.history.citation.unscanned` finding with path `.`, naming
     the count and the first three unscanned paths.

   Findings MUST be ordered by path, then the omitted summary, then the
   unscanned summary.
6. MUST produce byte-identical findings for two runs over the same repository.
7. MUST put its tests in a new `internal/baseline/history_citations_test.go`,
   over real temporary Git repositories, and MUST NOT edit any existing test
   file.

## Subtasks

- [ ] List the index safely and build the before, after and unit-directory
      sets.
- [ ] Extract repository-path tokens and Markdown link destinations, and
      resolve both sides.
- [ ] Read only safe, regular, tracked text files, and summarize unscanned
      ones.
- [ ] Compose, cap and order the findings.
- [ ] Add one test per acceptance criterion, with each negative case separate.

## Acceptance Criteria

- [ ] Each citation form is reported once, with its line: a repository path,
      an inline relative link, an image, a reference definition and a
      root-relative link.
- [ ] A relocated ADR's relative link to an ADR that stays is reported with its
      after-resolution.
- [ ] A link between two co-relocated ADRs is not reported.
- [ ] A citation already broken before the plan is not reported.
- [ ] A link inside a fenced code block, and a URL with a scheme, are not
      reported as links.
- [ ] An untracked file and an ignored file that hold citations and are made
      unreadable produce no finding and no unscanned entry.
- [ ] A tracked symbolic link, and a file under a symlinked directory, are not
      followed.
- [ ] A binary file is skipped, and an oversized tracked text file appears only
      in the unscanned summary.
- [ ] The three-citation cap and the 200-file cap hold, each with its summary.
- [ ] A citation of a relocated unit directory is reported, while a citation of
      a family root is not.
- [ ] With no moves, nothing is returned and Git is never called.
- [ ] Two runs return byte-identical findings.

## Context

- creates: `internal/baseline/history_citations.go`
- creates: `internal/baseline/history_citations_test.go`
- instruction: `internal/baseline/history_layout.go`
- instruction: `internal/baseline/repository.go`
- instruction: `docs/adr/0173-a-baseline-plan-reports-the-citations-its-history-relocations-break.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRelocationCitationsReportEachCitationForm|TestRelocationCitationsReportARelocatedFilesOutwardLink|TestRelocationCitationsSkipCoRelocatedLinks|TestRelocationCitationsSkipAlreadyBrokenCitations|TestRelocationCitationsSkipCodeFencesAndURLs|TestRelocationCitationsNeverOpenUntrackedOrIgnoredFiles|TestRelocationCitationsNeverFollowSymbolicLinks|TestRelocationCitationsSummarizeUnscannedFiles|TestRelocationCitationsCapEachFileAndThePlan|TestRelocationCitationsCountUnitDirectoriesNotFamilyRoots|TestRelocationCitationsDoNothingWithoutMoves|TestRelocationCitationsAreDeterministic)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRelocationCitationsReportEachCitationForm TestRelocationCitationsReportARelocatedFilesOutwardLink TestRelocationCitationsSkipCoRelocatedLinks TestRelocationCitationsSkipAlreadyBrokenCitations TestRelocationCitationsSkipCodeFencesAndURLs TestRelocationCitationsNeverOpenUntrackedOrIgnoredFiles TestRelocationCitationsNeverFollowSymbolicLinks TestRelocationCitationsSummarizeUnscannedFiles TestRelocationCitationsCapEachFileAndThePlan TestRelocationCitationsCountUnitDirectoriesNotFamilyRoots TestRelocationCitationsDoNothingWithoutMoves TestRelocationCitationsAreDeterministic; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf "missing PASS for %s\\n" "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the named tests exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1–2; User Story 2; Core Features 1, 2, 3 and 5;
  Success Metric 1; Recorded limits
- [_techspec.md](_techspec.md) — Interfaces; The resolution model; What the
  scan reads; Data Models; Testing Approach 1; Build Order 1
- ADR-0173; ADR-0120; ADR-0071
