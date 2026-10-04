---
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: active
created: 2026-10-04
surfaces: [backend, cli, docs]
---

# An archived retry that needs no recorded candidate

On 2026-10-04 the Delivery Queue parked Spec 0220 `run-unresolved`. Its newest
QA Report was `partial` with no finding-blocked row, and its two
environment-blocked rows both waited for an open Pull Request: one was the
Pull Request row, the other the Linux CI run that only a Pull Request starts.
The operator archived the Spec with the QA Archive Override, as the standing
authorization allows, and `roundfix deliver retry` refused with
`candidate head is missing`. No review had run, so no candidate head was
recorded, and only a `qa-environment-partial` park may use the Implement
start head of its Run in place of a candidate. The operator opened Pull
Request #367 and ran the pre-PR review by hand. Spec 0211 ended the same way
on 2026-10-02, before that exception existed.

This is a bug fix. It has two causes, both measured during authoring:

1. The Delivery Retry accepts an operator archive without a candidate only for
   the `qa-environment-partial` park. A `run-unresolved` item, which never
   has a candidate, cannot be resumed after an override at all.
2. The queue parks a zero-finding `partial` as `qa-environment-partial` only
   when its environment-blocked rows outnumber the rows blocked only because
   no Pull Request is open. In 0220's report both environment-blocked rows
   were such rows (2 against 2), so it parked `run-unresolved`, whose retry
   would run the QA gate again against the same missing Pull Request.

## Prerequisites

