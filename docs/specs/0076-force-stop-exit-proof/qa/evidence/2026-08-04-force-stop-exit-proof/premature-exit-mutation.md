# Premature-exit mutation probe

A local disposable clone was created at
`/tmp/rf-0076-mutation.KGWv4W/repo` from the exact report build
`c035ebb19dcb6eb81844f5195a0b89abbf99e4e1`. Only that disposable copy was
mutated: one `return` was inserted immediately after the `ignore` helper prints
`ready`.

The focused public test command then exited 1 as required:

```text
rtk proxy env GOCACHE=/tmp/rf-0076-mutation.KGWv4W/gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessControllerForceKillExitProof$' -v
=== RUN   TestOwnerProcessControllerForceKillExitProof
=== PAUSE TestOwnerProcessControllerForceKillExitProof
=== CONT  TestOwnerProcessControllerForceKillExitProof
    process_unix_test.go:65: owner process 82808 exited prematurely before controller force-kill escalation: <nil>
--- FAIL: TestOwnerProcessControllerForceKillExitProof (0.01s)
FAIL
FAIL  roundfix/internal/store  0.330s
FAIL
```

The failure names the premature exit and therefore distinguishes a clean
early exit from controller-caused `SIGKILL`.

The unchanged Run Worktree then ran the identical focused command and exited
0:

```text
--- PASS: TestOwnerProcessControllerForceKillExitProof (0.03s)
PASS
ok  roundfix/internal/store  0.238s
```

`git -c core.fsmonitor=false status --short
internal/store/process_unix_test.go` produced no output in the Run Worktree,
confirming the mutation never touched the report build.
