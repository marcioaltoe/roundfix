---
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
status: archived
created: 2026-09-29
surfaces: [backend, cli, docs]
archived: "2026-09-29"
source_slug: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
---


# Delivery that reviews and retries from where the item stands

`roundfix deliver` reviews a candidate and recovers a parked item from facts it
reads at that moment. On 2026-09-28 two of those facts were read from the
wrong place, and each cost a manual intervention:

- The pre-PR review diffed the tip of the default branch against the candidate
  head. Main had moved on after Spec 0175's item branch was cut: #265 retired
  Specs 0126 and 0127, and #266 authored Specs 0177–0180. The review therefore
  reported that the candidate deleted four active Specs and reactivated two
  retired ones. The same false "reversion" appeared in four reviews, and each
  needed a re-merge and a new review. A main-side edit of an archived Spec also
  reaches `archivedSpecs`, and the delivery engine parks such an item as
  `corrective-spec-required` for a Spec it never touched.
- Spec 0175's item was retried twice after `BudgetExceeded`. The second retry
  carried forward from the Run recorded on the item, the first Run
  `run_20260928T154312Z_63679c6540b34603`, whose Tasks were already on the
  item branch. It was refused with `declared input(s) moved`, and the newer Run
  `run_20260928T174529Z_c0cad0da5734a85f` kept Tasks 03–05. `implement` then
  refused to start while those proved Tasks waited, and the item parked
  `delivery-error` again. The operator carried the newer Run by hand with
  `roundfix reconcile --carry-forward` and retried.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a review record
  keeps its commits and a Delivery Queue item keeps its Spec slug, branch and
  Run ID. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git, the Run Database
  and the configured reviewer runtime only; no credential is read and no
  network call is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0169 (this Spec) makes the pre-PR
  review diff the candidate from its merge base, and ADR-0170 (this Spec) makes
  a Task already completed on its carry-forward target nothing to carry. Both
  are recorded before implementation. ADR-0153 keeps the pre-PR review an
  explicit provider policy that must review the current candidate, and ADR-0151
  keeps the configured review profile as the reviewer; neither changes.
  ADR-0165 parks a blocking review after archive for a corrective Spec, and
  ADR-0169 narrows that park to archived Specs the candidate itself changed.
  ADR-0057 and ADR-0014 keep Task status and Verification Daemon-owned, which
  is why a completed status on the target is authoritative. ADR-0053 keeps
  reconciliation proof-based, and every remaining carry-forward candidate keeps
  every proof. ADR-0158 keeps a BudgetExceeded Run recoverable through
  carry-forward. ADR-0026 integrates settled Tasks in completion order, which
  carry-forward still follows inside each Run. ADR-0090 forbids reusing a
  repository fact across a mutation, so the retry rereads the item's Task Graph
  after each carried Run. ADR-0052 makes completion compare-and-set, the
  model for the retry's guarded item transition, which is unchanged. ADR-0044
  reclaims an orphaned owner record only on proven owner death, which the
  retry's hand-off keeps. ADR-0139 keeps one Active Run per work target. This
  Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0104, ADR-0117,
  ADR-0155 and ADR-0156. ADR-0093 checks Spec consistency by citation, and
  ADR-0094 makes that consistency check artifact-presence-aware. ADR-0020,
  ADR-0038 and ADR-0127 cite ADR-0014 but govern the Agent prompt result, the
  Verification repair bound and process residue. ADR-0160 cites ADR-0014 but
  names the frozen authorization as the only authority that opens a red
  repository gate. ADR-0097 cites ADR-0080 but carries a QA row forward, not a
  Task. ADR-0161 cites ADR-0053 but releases a merged Spec's Runs on the merged
  head. ADR-0164 cites ADR-0158 but renews an Implement Run's budget at each
  Task settlement. ADR-0056 and ADR-0159 cite ADR-0038 but govern
  Verification capacity and independent Verification. This Spec changes none of
  them, so they do not apply. ADR-0166 (Spec 0181) has the Daemon record the
  paths a Task changed without declaring them; every Task here declares the
  paths it edits, so it changes nothing here. ADR-0167 (Spec 0181) keeps the
  pre-PR Pull Request row from deciding a qualifying partial; this Spec's gate
  aims at `pass`, so it does not apply. ADR-0168 (Spec 0181) opens a
  related-ADR gap only for ADRs that predate the Spec, which narrows this
  Spec's own check and holds. All the others hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved the Onda 4 plan in
  chat on 2026-09-29 ("Duas filas": Specs 0181 and 0182 in one
  `roundfix deliver start`); the Roundfix skill files ride the standing grant
  of 2026-09-18 for keeping the shipped skills true to the CLI, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A pre-PR review examines only what the candidate introduces, however far the
  default branch has moved since the candidate was cut.
- A delivery item is never parked for an archived Spec that only the default
  branch changed.
- Task Carry-Forward never refuses a Run because of Tasks its target already
  completed.
- A Delivery Retry recovers every proved Task of the item's Spec, whichever Run
  holds it and whichever Run the item recorded.

## Core Features

1. **The review diffs from the merge base.** `roundfix review` resolves the
   merge base of the current head and the selected base ref, the default
   branch when `--base` is omitted. That merge base is the before side of the
   diff sent to the reviewer and of the changed-Spec list. The prompt names it
   as the base commit, and the review record stores it as `baseCommit`. The
   record also stores the commit the base ref resolved to as `baseTipCommit`.
   A head that shares no history with the base ref is refused before any
   reviewer call.
2. **A review of an unchanged candidate stands while main moves.** A findings
   verdict is reused by repository, merge base, head and provider, so a new
   commit on the default branch does not force a fresh review of a candidate
   whose diff it did not change.
