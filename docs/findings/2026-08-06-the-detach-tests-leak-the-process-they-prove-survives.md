---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
---

# The detach tests leak the process they prove survives

**Date:** 2026-08-06
**Found by:** `ps` on a developer machine, while looking for stale Runs before
starting unrelated work. Nothing in Roundfix reported these.

Four orphaned processes were running on one machine, aged **1 day 22 hours to
3 days 15 hours**, two of them burning **1.9% CPU continuously**. Together they
had consumed **2h40m of CPU time doing nothing**. All four came from the
Roundfix Go test suite, not from any real Run.

| PID | Age | CPU consumed | Origin |
| --- | --- | --- | --- |
| 33864 | 3d 14h | 86 min | `TestRunImplementDetachSurvivesCallerProcessGroupKill` |
| 55351 | 3d 02h | 74 min | same test, a different execution |
| 54551 | 3d 02h | 8 s | `cli.test implement --spec 0001-widget-flow` |
| 87902 | 1d 22h | 4 s | same, with `--detach` |

Three were reparented to PID 1; `54551` was a session leader (`Ss`) with `55351`
as its child. The Go build directories that produced two of them had already
been garbage-collected — the owning test binaries were long gone.

## Why it leaks, and why it can never self-heal

`TestRunImplementDetachSurvivesCallerProcessGroupKill`
(`internal/cli/implement_test.go:1508`) does exactly what its name says:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
// ...
syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)   // kill the caller's whole group
```

The detached child survives that kill **by design** — that is the property under
test. But the test's only lever over the surviving child is a release file:

```go
dir := t.TempDir()
releasePrompt := filepath.Join(dir, "release")
t.Cleanup(func() { _ = os.WriteFile(releasePrompt, []byte("release\n"), 0o644) })
```

and the fake ACPX blocks on that file with a busy poll:

```sh
while [ ! -f "$ROUNDFIX_FAKE_ACPX_RELEASE" ]; do
  sleep 0.05
done
```

Twenty filesystem probes per second, forever, until the file appears. That is
the 1.9% CPU.

The trap is that **the release file lives inside `t.TempDir()`, which Go deletes
when the test ends.** Once the directory is gone, the file the loop waits for
can never be created by anyone. The loop is not merely unterminated — it is
*provably* infinite. Nothing in the test, the suite, or Roundfix itself will
ever end it. Only a human reading `ps` will.

Any path that ends the test before `mustWrite(t, releasePrompt, ...)` runs — a
`t.Fatalf`, a panic, a package timeout, a Ctrl-C, a CI job cancellation — leaks
a detached child plus its fake ACPX permanently. `t.Cleanup` covers the ordinary
failure path; it covers nothing when the test *binary* dies, because cleanups do
not run then.

## This failure mode was already on the record

Spec 0076's QA report already caught this test failing under full-suite load:

> `TestRunImplementDetachSurvivesCallerProcessGroupKill` exceeded its fixed 5s
> line deadline at 6.21s. The test passed 20/20 focused runs and the isolated
> `internal/cli` package passed, narrowing the nondeterminism to full-suite load.

That was filed as a flake. It is also, unnoticed, the **leak trigger**: the
timeout path is precisely the path that abandons a detached child. Every
full-suite sweep that trips this deadline strands one more pair of processes.
The two oldest orphans found here are three days apart, which is consistent with
accumulation across ordinary suite runs rather than a single bad day.

## Why no Roundfix command surfaces it

`roundfix runs list` reports `No Runs found` on this machine. These children
never registered a Run in the Run Database that the operator queries — they are
test-fixture Runs under throwaway `HOME` directories. So every diagnostic
Roundfix offers is blind to them:

- `runs list` is scoped to the current repository and to real Run records.
- `reconcile` classifies Run Worktrees, not stray processes.
- `gc` prunes journals and artifacts, not the processes that abandoned them.
- `doctor` checks readiness, not residue.

A tool whose central promise is detached execution has no command that answers
"what did I detach that is still running?" That gap is the reason these survived
three days on a machine whose owner actively maintains the project.

## Asks

1. **Terminate the detached child in the test, not just release it.** Record the
   detached PID (the Run ID is already parsed from stdout) and register a
   `t.Cleanup` that signals that process group directly. The release file is a
   cooperative shutdown; it needs a non-cooperative backstop, because the whole
   point of the test is that this child ignores the caller's death.

2. **Move the release file out of `t.TempDir()`.** As long as the sentinel lives
   in a directory Go deletes, a cleanup that writes it races the cleanup that
   removes it, and any late waiter is stranded by construction. `os.MkdirTemp`
   with an explicit removal ordered after the process kill would close this.

3. **Bound the fake ACPX wait.** A test fixture that polls forever converts every
   test-harness failure into a permanent CPU leak. A ceiling — say 120s, then
   exit non-zero — turns an invisible orphan into a visible test failure, which
   is the outcome a test suite is supposed to produce.

4. **Give Roundfix a way to see its own residue.** Something like
   `roundfix runs list --orphans`, or a `doctor` line that reports detached
   processes with no live Run record. Detached execution without an inventory
   command means the operator's only tool is `ps` and the knowledge that they
   should run it.

5. **Reclassify the 0076 flake.** The 6.21s-vs-5s deadline miss is not only
   nondeterminism under load; it is the leak's entry point. Fixing the deadline
   without fixing the teardown leaves the leak in place under any other failure.

## Evidence

- `internal/cli/implement_test.go:1508`
  (`TestRunImplementDetachSurvivesCallerProcessGroupKill`), and the fake ACPX
  release loop in the same file around line 220.
- `docs/specs/_archived/0076-force-stop-exit-proof/qa/qa-report-2026-08-04.md`,
  QA-11 and finding F-01 — the same test timing out under full-suite load.
- Process table captured 2026-08-06 before termination: PIDs 33864, 54551,
  55351, 87902, ages and CPU as tabulated above. All four ended on `SIGTERM`;
  none required `SIGKILL`, confirming they were idle waiters rather than wedged.
- `roundfix runs list` on the same machine at the same moment:
  `No Runs found. (119 terminal Run(s) hidden; use --state all)`.
