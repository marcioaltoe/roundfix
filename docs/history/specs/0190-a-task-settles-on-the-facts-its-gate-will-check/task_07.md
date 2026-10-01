---
task: task_07
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: completed
type: backend
complexity: low
---

# Task 07: The public event stream names the two Settlement Checks

## Overview

The QA gate of Run `run_20261001T024730Z_a33946bbe6ed3a5d` found F-01. `roundfix events <run> --filter verification` prints the Settlement Check events without their `command`, although the Run Event Journal holds both labels. The cause is the public projection in `internal/runevent/stream.go`, which keeps `command` only for an `unknown` classification. The stream redacts verification commands on purpose: `TestProjectStreamEventCoversStableCategoriesAndRedactsPayload` and `TestEventsReplayDefaultAndFilterJSONLRecordsOnly` require that a repository command such as `make verify` never appears. API Contract 2 asks only for the two labels Roundfix itself defines.

## Requirements

1. MUST make the verification projection in `internal/runevent/stream.go` set `command` when the payload's `command` is exactly `settlement check: spec consistency` or `settlement check: authorization`. Every other command MUST stay out of the public record, as today.
2. MUST add `TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot` in `internal/runevent/settlement_label_projection_test.go`. It covers `started`, `command-passed` and `failed` events for both labels, and a `make verify` event whose command stays absent.
3. MUST keep `TestProjectStreamEventCoversStableCategoriesAndRedactsPayload` and `TestEventsReplayDefaultAndFilterJSONLRecordsOnly` passing unchanged.
4. MUST NOT edit any other file.

## Subtasks

- [ ] Project the two labels.
- [ ] Prove the labels appear and repository commands stay redacted.

## Acceptance Criteria

- [ ] `roundfix events <run> --filter verification` shows both Settlement Check labels for a gated Task, and no repository command.

## Context

- interface: `internal/runevent/stream.go`
- creates: `internal/runevent/settlement_label_projection_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot|TestProjectStreamEventCoversStableCategoriesAndRedactsPayload)$' ./internal/runevent 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot TestProjectStreamEventCoversStableCategoriesAndRedactsPayload; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 -run '^TestEventsReplayDefaultAndFilterJSONLRecordsOnly$' ./internal/cli` — expected: exit 0; before this Task the new test does not exist.

## References

- task_02
- `_techspec.md` → API Contracts (API Contract 2)

## Result

The public verification projection now retains `command` for exactly
`settlement check: spec consistency` and `settlement check: authorization`
in ordinary verification events. Repository commands remain redacted, and
the existing vacuous and unknown classification projections are unchanged.

Acceptance evidence:

- `TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot` exercises
  both labels across `started`, `command-passed`, and `failed` through the
  verification filter and checks the serialized public `command` field.
  It also requires the field to be absent for `make verify`, a label with an
  appended shell command, and a label with leading whitespace in all three
  phases. Before the production edit, all six exact-label cases failed
  because the field was absent.
- `GOCACHE=/private/tmp/roundfix-task07-go-cache rtk proxy go test -count=1 -run 'TestASettlementCheck|TestProjectStreamEvent|TestEventsReplay' ./internal/runevent ./internal/cli`
  exited 0. This focused selection includes the new regression and both
  required existing redaction tests, which remain byte-identical.
- The initial focused check using the shared Go cache could not set up the
  CLI package because sandbox access to a cache entry was denied. The
  task-scoped cache rerun above passed.
- `GOCACHE=/private/tmp/roundfix-task07-go-cache rtk make verify-incremental`
  exited 0 on the rerun with process-table access and no concurrent repository
  edits. Formatting, vet, package tests, skill checks, and build passed.
  The first incremental attempt exited 2: CLI process-stop tests lacked
  process-table access, and the suite guard detected this Agent's concurrent
  Result edit. Both conditions were removed for the successful rerun.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The declared Verification commands were not run. No live Run was replayed;
the focused evidence covers the projection and existing CLI replay boundary.
Task status remains Daemon-owned. No other task or graph file was edited,
and no commit, push, or pull request was made.
