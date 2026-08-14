# Direct helper liveness

- Build: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`
- Binary build command: `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test -c -o /tmp/rf-0076-store-qa02.test ./internal/store`
- Binary build exit: `0`

The helper binary ran alone with
`ROUNDFIX_OWNER_PROCESS_HELPER=ignore` and
`-test.run='^TestOwnerProcessHelper$'`. The harness consumed its public stdout
protocol and sent Unix signals to the helper PID:

```text
readiness: ready
sigterm: sent
liveness: alive
after-sigterm: alive
cleanup-exit: 137
stderr: empty
```

The harness exited 0. Exit 137 followed the harness's explicit `SIGKILL` cleanup.
The process remained alive after `SIGTERM`, acknowledged `SIGUSR1`, and emitted
neither `fatal error` nor `all goroutines are asleep`.

Independent confirmation:

```text
$ rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessHelperIgnoreModeStaysAlive$' -v
=== RUN   TestOwnerProcessHelperIgnoreModeStaysAlive
--- PASS: TestOwnerProcessHelperIgnoreModeStaysAlive (0.01s)
PASS
ok  roundfix/internal/store  0.317s
```

An initial harness attempt launched the helper behind an `rtk proxy` child and
therefore signalled the wrapper PID; the wrapper exited before liveness could
be read. Running the already RTK-prefixed shell with the test binary as the
direct coprocess made the PID boundary truthful. The helper behavior above is
the resulting product evidence. zsh also reported a sandbox-denied background
`nice(5)` adjustment, but the direct helper journey completed and its exit was
observed, so no acceptance row was blocked.
