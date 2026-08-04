# Non-Goal and implementation-scope audit

- Build: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`
- Implementation range inspected: `9dadc83..c035ebb`

The implementation range changes one code file:
`internal/store/process_unix_test.go` (114 insertions, 6 deletions). Both
Task commits also change only their assigned Task file.

Focused diff inspection confirms:

- no production file under `internal/store` changed;
- `TerminateAndWait` and the process controller implementation are untouched;
- the existing force-kill test now calls a test-only causation assertion;
- the helper's `ignore` branch replaces `select {}` with `signal.Notify` on
  `SIGUSR1` and a channel range;
- the helper block introduces no sleep, timer, deadline, or retry;
- readiness scanning now completes before the test starts `cmd.Wait` in its
  observation goroutine;
- changes to other store tests are limited to the added direct liveness test
  required by the Spec.

The later F-01 repair changes only `internal/cli/implement_test.go`; it does not
change Force Stop product behavior or the store process controller.
