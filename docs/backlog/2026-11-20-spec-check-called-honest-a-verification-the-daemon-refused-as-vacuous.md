---
type: fix
status: open
created: 2026-11-20
spec: null
reason: null
---

# `roundfix spec check --run-verification` called honest a Verification command the Daemon then refused as vacuous

## Symptom

On 2026-11-19, the operator's agent added corrective `task_04` to Spec 0414 on its Item Branch (commit `a05ceae9`). Its `## Verification` held two commands; the second was:

```text
test -z "$(gofmt -l internal/cli/follow_up_queue_test.go)" && go vet ./internal/cli
```

That command exits 0 on the unchanged tree. `bin/roundfix spec check 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release --strict --run-verification`, built at the Item Branch head, reported `No findings. Authored Verification commands executed: 2`. The Delivery Queue then settled `task_04` failed with:

```text
Pre-work Verification vacuous: commands "test -z \"$(gofmt -l internal/cli/follow_up_queue_test.go)\" && go vet ./internal/cli" exited zero against the unchanged tree
```

and parked the item `run-unresolved`. Dropping the command (commit `54e8302d`, operator Correction C3) let the delivery finish (PR #930).

## Where

Both sides judge each command separately: the probe in `internal/cli/spec_check.go` (`specCheckVerificationVerdictVacuous`, about line 140) and the Daemon in `internal/daemon/task_engine.go:1704`. So the verdict differs by environment, not by rule. This is the class Spec 0396 addressed (the probe prepares its disposable checkout as the Daemon prepares a Task Worktree). Not measured: why `go vet ./internal/cli` or the `gofmt` subshell failed in the probe's checkout and passed in the Run Worktree. Candidates are the Go build cache or module cache seen by the probe, or the host's `core.fsmonitor=true` interacting with the disposable worktree.

## Expected

A command that exits zero against the unchanged tree in the Daemon's pre-work Verification is reported `vacuous` by `spec check --run-verification` in the same repository at the same head, or the probe names why it could not run the command (never `honest`).

## Source

Recorded by the operator, Cycle BG, from the 0414 delivery log and the probe output, on 2026-11-20.
