---
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: archived
created: 2026-09-29
surfaces: [backend, cli, docs]
archived: "2026-09-29"
source_slug: 0181-gates-that-refuse-only-what-someone-can-act-on
---


# Gates that refuse only what someone can act on

Waves 2 and 3 shipped on 2026-09-28 and 2026-09-29. Three of their gate
refusals had nothing to do with the work being wrong. Each cost reruns or hand
edits:

- **A file a Task created without declaring it failed the QA scope audit.**
  Spec 0180's `task_02` changed `internal/store/delivery_test.go`, and its
  `task_04` changed `internal/cli/deliver_revalidate_test.go`. Spec 0179's
  Task 06 changed three production lock files. Every one of those Tasks was
  completed and verified. The authored QA gate refused each Spec because the
  paths were not in the Task's `## Context`, and each Spec needed a corrective
  edit that only added the line. The implementing Agent could not have
  declared them: it may edit only its Task's `## Result`. Together this cost two
  extra QA Runs.
- **The Pull Request row blocked a qualifying `partial`.** The QA gate runs
  before any Pull Request exists, so the Pull Request row every matrix carries
  is always `blocked (environment: no open Pull Request)`. A qualifying declared
  `partial` requires `rows_blocked_environment: 0`. The v0.20.0 binary refuses
  Spec 0179's report `qa-report-2026-09-29.md` with
  `rows_blocked_environment is 1; expected 0`. Getting past it took a
  boilerplate fourth Unreachable Acceptance declaration, a new QA Requirement,
  and two extra QA Runs.
- **A new ADR forced edits to every other active Spec.** `SC-ADR-RELATED`
  reports an accepted ADR that cites an ADR a Spec lists. ADR-0161, ADR-0164
  and ADR-0165 each cite widely listed ADRs. When each one landed, every other
  active Spec listing the cited ADR failed the strict check, the Delivery Plan
  and queue revalidation. PRs #268, #272 and #274 each carried hand edits to
  other Specs' Project Constraints before they could merge.

## Project Constraints

- Identifier strategy: applicable — the only new identifiers are the
  Daemon-owned Task file heading `## Recorded paths` and the Task commit event
  payload key `recorded_paths`. Each follows the existing form of its family:
  the Daemon-written `## Carry-forward provenance` heading and the snake_case
  event payload keys. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential is read and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014 and ADR-0057 make the Daemon
  the writer of a Task's settled status in the Task file it commits, and the
  recorded paths join that write. ADR-0130 keeps a Governed Path under its
  authorization, so the Daemon never records one. ADR-0080 lets an
  environment-blocked row reach `pass` with equivalent evidence. ADR-0088
  authors the QA gate into the graph, and ADR-0091 makes the QA gate a Task
  node of its own type. The gate therefore runs before any Pull Request exists,
  which is why the Pull Request row is always environment-blocked.
  ADR-0155 makes the Pull Request row a coverage source every matrix carries.
  ADR-0097 carries a QA row forward only on unmoved evidence, and it holds
  unchanged. ADR-0093 checks Spec consistency by citation and ADR-0094 skips a
  detector whose artifact is absent. Both hold for the horizon, which reads only
  citations and commits. ADR-0116 keeps `SC-CITATION-UNSUPPORTED` reading every
  cited record, and ADR-0117 checks a defect at the stage that can produce it:
  a Spec's author is the one who can account for an ADR that existed when the
  Spec was written. This Spec adds ADR-0166, ADR-0167 and ADR-0168.
  ADR-0161, ADR-0164 and ADR-0165 are cited only as the ADRs whose landing
  caused the cascade. They govern the release of a merged Spec's Runs, the
  Run Budget and a blocking review after archive, which this Spec does not
  touch, so they do not apply. ADR-0020, ADR-0038, ADR-0127 and ADR-0160 cite
  ADR-0014 but govern the Agent prompt result, the Verification repair bound,
  process residue and the red repository gate. ADR-0056 and ADR-0159 cite
  ADR-0038 but govern Verification capacity and independent Verification. This
  Spec changes none of them, so they do not apply.
  ADR-0169 and ADR-0170 (Spec 0182) govern the pre-PR review diff and Delivery
  Retry carry-forward, which this Spec does not touch, so they do not apply.
  This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0104, ADR-0155
  and ADR-0156. All hold. ADR-0176 (added 2026-09-29 with the corrective
  task_06) makes citation checks read only the authored projection. ADR-0002
  does not apply: task_03's Result names it only as a test fixture's number,
  and this Spec's own gate runs on the v0.20.0 auditor, which still reads that
  section. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer expressly authorized edits to
  the `qa-gate` and `write-tasks` skills on 2026-09-29, answering "Autorizar"
  to a structured question, only to align them with this Spec's behavior. The
  Roundfix skill files ride the standing grant of 2026-09-18 for keeping the
  shipped skills true to the CLI. Both are recorded in
  [_authorization.md](_authorization.md);
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`, `make baseline-digests`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A completed, verified Task is never refused by the QA scope audit for an
  ordinary file it changed without declaring. The file is disclosed instead.
- A Spec whose only unmet rows are genuinely declared unreachable reaches a
  qualifying `partial`, even though the pre-PR Pull Request row is
  environment-blocked.
- Landing a new ADR opens no Spec Consistency gap in an active Spec that was
  committed before the ADR existed.
- Every refusal these gates keep still fires: a Governed Path outside its
  authorization, any other environment-blocked row in a `partial`, and a
  related ADR that predates the Spec.

## Core Features

1. **The Daemon records the paths a Task changed without declaring them.**
   When the Daemon commits a completed non-QA Task, it compares the paths the
   Task commit stages with the Task's declarations. It writes every staged path
   that is not the Task file, not declared under `interface:` or `creates:`, and
   not a Governed Path into a Daemon-owned `## Recorded paths` section at the
   end of the Task file. That section rides in the same commit, and the Task
   commit event lists the same paths as `recorded_paths`. The authored
   `## Context` is never edited. The QA scope audit counts a recorded path as
   declared and names it in its scope row.
