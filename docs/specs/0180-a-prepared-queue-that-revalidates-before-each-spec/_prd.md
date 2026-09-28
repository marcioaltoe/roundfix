---
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# A prepared queue that revalidates before each Spec

`roundfix deliver start <slug>...` records whatever slugs it is given and
starts working. Nothing tells the operator beforehand which of those Specs is
actually approved to run, and nothing re-checks a later Spec against the main it
will start from once earlier items have merged. The queue has no limit of its
own, and it presents every parked item at once. The owned `implement-spec`
skill still tells the Supervisor to run its own Task loop. Retired Spec 0127
carried these gaps as Core Features 1, 4, 5, 6 and 8. They are now four
Backlog Entries, and the efficiency waves hit each of them:

- Spec 0175 was authored against main `0160f70a` and queued after Spec 0173.
  Spec 0173's squash merge (`6fac37ea`) changed `internal/delivery/engine.go`,
  `internal/cli/deliver_workflow.go` and `internal/cli/carryforward.go`, and
  Spec 0175's Tasks declare all three as `interface:`. The queue would run
  Spec 0175 on the new main without noticing that its premise had moved.
- Wave-one and wave-two Specs repeatedly failed on governed-path declarations
  and on signature fallout from earlier merges. These are defects
  `roundfix spec check --strict` reports, but nothing runs it on the main an
  item actually starts from.
- A queue item whose authorization lacks `push`, `pull_request` or `merge`
  runs its whole Task Graph and only then parks as `unauthorized`.
