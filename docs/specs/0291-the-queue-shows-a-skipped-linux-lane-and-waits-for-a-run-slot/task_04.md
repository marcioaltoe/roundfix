---
task: task_04
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
status: pending
type: backend
complexity: high
---

# Task 04: The queue recognizes the ceiling refusal and shows the Run Slot Wait

## Overview

The queue turns every preflight refusal of its child `roundfix implement`
into a `delivery-error` park, the ceiling included (`_prd.md` → Sources 2).
This Task makes the queue's Implement step recognize a refusal whose only
reason is the Active Implement Run ceiling, read the Run slots the way
`implement` counts them, and hand the Queue Owner's engine the configured
`runs.slot_wait`, so the Run Slot Wait task_03 built runs in a real queue and
`roundfix deliver status` shows it. It also documents the wait for users and
in the Roundfix reference. It is verifiable on its own through compiled item
binaries, a real Run Database and the status command.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md`](references/2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md)
   (`_prd.md` → Sources 2), by implementing API Contracts 7 and 8 of
   `_techspec.md` in `internal/cli/deliver_workflow.go` and
   `internal/store/store.go`, with `ActiveImplementRunCeilingReason` and
   `commandDeliveryWorkflow.RunSlots` as `_techspec.md` → Interfaces sketches
   them.
2. MUST build `ActiveImplementRunCeilingError.Error()`'s first line from
   `ActiveImplementRunCeilingReason`, keeping its text byte for byte, and
   recognize the refusal only when all three hold: exit code `2`, no new
   Implement Run of the Spec, and the first non-empty line after the `Reason:`
   line of standard error, with ANSI color sequences removed, starting with
   that constant and a colon.
3. MUST state the hazard in the code's behavior, as API Contracts 7 and 8 do.
   The data sources are the child's exit code and standard error, the Run
   Database's Active Implement Runs and the Queue Owner's User Config; the
   destinations are the refusal the engine waits on and, through it, the
   item's blocker text that `roundfix deliver status` prints. `RunSlots`
   reads only: it reclaims no Run and signals no process, and it excludes a
   holder only when its recorded owner PID names a process that is not alive.
   A failed slot read returns the plain failure, which parks as today.
4. MUST wire the engine in `newCommandDeliveryEngine`: the workflow as
   `RunSlots` and `loaded.Config.Runs.SlotWait` as `RunSlotWait`. Nothing
   else in `internal/cli/deliver_workflow.go` changes beyond `Gate`
   (task_01), `deliveryRunResult`, `RunSpec`, `RunSlots` and that function.
5. MUST keep every existing test's name and assertions, including
   `TestRunCeilingRefusalNamesTheWayOut`, `TestImplementRunCeilingZeroDisables`,
   `TestRunCeilingIsCountedInsideTheCreateTransaction` and the item binary
   tests in `internal/cli/deliver_item_binary_test.go`; this Task declares no
   break.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   ending every child it starts in `t.Cleanup`, using temporary homes and
   repositories, and declaring no real-time budget under ten seconds:
   - in `internal/cli/deliver_run_slot_wait_test.go`:
     - `TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal`, rendering each
       standard error with the real `printPreflightFailure`, with the cases
       `_techspec.md` → Testing Approach lists;
     - `TestRunSpecReturnsARunSlotRefusalWithItsHolders`, through `RunSpec`
       with its own compiled item binary, built and copied the way
       `newDeliveryItemFixture` does, that answers `migrate --check` with exit
       `0` and `implement` with a ceiling refusal on standard error and exit
       `2`, and a Run Database holding one Active Implement Run of another
       repository: the error is a `*delivery.RunSlotRefusal` whose `Err` text
       equals today's failure and whose `Slots` name that holder and the
       configured ceiling;
     - `TestDeliverStatusShowsARunSlotWaitTranscript`, asserting Surface
       Transcript 2 of `_techspec.md` whole, with the blocker seeded from
       `delivery.RunSlotWaitText`;
   - in `internal/cli/deliver_run_slot_wait_unix_test.go`, under `//go:build
     unix`, `TestRunSlotsExcludeAHolderWhoseOwnerIsGone`, with one holder
     owned by the test process and one whose owner PID is a `sh` child the
     test started and waited for.

   The files compile for Windows, which `GOOS=windows go vet ./internal/cli`
   proves.
7. MUST document the wait:
   - in `docs/user-guide/commands/deliver.md`, a Run slot wait section: when
     the item waits, the blocker text `run-slot-wait: until <deadline>; <n>
     Run(s) hold <max> slot(s): <run id> in <repository> for <spec>` on the
     item's status line with the stage `running`, the two console lines, the
     poll, `runs.slot_wait`, the queue deadline that wins, and the park at the
     end of the wait with today's reason;
   - in `.agents/skills/roundfix/references/deliver.md`, the same behavior.

   MUST search the guides, the agent guides and the skill copies for wording
   that says the queue parks at the Active Implement Run ceiling, change each
   one in the files this Task declares, and name any other file found in the
   Task Result. No text inside any `### QA settlement` section may change.
