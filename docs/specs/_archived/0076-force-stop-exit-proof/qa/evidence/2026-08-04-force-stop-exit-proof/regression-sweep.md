# Regression sweep

The package sweep exited 0:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1
ok  roundfix/internal/store  0.347s
```

The immediately following full parallel sweep exited 1:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test -parallel 16 ./...
--- FAIL: TestRunImplementDetachSurvivesCallerProcessGroupKill (6.21s)
    implement_test.go:1510: timed out waiting for line
FAIL
FAIL  roundfix/internal/cli  69.805s
FAIL
```

Every other listed package passed, including `internal/store` in 7.479s. The
raw `go test` exit status was 1; no pipeline masked it.

## Isolation

The failing test is in `internal/cli/implement_test.go`, which Spec 0076 did
not change. Its first detach stdout line is bounded by a fixed five-second
`readLineWithTimeout` call.

The focused test passed 20 of 20 isolated invocations, each in 0.74–0.97s:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/cli -count=20 -parallel=16 -run '^TestRunImplementDetachSurvivesCallerProcessGroupKill$' -v
PASS
ok  roundfix/internal/cli  16.933s
```

The complete `internal/cli` package also passed when isolated:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/cli -count=1 -parallel=16
ok  roundfix/internal/cli  19.790s
```

The earlier `rtk make verify` invocation, which includes the same full parallel
sweep, passed. The observed failure is therefore nondeterministic and narrows
to full-suite load: the failing invocation took 6.21s to reach a line guarded
by a fixed 5s deadline while isolated invocations stayed below 1s. The exact
root cause inside detach startup is not established in this QA run.
