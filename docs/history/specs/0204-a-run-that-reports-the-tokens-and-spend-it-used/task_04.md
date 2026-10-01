---
task: task_04
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: completed
type: backend
complexity: medium
---

# Task 04: A Delivery Queue stops starting items at its token ceiling

## Overview

task_02 stores a queue's `MaxTokens` and sums its Runs' tokens, and nothing
sets or reads the ceiling. This Task adds `deliver start --max-tokens <n>`,
parks each queued item at or above the ceiling as `queue-token-ceiling`
before it has a worktree, refuses a retry there, and answers the Pending
Question (ADR-0199). It then describes everything this Spec ships in the
Roundfix Skill. It is verifiable on its own: with a ceiling below the queue's
recorded tokens, the next queued item parks and the item already running
continues.

## Requirements

1. MUST add `--max-tokens <n>` to `deliver start`, with the rules of API Contract 4: an integer of at least 1, stored with the queue, printed on the limits line as `tokens <n>`; a value below 1 or not an integer is a usage error with exit `2` that records no queue and creates no Run Database. The `deliver` help MUST name the flag and keep every string `TestRunCommandHelp` and `TestDeliverHelpNamesTheLimitFlags` require.
2. MUST check the ceiling in `Engine.Run` right after the deadline check, only for an item in stage `queued`, as `_techspec.md` → The ceiling states: at or above `MaxTokens`, the item parks with blocker `queue-token-ceiling` and no branch, worktree or workflow action. An item in any later stage, including `running`, MUST continue, and no Run is signalled, stopped or cancelled.
3. MUST refuse `Engine.Retry` for any parked item while the queue's tokens are at or above its ceiling, with the reason of Surface Transcript 6, before any workspace action.
4. MUST make `PendingQuestionFor` answer the blocker `queue-token-ceiling` with the text of Surface Transcript 4.
5. MUST leave a queue with no ceiling behaving exactly as before: an existing test of `internal/delivery` MUST pass unchanged, and one new test MUST prove that recorded tokens never park an item when `MaxTokens` is zero.
6. MUST put the new tests in `internal/delivery/token_ceiling_test.go`, with the real store and the existing fake workflow, and in `internal/cli/deliver_token_ceiling_test.go`, which reproduces Surface Transcripts 4 and 6 on a temporary Roundfix Home. `internal/cli/deliver_limits_test.go` changes only if a test in it pins the help text or limits line this Task changes, and the Result names it.
7. MUST document `--max-tokens`, the blocker, its Pending Question answer and the retry refusal in `docs/user-guide/commands/deliver.md`, and add to the `## Token usage` section of `docs/user-guide/usage.md` one paragraph on the ceiling and its boundary.
8. MUST describe, under a new heading `### Token usage` in each of `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/events.md`, `.agents/skills/roundfix/references/implement.md` and `.agents/skills/roundfix/references/runs.md`, what that command now reports: the `Usage:` line, `--max-tokens` and `queue-token-ceiling`; the `usage` stream category; the `Tokens:` line; and `roundfix runs show` with its schema. It MUST add no text inside the `### QA settlement` section of any skill.
9. MUST raise the Roundfix Skill's version in both front-matter fields of `.agents/skills/roundfix/SKILL.md` by one patch step from the value on this Task's starting commit, following the version rule in force on that commit, then run `make skills-sync`, record the version with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.
10. MUST prove each new gate can fail. The Result MUST record one sabotage of the start check (for example comparing with `>` instead of `>=`) and one of the running-item boundary (for example parking any unmerged item), each with the test that failed, and that the code was restored.

Spec 0194 creates the per-command guide and skill files this Task edits
(`_prd.md` → Prerequisites). They are declared under `creates:` because they do
not exist when this Spec is authored. This Task MUST NOT create one that Spec
0194 has not created: when such a file is absent, the Task stops and reports
that Spec 0194 has not landed.

## Subtasks

- [ ] Add the flag and its validation.
- [ ] Park queued items and refuse retries at the ceiling.
- [ ] Answer the Pending Question.
- [ ] Document the ceiling, then describe the Spec's surfaces in the skill and record its version.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] `deliver start --max-tokens 5000000` records the ceiling and prints `tokens 5000000`; `--max-tokens 0` is refused and records nothing.
- [ ] At or above the ceiling, a queued item parks as `queue-token-ceiling` without a worktree, and below it the item starts.
- [ ] An item in stage `running` continues above the ceiling.
- [ ] A retry at the ceiling is refused with the reason of Surface Transcript 6.
- [ ] `deliver status` reproduces Surface Transcript 4.
- [ ] The four skill reference files and their mirrors carry `### Token usage`, and the skill's version changed and is recorded.

