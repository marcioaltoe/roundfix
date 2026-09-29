---
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Review findings with evidence, and no unselected providers

A pre-PR review finding today ends in the operator's hands, and its
disposition lives nowhere Roundfix can read. On 2026-09-25 and 2026-09-28 the
Codex pre-PR review raised findings on Specs 0170, 0171 and 0172. The operator
fixed some with corrective Tasks and TechSpec edits, repaired a QA Report's
front matter by hand, and dismissed two findings only in pull request bodies:
the Windows portability finding on #258 and the stale-binary finding on #259.
No per-finding record exists, and four gaps follow from that:

- `roundfix review` writes `pre-pr-review.json` with the findings as one block
  of text, so a finding has no identity to carry a disposition.
- Nothing records that a finding was fixed by a commit or dismissed with
  evidence, and nothing ties that record to the head it reviewed.
- The Delivery Queue parks an item as `review-findings`. `roundfix deliver
  retry` then re-runs the reviewer at the same head, so a dismissed finding
  either parks the item again or clears when the reviewer happens to answer
  differently. Re-asking until the answer changes is not a review.
- When a review blocks on a candidate that archives a Spec, the correction
  path is ad hoc. An archived Spec must stay byte-identical, so a corrective
  Task inside it is not an option, and nothing says what is.

A fifth gap is older. The Roundfix skill still tells every agent that a pull
request "gets no review unless someone requests one" and how to request
CodeRabbit's, whatever Pre-PR Review Policy the repository selected, and the
generated config writes a `review_source` block with no word that it is the
legacy PR-feedback Review Source.

This Spec gives each finding an identity and an evidence-backed disposition
tied to its head. A findings verdict stands until the candidate changes, and a
blocking review after archive parks publication with the corrective Spec it
needs. An unselected provider is never asked for anything.

## Project Constraints

- Identifier strategy: applicable — a finding gets an ordinal identity `F1`,
  `F2`, … within the review record of one head, and a disposition is keyed by
  repository, head commit, finding identity and finding text. No other
  identifier changes; Run IDs, Delivery Queue items, blocker names already in
  use and the review record's existing fields keep their meaning. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review record, the disposition
  ledger and the Delivery Queue are local files and the Run Database; no
  credential, no network call and no provider request is added, and CodeRabbit
  is called by nothing new. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0153 makes pre-PR review an
  explicit provider policy where a failed, stale or missing enabled review is
  an error, never an omission; this Spec keeps that vocabulary and adds a
  dismissed-findings outcome that is never shown as a clean review. ADR-0151
  keeps `profiles.review` the reviewer's selection and names the legacy
  `review_source.name` a PR-feedback provider, which is the line this Spec
  draws in guidance and generated config. ADR-0142 binds Review Source
  Evidence to the pushed head for `watch`; the legacy commands and their
  `.coderabbit.yaml` coherence refusal keep working unchanged. ADR-0014 and
  ADR-0057 keep Verification and Task status Daemon-owned, so a disposition
  never settles a Task. ADR-0020, ADR-0038, ADR-0127 and ADR-0160 cite
  ADR-0014 but govern the Agent prompt result, the Verification repair,
  process residue and the red repository gate, which this Spec does not
  touch, so they do not apply; ADR-0056 and ADR-0159 cite ADR-0038 and govern
  Verification capacity and independent Verification, which this Spec does
  not touch either. ADR-0120 keeps retired Specs under the one history
  root, which is why a finding on an archived Spec needs a new Spec instead of
  an edit. ADR-0154 keeps an archive, overridden or not, separate from review
  and publication gates. This Spec adds ADR-0165, the late-correction
  policy. This Spec's gate is bound by ADR-0080, ADR-0091,
  ADR-0093, ADR-0096, ADR-0097, ADR-0104, ADR-0110, ADR-0117, ADR-0130,
  ADR-0155 and ADR-0156. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized continuing with
  the next wave after the v0.18.0 release in chat on 2026-09-28; the skill
  files ride the standing grant of 2026-09-18 for keeping the shipped skills
  true to the CLI, recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A repository is never asked, by Roundfix or its guidance, for a review
  provider it did not select.
- Every pre-PR review finding has an identity and at most one recorded,
  evidence-backed disposition tied to the head it reviewed.
- A findings verdict stands for its head: it clears only by dismissal with
  evidence at that head, and a changed candidate gets a fresh review.
- A blocking review of a candidate that archives a Spec parks publication and
  names the corrective Spec it needs; nothing inherits authority to author it.

