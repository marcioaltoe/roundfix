---
task: task_11
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
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

- `! grep -rq '0081-a-journal-cheap-to-write-and-keep' internal/store/` — expected: exits 0. The slug appears as a `filepath.Join` argument today, so the whole path never occurs as one literal; the slug alone is what has to disappear.
- `test -f internal/store/testdata/2026-08-11-prechange-roundfix.db` — expected: exits 0. `internal/store/testdata/` does not exist today, so this is what proves the fixtures moved rather than being copied by path alone.

Asserting that the corpus test passes is deliberately absent: it passes right
now, because the Spec is still at its active path. The proof that matters is
that it keeps passing once the Spec is archived, and no command can state that
before the archive exists. Requirement 2 carries it, and the Run-level gate
proves the test after this Task settles.

## References

- `_prd.md` → Success Metric R06.
- `task_10.md` → the corpus and expectations this Task relocates.

## Result

### Implementation

- Moved the recorded SQLite journal, consumer expectation JSON, Cockpit
  golden, and both Go overlay harness sources from the Spec baseline into
  `internal/store/testdata/`.
- Changed only the corpus input resolution: the test now builds each fixture
  path below `internal/store/testdata/`; its consumer commands, overlays,
  environment, and assertions remain unchanged.
- Left the measurement and consumer-replay Markdown records in the Spec
  baseline. The generator and measurement harness sources also remain there
  because this production test does not consume them.

### Acceptance-criterion evidence

- **Archive-safe corpus resolution:** source inspection shows every executable
  input is resolved below `internal/store/testdata/`, independent of the Spec
  location. The focused corpus replay passed from that directory. The Spec has
  not been archived in this Task worktree, so the Run-level post-archive proof
  remains owned by the Daemon gate.
- **No Spec path in test source:** the changed corpus function now joins only
  `internal`, `store`, and `testdata` below the repository root; it contains no
  Spec directory component. The Task's declared repository-wide search was
  not run in this Daemon-assigned turn.
- **Evidence records preserved:** `rtk ls
  docs/specs/0081-a-journal-cheap-to-write-and-keep/baseline` showed
  `2026-08-11-before.md`, `2026-08-11-repaired.md`, and
  `2026-08-11-consumer-replay.md` still present after the move.

### Focused checks

- Pre- and post-move `rtk sha256sum` results matched for all five files:
  `1c38ed391769fb2c3ba5e48875ed868737bf0bc5659c12bdbf7ee8b19f5c5dec`
  (database),
  `3c4eb9511ee29dda3116ac8863e8b7f0b5a74fcaf594470ce6f0788b7f5ab04b`
  (expectation JSON),
  `16eb1d01eb12a55728e0343c12768f04e40809359acc8214c368a42e5a8d9700`
  (Cockpit golden),
  `204d5a2d9eab3d57b17bd5ef9050cf167b6112e1ca52b4eee3178f80a9439260`
  (CLI harness), and
  `d7c6d7da9220ad863b084564c100c5ad3012ddd40879c682e23a476fd02cb5f8`
  (Cockpit harness).
- `GOCACHE=/private/tmp/roundfix-task11-gocache
  GOTMPDIR=/private/tmp/roundfix-task11-gotmp rtk go test ./internal/store
  -run '^TestJournalConsumerCorpusReplaysEveryConsumer$' -count=1`: passed;
  RTK reported three passing tests in one package.

The Task's declared Verification commands and repository-wide Verification
were not run; the Daemon owns those checks and settlement.