## Context

- creates: `internal/delivery/token_ceiling_test.go`
- creates: `internal/cli/deliver_token_ceiling_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/question.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/deliver_limits_test.go`
- creates: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `.agents/skills/roundfix/references/deliver.md`
- creates: `.agents/skills/roundfix/references/events.md`
- creates: `.agents/skills/roundfix/references/implement.md`
- creates: `.agents/skills/roundfix/references/runs.md`
- creates: `skills/roundfix/references/deliver.md`
- creates: `skills/roundfix/references/events.md`
- creates: `skills/roundfix/references/implement.md`
- creates: `skills/roundfix/references/runs.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/delivery/limits_test.go`
- instruction: `internal/cli/cli_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAQueuedItemParksAtTheTokenCeilingWithoutAWorktree|TestAQueuedItemStartsBelowTheTokenCeiling|TestARunningItemContinuesAboveTheTokenCeiling|TestARetryIsRefusedAtTheTokenCeiling|TestNoTokenCeilingNeverParks|TestPendingQuestionAnswersTheTokenCeiling|TestAQueuedItemParksAtTheDeadlineWithoutAWorktree|TestDeliverStartRecordsAndPrintsTheTokenCeiling|TestDeliverStartRefusesAMaxTokensBelowOne|TestDeliverStatusPrintsTheTokenCeilingQuestion|TestDeliverRetryIsRefusedAtTheTokenCeiling|TestDeliverHelpNamesTheLimitFlags|TestRunCommandHelp|TestEveryOwnedSkillVersionIsRecorded)$" ./internal/delivery ./internal/cli ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAQueuedItemParksAtTheTokenCeilingWithoutAWorktree TestAQueuedItemStartsBelowTheTokenCeiling TestARunningItemContinuesAboveTheTokenCeiling TestARetryIsRefusedAtTheTokenCeiling TestNoTokenCeilingNeverParks TestPendingQuestionAnswersTheTokenCeiling TestAQueuedItemParksAtTheDeadlineWithoutAWorktree TestDeliverStartRecordsAndPrintsTheTokenCeiling TestDeliverStartRefusesAMaxTokensBelowOne TestDeliverStatusPrintsTheTokenCeilingQuestion TestDeliverRetryIsRefusedAtTheTokenCeiling TestDeliverHelpNamesTheLimitFlags TestRunCommandHelp TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task ten of the fourteen named tests do not exist, so the command fails.
- `for pair in "docs/user-guide/commands/deliver.md|--max-tokens" "docs/user-guide/commands/deliver.md|queue-token-ceiling" "docs/user-guide/usage.md|--max-tokens" ".agents/skills/roundfix/references/deliver.md|### Token usage" ".agents/skills/roundfix/references/deliver.md|queue-token-ceiling" ".agents/skills/roundfix/references/events.md|### Token usage" ".agents/skills/roundfix/references/implement.md|### Token usage" ".agents/skills/roundfix/references/runs.md|roundfix runs show" ".agents/skills/roundfix/references/runs.md|roundfix/runs-show/v1"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- internal/delivery/token_ceiling_test.go)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/roundfix/SKILL.md" > "$tmp/before.md" || exit 1; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/version-old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/roundfix/SKILL.md > "$tmp/version-new"; if cmp -s "$tmp/version-old" "$tmp/version-new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no guide or skill file names `--max-tokens`, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 4; Core Features 10 and 11; Success Metric 5; Declared breaks; Recorded limits
- [_techspec.md](_techspec.md) — The ceiling; Surface Transcript 4; Surface Transcript 6; API Contract 4; API Contract 5; Testing Approach 5; Testing Approach 6; Build Order 4
- ADR-0199; ADR-0198

## Result

Implemented the token-ceiling slice for Daemon verification. The pre-existing
worktree change was this Task's Daemon-written `status: in_progress`; that
field remains unchanged. No other Task, Task Graph, commit, push or Pull
Request was changed or created. The prerequisite Spec 0194 guide and all
four canonical skill reference files existed before implementation.

