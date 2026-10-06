---
task: task_05
spec: 0241-retire-the-fool-autoresearch-and-council
status: completed
type: test
complexity: low
---

# Task 05: The special-entry fixture binds its socket under a short root

## Overview

Corrective Task for task_02's Verification failure on macOS. Its first
command requires a `--- PASS` line for
`TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering` in
`skills/repository_test.go`. That test's last case builds a ready repository
fixture under `t.TempDir()` and binds a unix socket at
`<root>/.agents/skills/agentic-cli-design/socket`. On macOS `t.TempDir()`
lives under `/var/folders/...`, so the socket path exceeds the 104-byte
`sun_path` limit, `net.Listen` fails with `bind: invalid argument`, and the
test skips. A skipped test prints `--- SKIP`, so task_02's Verification fails
for a reason unrelated to its work. On Linux CI the path fits and the test
passes.

This Task makes only the special-entry case build its fixture under a short
root, so the socket binds on macOS and the test passes. What the test proves
stays the same: a special filesystem entry inside an external skill's tree is
reported as an error that names its path. It runs before task_02, which edits
the same file.

## Requirements

1. MUST build the fixture of the special-entry case (the `specialRoot` block)
   under a short directory created by `os.MkdirTemp("/tmp", "rfsk")`, removed
   by `t.Cleanup`, so the socket path stays far below 104 bytes. The fixture
   writer MAY be split into a helper that takes the root, with
   `writeReadyRepositoryFixture` delegating to it with `t.TempDir()`.
2. MUST keep the assertion that `CheckRepository` returns an error naming the
   socket path, and keep closing the listener before the short root is
   removed.
3. MAY skip when the short directory cannot be created or the socket cannot
   be bound for another reason; MUST NOT skip because of the path length on
   macOS or Linux.
4. MUST NOT change the other cases of the test, any other test, or any
   production file. MUST NOT chdir, since the package's tests run in
   parallel.

## Subtasks

- [ ] Create the special-entry fixture under a short temporary root with cleanup.
- [ ] Run the test with `-v` on macOS and record its `--- PASS` line.

## Acceptance Criteria

- [ ] On macOS the test prints `--- PASS` instead of `--- SKIP`.
- [ ] The test still fails when `CheckRepository` accepts the socket.

## Context

- interface: `skills/repository_test.go`
- instruction: `skills/repository.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering " || { printf '%s\n' "$out"; printf 'missing pass: TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering\n' >&2; exit 1; }` — expected: exit 0; before this Task the socket bind fails on macOS with `bind: invalid argument`, the test prints `--- SKIP`, and the command fails.

## References

- task_02 → Verification command 1
- `_techspec.md` → Existing tests that change

## Result

The special-entry case now creates its repository root with
`os.MkdirTemp("/tmp", "rfsk")` and registers removal with `t.Cleanup`.
The existing listener cleanup is registered afterward, so it closes before
the directory is removed. The shared fixture writer now accepts a root;
`writeReadyRepositoryFixture` still delegates using `t.TempDir()`, preserving
the other cases and tests. The socket-error assertion is unchanged. A socket
path with a ten-digit temporary suffix is 60 bytes, below the macOS limit.

Focused evidence on macOS (`uname -s`: `Darwin`):

- Ran `GOCACHE=/tmp/roundfix-task05-gocache rtk proxy go test -count=1 -v -run 'TestCheckRepositoryHandlesNestedLinks|TestCheckRepositoryWithExternalUsesExplicitRequirement' ./skills`
  before editing. The target printed `--- SKIP` after binding its long
  `/var/folders/...` socket path failed with `invalid argument`.
- Acceptance criterion 1: the same focused selection after editing, with
  approved execution outside the sandbox, exited 0 and printed
  `--- PASS: TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering (0.06s)`.
  The companion test also passed. Inside the sandbox, the short socket path
  was denied with `operation not permitted`; this was a separate environment
  restriction, resolved by the approved rerun. The default Go build cache was
  also inaccessible, so these checks used the task-scoped cache above.
- Acceptance criterion 2: ran the focused selection with
  `-overlay=/tmp/rf-task05-mutation-_3_4b8xn/overlay.json`. The temporary
  overlay made the skill-folder walker ignore sockets, allowing
  `CheckRepository` to accept the socket. The target exited 1 and printed
  `--- FAIL: TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering`,
  with `expected external special-entry error naming "/tmp/rfsk1652707670/.agents/skills/agentic-cli-design/socket", got <nil>`.
  The overlay lived outside the repository and was removed afterward; no
  production file was edited.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exited 0.
- The first `GOCACHE=/tmp/roundfix-task05-gocache rtk make verify-incremental`
  run exited 2: the repository guard detected this Agent appending Result
  while the suite was running. Its test bodies reported `PASS`, but the
  guard rejected the mid-run change to this task file. The fresh rerun held
  the repository unchanged until exit and exited 0, including formatting,
  vet, tests, skill checks, and build.

The declared Verification command remains for the Daemon. Task status and
the Task Graph were not edited by this Agent. No follow-up changes are
included.

## Carry-forward provenance

- Source Run: `run_20261006T214048Z_131051a4f3c19cbe`
- Source commit: `cceb70017fb8f1c3bce8833750e88c8206b5d25e`
