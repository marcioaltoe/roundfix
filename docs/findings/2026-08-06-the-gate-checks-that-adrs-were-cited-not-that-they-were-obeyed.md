# 2026-08-06 — The gate checks that ADRs were cited, not that they were obeyed

status: pending

Source: a supervised autonomous session in the `fluxus` repository that delivered three specs
(0018, 0019, 0020) and two releases (`0.6.0`, `0.6.1`) between 2026-08-05 17:00 and
2026-08-06 09:45 America/Sao_Paulo, including a destructive production migration. **Eleven
Runs for three specs; exactly one ended `Clean`.** Ten ended Stopped, Unresolved, TimedOut,
ReviewSkipped or IntegrationPending — almost all on contract, not on broken code.

This finding reports only what is **not** already covered by
`2026-08-05-authored-verification-gates-are-untested-code.md`. Two of that finding's proposals
recurred verbatim in this session and are recorded here as recurrence evidence, not as new
proposals.

## 1. The Project Constraint audit is a citation check, and a citation check is not conformance

The strongest observation of the session. Spec 0020's Task 04 raised the outbound request
timeout from 60s to 120s. Two **accepted** ADRs in the target repository reject exactly that:

- ADR-0022: *"Considered and rejected: raising the global request timeout, ..."*
- ADR-0034: *"Considered and rejected: ... raising the request timeout instead of classifying
  the failure (already rejected by ADR-0022)."*

The QA gate ran its Project Constraint audit and reported **F-03: "Active ADR accounting omits
the governing own-deadline decision"** — that ADR-0034 was missing from the constraints list.
Its required correction was to *account for ADR 0034 and its applicability*.

Had the Supervisor done the minimum the gate asked — append "ADR 0034" to the list — the gate
would have passed, and a change that two accepted ADRs reject would have shipped to production
**wearing a conformance stamp**. The contradiction surfaced only because the Supervisor opened
the cited ADR and read it.

The audit verifies that applicable ADRs are *named with applicability and reasons*. It does not
verify that the Spec's own decisions do not contradict them. Those are different questions, and
only the second one protects the architecture.

Suggested behavior, in increasing cost:

- `write-techspec` requires, per applicable ADR, a line stating **how this Spec's decisions
  relate to that ADR's rejected alternatives** — not only that the ADR applies. An ADR whose
  "Considered and rejected" list intersects a Spec Decision must either not apply, or the Spec
  must supersede it explicitly with a new ADR.
- The QA gate's constraint audit reads each cited ADR's rejected-alternatives section and fails
  when a Spec Decision restates one. Mechanically feasible: these ADRs carry the rejections in
  prose, and a Spec's `## Decisions` section is short.
- Minimum viable version: the audit fails when an ADR that the changed files touch is **absent**
  from the constraints, which is what it already found here — but its `next:` should tell the
  Supervisor to *read the ADR and compare it against the Spec's Decisions*, not to *add it to a
  list*. The remediation text is what steered toward the cheap fix.

## 2. A declared-only `partial` can never archive through the supported path

The Archive Command's documented contract says an archive-eligible `partial` — one whose unmet
rows are covered by pre-run `## Unreachable Acceptance` declarations — does not prevent Archive.
Observed behavior contradicts it:

- The terminal QA Task settles `completed` only on `verdict: pass`. A `partial` settles `failed`.
- `roundfix archive` refuses while any Task is not `completed`:
  `Task "task_08" is "failed"; archive requires every Task to be "completed"`.

So the declared-unreachability mechanism, which exists precisely to let an honest `partial`
archive, cannot reach the Archive Command. In this session spec 0019 produced exactly that
shape: `verdict: partial`, `rows_blocked_declared: 1`, `rows_blocked_finding: 0`, zero failures
— and could not archive. It was archived only by the manual `archive-spec` skill procedure with
`qa_override: true`, which the skill itself says is the wrong instrument: *"declared
unreachability does not use or weaken that override."*

There is no supported path between "gate says partial for the reason the contract anticipated"
and "spec archives". Suggested behavior: settle the gate Task `completed` when the verdict is
`partial` **and** every unmet row is declared-covered, since that is already the archive
eligibility rule; or let Archive accept a `failed` gate Task whose newest report meets that same
rule. Either closes the gap without touching `qa_override`.

This is adjacent to the 2026-08-05 addendum "closing a failed gate needs a hand-edit the docs do
not name", but distinct: that one is about appending a corrective Task after a `fail`; this one
is about a `partial` that the contract explicitly calls archive-eligible.

## 3. A pull request that opens clean never gets a review requested

