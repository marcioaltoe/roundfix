---
task: task_05
spec: 0172-a-qa-gate-that-tells-the-truth
status: pending
type: backend
complexity: medium
---

# Task 05: The event stream survives

## Overview

`publishPreWorkProbeFindings` in `internal/daemon/task_engine.go` publishes the vacuous pre-work Verification event with `probed_commands`, while `projectVerificationRecord` in `internal/runevent/stream.go` requires `commands` for `verification_vacuous`. `encodeStreamEntry` in `internal/cli/events.go` returns the projection error, so `roundfix events` stops with `missing payload field "commands"` and exit `1`, and a Supervisor following the Run loses every later record, including its outcome. Journals already written before this fix hold only `probed_commands` and are read by the same command.

## Requirements

1. MUST make `publishPreWorkProbeFindings` add `commands`, the vacuous command strings in probe order, beside the retained `probed_commands`.
2. MUST make the vacuous projection read `commands`, and when it is absent derive it from the `probed_commands` entries whose verdict is `passed`, so journals written before the fix project; an event with neither key still fails projection.
3. MUST make `roundfix events` skip a record that `runevent.ProjectStreamEvent` cannot project, write one stderr warning naming its cursor, its event kind and the projection error, and continue, in replay and in follow; stdout stays JSONL records only, and a write or store error still exits `1`.
4. MUST replace `TestEventsMalformedRelevantPayloadFailsNoStdout` in `internal/cli/cli_test.go`, which pins the abort, with `TestEventsMalformedRelevantPayloadWarnsAndContinues`, which asserts exit `0`, the warning with its cursor, and the next well-formed record on stdout.
5. MUST state in `docs/user-guide/commands.md` (events), in the Supervisor Run Event Stream section of `.agents/skills/roundfix/SKILL.md` and in the Supervisor Run Event Stream entry of `CONTEXT.md` that the command skips a record it cannot project with a warning (the phrase `skips a record it cannot project`), replacing the skill's `malformed relevant Daemon payloads exit` wording, and name the vacuous record's `verification_vacuous` classification and `commands` field in the first two; regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
6. MUST put the new daemon test in `internal/daemon/vacuous_probe_projection_test.go`, projecting the payload the publisher actually wrote.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A vacuous pre-work event, as published, projects with its vacuous commands, and `roundfix events` replays a Run holding one with exit `0`.
- [ ] A journal written before the fix projects from `probed_commands`; an event with neither key is still malformed.
- [ ] A malformed record is skipped with a warning naming its cursor, in replay and in follow, and the next record is still emitted.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/runevent/stream.go`
- interface: `internal/cli/events.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/references/coverage-record.json`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestVacuousPreWorkEventProjectsItsCommands|TestProjectVacuousEventJournaledBeforeTheFixReadsProbedCommands|TestProjectVacuousEventWithoutCommandsIsMalformed|TestEventsMalformedRelevantPayloadWarnsAndContinues|TestEventsFollowWarnsOnAMalformedRecordAndContinues|TestEventsReplaysAVacuousPreWorkEvent)$" ./internal/daemon ./internal/runevent ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestVacuousPreWorkEventProjectsItsCommands TestProjectVacuousEventJournaledBeforeTheFixReadsProbedCommands TestProjectVacuousEventWithoutCommandsIsMalformed TestEventsMalformedRelevantPayloadWarnsAndContinues TestEventsFollowWarnsOnAMalformedRecordAndContinues TestEventsReplaysAVacuousPreWorkEvent; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && ! grep -q "TestEventsMalformedRelevantPayloadFailsNoStdout" internal/cli/cli_test.go && grep -q "skips a record it cannot project" docs/user-guide/commands.md && grep -q "skips a record it cannot project" CONTEXT.md && grep -q "skips a record it cannot project" .agents/skills/roundfix/SKILL.md && grep -q "verification_vacuous" .agents/skills/roundfix/SKILL.md && grep -q "verification_vacuous" docs/user-guide/commands.md && ! grep -q "malformed relevant Daemon payloads exit" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named cases exists and the abort test and wording remain, so the command fails.

## References

- [_techspec.md](_techspec.md) — The event stream
