---
date: 2026-08-03
surface: internal/cli
status: done
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# A 200 ms attach budget fails the Verification gate under CI load

`TestRunImplementDetachSurvivesCallerProcessGroupKill` failed the CI
Verification gate on PR #98, a documentation-only change touching fourteen
Markdown files and no Go source. Re-running the identical job on the identical
commit passed. The same test passed locally in `make verify` (3,136 tests) both
before and after the failure, and passed in isolation in 0.80 s.

## What was observed

```text
--- FAIL: TestRunImplementDetachSurvivesCallerProcessGroupKill (2.42s)
    implement_test.go:1527: expected attach to detach cleanly from active Run,
    got 2 stderr="roundfix attach failed: open Run Database reader
    \"/tmp/.../001/.roundfix/roundfix.db\": context deadline exceeded"
```

The test was the first CI Verification gate failure in the workflow's last ten
runs; the nine before it, across `main` and four branches, all succeeded.

## Where it comes from

`internal/cli/implement_test.go:1522` gives the `attach` invocation a
200-millisecond context:

```go
attachCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
attachCode := runCLIContext(t, attachCtx, []string{"attach", runID}, &attachStdout, &attachStderr)
```

That one deadline covers two different things. `store.OpenReader` uses it to
stat the file, open the SQLite connection, ping it, and read the migration
version — against a database an active writer is holding — and the attach
follow loop uses it to exit promptly. Only the second is what the test is
about. On an unloaded machine the open costs a few milliseconds and the
conflation is invisible; on a shared runner the open alone can exceed the whole
budget, and the failure surfaces as a context deadline in `OpenReader` rather
than anywhere near the behavior under test.

The line two statements above it already uses a named budget,
`waitForFile(t, promptStarted, implementWaitBudget)`, so the hardcoded literal
is the outlier at this call site rather than the house style.

## Why it matters

This is the fourth occurrence of the class Spec 0071 fixed in three packages:
a wait budget sized against an unloaded developer machine, which turns a
correct change red on a loaded runner. Its cost is not the minute of CI — it is
that a red gate on an unrelated change trains the next reader to re-run rather
than to read, and the run after that may be a real defect.

The shape of a fix is a deadline per concern: enough time to open the reader,
and a separate short bound on the follow loop that is the actual subject of the
assertion. Sizing the single budget upward would hide the conflation rather
than remove it.

## Evidence

- Failing run: GitHub Actions run 30862459339, job 91847174679, PR #98.
- Passing rerun: the same run id after `gh run rerun --failed`, same commit.
- Local: `make verify` exit 0, 3,136 tests in 24 packages;
  `go test ./internal/cli -count=1 -run 'TestRunImplementDetachSurvivesCallerProcessGroupKill'`
  PASS in 0.80 s.
