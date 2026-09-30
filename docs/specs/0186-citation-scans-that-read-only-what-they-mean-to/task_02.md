---
task: task_02
spec: 0186-citation-scans-that-read-only-what-they-mean-to
status: pending
type: backend
complexity: medium
---

# Task 02: The Relocation Citation scan reads through a repository root

## Overview

The Relocation Citation scan in `internal/baseline/history_citations.go` checks each tracked path's parent directories with `os.Lstat`, caches the result and then opens the file by its joined absolute path, with no-follow only on the last component. A tracked directory replaced by a symbolic link after the check redirects the open, possibly outside the repository, and the `os.SameFile` check still passes because both sides resolve through the link. This Task opens every tracked file through one `os.Root` opened on the repository root (ADR-0177). `os.Root` walks each component relative to a held handle and never leaves the root. The existing checks stay.

## Requirements

1. MUST open `os.OpenRoot(root)` once per `relocationCitationFindings` call that reads files, and close it before returning. It MUST open no root and touch no filesystem when there are no moves.
2. MUST change `openCitationFileNoFollow` in both platform files to take `(*os.Root, relative string)`:
   - on Unix it opens `root.OpenFile(filepath.FromSlash(relative), os.O_RDONLY|syscall.O_NONBLOCK, 0)`;
   - on Windows it opens `root.Open(filepath.FromSlash(relative))`.

   No absolute joined path may be opened for reading.
3. MUST keep the pre-open `lstatTrackedCitationPath` check, the post-open uncached re-check, the `os.SameFile` comparisons, the binary, oversize and FIFO handling, and every finding unchanged.
4. MUST declare a package-level test hook `citationBeforeOpen func(relative string)` in `history_citations.go`. It is nil in production, called after the pre-open check and immediately before the open, and assigned only by tests.
5. MUST NOT edit `go.mod`, any governed file or any existing test file, and MUST NOT add a dependency. It MUST rename or remove no top-level test and change no exported function signature.
6. MUST put the new tests in `internal/baseline/history_citations_root_test.go`, over real temporary Git repositories. The race test sets the hook to replace a tracked directory with a symbolic link to a directory outside the repository. That outside directory holds a file under the same relative name citing a relocated path the in-repository file does not cite.

## Subtasks

- [ ] Open the repository through `os.Root` in the scan.
- [ ] Route both platform opens through the root.
- [ ] Add the test hook and the root-containment tests.

## Acceptance Criteria

- [ ] A tracked directory swapped for a symbolic link to an outside directory between the check and the open yields no finding from the outside file's content.
- [ ] An unraced repository still reports its citations, and the existing `TestRelocationCitations*` tests, including the FIFO and symbolic-link cases, pass unchanged.
- [ ] The non-test code builds for `GOOS=windows`.

## Context

- instruction: `docs/adr/0177-the-relocation-citation-scan-reads-through-a-repository-root.md`
- instruction: `docs/adr/0173-a-baseline-plan-reports-the-citations-its-history-relocations-break.md`
- interface: `internal/baseline/history_citations.go`
- interface: `internal/baseline/history_citations_open_unix.go`
- interface: `internal/baseline/history_citations_open_windows.go`
- creates: `internal/baseline/history_citations_root_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRelocationCitationsNeverReadThroughASwappedDirectory|TestRelocationCitationsReadThroughTheRepositoryRoot|TestRelocationCitationsNeverFollowSymbolicLinks|TestRelocationCitationsNeverBlockOnAFIFO|TestRelocationCitationsReportEachCitationForm)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRelocationCitationsNeverReadThroughASwappedDirectory TestRelocationCitationsReadThroughTheRepositoryRoot TestRelocationCitationsNeverFollowSymbolicLinks TestRelocationCitationsNeverBlockOnAFIFO TestRelocationCitationsReportEachCitationForm; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && GOOS=windows GOARCH=amd64 go build -buildvcs=false -o /dev/null ./cmd/roundfix` — expected: exit 0; before this Task the two new named tests do not exist, so the command fails; after it, the non-test code also builds for Windows.

## References

- [_prd.md](_prd.md) — Goals 2 and 3; Core Feature 2; Success Metric 2; Success Metric 3
- [_techspec.md](_techspec.md) — System Architecture; Interfaces; Testing Approach 2 and 3; Build Order 2
- [references/2026-09-30-citation-scan-parent-directory-race.md](references/2026-09-30-citation-scan-parent-directory-race.md)
- ADR-0177; ADR-0173

## Result
