# Force-kill exit proof

- Build: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`

Normal public test-runner path:

```text
$ rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessControllerForceKillExitProof$' -v
=== RUN   TestOwnerProcessControllerForceKillExitProof
=== PAUSE TestOwnerProcessControllerForceKillExitProof
=== CONT  TestOwnerProcessControllerForceKillExitProof
--- PASS: TestOwnerProcessControllerForceKillExitProof (0.03s)
PASS
ok  roundfix/internal/store  0.241s
```

Race-detector confirmation:

```text
$ rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test -race ./internal/store -count=1 -run '^TestOwnerProcessControllerForceKillExitProof$' -v
--- PASS: TestOwnerProcessControllerForceKillExitProof (0.04s)
PASS
ok  roundfix/internal/store  1.445s
```

Both commands exited 0. The focused proof calls
`assertOwnerProcessForceKilled`, whose acceptance condition requires a
signaled `syscall.WaitStatus` with `status.Signal() == syscall.SIGKILL`; a
clean or other-signal exit fails. The direct helper evidence independently
shows that `SIGTERM` did not end the process before this escalation.
