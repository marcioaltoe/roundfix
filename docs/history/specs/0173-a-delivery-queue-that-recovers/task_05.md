---
task: task_05
spec: 0173-a-delivery-queue-that-recovers
status: completed
type: backend
complexity: high
---

# Task 05: `roundfix deliver retry` returns a parked item to its owner

## Overview

`runDeliverCommand` in `internal/cli/deliver.go` offers `start`, `status`, `resume` and `stop`; none returns a parked item to the queue, so the operator leaves `roundfix deliver` and finishes the Spec by hand. The owner started by `deliver resume` makes one pass and releases the queue through a deferred `ReleaseDeliveryQueueOwner`, so an item retried while it runs would be stranded. This Task adds the `retry` command on top of task_03's `Engine.Retry` and task_04's recovery, and makes the owner release the queue only when it is idle. The command reads and writes the Run Database shared with a detached owner process, proves that process through the owner process controller, and reports to the operator on stdout and stderr.

## Requirements

1. MUST add `retry` to `deliverUsage` (the usage line `roundfix deliver retry <slug>` and a Commands line) and to `runDeliverCommand`, and add `Retry(context.Context, string, string) (delivery.RetryResult, error)` to the `deliveryEngine` interface.
2. MUST make `runDeliverRetry` accept exactly one Spec slug through the deliver flag set, refusing a missing slug, an empty slug, an extra argument and an unknown flag with exit `2`.
3. MUST make `runDeliverRetry` build the engine with `newDeliveryEngine`, call `Retry`, and print to stdout `Carried forward from Run <run-id>: <task>, <task>` when Tasks were carried, then `Retried <slug>: <blocker> -> <stage>`.
4. MUST make the hand-off follow the owner the retry read: an owner that is alive and proven through `ownerProcesses.ProveOwner` gets the item and the command prints `Handed <slug> to Delivery Queue owner PID <pid>.` and exits `0` without starting a process; an owner that is dead or whose identity is unproven is released with `ReleaseDeliveryQueueOwner`, reported with the stderr notice `deliver resume` writes, and replaced by a started owner; with no owner, the command starts one through `startDeliveryOwner`.
5. MUST make every refusal and failure go through `printDeliverFailure` and exit `2`, starting no owner.
6. MUST make `runDeliveryOwner` run the engine, then call `ReleaseIdleDeliveryQueueOwner`, exit when it releases the queue, and run the engine again when it returns `false`; a release error fails the owner, and the deferred release stays for the error paths.
7. MUST change the top-level usage line in `internal/cli/cli.go` to `roundfix deliver <start|status|resume|retry|stop> [<slug> ...]`, and add `roundfix deliver retry <slug>` to the `deliver` row of `TestRunCommandHelp` in `internal/cli/cli_test.go`, adding, renaming or removing no top-level test there.
8. MUST document in the deliver section of `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md` the command (the string `roundfix deliver retry <slug>`), the re-entry stages, the carry-forward step, the owner hand-off, the refusals and the exit codes, and regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
9. MUST add to `CONTEXT.md` a **Delivery Queue** entry and a **Delivery Retry** entry, and make its Task Carry-Forward entry say that the act is reached through the Reconcile Command's `--carry-forward` switch or through a Delivery Retry and is never automatic.
10. MUST put the new tests in `internal/cli/deliver_retry_test.go`, driving the public CLI with the real Run Database, and injecting only the engine, the owner process controller and the owner starter as operating-system boundaries.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A retry with no owner starts one and prints the retried stage; a retry with a live proven owner hands the item over and starts nothing; a retry with a stale owner record reclaims it and starts one.
- [ ] A missing slug, an extra argument and an engine refusal exit `2` and start no owner.
- [ ] The owner runs the engine again for an item retried during its pass and releases the queue only after a pass that leaves every item merged or parked.
- [ ] `roundfix deliver --help` and `roundfix --help` name the retry command, and the guides and glossary describe it.

## Context

- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- creates: `internal/cli/deliver_retry_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliverRetryStartsAnOwnerWhenNoneIsRunning|TestDeliverRetryHandsTheItemToALiveOwner|TestDeliverRetryReclaimsAStaleOwnerAndStartsOne|TestDeliverRetryRequiresOneSlug|TestDeliverRetryRefusesAnExtraArgument|TestDeliverRetryRefusalStartsNoOwner|TestDeliveryOwnerRunsAgainForAnItemRetriedDuringItsPass|TestDeliveryOwnerReleasesAnIdleQueueAfterOnePass|TestTopLevelUsageNamesDeliverRetry|TestRunCommandHelp|TestDeliverCommandStopsAndResumesPersistedQueue|TestResumeReleasesAStaleOwner|TestDeliverCommandRefusesUnknownFlags)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverRetryStartsAnOwnerWhenNoneIsRunning TestDeliverRetryHandsTheItemToALiveOwner TestDeliverRetryReclaimsAStaleOwnerAndStartsOne TestDeliverRetryRequiresOneSlug TestDeliverRetryRefusesAnExtraArgument TestDeliverRetryRefusalStartsNoOwner TestDeliveryOwnerRunsAgainForAnItemRetriedDuringItsPass TestDeliveryOwnerReleasesAnIdleQueueAfterOnePass TestTopLevelUsageNamesDeliverRetry TestRunCommandHelp TestDeliverCommandStopsAndResumesPersistedQueue TestResumeReleasesAStaleOwner TestDeliverCommandRefusesUnknownFlags; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < internal/cli/cli_test.go | grep -qF -- "roundfix deliver retry <slug>" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "roundfix deliver retry <slug>" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "roundfix deliver retry <slug>" && grep -q "Delivery Retry" CONTEXT.md && grep -q "Delivery Queue" CONTEXT.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the nine new named tests exists and the command is documented nowhere, so the command fails.

## References

- [_techspec.md](_techspec.md) — The retry command and the owner hand-off

## Result

Implementation:

- Added `deliver retry` parsing, dispatch, engine invocation, carry-forward and
  retry output, and owner hand-off or replacement through the owner recorded by
  `Engine.Retry`.
- Changed the detached owner to release its claim through
  `ReleaseIdleDeliveryQueueOwner` after each pass and repeat while an item can
  still advance; the deferred release remains on failure paths.
- Added public-CLI tests backed by the real Run Database for owner start,
  proven-owner hand-off, stale-owner reclamation, each argument refusal,
  engine refusal, owner repetition, idle release, release failure, deferred
  cleanup, and help output.
- Documented Delivery Retry in the user guide, canonical and shipped Roundfix
  skills, and the domain glossary. `make skills-sync` regenerated the shipped
  skill; `make baseline-digests` reported `changed:false`.

Focused checks:

- Red signal: the initial focused test run failed because `retry` was an
  unknown deliver command, the owner ran once, and top-level help omitted
  `retry`.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run
  '^(TestDeliverRetry|TestDeliveryOwner|TestTopLevelUsageNamesDeliverRetry)'
  ./internal/cli`: passed.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -race -count=1 -run
  '^(TestDeliverRetry|TestDeliveryOwner)' ./internal/cli`: passed.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 ./internal/cli`:
  passed with process-table permission. The first sandboxed run reached two
  pre-existing force-stop integration tests and was blocked by `operation not
  permitted`; the permitted rerun passed.
- `GOCACHE=/tmp/roundfix-task05-gocache go vet ./internal/cli`: passed.
- `make fmt-check`: passed.
- `make skills-sync-check`: passed.
- Phrase inspection found `roundfix deliver retry <slug>` in the CLI help
  expectation, user guide, canonical skill and shipped skill, and found both
  Delivery Queue glossary entries.

Acceptance evidence:

- `TestDeliverRetryStartsAnOwnerWhenNoneIsRunning`,
  `TestDeliverRetryHandsTheItemToALiveOwner`, and
  `TestDeliverRetryReclaimsAStaleOwnerAndStartsOne` cover the three owner
  outcomes and their stdout or stderr reports.
- `TestDeliverRetryRequiresOneSlug`, `TestDeliverRetryRefusesAnEmptySlug`,
  `TestDeliverRetryRefusesAnExtraArgument`,
  `TestDeliverRetryRefusesUnknownFlag`, and
  `TestDeliverRetryRefusalStartsNoOwner` cover exit `2` without an owner start.
- `TestDeliveryOwnerRunsAgainForAnItemRetriedDuringItsPass`,
  `TestDeliveryOwnerReleasesAnIdleQueueAfterOnePass`,
  `TestDeliveryOwnerFailsWhenIdleReleaseFails`, and
  `TestDeliveryOwnerReleasesItsClaimAfterEngineFailure` cover idle-only release,
  repeated passes and failure cleanup.
- `TestTopLevelUsageNamesDeliverRetry` and the existing `TestRunCommandHelp`
  cover both help surfaces; the phrase and skill-parity checks cover the guides
  and glossary.

Daemon Verification was not run; it remains Daemon-owned.