2. **The pre-PR Pull Request row never decides a qualifying `partial`.**
   `roundfix qa-report accept`, QA settlement and `roundfix archive` share one
   eligibility policy. A Results row whose status is exactly
   `blocked (environment: no open Pull Request)` and whose provenance names the
   Pull Request row no longer counts against a qualifying declared `partial`.
   Every other environment-blocked row still refuses one. The Daemon's QA
   prompt and the `qa-gate` skill say so, and no Unreachable Acceptance
   declaration is needed for that row.
3. **A related-ADR gap has a horizon.** `SC-ADR-RELATED` reports an ADR for a
   committed Spec only when the commit that added the ADR is an ancestor of the
   commit that added the Spec's `_prd.md`. An uncommitted PRD, an ADR and PRD
   in different repositories, and unreadable history all keep the full check.
   `SC-ADR-UNLISTED` and `SC-CITATION-UNSUPPORTED` are unchanged.

4. **Citation checks read only what the Spec's authors wrote.** The Spec
   citation walk behind `SC-ADR-UNLISTED` skips `qa/` and, in Task files, the
   Agent-owned `## Result` and the Daemon-owned `## Recorded paths` and
   `## Carry-forward provenance` sections (ADR-0176). A citation in authored
   text is reported exactly as before. Added on 2026-09-29 as corrective
   task_06, approved by the maintainer after this Spec's first QA gate refused
   on a fixture number in task_03's Result.

## Non-Goals / Out of Scope

- Letting the Agent edit its Task's `## Context`, or changing what Task
  Carry-Forward, parallel path reservation or the fifty-entry Context limit
  read.
- Recording, or excusing, a Governed Path.
- Relaxing any other environment-blocked row for a qualifying `partial`, or
  changing the `pass` rule of ADR-0080.
- Rewriting existing QA Reports or the frontmatter counts they carry.
- Narrowing `SC-CITATION-UNSUPPORTED` or any other consistency detector beyond
  Core Features 3 and 4. `SC-ADR-UNLISTED` was listed here until 2026-09-29.
  This Spec's own QA gate then refused on a citation in an Agent's Result, and
  the maintainer approved narrowing its walk to authored text.
- A mechanical scope detector in the QA mechanical stage.

## Success Metrics

1. In the Daemon's Task cycle over a real Git repository, a Task that declares
   `interface: internal/store/delivery.go` and also changes
   `internal/store/delivery_test.go` settles `completed`. Its commit carries
   the Task file with `## Recorded paths` naming exactly
   `internal/store/delivery_test.go`, and the commit event lists the same path.
   A Task that changes only declared paths gets no section. A Governed Path, a
   QA Task and a failed Task are never recorded.
2. The built binary's `roundfix qa-report accept` exits `0` on Spec 0179's
   archived `qa-report-2026-09-29.md`, which the v0.20.0 binary refuses with
   `rows_blocked_environment is 1; expected 0`. It still exits `1` on a report
   whose environment-blocked rows include one outside the Pull Request row.
3. In a disposable clone at `ebeb997f` with ADR-0161 added in a later commit,
   the binary built from this Spec's starting main reports `SC-ADR-RELATED`
   for ADR-0161 on Spec 0180. The candidate binary reports none, and it still
   reports the gap when ADR-0161 is committed before the PRD.
4. The `qa-gate`, `write-tasks` and Roundfix skills, their mirrors, the user
   guide and `CONTEXT.md` describe the recorded paths, the pre-PR Pull Request
   row and the horizon, and `make skills-sync-check` exits `0`.

## Recorded limits

- A recorded path is disclosed, not planned. It reserves nothing against a
  parallel Task, so two Tasks can still collide on a file neither declared.
