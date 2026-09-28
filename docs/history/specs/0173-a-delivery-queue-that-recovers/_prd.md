---
spec: 0173-a-delivery-queue-that-recovers
status: archived
created: 2026-09-28
surfaces: [backend, cli, docs]
archived: "2026-09-28"
source_slug: 0173-a-delivery-queue-that-recovers
---


# A delivery queue that recovers

`roundfix deliver` exists so an operator does not orchestrate a Spec from Run
to merge by hand. In the first efficiency wave all three queued Specs parked
`run-unresolved`, and the queue had no way back: the owner loop skips every
parked item, `deliver resume` only restarts the owner, and `deliver start`
records a new queue. Each Spec was finished by hand: fast-forward the item
branch to the Run Branch, add a corrective Task, implement, review with an
explicit base, archive, gate, open the pull request and merge. Four defects
make that recovery necessary or make dispatch fail for reasons the operator
cannot see:

- A parked Delivery Queue item can never re-enter the queue, and the item does
  not even record the Run that parked it: the live Run Database holds all three
  wave-one items as `parked`, `run-unresolved`, with an empty Run ID.
- Task Carry-Forward refused Spec 0170's settled Tasks because `task_01` had
  changed `CONTEXT.md`, an input `task_04` declares. Both Tasks belong to the
  same Run. ADR-0026 integrates Tasks in completion order, and `task_04`
  settled first, but carry-forward stages the Tasks in Task Graph order, so it
  measured `task_04`'s input against a state its settlement never saw.
- The Daemon leaves `.roundfixrc.yml` out of every commit with no event and no
  reason, so a Task authorized to change Project Config settles completed and
  its change stays behind in the Run Worktree until the QA gate notices.
- Implement Preflight proves every configured tuple, fallbacks included, with a
  30-second setup limit and no retry. Under load an unused fallback's proof
  timed out and refused four dispatches, with advice to reconfigure a profile
  that works.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a Delivery Queue
  item keeps its Spec slug, item branch and Run ID. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git, the Run Database
  and ACP Runtime subprocesses only; no credential is read and no network call
  is added. Publication keeps using the existing GitHub CLI boundary. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0026 integrates settled Tasks in
  completion order, which carry-forward now follows; ADR-0053 keeps
  reconciliation proof-based, and Task Carry-Forward keeps every proof and its
  whole-set refusal; ADR-0158 keeps a BudgetExceeded Run recoverable through
  carry-forward; ADR-0052 makes completion compare-and-set, the model for the
  retry's guarded item transition; ADR-0044 reclaims an owner record only on
  proven death, which the retry's owner hand-off reuses; ADR-0139 keeps one
  Active Run per work target; ADR-0138 keeps one commit per verified Task;
  ADR-0014 and ADR-0057 keep Verification and Task status Daemon-owned, and
  carry-forward stays outside a Run; ADR-0153 keeps the pre-PR review an
  explicit provider policy, so a re-review runs the configured policy;
  ADR-0157 makes tracked source that cannot be committed a Task failure;
  ADR-0160 names the frozen authorization as the Daemon's only authority, which
  decides whether Project Config may be committed; ADR-0002 defines Project
  Config; ADR-0130 keeps a bounded path governed; ADR-0107 has Preflight prove
  every configured tuple and ADR-0050 substitutes no fallback before a Run
  exists, both unchanged by the retry; ADR-0140 proves the exact advertised
  tuple, and ADR-0114 says opening an Agent Session is not Agent work, so a
  retried proof opens a fresh session without starting Agent work.
  ADR-0097 carries a QA row forward, not a Task, so it does not apply. ADR-0020,
  ADR-0038 and ADR-0127 cite ADR-0014 but govern the Agent prompt result, the
  Verification repair bound and process residue; ADR-0056 and ADR-0159 cite
  ADR-0038 but govern Verification capacity and independent Verification; and
  ADR-0069 cites ADR-0050 but governs Baseline semantic analysis. This Spec
  changes none of them, so they do not apply. ADR-0093
  checks Spec consistency by citation, ADR-0104 accepts on evidence a Spec did
  not author, ADR-0155 makes the `qa` Task declare the matrix and ADR-0156
  makes a declared promise name a consuming Task. This Spec's gate is bound by
  ADR-0080, ADR-0091, ADR-0096 and ADR-0117. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved Onda 2 of the
  efficiency sequence in chat on 2026-09-28; the Go test file rides the
  standing grant of 2026-09-21 for governed source and the skill files ride the
  standing grant of 2026-09-18 for keeping the shipped skills true to the CLI,
  recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A parked item returns to the queue with one command, from the stage its
  recorded evidence names, without leaving `roundfix deliver`.
- A retried Run re-runs only the Tasks that did not finish.
- Tasks of one Run never refuse carry-forward over each other's changes.
- A Task never loses Project Config silently.
- A proof that timed out under load neither blocks dispatch on its first
  timeout nor advises reconfiguring a working profile.

## Core Features

1. **A parked item can be retried.** `roundfix deliver retry <slug>` returns
   one parked item to the queue. It decides the re-entry stage from the item's
   recorded evidence, never from the blocker text alone: while the Spec is
   still active on the item branch, the item re-enters the Run when any Task is
   unfinished and the pre-PR review of the current item head otherwise; after
   the Spec was archived on the item branch, the item re-enters the repository
   gate, or the current-head checks when a pull request is recorded. It refuses
   an item that is not parked, an item whose branch is gone, and an archived
   item whose head moved, and leaves the item unchanged when it refuses.
2. **A retried Run keeps what already passed.** Before an item whose Spec is
   still active re-enters the queue, the retry hands the settled Tasks of the
   item's Run back to the item branch through Task Carry-Forward, so only the
   unfinished Tasks run again. A parked `run-unresolved` item records the Run
   that parked it; an item parked before this Spec uses the newest Implement
   Run of its Spec on its item branch. A carry-forward that refuses the set
   refuses the retry.
