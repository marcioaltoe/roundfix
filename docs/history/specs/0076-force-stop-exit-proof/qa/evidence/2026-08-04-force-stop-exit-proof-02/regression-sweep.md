# Prior F-01 retest and regression sweep

- Build: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`
- Prior finding: F-01 from `qa-report-2026-08-04.md`
- Repair commit present: `c6d49d0a08b7ab8c3b2ae8b359e612e3ce975d8d`

## Store package

Command:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1
```

Exit: `0`

```text
ok  roundfix/internal/store  0.387s
```

## Full parallel sweep

Command:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test -parallel 16 ./...
```

Exit: `0`. The slowest packages were `internal/cli` at 71.285s and
`internal/baseline` at 70.270s. `internal/store` passed in 4.618s. No package
failed, and the prior fixed 5s detached-Run readiness diagnostic did not recur.

Independent confirmation: `rtk make verify` separately ran the same full
parallel sweep and reported 3,137 passing tests across 24 packages.

Scope confirmation: Task commits `d9be9f5` and `c035ebb` changed only their
assigned Task files and `internal/store/process_unix_test.go`. Repair commit
`c6d49d0` changed only `internal/cli/implement_test.go`. All are test or Spec
artifact paths; no production file under `internal/store` changed.
