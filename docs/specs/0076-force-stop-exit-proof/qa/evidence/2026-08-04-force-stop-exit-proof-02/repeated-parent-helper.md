# Repeated parent/helper proof

- Build: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`

Force-kill proof, 50 repetitions:

```text
$ rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=50 -run '^TestOwnerProcessControllerForceKillExitProof$'
ok  roundfix/internal/store  1.900s
exit: 0
```

Graceful and force-kill proofs together, 50 repetitions each:

```text
$ rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=50 -run '^TestOwnerProcessController(GracefulExitProof|ForceKillExitProof)$'
ok  roundfix/internal/store  1.969s
exit: 0
```

Each harness captured the test command's status before checking output and
exited nonzero on `FAIL`, `fatal error`, or `file already closed`. Both harnesses
exited 0. No pipeline masked either `go test` status.