3. **The retry reaches an owner.** The retry works whether or not the queue's
   owner is running: a live owner picks the item up before it releases the
   queue, a stale owner record is reclaimed, and with no owner the retry starts
   one.
4. **Carry-forward follows the Run's integration order.** Task Carry-Forward
   stages and proves candidates in the order the Run integrated their
   settlement commits, so an input changed by a Task the Run integrated earlier
   is not a moved input; an input changed on the checkout still refuses the
   whole set.
5. **Project Config is committed on authority or refused aloud.** A Task that
   changes `.roundfixrc.yml` commits it when the Spec's frozen authorization
   bounds that file, and otherwise settles failed with a lost-output reason
   naming it. A Batch or QA Report commit still never stages it, and reports
   the exclusion as a dropped path.
6. **A timed-out proof is retried once and called temporary.** A profile proof
   whose own setup deadline expires is retried once with a fresh disposable
   session. A second timeout fails with classification `temporary` and advice
   to rerun, not to reconfigure. A rejected selection and a cancelled command
   are never retried.

## Non-Goals / Out of Scope

- Retrying several items at once, or an automatic retry without an operator
  command.
- Re-running a Task that already passed, or carrying part of a refused set.
- Moving the head of an archived candidate, re-archiving a Spec, or replaying a
  recorded push, pull request creation or merge.
- A new Delivery Queue stage, a new blocker or a Run Database schema change.
- Staging Project Config from a Batch or QA Report commit.
- Deferring a fallback's proof until it is needed, or changing the 30-second
  setup limit.

## Success Metrics

1. Against a Run whose Tasks settled out of Task Graph order, one sharing an
   input the other changed, `roundfix reconcile <run> --carry-forward` exits
   `0` and carries both Tasks, while a checkout change to that input still
   exits `2` and leaves HEAD unchanged.
2. A `run-unresolved` item retried with `roundfix deliver retry <slug>` runs
   its next Implement Run with every Task the previous Run completed already
   `completed` on the item branch, and reaches the stage the table in Core
   Feature 1 names for each recorded state.
3. `roundfix deliver retry <slug>` exits `0` and either names the live owner it
   handed the item to or starts a detached owner; it exits `2` and leaves the
   item byte-for-byte unchanged for an item that is not parked, a missing item
   branch, a moved archived head or a refused carry-forward.
4. A Task that changes `.roundfixrc.yml` produces a commit containing it when
   the frozen authorization bounds the file, and otherwise settles `failed`
   with reason `Task commit lost output: .roundfixrc.yml (Project Config
   outside the Spec's authorization)` and a dropped-path Run Event.
5. A profile proof that times out once and then passes leaves `roundfix
   profiles validate` at exit `0` after two attempts; one that times out twice
   exits `2` with classification `temporary`; a rejection is attempted once.

## Recorded limits

- A retry that meets the same blocker parks the item again: retrying an
  `unauthorized` item, whose authorization is read from the archived
  candidate's parent, parks it `unauthorized` again, because the archived head
  cannot move.
- Carry-forward still requires a clean item worktree. An untracked file the
  worktree copy list places there without an ignore rule refuses the carry, as
  it refuses `reconcile --carry-forward` today.

## Decisions

- **Derive the stage from evidence, not from the blocker.** The item's stage
  becomes `parked`, so the stage it parked at is not stored, and one blocker
  (`review-stale`) parks from four stages. Whether the Spec is archived on the
  item branch, whether any Task is unfinished, and whether a pull request is
  recorded decide re-entry without a schema change.
- **Carry forward, never fast-forward.** A Run Branch after a failed QA gate
  holds the QA Report commit with a `fail` verdict, and an operator may have
  committed a corrective Task on the item branch; a fast-forward would bring
  the first and is impossible after the second. Carry-forward carries exactly
  the proved Tasks with their provenance.
- **Follow the integration order rather than weaken the proof.** Staging in the
  order the Run integrated its commits makes each candidate's staged state its
  settlement parent's state, so the proof "declared inputs unmoved since
  settlement" stays as strict and becomes true for sibling Tasks.
- **The owner releases the queue only when it is idle.** Releasing and checking
  for an advanceable item happen in one transaction, and the retry's item
  transition reads the owner in the same transaction, so a retry can never land
  between an owner's last look and its exit.
- **Commit Project Config only on a named grant.** The frozen authorization
  already bounds governed paths; a Spec that names `.roundfixrc.yml` there has
  the maintainer's authority to change it, and every other Task fails loudly.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a retry of an item that is not
parked, a retry that re-runs a completed Task, an owner that exits past a
retried item, a sibling input that still refuses, a checkout input that no
longer refuses, a silent Project Config drop, or a retried rejection would each
pass a happy-path test.

The outside-evidence rows rest on records this Spec did not produce: the Run
Event Journal of Run `run_20260925T182023Z_7e3e8c4b8bc65d93` in the live Run
Database, read through the built `roundfix events`, shows `task_04` settling
before `task_01`, the order carry-forward must follow; and the live Delivery
Queue rows, read through the built `roundfix deliver status` and a read-only
query, show the three wave-one items parked `run-unresolved` with no Run ID.

## Research basis

The four adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md). The Project Config drop was
captured in the secondbrain inbox
(`inbox/roundfix/_triaged/2026-09-24-daemon-silently-drops-project-config-from-task-commits.md`)
from Spec 0155's QA finding F-001. The proof timeouts were observed on
2026-09-24: four dispatches refused with `adapter error: context deadline
exceeded` for the `claude`/`opus`/`high` fallback.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
