---
task: task_01
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: test
complexity: medium
---

# Task 01: Measure what the journal actually costs

## Overview

The first artifact, and the one every later claim in this Spec cites. The
three-gigabyte measurement that opened the Spec is a symptom; the costs worth
fixing — a write lock held across a full-table aggregate at Run start, one
transaction and one fsync per agent output line — were only visible once the
write path was read rather than the disk.

Nothing here changes behaviour. It builds the harness and records the before.

## Requirements

1. MUST measure event-write latency, lock-wait time, and `SQLITE_BUSY`
   frequency as functions of journal size and concurrent writer count.
2. MUST measure Run-start cost against journal size, since the retention
   eligibility scan is the path this Spec claims is size-dependent.
3. MUST record the results as a committed artifact in the Spec folder, with
   the harness parameters, the machine facts that affect the numbers, and the
   commit measured.
4. MUST make the harness rerunnable, so an after can be produced from the same
   code and compared honestly against the before.
5. MUST seed its own journal fixtures rather than measuring the operator's live
   database, so the measurement is hermetic and repeatable.
6. MUST NOT change any production code path, configuration default, or schema.

## Subtasks

- [ ] Build the seeded-journal harness with its parameters.
- [ ] Record write latency, lock wait, `SQLITE_BUSY`, and Run-start cost.
- [ ] Commit the baseline artifact with the commit and machine facts.

## Acceptance Criteria

- [ ] The baseline artifact exists in the Spec folder and names the commit it
      measured.
- [ ] Rerunning the harness on the same commit reproduces comparable numbers.
- [ ] No production code path or default changed.

## Verification

- `output="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/*.md 2>/dev/null)"; [ -n "$output" ]`
  — expected: exit 0; the committed baseline artifact exists.
- `report="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/*.md | tail -1)"; grep -qiE 'latency' "$report" && grep -qi 'SQLITE_BUSY' "$report" && grep -qiE 'run start|run-start' "$report"`
  — expected: exit 0; the three measured dimensions are recorded.
- `output="$(go test -count=1 ./internal/store -run 'Baseline|Harness' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the harness is selected and runs.
  — expected: exit 0; nothing outside the Spec folder and store test files
  changed, so no production path moved.

## References

- `_prd.md` → Core Feature 1; Success Metrics; Decisions (measurement precedes
  design).
- `_techspec.md` → Testing Approach; Build Order 1.
- `references/2026-08-06-event-journal-payload-economics.md` → the intent this
  measurement serves.

## Result

### Implementation

- Added a rerunnable store-test harness that creates a fresh temporary Roundfix
  Home for each journal size, migrates and seeds its own Run Database, and
  never opens the operator database.
- Measured real `AppendRunEvent` wall time across 0, 1,000, and 10,000 seeded
  events with 1, 2, and 4 synchronized Store writers. The harness records p50
  and p95 latency, a same-size uncontended-baseline lock-wait estimate, and
  typed `SQLITE_BUSY` counts.
- Measured the no-candidate Run-start retention sequence at each journal size:
  `TerminalRunPruneCandidates` followed by the immediate-transaction
  `PruneTerminalRuns` path.
- Recorded the before in
  `baseline/2026-08-11-before.md`, including production commit
  `685d201b658cc46e944634a3c072da2a7d1d83c3`, fixed harness parameters,
  machine and SQLite facts, two fresh result sets, the lock-wait definition,
  repeatability evidence, and measurement limits.

### Focused checks

- Pre-change signal: `rtk rg --files internal/store
  docs/specs/0081-a-journal-cheap-to-write-and-keep | rtk grep -Ei
  'baseline|harness'` returned no matching path before implementation.
- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk proxy go test -count=2
  ./internal/store -run '^TestJournalMeasurementHarness$' -v` passed two fresh
  executions. Across both executions all 336 writes succeeded and typed
  `SQLITE_BUSY` frequency was zero. Four-writer p95 latency differed by less
  than 2% at each journal size; Run-start p50 differed by 1.6% at 1,000 events
  and 0.9% at 10,000 events.
- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk proxy go test -count=1
  ./internal/store -run
  '^TestJournalMeasurementHarnessRejectsInvalidParameters$' -v` passed all
  three negative parameter cases.
- `rtk git -c core.fsmonitor=false diff --check` passed.
- The first focused attempt without a task-local `GOCACHE` did not compile: the
  sandbox refused writes under `~/Library/Caches/go-build`. Repointing only the
  disposable Go build cache to `/private/tmp` removed that environment boundary.

### Acceptance evidence

1. **Baseline artifact exists and names the measured commit.**
   `baseline/2026-08-11-before.md` names the full production commit and records
   both measurement runs.
2. **The harness reproduces comparable numbers on the same commit.** The fresh
   `-count=2` run passed twice; the stable tail-latency and Run-start bands are
   quantified in the baseline report rather than inferred from one run.
3. **No production path or default changed.** The Task adds only
   `internal/store/journal_baseline_test.go`, the Spec-local baseline report,
   and this Result. The pre-existing `pending` to `in_progress` frontmatter
   change remains Daemon-owned.

### Daemon handoff

- The commands under this Task's `## Verification` were not run. The Daemon
  owns those commands, terminal status, and the Task commit.
