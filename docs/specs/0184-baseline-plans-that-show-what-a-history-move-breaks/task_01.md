---
task: task_01
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: completed
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
   `relocationCitationFindings(ctx context.Context, root string, moves []HistoryMove, refused map[string]bool) ([]Finding, error)`, where `refused` names by `From` the moves discovery reported as occupied-destination collisions.
   It returns nil with no Git call and no filesystem access when `moves` is
   empty.
2. MUST list the index with `listTrackedPaths(ctx, root)`. That function runs
   `git -C <root> -c core.fsmonitor=false ls-files -z --cached --full-name`
   with `GIT_OPTIONAL_LOCKS=0` and `GIT_TERMINAL_PROMPT=0`. It splits the
   output on NUL without trimming, drops empty and duplicate entries, keeps
   only paths for which `repositoryPathIsSafe` holds and that contain no
   control character (U+0000–U+001F, U+007F–U+009F), and sorts them. A path
   dropped this way is never opened or printed. A Git
   failure MUST be returned as an error.
3. MUST implement the resolution model in `_techspec.md` → Implementation
   Design → The resolution model:
   - files before and files after, where a move that history-layout discovery
     reports as an occupied-destination collision is left out, because
     `baseline apply` refuses it and leaves its source in place. The scan
     receives those collisions from `planHistoryMoves`; it never stats an
     untracked destination itself;
   - unit directories strictly below a move's `historyMoveSourceRoot`, with
     family and legacy roots never targets;
   - the citing file's after-location;
   - repository-path tokens in every scanned file;
   - Markdown link destinations (inline links, images and reference
     definitions, outside fenced code blocks) in `.md` and `.markdown` files;
   - one citation per `(line, before-resolution)`.

   A citation MUST count only when its before-resolution exists before and its
   after-resolution does not exist after. A link destination whose raw or
   percent-decoded form contains a control character MUST be skipped, so no
   finding ever prints one.
4. MUST read a tracked path only after `os.Lstat` shows that neither the path
   nor any of its parent directories is a symbolic link, and that the path is
   a regular file. It MUST then open the path with `O_RDONLY|O_NOFOLLOW|O_NONBLOCK`, so a path replaced by a FIFO or device after the check can never block the planner, and
   read it only when `os.SameFile` holds between the opened file's `Stat` and
   that `os.Lstat` and the opened file's mode is regular; a mismatch or a non-regular file skips the file, so a path swapped for a link
   between the check and the open is never read. It MUST skip a file whose first 8,000 bytes contain a NUL
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
- [ ] A tracked path replaced by a FIFO is skipped without blocking the plan.
- [ ] A tracked symbolic link, and a file under a symlinked directory, are not
      followed.
- [ ] A binary file is skipped, and an oversized tracked text file appears only
      in the unscanned summary.
- [ ] The three-citation cap and the 200-file cap hold, each with its summary.
- [ ] A citation of a relocated unit directory is reported, while a citation of
      a family root is not.
- [ ] With no moves, nothing is returned and Git is never called.
- [ ] Two runs return byte-identical findings.
- [ ] A tracked path and a link destination that contain a control character
      are never printed in any finding.
- [ ] A move whose destination is already occupied, which `baseline apply`
      refuses, produces no citation finding.

## Context

