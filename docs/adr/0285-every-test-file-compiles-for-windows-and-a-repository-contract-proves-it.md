---
status: accepted
created_at: 2026-10-17T20:40:00Z
updated_at: 2026-10-17T20:40:00Z
deprecated_at: null
superseded_by: null
---

# Every test file compiles for Windows, and a repository contract proves it

Windows is a release platform, and the coverage record lists each test with the
release platforms that build it, yet nothing compiled the tests for Windows. On
2026-10-17 `GOOS=windows go vet ./...` failed in three packages: Spec 0275
added a delivery probe test that reads `syscall.Stat_t` and `syscall.Getpgrp`
without a build constraint, and older test files of `internal/cli`,
`internal/agent` and `internal/baseline` already called `syscall.Kill`,
`syscall.Mkfifo`, `syscall.Wait4`, `SysProcAttr.Setpgid` and the Unix-only
Proven Signal helpers the same way. Production code built for Windows; only
the tests did not, so the coverage record listed those tests on Windows,
where they could not compile.

## Decision

**A test that needs Unix lives in a file that says so.** Every test function
or helper that uses a Unix-only symbol moves, unchanged and under its own name,
into a file with the `unix` build constraint, and the coverage record is
re-recorded with its sanctioned command so it lists those tests for Darwin and
Linux only.

**A repository contract compiles every test for Windows.** A Repository
Contract Test runs `go vet` with `GOOS=windows` and the `repocontract` and
`docscontract` tags over every package of the repository and fails with the
compiler's diagnostics. It keeps the default Contract Relevance, so the Full
Contract Run on main and before each release runs it, and a change to its own
package selects it; it costs about 5 s with a warm build cache and 12 s with a
cold one on the maintainer's machine.

## Considered Options

- Accepting the gap as a platform-only limit was rejected: the coverage record
  claims those tests for Windows, and a release ships a Windows build.
- A static scan of `syscall.<name>` selectors was rejected: it misses struct
  fields such as `Setpgid` and repository helpers that are Unix-only.
- Making the contract `always` was rejected: every Task's settlement would pay
  its cost for a defect class that the Full Contract Run catches before a
  release.
- Editing the Makefile or a CI workflow to add the check was not authorized.
