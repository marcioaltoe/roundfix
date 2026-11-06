---
task: task_03
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
status: pending
type: backend
complexity: medium
---

# Task 03: A full filesystem parks the item as a Disk-Full Park

## Overview

When the Item Binary build or the repository gate fails because a filesystem
is full, the Delivery Queue parks the item as `disk-full: <path>: <N> MiB
free` in the environment Park Class, with a next action that names what to
free, instead of a generic `delivery-error` or a `gate-failed` finding, and a
Delivery Retry from it is not counted against the retry limit. Every other
failure keeps its blocker. The `deliver` documents, the documentation
contract, `CONTEXT.md` and ADR-0334 describe the Disk-Full Park.

## Requirements

1. MUST add `internal/delivery/disk_full.go` with `BlockerDiskFull =
   "disk-full"` and `DiskFullError` (`Path`, `AvailableBytes`, `Err`;
   `Error` returns `Err`'s text, `Unwrap` returns `Err`, `Blocker` returns
   `disk-full: <Path>: <AvailableBytes / 1 MiB> MiB free`), as
   `_techspec.md` → Interfaces says.
2. MUST classify `disk-full` in `ClassifyPark` (`internal/delivery/park_class.go`)
   as the environment Park Class with the next action of API Contract 7 word
   for word:
   `free space on the filesystem that holds the path the blocker names: run roundfix gc to trim declared Shared Build Caches and reclaim ended Runs, and remove other files there; then run roundfix deliver retry <slug>; a retry from disk-full is not counted`.
3. MUST park a step error that carries a `DiskFullError` with its
   `Blocker()` in both places that park `delivery-error: <error>` today: the
   `Run` loop in `internal/delivery/engine.go` and `reviewSlots.review` in
   `internal/delivery/review_slots.go`, through one helper; any other error
   keeps `delivery-error: <error>` byte for byte.
4. MUST exempt `disk-full` from the retry limit in `Retry`
   (`internal/delivery/engine.go`) and from the increment in
   `RetryDeliveryQueueItem` (`internal/store/delivery.go`), beside
   `runtime-infrastructure`; the resumed stage keeps following the item's
   evidence.
5. MUST add `internal/cli/deliver_disk_full.go` with
   `diskFullFailure(ctx, loaded, err, logPath)` as API Contract 6 describes:
   `errors.Is(err, syscall.ENOSPC)` or `no space left on device` (any case)
   in the last 64 KiB of `logPath`; the path is the last absolute path
   written directly before `: no space left on device` in the log tail, else
   in the error text, else task_02's `inspectFreeSpace` location with the
   least free space; the free space read at the path's nearest existing
   directory through the `availableBytes` command dependency.
6. MUST make `stepExecutable` return `build item binary: <the same error
   text>` wrapping the `DiskFullError` when `diskFullFailure` finds one for
   the build's `item-binary-build.log`, and `Gate` return
   `GateResult{Passed: false, Reason: <error>, Blocker: <Blocker()>}` for a
   Verification failure or an `ENOSPC` error whose `repository-gate.log` or
   error names a full filesystem, as API Contract 6 says; every other build
   and gate result keeps its bytes. Keep the edits in
   `internal/cli/deliver_workflow.go` local to those two functions.
7. MUST add `TestAFullDiskParksTheItemAsDiskFull` and
   `TestADiskFullRetryIsNotCounted` in `internal/delivery/disk_full_park_test.go`,
   `TestTheStoreDoesNotCountADiskFullRetry` in
   `internal/store/delivery_disk_full_retry_test.go`, and
   `TestTheItemBuildOnAFullDiskIsADiskFullFailure` and
   `TestTheRepositoryGateOnAFullDiskParksDiskFull` in
   `internal/cli/deliver_disk_full_test.go`, with the rows `_techspec.md` →
   Testing Approach names, each package's rule on `t.Parallel()`, fake build
   and Verification commands that print the text, and the injected
   `availableBytes`; no test fills or reads a real disk. `TestEveryBlockerHasAParkClass`
   keeps passing with the new constant.
8. MUST write, word for word as `_techspec.md` → API Contract 8 gives them,
   in `docs/user-guide/commands/deliver.md` and the canonical
   `.agents/skills/roundfix/references/deliver.md`: `disk-full` in the
   `environment` row of the Park Class table, the paragraph starting "When
   the item binary build or the repository gate fails because a filesystem is
   full", the `disk-full` next action sentence exactly once, and the
   uncounted retry sentence; add `delivery.BlockerDiskFull` to the blockers
   whose documented next action `checkParkActionDocumentation`
   (`internal/docscontract/command_documentation_test.go`) requires exactly
   once in each `deliver` document; add through `domain-modeling` the
   `CONTEXT.md` entry **Disk-Full Park** with its text and `_Avoid_` line,
   beside **Free-Space Floor**, with `ADR-0370` written as the glossary's
   Markdown link, adding no other capitalized phrase; and write the ADR-0334
   amendment notice after its Decision's last bullet, with `ADR-0370` in the
   bold opening as a Markdown link to
   `0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`,
   and `updated_at: 2026-11-06T12:00:00Z`; no other line of ADR-0334 changes.