- creates: `internal/baseline/history_citations.go`
- creates: `internal/baseline/history_citations_open_unix.go`
- creates: `internal/baseline/history_citations_open_windows.go`
- creates: `internal/baseline/history_citations_test.go`
- instruction: `internal/baseline/history_layout.go`
- instruction: `internal/baseline/repository.go`
- instruction: `docs/adr/0173-a-baseline-plan-reports-the-citations-its-history-relocations-break.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRelocationCitationsReportEachCitationForm|TestRelocationCitationsReportARelocatedFilesOutwardLink|TestRelocationCitationsSkipCoRelocatedLinks|TestRelocationCitationsSkipAlreadyBrokenCitations|TestRelocationCitationsSkipCodeFencesAndURLs|TestRelocationCitationsNeverOpenUntrackedOrIgnoredFiles|TestRelocationCitationsNeverFollowSymbolicLinks|TestRelocationCitationsSummarizeUnscannedFiles|TestRelocationCitationsCapEachFileAndThePlan|TestRelocationCitationsCountUnitDirectoriesNotFamilyRoots|TestRelocationCitationsDoNothingWithoutMoves|TestRelocationCitationsAreDeterministic|TestRelocationCitationsNeverPrintAControlCharacter|TestRelocationCitationsIgnoreAMoveApplyWouldRefuse|TestRelocationCitationsNeverBlockOnAFIFO)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRelocationCitationsReportEachCitationForm TestRelocationCitationsReportARelocatedFilesOutwardLink TestRelocationCitationsSkipCoRelocatedLinks TestRelocationCitationsSkipAlreadyBrokenCitations TestRelocationCitationsSkipCodeFencesAndURLs TestRelocationCitationsNeverOpenUntrackedOrIgnoredFiles TestRelocationCitationsNeverFollowSymbolicLinks TestRelocationCitationsSummarizeUnscannedFiles TestRelocationCitationsCapEachFileAndThePlan TestRelocationCitationsCountUnitDirectoriesNotFamilyRoots TestRelocationCitationsDoNothingWithoutMoves TestRelocationCitationsAreDeterministic TestRelocationCitationsNeverPrintAControlCharacter TestRelocationCitationsIgnoreAMoveApplyWouldRefuse TestRelocationCitationsNeverBlockOnAFIFO; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf "missing PASS for %s\\n" "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the named tests exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1–2; User Story 2; Core Features 1, 2, 3 and 5;
  Success Metric 1; Recorded limits
- [_techspec.md](_techspec.md) — Interfaces; The resolution model; What the
  scan reads; Data Models; Testing Approach 1; Build Order 1
- ADR-0173; ADR-0120; ADR-0071

## Result

Implemented the Relocation Citation scan without wiring it into planning. The
scan lists and sorts safe Git-index paths, models files and unit directories
before and after applied moves, resolves repository-path tokens and supported
Markdown destinations from each citing file's before and after location, and
emits the capped citation, omitted and unscanned findings. Occupied-destination
moves named in `refused` leave both the source set and citing-file location
unchanged.

Tracked files are read from the working tree only after parent and final-path
`Lstat` checks. The final path is opened read-only with `O_NOFOLLOW` and
`O_NONBLOCK` on Unix (and the no-follow reparse-point equivalent on Windows),
then accepted only when the opened file is regular and
`os.SameFile` matches the checked file. Symlinks, changed file types, FIFOs and
binary files are skipped; oversized and genuinely unreadable regular files are
summarized without reading untracked or ignored paths.

Focused-check evidence:

- `GOCACHE=/private/tmp/roundfix-0184-task01-gocache go test -count=1 -run '^TestRelocationCitations' ./internal/baseline` exited 0 after the final implementation edit.
- `GOCACHE=/private/tmp/roundfix-0184-task01-gocache go test -count=1 ./internal/baseline` exited 0 in 56.796s after the final implementation edit.
- `GOCACHE=/private/tmp/roundfix-0184-task01-gocache go vet ./internal/baseline` exited 0.
- `go build ./cmd/roundfix` cross-compilation exited 0 for the five release
  targets: Darwin arm64/amd64, Linux arm64/amd64 and Windows amd64. A Windows
  cross-compile of the baseline test binary remains unavailable because the
  pre-existing `repository_test.go` uses the Unix-only `syscall.Mkfifo`.
- `make verify-incremental` reached and passed `go vet ./...` and the changed `internal/baseline` package, then the sandbox blocked another package's network access to `cafe.github.com`. An escalated rerun was denied because the broad suite could send unknown test data to that external endpoint. This environment block is not used as acceptance evidence.
- `git diff --check` exited 0.
- The Task's declared `## Verification` command was not run; the Daemon owns that gate.

Acceptance evidence:

- `TestRelocationCitationsReportEachCitationForm` covers a repository path,
  inline link, image, reference definition and root-relative link, each on its
  reported line.
- `TestRelocationCitationsReportARelocatedFilesOutwardLink` covers a relocated
  ADR's relative link to a file that stays, including its changed
  after-resolution.
- `TestRelocationCitationsSkipCoRelocatedLinks` covers a link whose citing and
  cited ADRs move together.
- `TestRelocationCitationsSkipAlreadyBrokenCitations` covers a target that did
  not exist before the move.
- `TestRelocationCitationsSkipCodeFencesAndURLs` covers fenced Markdown and a
  scheme URL.
- `TestRelocationCitationsNeverOpenUntrackedOrIgnoredFiles` makes both excluded
  files unreadable and observes no citation or unscanned finding.
- `TestRelocationCitationsNeverBlockOnAFIFO` replaces an indexed regular file
  with a FIFO and returns without a finding.
- `TestRelocationCitationsNeverFollowSymbolicLinks` covers an indexed symlink
  and an indexed file beneath a symlinked directory.
- `TestRelocationCitationsSummarizeUnscannedFiles` skips a binary and reports
  only the oversized tracked text file.
- `TestRelocationCitationsCapEachFileAndThePlan` observes three displayed
  citations plus `and 1 more`, 200 file findings, and one omitted-file summary.
- `TestRelocationCitationsCountUnitDirectoriesNotFamilyRoots` covers a legacy
  Spec unit and Review Artifact unit while excluding both family roots.
- `TestRelocationCitationsDoNothingWithoutMoves` supplies a nonexistent root
  and observes `(nil, nil)`, proving the zero-move path touches neither Git nor
  the filesystem.
- `TestRelocationCitationsAreDeterministic` compares byte-identical JSON from
  two scans of the same repository.
- `TestRelocationCitationsNeverPrintAControlCharacter` covers both an indexed
  path with a control character and a percent-decoded control in a link
  destination.
- `TestRelocationCitationsIgnoreAMoveApplyWouldRefuse` covers an occupied
  destination passed through `refused` and observes no citation finding.

Plan wiring, renderer/digest integration and apply invariance remain outside
this Task's slice.

## Carry-forward provenance

- Source Run: `run_20260929T222541Z_a2fa4eeca2973de1`
- Source commit: `421c5cccd7359be3f1fb31688515d208af7ebad0`
