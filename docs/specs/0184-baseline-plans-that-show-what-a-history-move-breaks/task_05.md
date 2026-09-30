---
task: task_05
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: pending
type: backend
complexity: low
---

# Task 05: The Windows citation open reads files synchronously

## Overview

task_01 opens each tracked file on Windows through `windows.CreateFile` with `FILE_FLAG_OVERLAPPED` in `internal/baseline/history_citations_open_windows.go`, then reads the handle through ordinary `os.File` reads. An overlapped handle needs asynchronous I/O, so the scan on Windows can fail its reads and report every regular file as unscanned. `roundfix` ships a `win32-x64` package. The non-blocking open exists only to keep a FIFO from blocking the planner, and Windows filesystems have no such FIFO. The pre-PR review of the 0184 candidate found this.

## Requirements

1. MUST open the file on Windows for synchronous reads: drop `FILE_FLAG_OVERLAPPED` and keep `FILE_FLAG_OPEN_REPARSE_POINT`, so a reparse point is still never followed.
2. MUST keep the Unix open (`O_RDONLY|O_NOFOLLOW|O_NONBLOCK`) and every scanner behavior unchanged.
3. MUST keep the package's non-test code building for `GOOS=windows`.
4. MUST NOT edit any test file. Existing Unix-only tests, such as the FIFO case in `internal/baseline/repository_test.go`, are out of scope and governed.
5. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Drop the overlapped flag from the Windows open.
- [ ] Prove the Windows build still passes.

## Acceptance Criteria

- [ ] `internal/baseline/history_citations_open_windows.go` no longer requests overlapped I/O and still opens reparse points without following them.
- [ ] `GOOS=windows GOARCH=amd64 go build ./internal/baseline` exits `0`, and no test file changes.

## Context

- interface: `internal/baseline/history_citations_open_windows.go`

## Verification

- `! grep -q "FILE_FLAG_OVERLAPPED" internal/baseline/history_citations_open_windows.go && grep -q "FILE_FLAG_OPEN_REPARSE_POINT" internal/baseline/history_citations_open_windows.go && GOOS=windows GOARCH=amd64 go build -buildvcs=false -o /dev/null ./internal/baseline && GOOS=windows GOARCH=amd64 go build -buildvcs=false -o /dev/null ./cmd/roundfix` — expected: exit 0; before this Task the file still requests `FILE_FLAG_OVERLAPPED`, so the first check fails.

## References

- task_01
- `_techspec.md` → What the scan reads

## Result
