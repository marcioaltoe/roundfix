---
task: task_11
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: test
complexity: low
---

# Task 11: Stop a production test from reading a Spec that will move

## Overview

`TestJournalConsumerCorpusReplaysEveryConsumer` resolves its fixtures from
`docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/`. Archiving the
Spec moves that directory to `docs/specs/_archived/...`, and the test fails
with `no such file or directory` for the recorded database, the expectation
JSON, the Cockpit golden, and both harness overlays.

This was proved on 2026-08-11 by archiving the Spec and running `make verify`:
the test passed before the move and failed immediately after. The archive was
reverted so this Task can land first.

Every Spec is archived when it closes, so a production test that reads its own
Spec directory is a test with an expiry date. The evidence records belong to the
Spec and stay there; the files the test *executes against* are test inputs and
belong beside the test.

## Requirements

1. MUST move the files this test consumes — the recorded journal database, the
   consumer expectation JSON, the Cockpit golden, and the harness overlay
   sources — under `internal/store/testdata/`, preserving their bytes exactly.
2. MUST resolve them by a path that does not name any Spec directory, so the
   test's result is identical before and after that Spec is archived.
3. MUST leave the Spec's evidence records — the measurement and consumer-replay
   Markdown — in the Spec. They are the record of what was measured, not inputs
   to a test.
4. MUST NOT change what the test asserts. This Task moves inputs; the pre/post
   comparison it performs stays byte-for-byte what Task 10 established.

## Subtasks

- [ ] Move the four consumed fixture kinds under `internal/store/testdata/`.
- [ ] Resolve them without naming a Spec directory.
- [ ] Leave the Markdown evidence in the Spec.

## Acceptance Criteria

- [ ] The test passes with the Spec directory at its archived path.
- [ ] No test source names `docs/specs/0081-a-journal-cheap-to-write-and-keep`.
- [ ] The Spec still holds its measurement and consumer-replay records.

## Bounded scope

This Task may create or modify only:

- `internal/store/journal_consumer_corpus_test.go`
- `internal/store/testdata/`
- `docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline/`
- `docs/specs/0081-a-journal-cheap-to-write-and-keep/task_11.md`

## Verification

- `! grep -rq 'docs/specs/0081-a-journal-cheap-to-write-and-keep' internal/store/` — expected: exits 0. The test names that path today.
- `GOCACHE="$PWD/.gocache" go test ./internal/store -run '^TestJournalConsumerCorpusReplaysEveryConsumer$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestJournalConsumerCorpusReplaysEveryConsumer'` — expected: exits 0, proving the move preserved the comparison.

## References

- `_prd.md` → Success Metric R06.
- `task_10.md` → the corpus and expectations this Task relocates.
