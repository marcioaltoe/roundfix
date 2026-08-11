---
task: task_10
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: test
complexity: high
---

# Task 10: Measure and record the before that the after is compared against

## Overview

Two of the Spec's six Success Metric rows failed for the same reason: the
evidence they compare against was produced by the changed build, so neither
comparison can establish what it claims.

- **R03 (F-02).** The current scenario reports 38,136 µs per 1,000 across 1,536
  Journal Sink events, and every batching rehearsal passes. The recorded
  baseline holds per-append percentiles and no comparable per-1,000 figure, and
  its same-harness "after" explicitly excludes batching. There is no before at
  the same boundary, so the improvement cannot be attributed to batching rather
  than to fewer events.
- **R06 (F-03).** `journal_consumer_corpus_test.go` was added in Task 05's own
  production commit, and its four cases reach full read, `events`, the header
  subset, and batch clocks. Attach, full Cockpit, reconcile replay detection,
  and `gc` have only current-fixture tests, which cannot establish pre/post
  identity.

This repository's acceptance standard is a characterization corpus captured
*before* the change with its breaks declared. Both rows were captured with the
change instead. The pre-Spec commit is `a2a4c86b`, and it is still buildable, so
the before is recoverable rather than lost.

## Requirements

1. MUST produce the per-1,000 measurement from a build of `a2a4c86b`, using the
   same harness, the same event count, and the same measurement boundary as the
   current figure, so the two numbers differ only in the code under test.
2. MUST record the before and after side by side under the Spec's `baseline/`
   directory, naming the commit each was measured from and the exact command
   that produced it, so a later reader can rerun both.
3. MUST state plainly whether the measured delta supports the write-amplification
   goal. If it does not, say so and stop; a Spec that reports an improvement it
   did not measure is worse than one that reports none.
4. MUST record a journal produced by the `a2a4c86b` build and replay it through
   all five named consumers — `events`, attach, Cockpit rendering, reconcile
   replay detection, and `gc` — comparing behaviour against the current build.
5. MUST NOT change production code. If a consumer's behaviour genuinely differs
   between the two builds, that is a finding to report, not a difference to
   absorb into the expectation.

## Subtasks

- [ ] Measure per-1,000 at the same boundary from `a2a4c86b`.
- [ ] Record both measurements with their commands and commits.
- [ ] Record a pre-change journal and replay it through all five consumers.

## Acceptance Criteria

- [ ] The baseline directory holds a before and an after at the same boundary,
      each naming its commit and command.
- [ ] The report states whether the delta supports the goal.
- [ ] A pre-change journal replays through all five consumers with the
      comparison recorded per consumer.

## Bounded scope

This Task may create or modify only:

- `docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/`
- `internal/store/journal_consumer_corpus_test.go`
- `docs/specs/0081-a-journal-cheap-to-write-and-keep/task_10.md`

## Verification

- `grep -rq 'a2a4c86b' docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/` — expected: exits 0. No baseline file names the pre-Spec commit before this Task.
- `grep -rqi 'per 1,000\|per-1,000' docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/2026-08-11-before.md` — expected: exits 0. The recorded before holds percentiles and no per-1,000 figure today.
- `GOCACHE="$PWD/.gocache" go test ./internal/store -run '^TestJournalConsumerCorpusReplaysEveryConsumer$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestJournalConsumerCorpusReplaysEveryConsumer'` — expected: exits 0. The case does not exist before this Task.

## References

- `_prd.md` → Success Metrics R03 and R06.
- `qa/qa-report-2026-08-11.md` → F-02 and F-03.