Spec 0078 made Roundfix publish the review request itself, with the success metric *"One Round
produces exactly one review request, asserted for a Round whose Final Push is followed by the
artifact-only docs commit."* The request is **Round-scoped**.

A `watch --until-clean` over a pull request that opens with **zero Review Issues** runs zero
Rounds, publishes zero requests, and waits for a review nobody asked for until the Review Source
timeout. Measured here on PR #89: `TimedOut` after 30 minutes with
`Review skipped: automatic reviews are disabled`, then a manual `@coderabbitai review`, then
`ReviewSkipped — Review limit reached`. **Thirty minutes of an unattended overnight budget, for
a pull request that had nothing to fix.**

The failure mode is systematic: every clean pull request pays it. Suggested behavior: when the
Review Source has automatic review disabled and the expected head carries no review evidence,
publish the request at watch start rather than only after a Round's Final Push. The command is
already computed and carried in `ManualReviewCommand`.

## 4. `release plan` is unusable against a repository whose tags carry no `v`

`roundfix release plan` refuses with `no stable release tag; create or pass a stable
vMAJOR.MINOR.PATCH base tag`. The target repository's convention is bare `0.5.7`, `0.6.0`,
`0.6.1`, and its tag-triggered CD fires on `*`, so the convention is load-bearing and correct
for that repo.

The command is therefore unavailable exactly where the autonomous loop most wants it: it is the
step that classifies impact and forces the `approval_required` decision for a minor or major
release. In this session the release had to be planned by reading the changelog and the tag
history by hand, and the human approval that `release plan` would have demanded was obtained
conversationally instead.

Suggested behavior: accept an optional configured tag pattern (or infer the prefix from the
repository's existing stable tags) rather than requiring the `v` prefix. Refusing to plan is a
much worse outcome than planning against a non-`v` tag.

## 5. SQLITE_BUSY at two concurrent writers, and GC that cannot reclaim

Two Tasks in one wave, both `data`, both died with
`Agent Batch failed: publish Run Events: begin Run Event append: database is locked (5)
(SQLITE_BUSY)`. Concurrency was 3; **the wave contained exactly 2 Tasks**, so two concurrent
writers were enough. The Run Database was **2.6 GB** with an 11 MB WAL.

`roundfix gc` could not help: it pruned 294 Runs and removed **0 journal rows and 0 artifact
bytes**, because the volume is recent Agent payloads inside the 336h retention window. The
operator's only lever was dropping `worktree.concurrency` to 1, which cost roughly half the
session's wall-clock.

Suggested behavior: a busy-timeout and/or write serialization on Run Event append, so two
Agents settling near-simultaneously queue instead of failing their Batches; and a `gc` report
that names *why* nothing was reclaimable (retention window versus eligibility) rather than
printing zeros, which reads as a no-op bug.

## 6. Recurrence evidence for two 2026-08-05 proposals

Both already proposed in `2026-08-05-authored-verification-gates-are-untested-code.md`; this
session is independent evidence that they were not theoretical.

- **§2 mandatory Verification lint.** **Nine** authored Verification commands across two specs
  declared `expected: no match` with a bare `grep`. Roundfix gates on exit status and `grep`
  exits 1 when it finds nothing, so the intended success was scored as failure. The class is
  *never-green*, which the proposed red-pre-work lint catches. One had already failed a Task
  before the Supervisor found and fixed all nine at once.
- **§4 serialize ordinal-artifact generators.** Two `data` Tasks in one wave each generated
  drizzle migration `0035` from the same base. Cherry-pick integration cannot merge a drizzle
  snapshot, so one Task's completed, verified work had to be discarded and redone. The proposal
  to serialize ordinal-artifact generators would have prevented it exactly.

## 7. What worked — keep

- **The gate audits the Supervisor.** It caught a protected-tooling mutation the Supervisor had
  committed (`.roundfixrc.yml`), naming the exact commit and that no Spec artifact authorized
  it. A gate that only audits Agents would have missed it. This is the single most valuable
  behavior observed.
- **`roundfix settle` as real recovery.** Recovered two Tasks' completed work — 25 minutes and
  a full data Task — after infrastructure failures, with Verification re-run and integration
  performed. Nothing was redone.
- **`reconcile` refusing to guess.** It preserved an orphan Task Branch it could not prove safe.
  The Supervisor verified supersession by hand and the skill's own warning — that the suggested
  `merge --ff-only` can be wrong when it would overwrite newer work — was exactly right: the
  orphan carried a superseded migration `0035` that would have clobbered `0036`–`0038`.
- **ADR-0020 classification and Verification Feedback.** No spurious Batch failures from
  transport noise across eleven Runs.
