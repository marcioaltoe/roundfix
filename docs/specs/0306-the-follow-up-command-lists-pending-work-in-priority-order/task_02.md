---
task: task_02
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
status: completed
type: backend
complexity: medium
---

# Task 02: The judge answers four yes/no questions per Follow-up Item, after a secret scan

## Overview

The judge asks one question per request for five fixed judgment kinds. This
Task adds the `follow-up-priority` judgment: one request per Follow-up Item
carrying four Noul questions, a verdict per question at the thresholds of
ADR-0319, a secret scan before anything leaves the machine, and one Judge Log
row per call that names the model that answered (`_prd.md` → Core Feature 4;
`_techspec.md` → API Contract 5, Data Models). It is testable alone with the
package's fake transports.

## Requirements

1. MUST add the judgment `follow-up-priority` to `internal/judge/questions.json`
   after `archive-value`, exactly as `_techspec.md` → Data Models shapes it:
   `question_id` `follow_up_priority`, the four Noul questions `is_fix`,
   `blocks_other_project`, `is_small` and `affects_production` with the
   instructions written there verbatim and no `criteria` key, `yes_at_least`
   `0.7`, `no_at_most` `0.3`, `item_max_chars` `4000` and `max_items` `100`.
   It MUST NOT change any other judgment, the pin, the transports, the price
   or the ceiling, so `TestQuestionFileLoads` keeps passing unchanged.
2. MUST decode it into a new `Questions` field and add
   `PrioritizeFollowUp(ctx, q, req FollowUpRequest) (FollowUpReport, error)`
   with the types of `_techspec.md` → Interfaces, in a new
   `internal/judge/follow_up.go`, following `AdviseArchive` for the Run
   Marker skip, the key selection, the month's ceiling before the run and
   before every attempt, the 429/529 retry, redaction of key values, and
   `Stopped` on a service failure or a refused key.
3. MUST send, for each item in order, one request whose `questions` map holds
   the four questions under their ids, and whose `state` is
   `{"project": <req.Project>, "item": {"kind", "type", "origin", "title",
   "text"}}`, where `text` is the item text cut to `item_max_chars` runes on a
   rune boundary. The `client.ask` and `evaluate` paths MUST keep sending one
   question for every other kind.
4. MUST scan the text the state will carry (the title, the origin and the
   cut text, joined by newlines, before JSON encoding) with
   `secretscan.ScanContent(<item ref>, <that text>)` before sending; on a match, send nothing for that item and mark
   it `skipped` with the reason `not sent: secret pattern <kind>`, naming the
   first match's kind and never its value.
5. MUST mark an item `answered` only when the reported model matches the
   accepted-model pin and all four answers are Nouls between 0 and 1;
   otherwise `skipped` with the existing reasons (`answered by <model>,
   thresholds belong to <pin>`, `unreadable answer`). For an answered item,
   each verdict is `yes` at or above `yes_at_least`, `no` at or below
   `no_at_most`, and `unsure` otherwise.
6. MUST mark the items after the first `max_items` as `skipped` with
   `item limit reached (100)` without sending them, and return one
   `FollowUpAnswer` per input item, in input order.
7. MUST append one Judge Log row per call with `judgment`
   `follow-up-priority`, `question_id` `follow_up_priority`, `artifact` the
   item reference, `target` the item kind, an empty `spec`, `probabilities`
   holding the four Noul values by question id, `noul` null, `model` the
   model the service reported and `requested_model` the transport's request
   model, using the existing log path and file modes.
