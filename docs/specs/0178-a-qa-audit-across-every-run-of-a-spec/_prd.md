---
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: active
created: 2026-09-28
surfaces: [backend, docs]
---

# A QA audit across every Run of a Spec

The terminal QA gate must audit what a Spec delivers, not what its last Run
happened to commit, and the auditor identity a QA Report records must mean one
thing and be acted on. Two defects break that today:

- The mechanical stage reads Task commits from `git log <Run start head>..HEAD`
  (`mechanicalTaskCommits` in `internal/daemon/task_engine.go`) and keeps only
  the newest commit of each Task. Specs 0170 to 0176 each took two to four Runs,
  so Task commits of earlier Runs are the common case, and they are already
  ancestors of the Run start head. Spec 0172's authorization audit caught its
  out-of-grant file only because that file was in the current Run. The same
  request reads each grant at the Run start head, so a grant widened on the Spec
  branch inside the consuming pull request authorizes the Tasks after it, which
  the documented rule forbids.
- `auditor_staleness` compares the Daemon's build commit with the audited head.
  In a Roundfix self-audit the Daemon is built from main and the audited head
  carries the candidate's commits, so the field reads `stale` by construction,
  and nothing reads it. The QA Agent runs public-CLI rows with whichever
  `roundfix` it finds, and four reports of 2026-09-25 to 2026-09-28 rewrote
  `auditing_binary` to name that build, so the field named two binaries.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; the Delivery Base is
  a Git revision, and the new `user_flow_binary` key carries a build identity
  that `roundfix --version` already prints. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads and files only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0096 makes the mechanical stage
  compute the tooling authorization's bounded files against each Task commit's
  changed paths, withhold the Agent Session on a blocking fact, and never make a
  verdict more permissive; ADR-0130 judges only governed paths; ADR-0132 makes a
  refused gate record the refusal as its terminal row; ADR-0117 places a check
  with the stage that can produce the defect; ADR-0014 and ADR-0057 keep
  Verification and Task status Daemon-owned; ADR-0015 reads the verdict from the
  QA Report and ADR-0080 owns verdict semantics; ADR-0091 keeps the gate a Task
  node; ADR-0023 runs a Spec Run in its own Run Worktree; ADR-0138 commits each
  Task on the user's non-default branch, which is what makes the repository
  default branch the delivery target; ADR-0089 makes code under test take its
  environment explicitly, so the Auditing Binary reaches the Daemon as a
  dependency. ADR-0020, ADR-0038, ADR-0127, ADR-0141 and ADR-0160 cite ADR-0014
  or ADR-0023 but govern the Agent prompt result, the Verification repair,
  process residue, review Runs and the red repository gate's entry, which this
  Spec does not touch, so they do not apply; ADR-0056 and ADR-0159, which cite
  ADR-0038, govern Verification capacity and independent Verification and do
  not apply either. ADR-0093
  checks Spec consistency by citation, ADR-0104 accepts on evidence a Spec did
  not author, ADR-0155 makes the `qa` Task declare the matrix and ADR-0156
  makes a declared promise name a consuming Task. This Spec's gate is bound by
  ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized continuing with the
  next wave after the v0.18.0 release in chat on 2026-09-28 ("Após o release,
  pode continuar com as implementações da onda seguinte"); the skill files ride
  the standing grant of 2026-09-18 for keeping the shipped skills true to the
  CLI, recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- The mechanical stage audits every Task commit a Spec accumulated across all
  its Runs.
- A grant authorizes a Task commit only once it is in the delivery target's
  ancestry.
- The gate audits by command only what the mechanical stage did not audit.
- `auditor_staleness` answers a question whose answer is not fixed by
  construction, and a `stale` auditor is published as a warning while the gate
  proceeds.
- `auditing_binary` names one binary, and the binary the public rows ran is
  recorded and checked.

## Core Features

1. **Every Run's Task commits reach the audit.** The Daemon resolves the
   Delivery Base, the merge base of the audited head and the repository default
   branch, and the mechanical stage reads every commit of every non-QA Task of
   the Spec between it and the audited head, not only the newest commit of each
   Task. Each grant is read at the Delivery Base, as it stands before the
   consuming delivery. When the default branch cannot be determined, the stage
   keeps today's Run start head and records the skip
   `Task commits of earlier Runs`.
2. **The audit table names each commit.** The report's authorization audit
   table gains a `Commit` column, and the qa-gate skill audits by command only a
   Task commit that table does not list, or every Task commit when the skip is
   recorded.
3. **Staleness is measured against the Delivery Base, and stale warns.** The
   Daemon compares its build commit with the Delivery Base instead of the
   audited head. A `stale` auditor publishes one `daemon.qa` Run Event naming
   the staleness line and the rebuild instruction, and the seeded report records
   the same warning; the gate then proceeds as for `current` and `unknown`, with
   repository Verification and the Agent turn.
4. **The report names the user-flow binary.** `auditing_binary` and
   `auditor_staleness` stay as the Daemon seeded them, and a `pass` or `partial`
   that rewrites either is not accepted. The QA Agent records the binary its
   public rows ran as `user_flow_binary`; in a Roundfix self-audit, settlement
   refuses a `pass` or `partial` whose `user_flow_binary` is missing or was not
   built from the audited head, and the QA prompt says so.

## Non-Goals / Out of Scope

- Changing `QAReportEligibility`, so `roundfix archive`, `roundfix settle` and
  `roundfix qa-report accept` keep today's acceptance and every archived report
  stays eligible.
- Rebuilding or restarting the Daemon automatically when it is stale.
- Changing how a Task commit is authorized once its grant is read: the bounded
  files, the self-approval rule and sanctioned regeneration are unchanged.
- Rewriting archived QA Reports or any archived Spec.

## Success Metrics

1. In a disposable repository where task_01 was committed by an earlier Run and
   task_02 by the current Run, a governed out-of-grant change in task_01 yields a
   `QA-AUTH-PATHS` finding naming task_01's commit; with the Run start head as
   the range start, as today, it yields none.
2. A Task with two commits, whose older commit changes an ungranted governed
   path, yields a finding naming the older commit.
3. A grant widened on the Spec branch after the Delivery Base does not
   authorize the widened path; the same widening landed on the default branch
   and merged into the Spec branch does.
4. A Daemon built from the Delivery Base seeds `auditor_staleness` as
   `current: commit ancestry: build commit does not predate the delivery base`
   while the audited head carries candidate commits and publishes no staleness
   warning; a Daemon built from an ancestor of the Delivery Base publishes
   exactly one `daemon.qa` event with phase `auditor_staleness`, its seeded
   report carries the `## Auditor staleness warning` section, and repository
   Verification and the Agent turn still run.
5. In a self-audit, a `pass` whose `user_flow_binary` names the audited head
   settles `completed`; a missing `user_flow_binary`, one built from another
   commit, and a rewritten `auditing_binary` each settle the QA Task `failed`
   with a reason naming the cause.

## Recorded limits

- A report that deletes the seeded `auditing_binary` or `auditor_staleness`
  line is not refused; only a changed value is. Reports written in code and by
  earlier writers carry neither line, and the qa-gate skill forbids deleting
  them.
- A self-audit is recognized only when the Daemon carries a build commit, which
  `make build` stamps; a released Roundfix binary auditing the Roundfix
  repository skips the `user_flow_binary` check.
- A repository whose default branch cannot be determined keeps the Run start
  head as the range start and records the skip `Task commits of earlier Runs`
  instead of blocking the gate; the qa-gate skill then audits every Task commit
  by command.
- A stale Daemon is not refused. The warning is recorded and published, and the
  operator rebuilds before the next Spec.

## Decisions

- **The Delivery Base is the merge base with the default branch.** A Spec's
  Task commits land on a non-default work branch (ADR-0138) that is delivered
  to the repository default branch, so every commit between that merge base and
  the audited head is the Spec's delivery, whichever Run made it. The
  remote-tracking ref is read first because it is the pull request base.
- **Read each grant at the Delivery Base.** `docs/agents/spec-routing.md`
  requires a new or widened grant to land in the delivery target's ancestry
  before the consuming delivery. Reading it at the Run start head let Spec
  0172's in-PR widening pass. A mid-Spec widening now lands through its own
  pull request, as Spec 0175's phrase-check amendment did (#264), and reaches
  the Spec branch by merging the default branch.
- **The auditor is the delivery target's build, by design.** The Daemon judges
  a candidate with code the delivery target already accepted, never with code
  the candidate changes, so its staleness is measured against the Delivery Base.
  `stale` then means the Daemon lacks rules the delivery target holds, which a
  rebuild fixes.
- **A stale auditor warns; it does not refuse.** A multi-item
  `roundfix deliver` queue keeps one binary while each later item starts from a
  newer main, so a refusal would stop every item after the first at QA and
  defeat the queue. The protection that matters for the product is that the
  public rows ran a build of the audited head, which the `user_flow_binary`
  check enforces; the Daemon's own staleness is published so an operator
  rebuilds before the next Spec.
- **Two binaries, two keys.** The QA Agent's public rows need the candidate's
  build and the mechanical stage needs the delivery target's; one key cannot
  name both, which is why the reports disagreed.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: an earlier Run's commit left
out, an older commit of a Task left out, an in-PR grant accepted, a self-audit
reading `stale`, a stale Daemon that stays silent or stops the gate, or a rewritten
auditor field accepted would each pass a happy-path test.

The outside-evidence row rests on history this Spec did not write: Spec 0172's
retained commits `256ad156..ebac7cbe` and the QA Reports of Specs 0171, 0172 and
0174. Replayed in a disposable clone, the built tree's commit selection must
list every Task commit of Spec 0172 that changes a governed path (the
authorization audit selects only those, by design), including task_01, task_02
and task_04 from before its `chore: merge main into 0172` commit. The staleness those reports recorded
as `stale` is recomputed against each audited head's Delivery Base: measured on
2026-09-28, the 0171 heads have Delivery Base `256ad156` and the 0174 head
`3e6cfe9f` has `f0faa780`, each equal to the Daemon build that audited it, so
both read `current`; the 0172 head `ebac7cbe` has Delivery Base `fc296df0`,
which the branch merged after the Daemon `256ad156` was built, so it still reads
`stale`. The row records all three results. When those objects are no longer in
the repository, it is recorded as blocked with that reason.

## Research basis

Spec 0138 routed the commit-range gap to a later Spec
(`references/2026-09-14-the-mechanical-stage-audits-only-the-current-runs-task-commits.md`);
Spec 0119's report `qa-report-2026-09-10-01.md` reconstructed thirteen Task
commits by hand and found a real out-of-grant change. The staleness gap was
triaged from the secondbrain capture
`inbox/roundfix/2026-09-19-qa-gate-passes-on-self-declared-stale-binary.md`
(`references/2026-09-25-the-auditor-staleness-signal-is-recorded-but-never-acted-on.md`).
The adopted sources are indexed in [references/_index.md](references/_index.md).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