- The QA scope audit is a disclosure check for ordinary paths. The pre-PR
  review of the diff against the Spec stays the judge of whether an undeclared
  change belongs to the slice.
- A report whose Results table has no provenance column gets no pre-PR
  exception.
- A PRD revised after a later ADR landed keeps the horizon of its first commit.
  Its author is not told about that ADR.
- This Spec's own Run executes on the v0.20.0 Daemon, which records nothing.
  Its Tasks therefore declare every path they expect to change, and its QA
  scope row applies the policy this Spec introduces: a test file a Task's
  `## Result` names as invalidated by the change is disclosed, not refused.

## Decisions

- **Record outside the authored Context.** Appending to `## Context` would
  change authored bytes. An appended `creates:` entry can collide with another
  Task or pass the fifty-entry limit, and either makes the Task file
  unparseable. A Daemon-owned section, as `## Carry-forward provenance`
  already is, changes nothing that carry-forward, path reservation or the
  consistency check reads. See ADR-0166.
- **The Daemon records, not the Agent.** An Agent asked to append its own
  entries forgets exactly when it matters. The Daemon already computes the
  staged set and already writes the Task file it commits.
- **Disclose ordinary paths, refuse governed ones.** Published scope checkers
  block only protected surfaces outright and report ordinary out-of-scope files
  as warnings (see Research basis). Roundfix already refuses a Governed Path
  outside its authorization, and that refusal stays.
- **Derive the pre-PR row from the Results rows, not a new key.** The
  mechanical stage already validates each row's status and provenance, and the
  frontmatter counts keep their meaning. A new key would be one more number an
  Agent can get wrong. See ADR-0167.
- **Only the Pull Request row.** Its block is structural: ADR-0088 fixes the
  order, and every Spec pays for it. Other environment blocks depend on the Run
  and remain declarable.
- **A commit-ancestry horizon, not a date.** ADR-0161 was created at
  `2026-09-28T16:02Z` and Spec 0180's PRD on `2026-09-28`, so a date comparison
  still fires on a same-day tie. Commit ancestry is exact, and the
  ADR-0161 cascade on 2026-09-28 is what it must stop. See ADR-0168.

## Acceptance evidence

Each Core Feature needs positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight:

- a recorded Governed Path;
- a recorded path in a QA Task or a failed Task;
- an edited `## Context`;
- a non-PR environment row excused;
- a PR-status row without the Pull Request source excused;
- a predating ADR skipped;
- an uncommitted Spec checked with a horizon.

Each of these would pass a happy-path test.

The outside-evidence row rests on sources this Spec did not produce:

- **Oraculum Spec 0027's QA report.** Its Secondbrain mirror is at
  `~/dev/secondbrain/projects/oraculum/mirror/docs/history/specs/0027-agregadores-de-vendas/qa/qa-report-2026-08-07-02.md`.
  The report records `blocked (environment: no open Pull Request)` beside
  three other environment blocks. The candidate binary's `roundfix qa-report
  accept` must still refuse it.
- **The fleet QA reports carrying the same row.** Reports under the Conexus,
  Fiscus, Fluxus, Oraculum and Vortex mirrors carry it, and the gate records
  how many.
- **The ADR-0161 cascade.** It is reproduced on commits `ebeb997f` and
  `06afa835`, which other sessions wrote.

When the Secondbrain is unavailable, the row is recorded as blocked with that
reason.

## Research basis

The three adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md).

Secondbrain, consulted through `wiki/index.md` and
`qmd query ... --all --files --min-score 0.3`:

- `projects/roundfix/mirror/docs/history/findings/2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md`
  records the structural cause of the Pull Request block that ADR-0080 answered
  for `pass`.
- `projects/roundfix/mirror/docs/history/findings/2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md`
  calls `SC-ADR-RELATED` "a suggestion to the author of an active Spec", the
  principle the horizon extends from archived Specs to Specs committed before
  the ADR.
- The Oraculum, Conexus, Fiscus, Fluxus and Vortex QA report mirrors show the
  Pull Request row blocked the same way across the fleet.

Exa, one search on how coding-agent harnesses treat files changed outside a
declared scope:

- TaskBound (https://github.com/Conalh/TaskBound) grades drift by severity.
  An ordinary out-of-scope file is "medium" and reported for review, while
  excluded or sensitive surfaces are "critical".
- AgentScope (https://abdouloued.github.io/agentscopev2/) hard-blocks only
  configured paths and warns on the rest.
- agentdiff (https://github.com/iselur/agentdiff) reports out-of-scope files as
  MED findings for review.
- agent-guardrails
  (https://github.com/logi-cmd/agent-guardrails/blob/HEAD/README.md) keeps a
  `violationBudget` for minor scope slips.

That published practice supports disclosing an ordinary undeclared path and
refusing only the governed ones. The tools' numbers are their authors' and are
not re-measured here.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