None. Specs 0222 and 0223 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the operator orders the queue so that the
later Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes. Queue items
  keep their Spec slug and Runs their Run identifiers; the blocker names are
  unchanged. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the retry and the park
  classification read only local Git, the Run Branch's QA Report and the Run
  Database; no credential or network call is added, and no test reaches
  GitHub. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0229 (this Spec) decides both
  rules, ADR-0229: "accepts the item head when Git proves it descends from
  the Implement start head of the item's recorded Run, whatever the park".
  ADR-0223 keeps a delivery's work across archive, requeue and review with
  its descendant rule for an item that has a candidate,
  ADR-0223: "A park for a corrective Spec still refuses a moved head".
  ADR-0167 still holds that the pre-PR Pull Request row never decides a
  qualifying partial for archive, ADR-0167: "That row no longer
  counts against a qualifying". ADR-0154 keeps the override a record of user
  authority, ADR-0154: "Preserve its approval source/date and actual QA
  outcome or absence". ADR-0165 keeps a review finding after archive parked
  for a corrective Spec, ADR-0165: "an archived Spec is never edited to absorb
  a finding". ADR-0187 and ADR-0189 govern the Roundfix Skill edit
  and its version, and ADR-0184 has the TechSpec state changed command
  surfaces, ADR-0184: "A TechSpec now declares numbered Surface Transcripts". The gate is bound by ADR-0080, ADR-0088,
  ADR-0091, ADR-0104, ADR-0155 and ADR-0156, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by citation
  and receipt. ADR-0178 and ADR-0179 decide the grant each Task commit runs
  under, and ADR-0166 records undeclared paths; every Task declares its
  paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage
  and row carry, ADR-0182 cites ADR-0117 but settles a Task on the facts its
  gate will check, ADR-0169 cites ADR-0165 but decides the review's merge-base
  diff, ADR-0196 cites ADR-0169 but decides when a review finding parks,
  ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row
  records, when it is observed again and its evidence snapshot, and ADR-0192
  cites ADR-0178 but decides derived-path conflicts; this Spec changes none of
  them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and
  its `SKILL.md` mirror are Governed Paths, and the maintainer authorized
  skill edits ("considere autorizado a ajustar todas as skills se
  necessário") and this Spec (2026-10-04, "deliver M"). No other Governed
  Path changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0224-an-archived-retry-that-needs-no-recorded-candidate/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- An item the operator archived with the QA Archive Override resumes at
  `reviewing` after `roundfix deliver retry`, whatever its park, when its
  head descends from the Implement start head of its recorded Run and no
  candidate is recorded.
- A zero-finding QA `partial` blocked only by rows waiting for an open Pull
  Request parks `qa-environment-partial`, so status prints the override
  recovery instead of a Run retry that would park again.
- Every existing refusal holds: no override, a non-descendant head, missing
  history, `corrective-spec-required` and `pull-request-conflict`.

## Core Features

1. **Retry from the Run start after any override.** A Delivery Retry of an
   archived item whose archive records `qa_override: true` and that has no
   candidate head finds the Implement start head of the item's recorded Run.
   When Git proves the item head descends from it, the retry records that
   head as the candidate and returns the item to `reviewing`, for every park
   that reaches the archived branch of the retry (ADR-0229).
2. **An environment-only partial parks for the operator.** After an
   unresolved Run, a newest Run-Branch QA Report that is `partial`, has no
   finding-blocked row and has at least one environment-blocked row parks
   `qa-environment-partial`, including when every such row waits for an open
   Pull Request (ADR-0229).
3. **The skill and the guide say so.** The Roundfix Skill's `deliver`
   reference and the `deliver` command guide state both rules and drop the
   sentences that limit the Run start head to the environment park and keep a
   Pull-Request-only partial `run-unresolved`.

## Non-Goals / Out of Scope

- Archiving a `partial` without the operator's override, or changing the
  Archive Command's eligibility (ADR-0167).
- Accepting an archive without `qa_override: true` when no candidate exists.
- Changing the retry of `corrective-spec-required` or `pull-request-conflict`
  items, Park Class names, blocker names, retry limits or status wording.
- Changing what the QA gate writes for the Pull Request row.

## Success Metrics

1. Success Metric: in a disposable repository with a real archive made by
   `roundfix archive --qa-override`, a `run-unresolved` item with no
   candidate resumes at `reviewing` with the archived head as its only
   candidate, where today the retry fails with `candidate head is missing`.
2. Success Metric: a Run-Branch QA Report shaped like 0220's, `partial` with
   two environment-blocked rows that both wait for an open Pull Request and
   none blocked by a finding, sets the environment partial; a partial with
   only declared rows or with a finding row does not.
3. Success Metric: every existing operator-archive, Run start and park
   classification test passes, with one declared change: the `pre-PR only`
   case of `TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch` now
   expects an environment partial.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The operator's intervention log (`~/.roundfix-operator/queue-interventions.md`),
  entries 96 (0211, 2026-10-02: "PR by hand (retry refused: candidate head
  empty)") and 140 to 142 (0220, 2026-10-04: "retry refused 'candidate head
  is missing' (run-unresolved park, no candidate) → PR by hand + review by
  hand").
- The live Run Database, read with `sqlite3 -readonly` during authoring on
  2026-10-04: the 0220 item is `parked`, blocker `run-unresolved`,
  `candidate_commits` `[]`, Run `run_20261004T010416Z_bdbf1cc82a3473ae`, an
  Implement Run with start head `9c98d516`.
- The archived 0220 QA Report, read through the shipped QA Report reader on
  2026-10-04: verdict `partial`, `rows_blocked_environment: 2`, two rows
  counted as waiting for a Pull Request, `rows_blocked_finding: 0`, so the
  queue's test `2 > 2` was false; the Archive Command's eligibility for the
  same report is "newest QA Report verdict is "partial"; expected "pass"".
- Git's documentation of `git merge-base --is-ancestor`
  (<https://git-scm.com/docs/git-merge-base>), the ancestry proof the retry
  uses: it will "exit with status 0 if true, or with status 1 if not. Errors
  are signaled by a non-zero status that is not 1".
- Temporal's guide to recovering a process without restarting
  (<https://docs.temporal.io/guides/recover-without-restart>), a published
  account of the same recovery shape: after an operator's correction "The
  process picks up where it left off: completed steps do not re-execute".
- The Secondbrain holds no independent entry on this defect; its search on
  2026-10-04 returned only mirrors of Roundfix's own guide and archived Specs
  0201, 0211 and 0219.

## Decisions

- The override, not the park, authorizes the Run start anchor. See ADR-0229.
- A zero-finding partial with any environment-blocked row is an environment
  park. See ADR-0229.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
