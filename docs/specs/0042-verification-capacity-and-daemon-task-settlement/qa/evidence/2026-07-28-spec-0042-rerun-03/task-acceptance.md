# Task acceptance evidence

Current focused suites passed on build `1b1bfc3`:

```text
go test ./internal/config ./internal/agent ./internal/cli \
  -run 'Test(Load.*Verification|DefaultConfigYAML.*Verification|Validate.*Verification|Build(Task|VerificationRepair)Prompt|CommandUsage|DocumentationContract|Settle.*TaskStatus|RunImplement.*VerificationCapacit|TaskCycle.*Capacit)' \
  -count=1
```

```text
go test ./internal/daemon ./internal/runevent ./internal/tui \
  -run 'Test(TaskCycle.*(IntegratedVerificationCapacit|VerificationCapacit|WaitingForVerification|RepairReacquires|StopRequestWhileQueued|DaemonStatus|AgentStatus|AgentFailure|TemporaryVerification|ExclusiveRetry|RetryBudget|DeterministicRetry)|ExecVerifier.*(Temporary|Exit|Diagnostic|Cancel)|WaitingForVerificationReplay|VerificationTemporaryRetryExclusiveProjection|RenderLiveRunView.*(Capacit|Review)|Cockpit.*(Task|Review|Key|Detail|Viewport|NoColor|Capacit|WaitingForVerification))' \
  -count=1
```

These commands cover Task 01 Config defaults, hierarchy, strict validation,
generated YAML, planning, event compatibility, and absence of a flag; Task 02
status normalization, Result preservation, prompt ownership, Agent failure,
Stop, and Settle regression; Task 03 capacity 1/2, event ordering,
release/reacquire, exclusive fairness, and cancellation; and Task 04 typed
exit 75, distinct diagnostics, one retry, repair interaction, exhaustion, and
bounded public fields.

The full gate, race suite, public macro suite, and 20 repeated concurrency
runs provide the current assembled confirmation. Task 03's historical
`TestVerificationGate.*` selector remains stale after later test
consolidation, so this gate did not credit that old selector; current named
TaskCycle, public Stop, repeated, and race cases prove its criteria.