8. MUST add `internal/judge/follow_up_test.go` with these tests, all on
   `roundTripFunc` fakes and a temporary home, never the network and never a
   real key (build any key-shaped or secret-shaped value at runtime by
   concatenation, obviously fake):
   - `TestFollowUpAsksFourQuestionsInOneRequest`: one request per item, each
     body naming exactly the four ids and the item's state;
   - `TestFollowUpVerdictsUseTheThresholds`: answers of exactly `0.7`,
     `0.69`, `0.31` and exactly `0.3` give `yes`, `unsure`, `unsure` and `no`;
   - `TestFollowUpLogsTheModelOfEachAnswer`: one row per call whose `model`
     is the versioned id the fake reports (for example `jev-1.13.0`) and
     whose `probabilities` hold the four values;
   - `TestFollowUpNeverSendsASecretShapedItem`: an item whose text holds a
     value `secretscan` recognizes produces no request, a `skipped` answer
     naming the pattern's kind, and no copy of the value in any request body,
     log row or reason, while the next item is still asked;
   - `TestFollowUpCutsTextOnARuneBoundary`: an item of 3,999 ASCII
     characters followed by `ção` sends a `text` of exactly 4,000 runes that
     ends with `ç` and holds no U+FFFD;
   - `TestFollowUpFailsOpen`, with one subtest each for a missing key, a set
     `RunID`, the month's ceiling already reached, an unreadable Judge Log, a
     503 that stops the remaining items, and a model outside the pin; none
     returns an error, and none but the last two sends a request;
   - `TestFollowUpStopsAtTheItemLimit`: with 101 items, 100 requests and the
     101st `skipped` with `item limit reached (100)`.
9. MUST NOT change `internal/cli`, `roundfix spec judge`, `AdviseArchive`,
   the Judge Key Variables or `WithoutKeyVariables`.

## Subtasks

- [ ] Add the judgment to `questions.json` and decode it.
- [ ] Teach `client.ask` and `evaluate` the four-question kind without changing the others.
- [ ] Add `PrioritizeFollowUp` with the scan, the cut, the limit, the verdicts and the log rows.
- [ ] Add the seven tests.

## Acceptance Criteria

- [ ] Each Follow-up Item is asked in one request of four Nouls, and each answer becomes `yes`, `no` or `unsure` at 0.7 and 0.3 inclusive.
- [ ] A secret-shaped item, an item beyond the limit, a missing key, a Run, the ceiling and a failed service send nothing for the affected items and say why.
- [ ] Every call leaves one Judge Log row naming the model that answered.

## Consumers

- checkpoint: task_03

## Context

- instruction: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- instruction: `internal/judge/archive_advice.go`
- instruction: `internal/judge/client_test.go`
- instruction: `internal/secretscan/scan.go`
- interface: `internal/judge/questions.json`
- interface: `internal/judge/questions.go`
- interface: `internal/judge/client.go`
- creates: `internal/judge/follow_up.go`
- creates: `internal/judge/follow_up_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestFollowUpAsksFourQuestionsInOneRequest|TestFollowUpVerdictsUseTheThresholds|TestFollowUpLogsTheModelOfEachAnswer|TestFollowUpNeverSendsASecretShapedItem|TestFollowUpCutsTextOnARuneBoundary|TestFollowUpFailsOpen|TestFollowUpStopsAtTheItemLimit|TestQuestionFileLoads|TestArchiveAdviceJudgesOnlyCandidateFiles)$' ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpAsksFourQuestionsInOneRequest TestFollowUpVerdictsUseTheThresholds TestFollowUpLogsTheModelOfEachAnswer TestFollowUpNeverSendsASecretShapedItem TestFollowUpCutsTextOnARuneBoundary TestFollowUpFailsOpen TestFollowUpStopsAtTheItemLimit TestQuestionFileLoads TestArchiveAdviceJudgesOnlyCandidateFiles; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q '"follow-up-priority"' internal/judge/questions.json || { echo 'missing judgment' >&2; exit 1; }; GOOS=windows go vet ./internal/judge` — expected: exit 0. Before this Task the seven follow-up tests do not exist, so the command fails at its first pass check. Afterwards they pass beside the unchanged question-file and archive-advice tests, the judgment is in the embedded catalog, and the package's tests compile for Windows.