3. **A completed Task is nothing to carry.** Task Carry-Forward records a Task
   whose status is `completed` on its target as already completed. It is not
   staged, proved or refused, and it never refuses the rest of the set. The
   remaining Tasks pass or refuse as one whole set, as before. A Run whose
   Tasks are all completed on the target carries nothing and succeeds. This
   holds for `roundfix reconcile <run-id> --carry-forward`, for the Implement
   Command's carry-forward Preflight and for a Delivery Retry.
4. **A Delivery Retry carries from every Run of its item.** A retry of an item
   whose Spec is still active considers every terminal Implement Run of that
   Spec on the item branch, newest first, together with the Run the item
   records. It carries each Run's remaining proved Tasks and rereads the item's
   Task Graph between Runs. It records the newest Run on the item and prints one
   `Carried forward from Run <run-id>: <task>, …` line per Run that moved Tasks.
   A Run whose Tasks are all completed on the item needs no Run Worktree and is
   never refused.

## Non-Goals / Out of Scope

- Merging the default branch into a candidate, rebasing it, or changing what
  the delivery queue pushes.
- Changing the reviewer prompt's wording, the review policy or its providers.
- Carrying part of a set in which a remaining Task refuses, or re-proving a
  Task its target already completed.
- Retrying several items at once, an automatic retry, or a Run Database schema
  change.
- Reconciling or releasing the Run Worktrees of Runs that carried nothing.

## Success Metrics

1. In a disposable repository where the default branch gains commits after the
   candidate was cut, one of them editing an archived Spec, the review prompt's
   diff contains only the candidate's changes. The record's `baseCommit` equals
   `git merge-base <base-tip> HEAD`, its `baseTipCommit` equals the base tip,
   and its `archivedSpecs` is empty. A candidate that itself edits an archived
   Spec still lists it.
2. A findings record made before the default branch moved is reused after it
   moves, with no reviewer call. A base ref that shares no history with `HEAD`
   exits `2` before any reviewer call.
3. `roundfix reconcile <run-id> --carry-forward` over a Run whose `task_01` is
   already completed on the checkout exits `0` and carries only `task_02`.
   Repeating it exits `0` with `HEAD` unchanged. A moved input of `task_02`
   still refuses the remaining set with `HEAD` unchanged.
4. With two BudgetExceeded Runs of one Spec on the item branch, where the item
   records the older Run and the older Run's Tasks are already on the item
   branch, a retry carries the newer Run's Tasks. It records the newer Run on
   the item and prints one `Carried forward` line naming it.

## Recorded limits

- A review of a candidate that merged the default branch into itself has the
  default branch's tip at that merge as its merge base, as a pull request has;
  the review then covers the merge's conflict resolutions too.
- With several merge bases, the one `git merge-base` reports is used, as
  `git diff A...B` does.
- A retry carries Runs one at a time. When a later Run refuses, the Tasks of
  Runs carried before it stay on the item branch. They are proved and complete,
  and the refusal names them.
- A Task whose target status is `completed` is trusted as completed. A status
  hand-edited to `completed` on the target is out of reach of this proof, as it
  is of every other reader of the Task Graph.

## Decisions

- **Resolve the merge base once, where the base is resolved.** One commit
  feeds the diff, the changed-Spec list, the prompt and the record, so no
  consumer can disagree with another. Switching each `git diff` call to
  `base...head` would leave the prompt, the record and the reuse key naming
  the tip of main. See ADR-0169.
- **Keep the tip for provenance, key reuse on the merge base.** Two reviews
  with one merge base and one head examine the same diff, so reuse must not
  depend on where main has moved. The tip stays in the record so a reader can
  see how far main had moved.
- **Filter completed Tasks in the shared inspection.** Reconcile, the
  implement Preflight and the retry refused for the same reason. Fixing only
  the retry would leave `reconcile --carry-forward` refusing the operator's own
  second carry and the Preflight noting a fully carried Run as unavailable.
  See ADR-0170.
- **Newest Run first.** The newest Run started from the most complete item
  head, so carrying it first turns the older Runs' Tasks into already-completed
  ones instead of moved inputs.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight. A review that still shows
main's commits, an archived Spec listed only because main edited it, a verdict
discarded because main moved, a completed Task that still refuses its set, a
remaining Task that no longer refuses on a moved input, a retry that carries
from the recorded Run only, and a fully carried Run refused for a missing
worktree would each pass a happy-path test.

The outside-evidence rows rest on sources this Spec did not produce:

- This repository's history. Spec 0175's reviewed head `74a51f23` and the main
  tip `159b41db` it was reviewed against were written by other sessions.
  `git diff --name-status 159b41db 74a51f23` shows Specs 0177–0180 deleted and
  0126/0127 reactivated. The diff from their merge base shows only 0175's own
  changes.
- The live Run Database, read through the built `roundfix runs list --all
  --state all --limit 0`. It lists Spec 0175's two BudgetExceeded Runs on one
  item branch, the shape the retry now carries.
- Git's own documentation. `git diff A...B` "is equivalent to
  `git diff $(git merge-base A B) B`" (<https://git-scm.com/docs/git-diff>),
  and GitHub states that a pull request shows this three-dot comparison
  because it shows "what a pull request introduces"
  (<https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-comparing-branches-in-pull-requests>).

## Research basis

The two adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md). The Secondbrain was consulted
through `wiki/index.md` and
`qmd query "pre-PR review diff against moved main merge base false reversion"`
and `qmd query "deliver retry carry forward stale run newest"`. It holds no
knowledge beyond this repository's own mirrors, including the adopted entries
and Spec 0173's carry-forward references. Exa located the two primary sources
cited above. They confirm that the merge-base comparison is the established
meaning of "what a candidate introduces", which is the premise of ADR-0169.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
