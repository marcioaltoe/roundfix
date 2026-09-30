---
task: task_04
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: pending
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
