# Repeated parent/helper rail

The exact force-kill parent/helper proof ran 50 times as one isolated command:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=50 -run '^TestOwnerProcessControllerForceKillExitProof$'
ok  roundfix/internal/store  1.943s
```

Exit status: 0. A child runtime fatal error, `file already closed`, or test
failure would make the raw `go test` command nonzero; no pipeline or pager
masked its status.

Independent confirmation ran both graceful and force-kill paths 50 times:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=50 -run '^TestOwnerProcessController(GracefulExitProof|ForceKillExitProof)$'
ok  roundfix/internal/store  2.096s
```

Exit status: 0.
