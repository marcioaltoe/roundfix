---
task: task_04
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: pending
type: backend
complexity: medium
---

# Task 04: `deliver start` takes explicit limits and `deliver status` presents one Pending Question

## Overview

task_02 records and enforces queue limits, but the operator has no way to set them, and `runDeliverStatus` in `internal/cli/deliver.go` prints every item row with equal weight. This Task adds the `--max-duration` and `--max-retries` flags to `deliver start`, prints the recorded limits from `deliver start` and `deliver status`, and presents one Pending Question: the lowest-position parked item and the explicit action that answers it. The limits come from the operator's flags and are persisted in the Run Database. The question is derived from the persisted item records, read by the operator or a Supervisor, and changed only by an operator command.

## Requirements

1. MUST make `parseDeliverStart` accept `--max-duration <duration>`, a positive Go duration, and `--max-retries <n>`, an integer of at least `1`, both declared as value flags to `hoistCommandFlags`. A zero, negative or malformed value exits `2` and records no queue.
2. MUST make `runDeliverStart` set the deadline to `time.Now().UTC()` plus the duration, truncated to whole seconds. It records the queue with `CreateDeliveryQueueWithLimits`, then prints to stdout, before the detached owner report, `Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, spend not measured`.
3. MUST keep the `deliverUsage` line `roundfix deliver start <slug>...` and add a Flags block naming `--max-duration <duration>` and `--max-retries <n>`.
4. MUST add `internal/delivery/question.go` with `PendingQuestion` and `PendingQuestionFor(queue store.DeliveryQueue) (PendingQuestion, bool)` as the TechSpec states. The question is the parked item with the lowest position, `Waiting` counts the other parked items, and each blocker class receives the answer from the TechSpec's table.
5. MUST make `runDeliverStatus` print the item rows unchanged, then the `Limits:` line. When an item is parked, it MUST also print `Pending question: <slug> parked <blocker>` and `Answer: <answer>`, and `Waiting behind it: <n> parked item(s)` when `n` is greater than zero.
6. MUST NOT let the engine, the owner loop or any elapsed time change a parked item. Only a Delivery Retry or a new queue answers the question.
7. MUST update only the tests this change invalidates: the exact status expectations in `TestDeliverStatusPrintsTheItemWorktree` in `internal/cli/deliver_test.go` gain the `Limits:` line and, for its parked item, the question lines.
8. MUST document in the deliver section of `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md`:
   - both flags, by the strings `--max-duration` and `--max-retries`;
   - the `Limits:` line;
   - the `queue-deadline` blocker and the retry limit refusal;
   - the Pending Question, by the string `Pending question:`.

   Then MUST regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
9. MUST add to `CONTEXT.md` the entries **Delivery Plan**, **Delivery Revalidation**, **Delivery Queue Limit** and **Pending Question**, each with an `_Avoid_` line, and make its **Delivery Queue** entry name the queue's limits.
10. MUST put the question tests in `internal/delivery/question_test.go` and the command tests in `internal/cli/deliver_limits_test.go`. The command tests drive the public CLI with the real Run Database and inject only the owner starter, as the existing deliver tests do. The owner-pass test advances the injected `Clock` with no wall-clock wait.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] `deliver start --max-duration 2h --max-retries 2` records a deadline within the start window and a retry limit of 2, and prints them in the `Limits:` line. Omitted flags record and print `none`.
- [ ] A zero, negative or malformed `--max-duration`, and a `--max-retries` below `1`, each exit `2` and record no queue.
- [ ] For a queue with two parked items, `deliver status` prints exactly one Pending Question, naming the lower-position item, with one item waiting. A queue with no parked item prints no question.
- [ ] Each blocker class receives its answer. The question is unchanged after an owner pass with the clock advanced past the deadline.
- [ ] `roundfix deliver --help` names both flags, and the guides and glossary describe the limits and the Pending Question.

## Context

- interface: `internal/cli/deliver.go`
- creates: `internal/cli/deliver_limits_test.go`
- interface: `internal/cli/deliver_test.go`
- creates: `internal/delivery/question.go`
- creates: `internal/delivery/question_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPendingQuestionIsTheLowestPositionParkedItem|TestPendingQuestionAnswersEachBlockerClass|TestNoPendingQuestionWithoutAParkedItem|TestAPendingQuestionSurvivesOwnerPassesAndTime|TestDeliverStartRecordsAndPrintsItsLimits|TestDeliverStartRecordsNoneForOmittedLimits|TestDeliverStartRefusesANonPositiveMaxDuration|TestDeliverStartRefusesAMalformedMaxDuration|TestDeliverStartRefusesAMaxRetriesBelowOne|TestDeliverStatusPrintsOnePendingQuestion|TestDeliverStatusPrintsNoQuestionWithoutAParkedItem|TestDeliverHelpNamesTheLimitFlags|TestDeliverStatusPrintsTheItemWorktree|TestATerminalQueueIsReplacedByANewStart)$" ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPendingQuestionIsTheLowestPositionParkedItem TestPendingQuestionAnswersEachBlockerClass TestNoPendingQuestionWithoutAParkedItem TestAPendingQuestionSurvivesOwnerPassesAndTime TestDeliverStartRecordsAndPrintsItsLimits TestDeliverStartRecordsNoneForOmittedLimits TestDeliverStartRefusesANonPositiveMaxDuration TestDeliverStartRefusesAMalformedMaxDuration TestDeliverStartRefusesAMaxRetriesBelowOne TestDeliverStatusPrintsOnePendingQuestion TestDeliverStatusPrintsNoQuestionWithoutAParkedItem TestDeliverHelpNamesTheLimitFlags TestDeliverStatusPrintsTheItemWorktree TestATerminalQueueIsReplacedByANewStart; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "--max-duration" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "--max-retries" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "Pending question:" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "--max-duration" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "Pending question:" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Delivery Plan**" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Delivery Revalidation**" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Delivery Queue Limit**" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Pending Question**" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the twelve new named tests exists and neither the flags nor the glossary entries are documented, so the command fails.

## References

- [_techspec.md](_techspec.md) — Limit flags, status and the Pending Question
- `_prd.md` → Goals 4-5; Core Features 4-5; Success Metrics 3-4
- `_techspec.md` → API Contract 5; API Contract 7; Testing Approach 4

## Result