## References

- `_prd.md` → Core Feature 4; User Story 2; Success Metric 2; Success Metric 3
- `_techspec.md` → Interfaces; Data Models; API Contract 5; Invariants 3-4; Testing Approach; Build Order 2
- ADR-0209
- ADR-0319

## Result

Implemented the Task 02 judge slice. `Questions.FollowUp` decodes the new
`follow-up-priority` catalog entry with the four verbatim instructions,
inclusive 0.7/0.3 thresholds, 4,000-rune excerpt and 100-item cap.
`PrioritizeFollowUp` returns one answer per input in order, sends four Nouls
in one request per eligible item, scans the unencoded title, origin, excerpt,
type and kind, and uses the existing keys, transports, pin and private Judge
Log. No live judge request or real credential was used.

Acceptance evidence from focused fake-transport checks:

- Four questions and verdicts: `TestFollowUpAsksFourQuestionsInOneRequest`
  asserts the exact question map and item state, including Portuguese text;
  `TestFollowUpVerdictsUseTheThresholds` asserts yes/unsure/unsure/no for
  0.7/0.69/0.31/0.3. `TestFollowUpCutsTextOnARuneBoundary` checks 4,000
  runes ending with `ç`, without U+FFFD. `TestFollowUpRejectsUnreadableNouls`
  covers missing, null, negative, above-one and wrong-type Nouls and broken
  JSON.
- No-send and stop behavior: `TestFollowUpNeverSendsASecretShapedItem`
  checks text, title, origin, kind and type independently, verifies the
  scanner recognizes the runtime-built fake, and checks the value never
  appears in requests, reports or log rows while the next item is answered.
  `TestFollowUpFailsOpen` covers missing keys, a Run Marker, exhausted
  monthly budget, malformed Judge Log, HTTP 503 and a model outside the pin.
  `TestFollowUpStopsAtTheItemLimit` asserts 100 requests and the 101st
  answer's exact limit reason. `TestFollowUpRetriesAndChecksEachAttempt`
  covers 429/529 recovery, retry exhaustion, budget consumption before a
  retry or the next item, a refused key and error-message key redaction.
- Judge Log: `TestFollowUpLogsTheModelOfEachAnswer` asserts distinct reported
  versioned models, requested model, judgment, question id, reference, kind,
  empty Spec, null Noul and all four probabilities, plus accumulated cost.
  Retry tests assert a separate row and attempt number for each call.

Focused checks run:

- `go test ./internal/judge -run '^TestFollowUp' -count=1`: exit 0 after
  correcting a fake-token fixture length to match the scanner's exact
  GitHub-token pattern. The initial check caught the nonmatching fixture.
- `go test ./internal/judge -run '^(TestFollowUp|TestAsk|TestModelPin|TestRunNever)' -count=1 -v`:
  exit 0; all nine new follow-up tests and the selected existing transport,
  model-pin and redirect checks passed after the final code edit.
- `git diff --check`: exit 0.
- Parsed the current and HEAD question catalogs, removed only the new
  judgment from the current copy, and compared them: identical. The new
  entry follows `archive-value`; all existing judgments, model settings,
  transports, price and ceiling are unchanged.

Environment: the sandbox could not initialize with the symlinked writable
root; commands used authorized escalation in the resolved assigned checkout.
The only pre-existing changed path was this Task file's daemon-owned
`pending` to `in_progress` status change, which was preserved.

Declared Task Verification, Windows vet and repository Verification were not
run; they remain daemon-owned. Task status and checkboxes are unchanged.
No other Task file, Task Graph manifest, CLI implementation, archive advice,
key-variable contract or `WithoutKeyVariables` was changed. No follow-up
outside this slice was identified.

## Carry-forward provenance

- Source Run: `run_20261024T004042Z_a30dec8c359819a6`
- Source commit: `b9c62a1ad390fa50b0dfdceb24f43050b9379deb`
