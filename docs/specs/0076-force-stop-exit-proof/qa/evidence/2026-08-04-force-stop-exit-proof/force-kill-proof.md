# Force-kill causation proof

The focused proof was exercised through the public Go test runner:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessControllerForceKillExitProof$' -v
=== RUN   TestOwnerProcessControllerForceKillExitProof
=== PAUSE TestOwnerProcessControllerForceKillExitProof
=== CONT  TestOwnerProcessControllerForceKillExitProof
--- PASS: TestOwnerProcessControllerForceKillExitProof (0.04s)
PASS
ok  roundfix/internal/store  0.334s
```

The identical focused proof under the race detector also exited 0:

```text
--- PASS: TestOwnerProcessControllerForceKillExitProof (0.05s)
PASS
ok  roundfix/internal/store  1.463s
```

Source inspection ties the public result to its intended observable:
`assertOwnerProcessForceKilled` accepts only an `exec.ExitError` whose
`syscall.WaitStatus` is signaled by `SIGKILL`; a clean exit or another signal
fails. The standalone helper evidence independently proves that the helper
survives `SIGTERM` before this assertion is credited.