### Implementation and acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Start records and prints `5000000`; zero records nothing | `deliver start` parses an explicitly supplied `--max-tokens` as an int64 of at least 1 before opening the Run Database, hoists it alongside the existing value flags, and stores it in the queue limits. `TestDeliverStartRecordsAndPrintsTheTokenCeiling` checks the stored value and exact limits line. `TestDeliverStartRefusesAMaxTokensBelowOne` checks zero and negative values, exit 2, no stdout, no owner launch and no Run Database. `TestDeliverStartRefusesANonIntegerMaxTokens` checks fractional, nonnumeric and overflowing values with the same no-write boundary. |
| At/above the ceiling parks without workspace; below starts | `Engine.Run` reads real queue usage immediately after the deadline check, only for a queued item with a positive ceiling, and parks before `advanceItem`. `TestAQueuedItemParksAtTheTokenCeilingWithoutAWorktree` checks both equality (5000000) and excess (5639755), empty branch/worktree, no provisioning and no workflow events. `TestAQueuedItemStartsBelowTheTokenCeiling` checks 4999999 tokens starts and merges. |
| A running item continues above the ceiling | `TestARunningItemContinuesAboveTheTokenCeiling` seeds a running item with a recorded workspace and queue usage of 5639755 against a 5000000 ceiling, then checks its Run starts and advances through merge. The new gate neither signals nor cancels Runs. |
| Retry refuses at the ceiling with Transcript 6 reason | `Engine.Retry` reads usage before deadline-blocker handling, retry-count enforcement, prerequisite re-entry or workspace work. `TestARetryIsRefusedAtTheTokenCeiling` checks the exact reason at equality for token-ceiling, deadline, unresolved-Run and prerequisite blockers, unchanged persisted item, and zero workspace/recovery/workflow actions. `TestDeliverRetryIsRefusedAtTheTokenCeiling` uses the real CLI and Engine with a temporary Home to check Transcript 6's stderr prefix/reason, exit 2, empty stdout, unchanged queue and no owner launch above the ceiling. |
| Status reproduces Transcript 4 | `PendingQuestionFor` supplies the exact new-queue answer for `queue-token-ceiling`. `TestPendingQuestionAnswersTheTokenCeiling` checks that answer; `TestDeliverStatusPrintsTheTokenCeilingQuestion` compares all stdout bytes, empty stderr and exit 0 against Transcript 4. Existing blocker `Park:` lines remain; this new blocker uses its item row and Pending Question without an extra `Park:` line, matching the authored transcript. |
| References and mirrors carry Token usage; version recorded | Added `### Token usage` to deliver/events/implement/runs references with their command's shipped surface, including the `Usage:` line and ceiling, default `usage` stream category, `Tokens:` line, and `roundfix runs show` schema `roundfix/runs-show/v1`. Both version fields advance from starting HEAD's 0.1.5 to 0.1.6; mirrors and the version digest were generated. Byte inspection confirms the full canonical and mirrored skill trees match and QA settlement sections are unchanged. |

`TestNoTokenCeilingNeverParks` proves a queue with 5639755 recorded tokens and
`MaxTokens == 0` still starts and merges its queued item. The existing
`TestAQueuedItemParksAtTheDeadlineWithoutAWorktree`,
`TestDeliverHelpNamesTheLimitFlags` and `TestRunCommandHelp` passed unchanged.
`TestDeliverHelpNamesTheTokenCeiling` checks the new help flag separately.
`internal/cli/deliver_limits_test.go` was not changed: its pinned limits line
already included `tokens none`, and all its required help strings remain.

The deliver guide now documents validation, the ceiling, Pending Question,
retry refusal and the item-start boundary. The usage guide's Token usage
section adds the ceiling paragraph. No unrelated implementation slice was
added.

### Focused checks

Commands used a task-scoped `GOCACHE=/private/tmp/roundfix-task04-cache` and
`rtk proxy` for exact Go output. No command from this Task's Verification
section was run.

- Before implementation, the new ceiling park test failed at both equality
  and excess because the item merged. Retry and Pending Question tests failed
  for their missing behavior, and CLI tests failed on the missing flag and
  transcript differences.
