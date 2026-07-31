---
task: task_01
spec: 0055-owner-identity-without-fork
status: completed
type: backend
complexity: high
---

# Task 01: Read owner start identity from the kernel, never a fork

## Overview

`processStartIdentity` shells out to `ps -o lstart=` on every ownership proof, so
Force Stop fails exactly when the host cannot fork — the condition it exists to
escape. Replace it with a direct kernel read per platform: procfs on Linux,
`sysctl` on macOS. The token stays opaque and equality-compared.

## Requirements

1. MUST remove `processStartIdentity` from `internal/store/process_unix.go` and
   provide it in `process_linux.go` (`//go:build linux`) and
   `process_darwin.go` (`//go:build darwin`).
2. MUST spawn no subprocess on either platform and MUST NOT use cgo.
3. MUST read the Linux start time as field 22 of `/proc/<pid>/stat`, counting
   fields from the **last** `)` so a comm containing spaces or parentheses
   cannot shift the index.
4. MUST read the macOS start time from `KERN_PROC_PID` through the already
   required `golang.org/x/sys/unix`; no new module may be added.
5. MUST prefix the token with its platform (`linux:` / `darwin:`) so a token
   this platform could not have produced is recognizable as such.
6. MUST return the same "process is gone" signal the caller already handles when
   the process does not exist (ENOENT on Linux, ESRCH on macOS), distinct from a
   read failure.
7. MUST add `process_unix_other.go` (`//go:build unix && !linux && !darwin`)
   returning an unreadable-identity error, so an unsupported Unix degrades
   instead of failing to build.
8. MUST leave `processAbsent`, `signalOwnerProcess`, `process_windows.go`, and
   `process_other.go` unchanged.

## Subtasks

- [ ] Add the two kernel implementations and the unsupported-Unix stub.
- [ ] Remove the `ps` implementation and its `os/exec` dependency from
      `process_unix.go`.
- [ ] Prove no subprocess is spawned, and prove the token is stable across
      repeated reads of one live process.

## Acceptance Criteria

- [ ] The identity read spawns no process: a test that shadows `ps` on `PATH`
      with an executable that fails the test if invoked stays green.
- [ ] Two consecutive reads of the same live process return the same token.
- [ ] Reading a nonexistent PID reports the process-gone condition, not a read
      failure.
- [ ] A `/proc/<pid>/stat` comm containing a space and a `)` still yields the
      correct field (covered with a synthetic stat payload).
- [ ] `go build` succeeds for linux, darwin, and windows targets.

## Context

- interface: `internal/store/process_unix.go`
- interface: `internal/store/process_linux.go`
- interface: `internal/store/process_darwin.go`
- interface: `internal/store/process_unix_other.go`
- interface: `internal/store/process_unix_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `GOOS=linux go build ./internal/store/` and `GOOS=windows go build ./internal/store/`
  — expected: both compile.
- `go test -count=1 ./internal/store/` — expected: pass, including the
  no-subprocess and comm-parsing cases.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 1, Story 1, Feature 1; `_techspec.md` → Build Order 1,
Interfaces, Risks (procfs field parsing).

## Result

Implemented direct, platform-tagged owner start-identity reads. Linux now reads
field 22 from `/proc/<pid>/stat` after locating the comm field's last `)`,
macOS reads `kern.proc.pid` through `golang.org/x/sys/unix`, and unsupported
Unix targets return an unreadable/unsupported identity error. Tokens are
prefixed with `linux:` or `darwin:`. The common Unix process absence and signal
functions, Windows implementation, and non-Unix implementation are unchanged.

Focused checks and acceptance evidence:

- No subprocess: before the implementation,
  `GOCACHE=<repo>/.gocache rtk go test ./internal/store -run '^TestOwnerProcessIdentityDoesNotSpawnPS$'`
  failed because the PATH-shadowed `ps` exited 99. After the implementation,
  the same case passed within the 15-case focused owner-process run. A static
  absence check over the four Unix production files also found no `os/exec`,
  `exec.Command`, or `import "C"`.
- Stable token: `TestOwnerProcessIdentityIsStableForOneProcess` passed twice in
  focused runs and also asserted the current `darwin:` platform prefix.
- Process gone: `TestOwnerProcessIdentityFailsForAbsentProcess` passed on the
  Darwin host and matched `ESRCH`; the Linux read wraps `os.ReadFile` errors
  with `%w`, preserving `ENOENT` for `/proc/<pid>/stat`.
- Proc stat parsing: the Linux-tagged synthetic payload includes spaces and a
  `)` inside comm; `GOCACHE=<repo>/.gocache rtk go test -count=1
  ./internal/store/process_linux.go ./internal/store/process_linux_test.go`
  passed both the field-22 and malformed-payload cases.
- Target compilation: cgo-disabled focused compilation/execution passed on
  Darwin (`15` owner-process cases), and cgo-disabled `go test -c` succeeded
  for Linux, Windows, and unsupported-Unix FreeBSD package test binaries.
- `rtk git diff --check` exited 0.

The commands under `## Verification` were not run; Daemon Verification owns
those commands and the Task verdict.