- `deliver status` lists every parked item with equal weight, and no bound
  stops a queue that runs past the time the operator intended.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a Delivery Queue
  item keeps its Spec slug, item branch, Run ID and merge commit, and new
  blockers are plain strings like the existing ones. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential is read and no network call is added.
  Publication keeps using the existing GitHub CLI boundary. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 and ADR-0094 make the Spec
  Consistency Check read citations only and run at every stage, which the
  preparation and the revalidation reuse unchanged. ADR-0117 checks a defect
  at the stage that can produce it, and an earlier item's merge is what
  produces a later item's stale premise. ADR-0130 keeps a bounded path
  governed, and ADR-0160 names the frozen authorization as the only authority;
  the preparation reads that authorization and grants nothing. ADR-0052 makes
  completion compare-and-set, the model for the retry limit checked in the
  same transaction as the guarded item transition. ADR-0044 reclaims an owner
  record only on proven death, which the retry keeps. ADR-0139 keeps one
  Active Run per work target, and the queue stays one Spec at a time.
  ADR-0137 and ADR-0158 keep the Run Budget a per-Run bound, so the queue
  deadline is a separate limit and never reinterprets it. ADR-0014 and
  ADR-0057 keep Verification and Task status Daemon-owned, which the rewritten
  `implement-spec` skill defers to. ADR-0153 keeps the pre-PR review an
  explicit provider policy, unchanged by this Spec. ADR-0020, ADR-0038 and
  ADR-0127 cite ADR-0014 but govern the Agent prompt result, the Verification
  repair bound and process residue; ADR-0053 cites ADR-0052 but governs
  terminal Run Worktree reconciliation; and ADR-0097 carries a QA row forward,
  not a queue item; ADR-0056 and ADR-0159 cite ADR-0038 but govern
  Verification capacity and independent Verification. This Spec changes none
  of them, so they do not apply. This
  Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0104, ADR-0155 and
  ADR-0156. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized continuing with
  the next wave after the v0.18.0 release in chat on 2026-09-28; the Go test
  file rides the standing grant of 2026-09-21 for governed source, the
  Roundfix skill files ride the standing grant of 2026-09-18 for keeping the
  shipped skills true to the CLI, and the maintainer expressly authorized the
  `implement-spec` rewrite on 2026-09-28, all recorded in
  [_authorization.md](_authorization.md);
  bounded files: `internal/cli/cli_test.go`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/implement-spec/SKILL.md`, `skills/implement-spec/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`, `make baseline-digests`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- Before a queue runs, the operator sees which Specs are approved to run, what
  blocks the others and which files a later Spec shares with an earlier one.
- A queue never records a Spec that has no delivery authority.
- Each queued Spec is re-checked, before its first Run, against the main it
  actually starts from.
- A queue carries its own explicit limits and enforces them.
- The maintainer faces one queue decision at a time, and only an explicit
  command answers it.
- The `implement-spec` entry point hands implementation to Roundfix and keeps
  no loop of its own.

## Core Features

1. **A read-only Delivery Plan.** `roundfix deliver plan [<slug>...]` reports,
   for each named Spec (every active Spec when none is named), whether it is
   approved to run and why not. It is approved when its authorization grants
   `implement`, `commit`, `push`, `pull_request` and `merge` and its strict
   Spec Consistency Check reports no error. The plan also names the production
   Go files each Spec declares as `interface:` that an earlier Spec in the
   given order declares too. It lists the open Backlog Entries, unresolved
   Findings and repository inbox notes as intent that is not approved to run.
   It writes nothing, and `--json` prints the same facts.
2. **A queue records only authorized Specs.** `roundfix deliver start` refuses
   the whole queue when any slug's authorization does not grant all five
   delivery operations. It names each such Spec and its reason, and records
   nothing. A strict consistency finding does not refuse the start: the
   Delivery Plan reports it, and the item is re-checked against its own
   starting main.
3. **Delivery Revalidation before the first Run.** After an item's worktree is
   created from the refreshed default branch, and before its Run starts, the
   item passes the strict Spec Consistency Check there. A finding parks the
   item as `revalidation-failed: <codes>`, and a Delivery Retry of it re-runs
   the check and refuses while findings remain. The revalidation also compares
   the production Go files the Spec's Tasks declare as `interface:` with the
   merge commits of earlier items of the same queue. An overlap never stops the
   item: it records a `premise-changed` warning on the item, naming the files
   and the merge, prints it to the delivery console log, and the item continues
   to its Run. `deliver status` shows the warning.
4. **Explicit, enforced queue limits.** `deliver start` accepts
   `--max-duration <duration>` and `--max-retries <n>`. They are recorded with
   the queue as a deadline and a per-item retry limit, and every omitted limit
   is recorded and printed as `none`. After the deadline no new item starts:
   a queued item parks as `queue-deadline` without a worktree, while an item
   already past `queued` continues. A retry beyond the limit is refused.
   Whole-Spec concurrency is one, and spending is reported as not measured.
   `deliver start` and `deliver status` print every limit.
5. **One Pending Question.** `deliver status` presents one Pending Question:
   the lowest-position parked item, its blocker and the explicit action that
   answers it, plus how many parked items wait behind it. Only an operator
   command changes a parked item. No owner pass, elapsed time or recommended
   option answers it.
6. **`implement-spec` delegates to Roundfix.** The owned skill prepares with
   `roundfix deliver plan`. It hands the Spec to `roundfix implement --spec`,
   or a merge-through sequence to `roundfix deliver start`, and monitors
   through `roundfix deliver status`. It asks the maintainer only the Pending
   Question, and it no longer tells the Supervisor to run Tasks, write code
   or tests, or run the QA gate itself.

## Non-Goals / Out of Scope

- A second implementation loop, a window script or any Supervisor-written code
  or tests.
- Parallel whole-Spec delivery or a concurrency setting.
- Measuring or capping API or subscription spend.
- Refusing a retry because another parked item is the Pending Question.
- Detecting a changed signature in a file the Spec does not declare, or a
  merge that did not come from the same queue.
- Inventorying the Secondbrain inbox, which Roundfix has no configured location
  for.
- Changing the Run Window, the Run Budget, the archive boundary or the pre-PR
  review policy.

## Success Metrics

1. In a disposable repository holding one approved Spec and one Spec whose
   authorization lacks `merge`, the built `roundfix deliver plan` exits `1`,
   reports the first `approved` and the second `blocked` with the missing
   operation. It creates no Run Database and leaves the checkout's HEAD and
   status unchanged. `roundfix deliver start` with both slugs exits `2` and
   records no queue.
2. Engine and workflow tests show that an item whose starting main makes its
   strict consistency check fail parks as `revalidation-failed` before any Run.
   An item whose declared production Go file an earlier item's merge commit
   changed reaches `running` with a `premise-changed` warning naming that file
   and that merge, and an item with no overlap reaches `running` with no
   warning. Against this repository's history, the paths Spec 0173's merge
   commit `6fac37ea` changed intersect Spec 0175's declared production Go
   files in exactly `internal/cli/carryforward.go`,
   `internal/cli/deliver_workflow.go` and `internal/delivery/engine.go`.
3. With `--max-duration 1ns`, the built binary's queue parks its first item as
   `queue-deadline` without a worktree. `deliver retry` refuses it with exit
   `2`. With `--max-retries 1`, a second retry of the same parked item exits
   `2` and leaves the item unchanged.
4. `roundfix deliver status` prints exactly one `Pending question:` line for a
   queue with two parked items, names the lower-position item and reports one
   item waiting. The question is unchanged after another owner pass and after
   the clock advances.
5. The `implement-spec` skill and its mirror name `roundfix deliver plan`,
   `roundfix implement --spec` and `roundfix deliver start`. Neither tells the
   Supervisor to run the implement-task cycle, and `make skills-sync-check`
   exits `0`.

## Recorded limits

- The premise check only sees merges of earlier items in the same queue and
  only files a Spec declares as production Go `interface:` paths. A signature
  changed elsewhere still reaches the Run, where the Task's own Verification is
  the detector.
- A `revalidation-failed` item is amended on its item branch, in the worktree
  `deliver status` prints, because the item was cut from the main that failed
  the check. An amendment merged to main afterwards does not reach that branch.
- A `premise-changed` warning is a record, not a question: nothing answers or
  clears it, and it stays on the item until the queue is replaced.
- Spending is not measured, so no spending limit exists; `deliver status`
  says so instead of implying zero.
- The Pending Question orders the operator's decisions. It does not stop the
  owner from advancing independent items, and it does not refuse a retry of
  another parked item.
- A limit is fixed when the queue is recorded; changing one means recording a
  new queue once the current one has no unfinished item.

## Decisions

- **Start refuses on authority, not on consistency.** An authorization is
  read where the Spec lives and cannot be supplied by an earlier merge in
  normal authoring, so refusing it at start saves a whole Run. A consistency
  finding can be fixed or caused by an earlier item's merge, so the authoritative
  check is the revalidation on the item's own starting main.
- **Revalidate inside the item worktree, after creating it.** The item worktree
  is the only place where the main the item starts from exists. Parking there
  keeps the item branch and worktree, so a Delivery Retry can resume the item.
- **A changed premise warns; it does not park.** Two queued Specs that edit
  the same production Go file are common. Spec 0175 queued after Spec 0173 is
  one such pair, and the overlap did not break it. A park would force a manual
  `deliver retry` on most multi-Spec queues. The warning keeps the fact durable
  and visible for the reviewer, while the Task's own Verification remains the
  detector for a real break.
- **Production Go only for the premise check.** Guides and tests are shared by
  almost every Spec. A merge that touches them is re-read by the Agent and
  checked by the Task's own Verification. A declared production Go file changed
  under a Spec is where signature fallout lives.
- **Derive the Pending Question; do not store it.** It is the lowest-position
  parked item, read from the item records that already persist, so it cannot
  disagree with them after a restart.
- **Limits are recorded with the queue, in one new schema version.** Limits
  enforced by a detached owner have to survive restart. They live next to the
  owner record, and the retry count lives on the item, where the retry
  transaction already reads. The same version gives each item a warning field,
  so a `premise-changed` warning survives restart and every later stage.
- **Omitted limits are explicit `none`.** Requiring both flags would break
  every existing `deliver start` invocation. The queue therefore records and
  prints `none` and never infers a limit from the Run Window or the Run Budget.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight: a plan that writes
anything, a start that records an unauthorized Spec, a Run started for an item
whose starting main fails its check, a premise change that goes unnamed or
stops the item, an
item started after the deadline, a retry past its limit, a second question
presented at once, and a skill that still runs Tasks would each pass a
happy-path test.

The outside-evidence row rests on history this Spec did not produce. Spec
0173's squash merge `6fac37ea` and Spec 0175's Task declarations, as
committed at `b92aefda`, were written by other sessions. The QA gate measures
their intersection with `git diff --name-only 6fac37ea^ 6fac37ea` and records
it beside the `premise-changed` warning the revalidation records for the same
shape.

## Research basis

The four adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md). Each was carried from retired
Spec 0127 (`docs/history/specs/0127-durable-unattended-spec-workflow/_prd.md`,
Core Features 1, 4, 5, 6 and 8) when that portfolio Spec was retired on
2026-09-28. Its TechSpec's "Limits and cancellation" section proposed a
serial queue, no paid call without a recorded ceiling and a bounded number of
corrections. It recorded them as proposals, not answers, and this Spec records
those proposals as limits the operator states when starting a queue.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
