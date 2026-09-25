---
task: task_07
spec: 0172-a-qa-gate-that-tells-the-truth
status: pending
type: backend
complexity: low
---

# Task 07: The journal consumer corpus replays the new event stream signature

## Overview

Corrective Task from the second QA gate of 2026-09-25, whose repository Verification precondition failed: task_05 added an `io.Writer` (stderr) parameter to `replayEventStream` in `internal/cli/events.go`, but the journal consumer corpus harness `internal/store/testdata/task10-cli-consumer-harness_test.go.txt`, which `TestJournalConsumerCorpusReplaysEveryConsumer` in `internal/store/journal_consumer_corpus_test.go` copies into `internal/cli` and compiles, still calls the old six-argument form, so `internal/cli` fails to build inside the corpus replay.

## Requirements

1. MUST update the harness's `replayEventStream` call to the current signature, passing a writer whose content the harness either discards or asserts empty for a well-formed journal, and keep every observation the harness records unchanged.
2. MUST search every `internal/**/testdata` harness for other calls to functions whose signatures this Spec changed and update them the same way.
3. MUST NOT change `replayEventStream` or any production file.

## Subtasks

- [ ] Implement the requirements above.

## Acceptance Criteria

- [ ] `TestJournalConsumerCorpusReplaysEveryConsumer` passes, including the `events` consumer subtest.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/testdata/task10-cli-consumer-harness_test.go.txt`
- interface: `internal/store/journal_consumer_corpus_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestJournalConsumerCorpusReplaysEveryConsumer$" ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestJournalConsumerCorpusReplaysEveryConsumer"` — expected: exit 0; before this Task the harness calls the six-argument form and `internal/cli` fails to build in the replay, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
