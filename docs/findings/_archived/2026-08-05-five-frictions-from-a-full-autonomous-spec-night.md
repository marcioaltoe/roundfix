---
status: done
created_at: 2026-08-05
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# 2026-08-05 — Five frictions from a full autonomous spec night

Source: a real overnight Supervisor session in the `fiscus` repository driving spec
`0001-fundacao-de-contextos` (10-task graph with authored QA gate, docs-heavy) through three
Implement Runs (`run_20260805T011231Z_fdc643c4319c4d39` Failed,
`run_20260805T011820Z_ca873dfebd68b418` Unresolved after a correct QA fail,
`run_20260805T025316Z_a62a0355b8033937` Clean) and one watch Run
(`run_20260805T030642Z_708d41a5ad9ed303` MaxRoundsReached). The spec shipped; every friction
below cost Supervisor diagnosis time or left a Run in a state only trial-and-error resolved.

## 1. An invalidated QA gate refuses preflight without printing the recovery

What was observed: after the gate failed with a real finding, the Supervisor authored a
corrective `task_11`, added it to the gate's `needs`, and re-ran implement. Preflight refused:

```
validate qa gate: QA gate result is invalidated for Task "task_10" because these
dependencies are not completed: task_11
```

The message names the state but not the exit. The working recovery — reset the gate task's
stale `status: failed` to `pending` by hand — was discovered by experiment, and hand-editing a
Daemon-owned status contradicts the documented ownership model.

Suggested behavior: when the loader detects a gate whose result is invalidated by a graph
insertion, either auto-reset the gate to `pending` (journaling the transition, since the
insertion is already "named, never silent"), or have preflight print the exact remediation as
its `next:` action.

## 2. Doctor passes on a machine where the first Daemon commit cannot succeed

What was observed: `roundfix doctor` reported ok on every line (after the skills remediation),
yet Run 1 died three minutes in — at the Daemon's first Task commit — because the fresh Task
Worktree had no `node_modules` and the repository's commit-msg hook crashed:

```
Error: Cannot find module "@commitlint/config-conventional" from
"/Users/marcio/.roundfix/worktrees/fiscus-2f7556b8/run_...task_01"
```

`worktree.bootstrap` was empty. Doctor proves adapters, profiles, and skills; nothing proves a
worktree can *commit* in a repository whose hooks need installed devDependencies.

Suggested behavior: doctor (or implement preflight) warns when the repository declares commit
hooks (`core.hooksPath`, `.githooks/`, husky) and `worktree.bootstrap` is empty; ideally setup
offers a bootstrap suggestion for lockfile ecosystems it recognizes. A disposable
worktree-commit smoke as an opt-in doctor check would have caught this in minutes, during the
day.

## 3. A Daemon commit failure after settlement leaves a Task no command can recover

What was observed: in Run 1 the event stream recorded `task_01 settled completed`
(cursor 781) and then the Run failed on the commit itself (cursor 784). Resulting state: the
Task Worktree's task file said `completed`, the checkout's said `pending`, and no commit
existed. `settle` cannot act — it selects only surfaces whose task file is `failed` — so the
completed work was unreachable by any supported command; the only path was
`reconcile --apply` plus a full re-run of the Task.

Suggested behavior: a Daemon commit failure should settle the Task `failed` (with the commit
error as reason), which makes the kept worktree a valid `settle` surface and preserves the
Agent's work instead of discarding it.

## 4. The verified-evidence wait consumes Rounds and mislabels the outcome

What was observed: after Batch 001 resolved all three Review Issues and CodeRabbit re-reviewed
the pushed head with zero new findings, the watch Run looped on evidence polling — the console
shows `Reused Round 002 with 0 Review Issue(s)` at ~35s intervals between 00:23:48 and
00:26:10 — each iteration consuming a Round until `MaxRoundsReached after 6 Round(s)` with a
final report of `3 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved`. Six Rounds were
exhausted in about seven minutes, five of them without any Agent work, while
`WaitingForReviewCheck` simply hadn't flipped to `verified`.

Suggested behavior: only Rounds that dispatch Agent work should count toward `max_rounds`;
evidence waiting is already bounded by `review_timeout` and the check grace period. A Run that
ends with zero unresolved issues and only missing verified evidence should end
`CleanUnverified` (its documented meaning) rather than `MaxRoundsReached`, whose next action
misdirects the operator toward "remaining Review Issues" that do not exist.

## 5. Excluded nitpicks deadlock merges under required conversation resolution

What was observed: with `review_source.include_nitpicks: false` (the default), CodeRabbit's two
trivial threads never became Review Issues, so the watch Run correctly reported zero unresolved
— while GitHub's `required_conversation_resolution` counted two open threads and blocked the
merge. No actor owns those threads: Roundfix filtered them out by design, the Supervisor is
barred from manual thread resolution, and the human wakes up to resolve two trivia clicks.

Suggested behavior: an opt-in courtesy pass — when a watch Run reaches its terminal state with
below-threshold threads still open, post the existing idempotent Outcome Comment shape
("acknowledged, below the configured severity threshold") and resolve them, so branch
protections based on conversation resolution can pass without a human. Alternatively, document
the incompatibility: `include_nitpicks: false` + GitHub required conversation resolution cannot
reach a mergeable state autonomously.

## Companion observation

The QA Agent's managed sandbox denies loopback TCP and the Docker socket (`connect EPERM
127.0.0.1:9433` during focused checks; three QA rows settled `blocked (environment...)`), so
every database-backed acceptance row in this repository will settle environment-blocked
forever. The spec-side mitigation (`## Unreachable Acceptance` declarations) exists, but a
per-category sandbox relaxation for `qa` would remove the recurring gap — related to the
existing finding `2026-08-05-agent-full-access-passes-config-validation-and-fails-every-task.md`.
