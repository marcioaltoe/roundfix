---
task: task_05
spec: 0241-retire-the-fool-autoresearch-and-council
status: pending
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