## Core Features

1. **An unselected provider receives no request.** The Roundfix skill,
   `docs/user-guide/commands.md`, `docs/user-guide/usage.md` and
   `docs/user-guide/configuration.md` scope every CodeRabbit request to a
   repository whose policy selects `coderabbit` or to the legacy PR-feedback
   commands `fetch`, `watch` and `resolve`. The generated config labels
   `review_source` as that legacy Review Source. `roundfix review` and the
   Delivery Queue never read `.coderabbit.yaml`, never publish a CodeRabbit
   request and never fail because CodeRabbit is absent. When the legacy
   commands run, CodeRabbit stays fully available, coherence refusal included.
2. **Every finding gets a recorded disposition.** A findings record lists each
   finding as `F1`, `F2`, … with its text. `roundfix review dispose <id>`
   records exactly one disposition per finding in an append-only ledger beside
   the review record: `--dismiss --evidence <text>` at the reviewed head, or
   `--fixed-by <commit>` naming a descendant of it. Every refusal writes
   nothing.
3. **A findings verdict stands for its head.** `roundfix review` does not ask
   the reviewer again for a repository, base, head and provider whose recorded
   verdict is findings. When every finding is dismissed with evidence at that
   head, it reports `findings-dismissed` and exits `0`; otherwise the findings
   stand and it exits `1`. A new head, base or provider gets a fresh review. The
   Delivery Queue advances a `findings-dismissed` item to archive and parks a
   standing one as `review-findings`, as before.
4. **A blocking review after archive parks publication.** The review record
   names the Specs the candidate archives. A findings verdict on such a
   candidate parks its Delivery Queue item as `corrective-spec-required:
   <slug>`, and the stderr of `roundfix review` names the corrective Spec it
   needs. `roundfix deliver retry` refuses that item once its head moved, and
   returns it to `reviewing` when the head is unchanged, so dismissals decide.

## Non-Goals / Out of Scope

- A local CodeRabbit pre-PR review surface: `coderabbit` stays refused by
  `roundfix review` as it is today.
- Changing the Delivery Queue's stage order, the legacy `watch`, `resolve` and
  `fetch` behavior, or which Agent Selection Profiles `roundfix doctor` proves.
- Authoring, starting or authorizing a corrective Spec automatically, or
  editing an archived Spec.
- Recording a disposition inside the reviewed tree or in a pull request body.
- Rewriting existing review records, archived Specs or Delivery Queue items.

## Success Metrics

1. With a `.coderabbit.yaml` that disables automatic review and
   `review_source.request_review: false` — the pair Preflight refuses for
   `watch` — `roundfix review` exits `0` under `codex` and under `none`. The
   Delivery Queue's publication carries no CodeRabbit request under any policy,
   the generated config scopes `review_source`, and
   `TestRunEnforcesReviewRequestCoherence` stays green.
2. A findings record lists `F1..Fn`; `roundfix review dispose` appends exactly
   one ledger entry per finding. An unknown identity, both or neither
   disposition flags, blank evidence, a dismissal at a moved head, a
   `--fixed-by` commit that does not descend from the reviewed head, and a
   second disposition each exit `2` and append nothing.
3. With every finding dismissed at the head, `roundfix review` exits `0` with
   `findings-dismissed` and makes no reviewer call. With one finding left, it
   exits `1` with `findings` and makes no reviewer call. After a new commit, a
   different base or a different provider, the reviewer is called. The Delivery
   Queue moves a `findings-dismissed` item to `archiving`.
4. A findings verdict on a candidate that archives a Spec records that Spec in
   `archivedSpecs` and parks the item as `corrective-spec-required: <slug>`.
   `roundfix deliver retry` refuses the item with exit `2` when its head moved
   and returns it to `reviewing` when unchanged.
5. The dispositions recorded by hand in the bodies of #257, #258 and #259
   replay through the built binary: the two dismissals clear their head, and
   the fix from #257 leaves its head standing until a fresh review.

## Unreachable Acceptance

- criterion: Success Metric 2 — `roundfix review dispose` exercised through the
  built binary against a findings record a reviewer produced
  reason: the QA sandbox has no authorized ACP reviewer stand-in that reaches
  the built CLI, and a findings record may never be written by hand, so the
  gate cannot create the record the command disposes
  satisfied-by: the task_02 tests (every ledger field and each refusal)
  executed against the built tree, and the first `roundfix review dispose`
  run on a real findings record after this Spec merges