8. MUST add **Run Slot Wait** to `CONTEXT.md` through the `domain-modeling`
   skill: the bounded wait of a running Delivery Queue item whose Implement
   step was refused only by the Active Implement Run ceiling, shown as its
   blocker text, polled against the Run Database, ended by `runs.slot_wait`
   or the queue deadline, linking ADR-0304.
9. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. Before them, raise both version fields of the Roundfix
   skill's canonical `SKILL.md` (`metadata.version` and the top-level
   `version`) by one patch over the value the tree holds when this Task
   starts. The Agent runs these commands itself before handing back; no
   record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.

## Subtasks

- [ ] Export the ceiling message's prefix and build the error from it.
- [ ] Recognize the ceiling-only refusal and fill its Run slots.
- [ ] Answer `RunSlots` and wire the engine.
- [ ] Add the recognition, item binary, status and unix tests.
- [ ] Document the wait in the guide, the reference and the glossary; sync, record and re-record.

## Acceptance Criteria

- [ ] A child refused only by the ceiling returns a `*delivery.RunSlotRefusal` with the live holders; every other refusal keeps today's error.
- [ ] A holder whose owner process is gone does not count.
- [ ] `roundfix deliver status` reproduces Surface Transcript 2.
- [ ] The guide, the reference and the glossary describe the Run Slot Wait.

## Consumers

- `commandDeliveryWorkflow.RunSlots` and the refusal recognition — consumers: `./internal/cli`, `./internal/delivery`
- `store.ActiveImplementRunCeilingReason` — consumers: `./internal/store`, `./internal/cli`

## Context

- instruction: `docs/adr/0304-a-queue-item-waits-for-a-run-slot-up-to-runs-slot-wait.md`
- instruction: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- instruction: `internal/cli/implement.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/cli/deliver_item_binary_test.go`
- instruction: `internal/cli/deliver_park_status_test.go`
- instruction: `internal/store/store_test.go`
- instruction: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/store/store.go`
- creates: `internal/cli/deliver_run_slot_wait_test.go`
- creates: `internal/cli/deliver_run_slot_wait_unix_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal|TestRunSpecReturnsARunSlotRefusalWithItsHolders|TestDeliverStatusShowsARunSlotWaitTranscript|TestRunSlotsExcludeAHolderWhoseOwnerIsGone|TestRunCeilingRefusalNamesTheWayOut|TestImplementRunCeilingZeroDisables|TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments|TestDeliverStatusPrintsAParkLinePerParkedItem)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal TestRunSpecReturnsARunSlotRefusalWithItsHolders TestDeliverStatusShowsARunSlotWaitTranscript TestRunSlotsExcludeAHolderWhoseOwnerIsGone TestRunCeilingRefusalNamesTheWayOut TestDeliverStatusPrintsAParkLinePerParkedItem; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestRunCeilingIsCountedInsideTheCreateTransaction$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestRunCeilingIsCountedInsideTheCreateTransaction ' || { echo 'missing store ceiling pass' >&2; exit 1; }; grep -qF -- 'ActiveImplementRunCeilingReason' internal/store/store.go || { echo 'missing exported ceiling prefix' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'run-slot-wait: until <deadline>' || { printf 'missing wait text in %s\n' "$f" >&2; exit 1; }; grep -qF -- 'runs.slot_wait' "$f" || { printf 'missing runs.slot_wait in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Run Slot Wait**:' || { echo 'missing glossary term' >&2; exit 1; }; go test -count=1 ./internal/delivery || exit 1; GOOS=windows go vet ./internal/cli || exit 1; GOOS=windows go build -buildvcs=false ./internal/cli ./internal/store || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards recognition, the item binary refusal, the status transcript and the dead-holder exclusion pass with the existing ceiling, item binary and status tests and the store's ceiling race test, the delivery package passes, the store exports the prefix, the guide and the reference describe the wait and its key, the glossary defines the term, the CLI's tests vet and its packages build for Windows, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 4; Core Feature 5; Success Metric 4; Success Metric 5; Glossary
- `_techspec.md` → API Contract 6; API Contract 7; API Contract 8; API Contract 9; Surface Transcript 2; Invariants 4-5; Interfaces; Current behavior; Testing Approach; Build Order 4; Glossary
- ADR-0304; ADR-0279; ADR-0199; ADR-0184; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280
