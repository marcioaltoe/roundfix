# Direct helper liveness

The helper test binary was built from the report build:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test -c -o /tmp/rf-0076-store.test ./internal/store
```

The direct harness started that binary alone with
`ROUNDFIX_OWNER_PROCESS_HELPER=ignore`, read its stdout, sent `SIGTERM`, sent
the registered `SIGUSR1` liveness probe, checked the process remained alive,
then sent `SIGKILL` and waited for exit. `rtk zsh
/tmp/rf-0076-direct-helper.sh` exited 0 with:

```text
readiness=ready
post_sigterm_liveness=alive
alive_after_sigterm=true
termination=SIGKILL
wait_status=137
fatal_runtime_diagnostic=false
```

macOS zsh also printed `nice(5) failed: operation not permitted` when starting
the background helper in the sandbox. The helper still started, produced both
protocol lines, survived `SIGTERM`, and was reaped after `SIGKILL`; this is an
environment warning, not a helper failure.

Independent confirmation:

```text
rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessHelperIgnoreModeStaysAlive$' -v
=== RUN   TestOwnerProcessHelperIgnoreModeStaysAlive
--- PASS: TestOwnerProcessHelperIgnoreModeStaysAlive (0.01s)
PASS
ok  roundfix/internal/store  0.305s
```

The independent command exited 0.
