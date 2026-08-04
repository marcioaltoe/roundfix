# Non-Goal and scope audit

The complete implementation diff from planning commit `9dadc83` to report
build `c035ebb` changes one `internal/store` path:

```text
internal/store/process_unix_test.go
```

No production file under `internal/store` changed. In particular,
`TerminateAndWait` remains in `internal/store/process.go`, outside the diff.

The `ignore` helper block uses this event-driven shape:

```go
signal.Ignore(syscall.SIGTERM)
defer signal.Reset(syscall.SIGTERM)
signals := make(chan os.Signal, 1)
signal.Notify(signals, syscall.SIGUSR1)
defer signal.Stop(signals)
fmt.Fprintln(os.Stdout, "ready")
for range signals {
    fmt.Fprintln(os.Stdout, "alive")
}
```

It contains no empty `select`, sleep, timer, or deadline. Diff inspection shows
the remaining changes only reorder readiness before `cmd.Wait`, add direct
test-process observations, and assert `SIGKILL` causation. Force Stop product
behavior, owner-identity rules, the last-resort flag, unrelated production
symbols, and protected tooling remain unchanged.