9. MUST run `make skills-sync`, then raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide and skill edit; never
   hand-edit a mirror, a version field or a record. Of the Governed Paths the
   Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md` and
   `skills/roundfix/SKILL.md`.
10. SHOULD add the `disk-full` row to `TestParkClassesAndNextActions` with
   its exact next action, and update every test that pins a blocker this Task
   changes on the default branch when the Task runs.

## Subtasks

- [ ] Add `DiskFullError`, the blocker and its Park Class.
- [ ] Park a `DiskFullError` at both `delivery-error` sites.
- [ ] Exempt `disk-full` from the retry limit and count.
- [ ] Detect a full filesystem in the item build and the gate.
- [ ] Add the delivery, store and CLI tests.
- [ ] Document the park, extend the documentation contract, add the glossary entry and the ADR-0334 notice, sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A build or gate log ending in `write /tmp/x/_pkg_.a: no space left on device` parks `disk-full: /tmp/x/_pkg_.a: <N> MiB free`, environment class, with the next action above.
- [ ] A log with no absolute path names the inspected location with the least free space.
- [ ] Any other failure keeps `delivery-error` or `gate-failed` byte for byte.
- [ ] A retry from `disk-full` leaves the retry count unchanged and is accepted at the limit.
- [ ] The `deliver` documents state the park, its next action once and its uncounted retry; `CONTEXT.md` and ADR-0334 describe it.

## Context

- instruction: `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/delivery/runtime_infrastructure_test.go`
- instruction: `internal/cli/deliver_item_binary_test.go`
- instruction: `internal/runcause/signatures.json`
- creates: `internal/delivery/disk_full.go`
- creates: `internal/delivery/disk_full_park_test.go`
- creates: `internal/store/delivery_disk_full_retry_test.go`
- creates: `internal/cli/deliver_disk_full.go`
- creates: `internal/cli/deliver_disk_full_test.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/review_slots.go`
- interface: `internal/store/delivery.go`
- interface: `internal/cli/deliver_workflow.go`
- instruction: `internal/cli/doctor_disk_space.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `CONTEXT.md`
- interface: `docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `disk-full` blocker and its Park Class — consumers: `./internal/delivery`, `./internal/cli`
- The uncounted retry — consumers: `./internal/delivery`, `./internal/store`
- The Item Binary build and gate errors — consumers: `./internal/cli`
- The `deliver` documents' Park Class table and actions, the Roundfix skill version and the glossary — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAFullDiskParksTheItemAsDiskFull|TestADiskFullRetryIsNotCounted|TestEveryBlockerHasAParkClass)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAFullDiskParksTheItemAsDiskFull TestADiskFullRetryIsNotCounted TestEveryBlockerHasAParkClass; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^TestTheStoreDoesNotCountADiskFullRetry$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheStoreDoesNotCountADiskFullRetry ' || { echo 'missing pass: TestTheStoreDoesNotCountADiskFullRetry' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/delivery ./internal/store && GOOS=windows go build -buildvcs=false ./internal/delivery ./internal/store` — expected: exit 0. Before this Task none of the new tests exists, so the first pass line is missing; afterwards a `DiskFullError` parks `disk-full` in the environment Park Class, its retry is not counted by the engine or the store, and both suites pass.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheItemBuildOnAFullDiskIsADiskFullFailure|TestTheRepositoryGateOnAFullDiskParksDiskFull|TestDeliveryStepParksWhenTheItemBuildFails)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheItemBuildOnAFullDiskIsADiskFullFailure TestTheRepositoryGateOnAFullDiskParksDiskFull 'TestTheItemBuildOnAFullDiskIsADiskFullFailure/path_in_the_log' 'TestTheItemBuildOnAFullDiskIsADiskFullFailure/no_path_in_the_log' TestDeliveryStepParksWhenTheItemBuildFails; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/cli` — expected: exit 0. Before this Task the CLI tests do not exist; afterwards the item build and the gate name a full filesystem, another build failure keeps its error, and the CLI suite passes and builds for Windows.

- `grep -q 'delivery.BlockerDiskFull' internal/docscontract/command_documentation_test.go || { echo 'park contract misses disk-full' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- 'parks as .disk-full: <path>: <N> MiB free. in the environment Park Class' || { echo "missing Disk-Full Park in $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -q -- 'A Delivery Retry from .disk-full. is not counted against the retry limit either.' || { echo "missing uncounted retry in $f" >&2; exit 1; }; grep -q '^| .environment. | .*.disk-full.' "$f" || { echo "no disk-full in the environment row: $f" >&2; exit 1; }; done; grep -F -- '**Disk-Full Park**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Disk-Full Park\*\*(:| —( |$))' || { echo 'missing term: Disk-Full Park' >&2; exit 1; }; f=docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md; grep -q '^updated_at: 2026-11-06T12:00:00Z$' "$f" && tr -s '[:space:]' ' ' < "$f" | grep -qF -- '**Amended by [ADR-0370](0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md) (2026-11-06).**' || { echo "missing amendment: $f" >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the contract does not name the blocker, so the first check fails; afterwards the three `deliver` documents state the park and its uncounted retry, the glossary and ADR-0334 describe it, and the documentation contracts (the park actions and the Behavior Surface fingerprint among them), the mirrors and the recorded skill version pass.

## References

- `_prd.md` → Goal 5; Core Features 4 and 5; Success Metric 5; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 5 to 7; Interfaces; API Contracts 6 and 7; API Contract 8; Invariant 6; Declared breaks; Testing Approach; Build Order 3
- ADR-0370; ADR-0334
