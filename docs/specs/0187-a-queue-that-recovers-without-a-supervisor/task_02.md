---
task: task_02
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: pending
type: backend
complexity: medium
---

# Task 02: A retry refused by a Spec amendment prints how to recover it

## Overview

A Delivery Retry carries forward the proved Tasks of the item's Runs. When a later non-Task commit on the item branch amends the Spec, every Task's declared inputs move, so the carry-forward correctly refuses the whole set. The refusal's next action names only `roundfix reconcile <run> --carry-forward`, and that refuses the same way while the amendment stays under it. On 2026-09-29 recovering Spec 0181 took a reset, the carry-forward and a cherry-pick of the amendment, worked out by hand. This Task names the amending commits in the refusal and prints the recovery commands in order. Roundfix prints them and never runs them.

## Requirements

1. MUST add `carryForwardAmendments(ctx, repository, run, candidates)` to `internal/cli/carryforward.go`. It returns the amending commits, oldest first, with `ok` true only when all of these hold:
   - every candidate with a refusal reason was refused for moved inputs and lists at least one;
   - `run.HeadSHA` resolves in `repository`;
   - `git log <run.HeadSHA>..HEAD -- <moved inputs>` lists at least one commit, and none of them carries a `Roundfix-Task` trailer.
2. MUST give `deliveryCarryForwardRefusal` in `internal/cli/deliver_workflow.go` an `amendments` list, filled from that helper when `CarryForward` refuses. With amendments, `Error()` MUST append `; amended by <sha>, <sha>`, and `NextAction()` MUST return these five commands in order, with real values substituted:
   - `git -C <worktree> branch roundfix-amended-<run-id> HEAD`
   - `git -C <worktree> reset --hard <first-amendment>^`
   - `(cd <worktree> && roundfix reconcile <run-id> --carry-forward)`, an executable line scoped to the worktree
   - `git -C <worktree> cherry-pick <amendment> …`
   - `roundfix deliver retry <slug>`
3. MUST leave both texts byte-identical to today's without amendments, including when a Task commit touched a moved input.
4. MUST POSIX single-quote every printed argument that comes from state (the worktree path, the Run ID, commit IDs and the slug) through one helper, so a path with spaces or shell metacharacters stays one literal argument. It MUST NOT run any of the printed commands, change what carry-forward proves, or change the exit code of a refused retry.
5. MUST rename or remove no top-level test, change no exported function signature, and put the new tests in `internal/cli/deliver_retry_amendment_test.go` over real Git repositories and the package's Run Database fixtures.
6. MUST describe the Delivery Retry refusal and its recovery commands in `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md` with the phrase `amended by`, then run `make skills-sync`.

## Subtasks

- [ ] Detect amending commits for a moved-input refusal.
- [ ] Name them in the refusal and print the ordered recovery.
- [ ] Describe the behavior in the guide and the skill, then sync the mirror.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A Run with a settled Task, followed by a non-Task commit that edits the Spec's `_prd.md`, refuses. The reason names that commit, and the next action lists the five commands in order.
- [ ] A moved input changed by a commit carrying a `Roundfix-Task` trailer keeps today's reason and next action.
- [ ] A refusal that is not a moved-input refusal keeps today's reason and next action.
- [ ] The guide, the skill and its mirror carry `amended by`, and `make skills-sync-check` passes.

## Context

- interface: `internal/cli/carryforward.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_retry_amendment_test.go`
- instruction: `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRetryRefusalNamesTheAmendingCommitAndTheRecovery|TestRetryRefusalByATaskCommitKeepsTodaysText|TestRetryRefusalWithoutMovedInputsKeepsTodaysText)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryRefusalNamesTheAmendingCommitAndTheRecovery TestRetryRefusalByATaskCommitKeepsTodaysText TestRetryRefusalWithoutMovedInputsKeepsTodaysText; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the three named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands.md|amended by" ".agents/skills/roundfix/SKILL.md|amended by" "skills/roundfix/SKILL.md|amended by"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && make skills-sync-check` — expected: exit 0; before this Task the phrase `amended by` is in neither the guide nor the skill, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; Core Feature 2; Success Metric 2
- [_techspec.md](_techspec.md) — Retry hint; API Contract 2; Testing Approach 2; Build Order 2
- [references/2026-09-30-retry-refusal-does-not-say-carry-forward-first.md](references/2026-09-30-retry-refusal-does-not-say-carry-forward-first.md)
- ADR-0170; ADR-0053; ADR-0158

## Result
