---
task: task_04
spec: 0156-a-delivery-loop-that-outlives-the-session
status: completed
type: backend
complexity: medium
---

# Task 04: The command family

## Overview

Expose the queue as `roundfix deliver start <slug>...`, `status`, `resume` and `stop`, starting the owner detached with the mechanism `implement --detach` already uses.

## Requirements

1. MUST record the queue and start a detached owner on `deliver start`.
2. MUST print each item's stage and blocker on `deliver status`.
3. MUST end the owner on `deliver stop` and restart it from the persisted queue on `deliver resume`.
4. MUST refuse unknown flags and unknown Spec slugs, matching the surrounding commands.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The command family appears on the public surface.
- [ ] `deliver status` prints each item's stage and blocker.
- [ ] Unknown slugs are refused before any queue is recorded.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/detach.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestDeliverCommandRecordsAQueueAndReportsIt$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDeliverCommandRecordsAQueueAndReportsIt"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `go run -buildvcs=false ./cmd/roundfix --help 2>&1 | grep -q "roundfix deliver"` — expected: exit 0; the command family appears on the public surface. Before this Task the usage has no deliver line, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components

## Result

Implemented the `deliver` command family and its detached owner. `deliver
start` validates every Spec before opening the Run Database, records the
ordered queue, and launches `deliver resume` through the existing detach
handshake. The owner records its PID and process-start identity on the queue,
runs the Delivery Engine through the existing Implement, review, archive,
repository-gate, authorization, and GitHub boundaries, and releases its owner
claim on exit. `deliver stop` proves and terminates that recorded process tree;
`deliver resume` preserves the queue, refuses a live owner, and reclaims only a
proved-absent owner record.

`deliver status` reads the durable queue and prints one tab-separated line per
item with Spec slug, stage, and blocker (`-` when no blocker is present). The
top-level help and command help expose `start`, `status`, `resume`, and `stop`.
The Run Database migrates from schema 14 to 15 by adding only the queue-owner
columns; existing Run and Delivery Queue data remain readable.

Focused checks:

- Pre-change inspection found no `internal/cli/deliver.go`, no `deliver`
  dispatch or top-level help entry, and no
  `TestDeliverCommandRecordsAQueueAndReportsIt` case.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run
  '^TestDeliverCommand(RecordsAQueueAndReportsIt|RejectsUnknownSlugBeforeRecordingQueue|StopsAndResumesPersistedQueue|RefusesUnknownFlags)$'
  ./internal/cli` passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run
  '^TestDeliveryQueue(RoundTripsItemsAndReceipts|OwnerMigrationPreservesExistingQueue)$'
  ./internal/store` passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1
  ./internal/cli ./internal/store ./internal/delivery` passed after the
  existing CLI suite's network-backed check was allowed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go vet ./internal/cli
  ./internal/store ./internal/delivery` passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache make verify-incremental`
  passed after the final code edit, including repository-wide vet and tests.
- `git diff --check` passed.

Acceptance evidence:

- The command family appears on the public surface: `TestRunHelp` asserts the
  top-level `roundfix deliver` entry, and `TestRunCommandHelp/deliver` asserts
  all four subcommands; both passed in the CLI package and incremental suites.
- `deliver status` prints each item's stage and blocker:
  `TestDeliverCommandRecordsAQueueAndReportsIt` records a real SQLite queue,
  parks its item as `review-stale`, and observes the exact
  `<slug>\tparked\treview-stale` public output.
- Unknown slugs are refused before any queue is recorded:
  `TestDeliverCommandRejectsUnknownSlugBeforeRecordingQueue` observes exit 2,
  no owner launch, and no Run Database file.

Additional requirement evidence:

- `TestDeliverCommandStopsAndResumesPersistedQueue` proves `stop` uses the
  stored PID and identity, clears only that owner claim, and `resume` launches
  from the unchanged persisted item.
- `TestDeliverCommandRefusesUnknownFlags` covers every subcommand, including a
  flag placed after the `start` slug.
- `TestDeliveryQueueRoundTripsItemsAndReceipts` now also proves owner identity
  persistence, refusal of a competing owner, and identity-bound release across
  a database reopen.
- `TestDeliveryQueueOwnerMigrationPreservesExistingQueue` proves schema 14
  queue items, stages, and blockers survive the schema 15 owner-column
  migration.

Task 05 remains the owner of shipped-skill and user-guide documentation; no
Task 05 files were changed here. The Daemon-owned commands under
`## Verification` were not run.

## Carry-forward provenance

- Source Run: `run_20260924T121521Z_744099e82ec81b88`
- Source commit: `a12d560216452e51831234c02895d597eb2bcce6`
