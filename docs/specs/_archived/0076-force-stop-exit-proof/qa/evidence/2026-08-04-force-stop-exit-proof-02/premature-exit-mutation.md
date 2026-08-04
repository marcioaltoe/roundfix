# Premature-exit mutation sensitivity

- Build copied: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`
- Disposable copy: `/private/tmp/rf-0076-qa02.qNS6Q9`
- Mutation: add `return` immediately after the `ignore` helper emits `ready`
- Run Worktree source mutation: none

The disposable copy ran the public focused proof through the normal Go test
runner. The QA harness required both a nonzero test exit and the named
premature-exit diagnostic:

```text
=== RUN   TestOwnerProcessControllerForceKillExitProof
=== PAUSE TestOwnerProcessControllerForceKillExitProof
=== CONT  TestOwnerProcessControllerForceKillExitProof
    process_unix_test.go:65: owner process 21193 exited prematurely before controller force-kill escalation: <nil>
--- FAIL: TestOwnerProcessControllerForceKillExitProof (0.01s)
FAIL
FAIL  roundfix/internal/store  0.333s
FAIL
observed-test-exit: 1
```

The evidence harness exited 0 only because it observed the required failing
test and diagnostic. The unchanged Run Worktree then ran the identical focused
test command:

```text
--- PASS: TestOwnerProcessControllerForceKillExitProof (0.03s)
PASS
ok  roundfix/internal/store  0.288s
```

The unchanged command exited 0. This proves the regression is mutation-sensitive
and the positive build remains green.