- criterion: Success Metric 3 — the reuse and fresh-review paths of
  `roundfix review` exercised through the built binary
  reason: the same missing reviewer stand-in prevents the gate from producing
  the findings record the reuse path reads
  satisfied-by: the task_03 tests executed against the built tree, and the
  first reuse of a real findings record after this Spec merges

- criterion: Success Metric 5 — the #257, #258 and #259 dispositions replayed
  through the built binary
  reason: replaying them needs findings records produced by a live reviewer at
  those heads, which the sandbox cannot create; the pull request bodies are
  captured locally in `qa/evidence/2026-09-28-pr-bodies-257-259.md`
  satisfied-by: the captured bodies checked against the dispositions the
  ledger format can express, and a live replay after this Spec merges

## Recorded limits

- A finding's identity is ordinal within one head's record. A fresh review of
  a new head may number the same concern differently, so a dismissal never
  carries across heads, by design.
- The split into findings follows list markers at the start of a line. A
  reviewer answer with no marker is one finding, which can make one
  disposition cover several concerns; the review prompt asks for one list item
  per finding to keep that rare.
- `deliver start` accepts only active Specs, and the queue reviews before it
  archives, so a queue item reaches `corrective-spec-required` only when its
  candidate archives some Spec before review. The same check covers the
  Supervisor order in `docs/agents/autonomous-work.md`, which archives before
  `roundfix review`.
- How a corrective Spec's work reaches the parked candidate, on the same item
  branch or as a new queue item, is left to a future Spec. Until then a parked
  `corrective-spec-required` item whose head moved stays refused by `roundfix
  deliver retry`, and ADR-0165 records the limit.

## Decisions

- **The late-correction policy is conservative, and it is an ADR.** The
  review policy left the archive-first late correction pending. The maintainer
  confirmed this policy on 2026-09-28 and asked for it as an ADR, so task_04
  records it as ADR-0165. An archived Spec is never edited to absorb a
  finding. Publication parks as `corrective-spec-required: <slug>`, and the
  correction is a new corrective Spec with its own `_authorization.md` and QA
  gate. No Run budget, corrective-Task ceiling or queue grant authorizes that
  Spec, and Roundfix never authors or starts it. After the correction, the
  candidate is reviewed afresh. How the corrective work reaches the parked
  candidate is a recorded limit for a future Spec.
- **A dismissal is the only way a findings verdict clears at its head.** A fix
  changes the candidate, so `--fixed-by` records history and never clears
  anything; the fresh review of the new head does.
- **`findings-dismissed` is its own outcome.** It exits `0` so publication can
  proceed, carries the findings and their dismissals, and is never written as
  `reviewed`. The disabled-review rule applies here too: dismissed findings are
  not a clean review.
- **Dispositions live beside the record, never in the tree.** Recording one
  must not change the candidate it describes. The append-only
  `pre-pr-review-dispositions.jsonl` in the Artifact Directory keeps each
  entry durable across later reviews.
- **The recorded verdict is reused only for findings.** Re-asking at an
  unchanged head is the permissive direction a gate must not take; a
  `reviewed`, `blocked` or `omitted` record is never reused.
- **The blocker is `corrective-spec-required`.** The Backlog Entry's shape
  proposed `review-blocked`, which already means an incomplete review, so a
  distinct name keeps both meanings readable.
- **CodeRabbit loses no capability.** The fix is scope, not removal: the legacy
  Review Source, its request command and its coherence refusal stay as they
  are when the legacy commands run.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight: a review blocked by an
unselected provider's configuration, a disposition written on refusal, a
reviewer re-asked at an unchanged head, a fix that clears its own head, or a
retry that moves a corrected archive forward would each pass a happy-path
test.

The outside-evidence row rests on dispositions this Spec did not write: the
Review dispositions sections of #258 (Windows portability, dismissed; TechSpec
premise, fixed) and #259 (stale build, dismissed), and the corrective Task of
#257 (round-one findings, fixed by task_06). The QA gate replays each through
the built binary and records the pull request as the origin.

## Research basis

The bodies of #257, #258 and #259 record every pre-PR review outcome of the
first efficiency wave as prose, and they are the only record of those
dispositions. The
retired review policy, now in
`docs/history/specs/0126-agent-review-before-pull-request/`, carried Core
Features 4 and 5, and the retired durable workflow in
`docs/history/specs/0127-durable-unattended-spec-workflow/` carried Core
Feature 9. Both left the late-correction authority pending. The adopted
sources are indexed in [references/_index.md](references/_index.md).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