- Final focused check:
  `go test ./internal/delivery ./internal/cli -run 'TokenCeiling|MaxTokens|NoTokenCeiling|TestAQueuedItemParksAtTheDeadlineWithoutAWorktree|TestDeliverHelpNamesTheLimitFlags|TestRunCommandHelp' -count=1`
  exited 0; both packages reported `ok` after restoring sabotage and adding
  deadline-blocker retry coverage.
- `make skills-sync` exited 0. It rewrote exactly
  `skills/roundfix/SKILL.md`, `skills/roundfix/references/deliver.md`,
  `skills/roundfix/references/events.md`,
  `skills/roundfix/references/implement.md` and
  `skills/roundfix/references/runs.md`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and rewrote `skills/testdata/owned-skill-versions.json`, recording
  Roundfix 0.1.6 and its content digest.
- `make baseline-digests` exited 0 and reported
  `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}`.
  It rewrote no files.
- Python byte comparisons checked every file in the canonical and mirrored
  skill trees, all four Token usage headings, both 0.1.6 fields in both
  SKILL.md files, and preserved QA settlement sections against starting HEAD.
  All assertions passed. `git -c core.fsmonitor=false diff --check` exited 0.

### Sabotage evidence

1. Changed only the queued-start comparison from `>=` to `>`.
   `go test ./internal/delivery -run '^TestAQueuedItemParksAtTheTokenCeilingWithoutAWorktree$' -count=1`
   exited 1: the `5000000` subtest observed a merged item instead of a park.
   Restored the original comparison before the next probe.
2. Changed the queued-only boundary to all unmerged items
   (`item.Stage != store.DeliveryStageMerged`).
   `go test ./internal/delivery -run '^TestARunningItemContinuesAboveTheTokenCeiling$' -count=1`
   exited 1: the running item parked as `queue-token-ceiling` with no Run
   events. Restored the queued-only boundary; the final focused check above
   passed with both mutations absent.

Daemon-owned authored Verification and Task settlement remain pending. There
are no follow-up implementation changes required by this slice.

### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261001T094428Z_f4d3047c47b08ac0/verification/batch-004-attempt-1.log`.
The configured `make verify-changed` failed on
`TestEveryBlockerHasAParkClass`: the new exported token-ceiling blocker was
missing from the shared classification table, including its colon-detail
form. That command and the authored Verification commands were not rerun.

Added `TestTokenCeilingHasABudgetParkClass` in this Task's
`internal/delivery/token_ceiling_test.go`. Before the repair,
`go test ./internal/delivery -run '^(TestTokenCeilingHasABudgetParkClass|TestEveryBlockerHasAParkClass)$' -count=1`
exited 1, reproducing the missing classification and inconsistent recovery
action for the detailed blocker.

`internal/delivery/park_class.go` now classifies `queue-token-ceiling` as
`budget` and owns its exact new-queue answer. Removed the special answer
branch from `question.go`, restoring Pending Question's existing delegation
to `ClassifyPark`; that file now matches starting HEAD. The classifier is
part of this Task's new blocker contract. The existing classification sweep
was not changed or weakened. Transcript 4's status output remains unchanged.

After the repair, with `GOCACHE=/private/tmp/roundfix-task04-cache` and
`rtk proxy` for exact output,
`go test ./internal/delivery ./internal/cli -run 'TokenCeiling|MaxTokens|NoTokenCeiling|ParkClass|EveryBlocker|ExistingGenericParkAnswers|TestAQueuedItemParksAtTheDeadlineWithoutAWorktree|TestDeliverHelpNamesTheLimitFlags|TestRunCommandHelp' -count=1`
exited 0 and both packages reported `ok`. This includes the new bare/detailed
blocker regression, the unchanged exported-blocker sweep, existing Park Class
and generic-answer contracts, and the Task's token-ceiling checks and CLI
transcripts. `git -c core.fsmonitor=false diff --check` also exited 0.

The skill source and mirrors were unchanged by this repair, so the recorded
0.1.6 content digest and prior regeneration evidence remain applicable.
Task status remains Daemon-owned and unchanged; no other Task or Task Graph
was edited, and no commit, push or Pull Request was created. The Daemon's full
configured Verification rerun and settlement remain pending.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/delivery/park_class.go`

## Carry-forward provenance

- Source Run: `run_20261001T094428Z_f4d3047c47b08ac0`
- Source commit: `946e11c283f7c2cf93341384f01a1afe0a269e89`
