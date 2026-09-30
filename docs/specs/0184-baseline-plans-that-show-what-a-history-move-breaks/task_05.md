---
task: task_05
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: completed
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

The Windows citation opener now creates a synchronous read handle by omitting
`FILE_FLAG_OVERLAPPED`. It retains `FILE_FLAG_OPEN_REPARSE_POINT`, so the
existing no-follow behavior remains in place. The Unix opener, scanner logic,
exported signatures and tests are unchanged.

Focused-check evidence:

- `rtk rg -n 'FILE_FLAG_(OVERLAPPED|OPEN_REPARSE_POINT)' internal/baseline/history_citations_open_windows.go internal/baseline/history_citations_open_unix.go` exited 0 and printed only the Windows `FILE_FLAG_OPEN_REPARSE_POINT` use.
- `rtk env GOCACHE=/private/tmp/roundfix-0184-task05-gocache GOOS=windows GOARCH=amd64 go list -f '{{range .GoFiles}}{{$.Dir}}/{{.}} {{end}}' ./internal/baseline` exited 0. Compiling that production-file list with `go test -c -o /private/tmp/roundfix-0184-task05-baseline.test.exe` under the same Windows environment exited 0 with `[no test files]`. This focused check compiles the package's non-test Windows source without invoking the declared Verification command.
- `rtk env GOCACHE=/private/tmp/roundfix-0184-task05-gocache make verify-incremental` exited 0 after rerunning with process-table access. The sandboxed attempt had failed only in two `internal/cli` force-stop integration tests because process-table access returned `operation not permitted`; the unrestricted rerun passed those tests and the complete incremental target.
- `rtk git diff --name-only -- '*_test.go'` exited 0 with no output, and `rtk git diff --check` exited 0.
- The Task's declared `## Verification` command was not run; the Daemon owns that gate.

Acceptance evidence:

- The Windows flag inspection shows synchronous I/O with the reparse-point guard retained, and the one-line production diff changes no other opener or scanner behavior.
- The focused Windows production compile exits 0, while the empty test-file diff proves no test was edited. The Daemon still owns the exact `GOOS=windows GOARCH=amd64 go build` acceptance command.

## Carry-forward provenance

- Source Run: `run_20260930T011500Z_f0e23f495ac8b795`
- Source commit: `381761bc2e7c0c733d7ec809d6a01b2e968a7fcd`
