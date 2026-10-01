---
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: archived
created: 2026-09-30
surfaces: [backend, cli, docs]
archived: "2026-10-01"
source_slug: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
---


# A reviewer that validates its findings and remembers its rounds

The Pre-PR Review Command runs the configured reviewer over the candidate
before a Pull Request exists. With the Delivery Queue, each finding parks the
item until an operator acts. Three problems make that costly.

- **The reviewer does not know the delivery's own order.** On 2026-09-30 the
  review of Spec 0199 raised two findings. One said the QA Report audits a
  commit other than the candidate head, which is true because a report cannot
  name the commit that records it. The other said a Task was `completed` while
  its Result says the Daemon owns status, which is how the Daemon settles a
  Task. Both were dismissed with evidence. The Finding of 2026-09-30 about the
  v0.22.0 queue records it as its fifth class of intervention. This
  repository's disposition ledger holds 34 dispositions written between
  2026-09-29 and 2026-09-30: 23 `fixed` and 11 `dismissed`. Five of the 11
  restate the delivery's order: two QA Report commits (Specs 0182 and 0199),
  one Daemon settlement (0199) and two planning candidates whose Tasks were
  pending by design (0188).
- **Nothing checks a finding before it parks the item.** A finding is
  whatever the reviewer's final message lists. The prompt asks for a file and
  a line, but nothing checks that the line is part of the candidate, or that
  the finding states what breaks.
- **Every round starts from nothing.** Each review opens a new reviewer
  session, reads the whole merge-base diff again and knows nothing of the last
  round's findings or how each was disposed. The maintainer's rule of
  2026-09-09 allows at most two review rounds per candidate. No code enforces
  it, so it lives only in operating notes.

This Spec tells the reviewer the delivery's conventions, and validates every
finding before it can park an item. A dismissed finding is recorded with its
reason and is never dropped. Reviews of one candidate become one Reviewer
Lineage of at most two rounds: the second round reviews only what changed, and
it carries the first round's findings and dispositions.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. A Delivery
  Convention is named by a fixed ordinal (`C1` to `C4`), and a Reviewer
  Lineage by the record that already names its checkout, merge base and
  provider. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git, local files and the
  configured local Agent runtime through acpx. No credential is read, and no
  network call or remote service is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0196 (this Spec) validates each
  finding against the candidate diff and the Delivery Conventions before it
  parks, and ADR-0197 (this Spec) bounds a Reviewer Lineage at two rounds.
  Both are recorded before implementation. ADR-0153 keeps the Pre-PR Review
  Policy explicit. ADR-0153: "A failed, unavailable, incomplete, stale or
  missing enabled review is an error, never implicit permission to select
  none." Validation never selects `none`, and a validator that cannot run
  leaves every anchored finding standing. ADR-0169 makes the merge base the
  diff's before side, and an anchor is checked against that diff. ADR-0169:
  "Roundfix now resolves the merge base of the candidate head and the selected
  base ref once." ADR-0174 reads the verdict from the final message, and the
  finding grammar reads that message. ADR-0174: "The reviewer's verdict is
  read from the last non-blank message". ADR-0018 keeps one Agent Session per
  Run, and a pre-PR review is not a Run. ADR-0018: "sessions are never reused
  across Runs". ADR-0051 gives Tasks and QA their own Agent Sessions, which a
  reviewer lineage does not touch. ADR-0165 parks a candidate that changed an
  archived Spec, and a standing finding still takes that park. ADR-0017 drives
  every runtime through acpx, and a continued session is an acpx named
  session. ADR-0151 keeps the review profile, and the validator uses its
  selection. This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156. ADR-0093 and ADR-0094 check its
  consistency, ADR-0166 records undeclared Task paths, ADR-0167 keeps the
  pre-PR Pull Request row from deciding a qualifying partial, ADR-0168 bounds
  the related-ADR check and ADR-0176 reads only authored text. All of them
  hold. ADR-0097 cites ADR-0080 but carries a QA row forward on unmoved
  evidence, and ADR-0182 cites ADR-0096 but moves mechanical facts to Task
  settlement. This Spec changes neither, so they do not apply. ADR-0194 and ADR-0195,
  written by Specs of the same wave, cite ADR-0097 but decide how the Daemon
  records QA evidence and which QA rows are observed on every pass. This
  Spec's gate uses neither, so they do not apply. ADR-0183 and ADR-0184 add receipts and Surface Transcripts, which the
  TechSpec provides. ADR-0189 versions an owned skill by its content, and the
  Roundfix skill's version is raised with its text. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário"), recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Prerequisites

- Spec 0194 is merged first. It moves the review section of the Roundfix
  skill and the command guide into per-command reference files, and this
  Spec's documentation Task edits those files.
- Spec 0191 is merged first. It adds the Surface Transcript and receipt checks
  that this Spec's TechSpec is written for.

The Delivery Queue does not enforce this order. The operator orders the queue.

## Goals

- A finding that restates the delivery's designed order never parks a
  delivery, and the record says which convention it restated.
- A finding parks only when it names a line the candidate diff shows and
  states a failure.
- A dismissed finding is always visible in the record, with the rule and the
  reason that dismissed it.
- A second review round reads only the change since the first, and knows the
  first round's findings and dispositions.
- The two-round ceiling is enforced by the product, and closing a candidate
  at the ceiling takes a recorded operator disposition for every standing
  finding.

## User Stories

1. As an operator running the Delivery Queue, I want the reviewer to know that
   the QA Report commit, the Daemon's settlement and the archive commit are
   designed, so that the queue does not park on them.
2. As an operator, I want a finding that cites no line of the candidate to be
   set aside with its reason, so that I read it once and do not have to
   dismiss it by hand.
3. As a maintainer, I want every dismissal kept in the record, so that I can
   audit what validation set aside.
4. As an operator retrying a parked item after a fix, I want the second review
   to look at the fix and remember what it raised before, so that it neither
   rereads the whole candidate nor repeats a dismissed finding.
5. As a maintainer, I want the review to stop at two rounds and close only on
   my recorded dispositions, so that a third round is never spent.

## Core Features

1. **Delivery Conventions in the prompt.** The review prompt carries a closed,
   numbered list of what a delivery writes by design:
   - `C1`: the QA Report records the head it audited, and is committed after
     it.
   - `C2`: the Daemon writes a Task's status and its Daemon-owned sections
     after its Verification passes.
   - `C3`: the archive commit moves a Spec to the archive root and stamps it.
   - `C4`: a planning candidate authors a Spec whose Tasks are all pending,
     with no QA Report.

   The prompt tells the reviewer not to report these as defects. Each
   convention also owns a region of the candidate, and only a finding anchored
   in that region can be dismissed as restating it.
2. **A finding grammar.** The prompt keeps today's verdict grammar and adds to
   it. Each finding opens with an anchor, a path and a line or a line range,
   and states its failure in a clause that starts with `Failure:`.
3. **Anchor validation.** Roundfix reads each finding's anchor. A finding whose
   anchor names no line the candidate diff shows, with its context lines, is
   dismissed as `unanchored`. An answer with findings in which no finding
   carries a readable anchor is blocked, because the reviewer did not follow
   the grammar.
4. **Convention validation.** A finding anchored in a convention's region, or
   one without a failure clause, goes to a validator. The validator is a
   sealed prompt on the review profile's selection, and it may use no tool. It
   may dismiss the finding as restating that convention, or, when the finding
   has no failure clause, as stating no failure. Every other finding stands.
   When the validator cannot run, times out, uses a tool or answers outside
   its grammar, every finding it was asked about stands, and the record says
   why.
5. **Recorded dismissals.** Each finding in the record carries its anchor and
   its validation: `stands`, or `dismissed-by-validation` with its rule and
   reason. Only a standing finding counts toward the verdict. When no finding
   stands after validation and the operator's evidence-backed dismissals, the
   record reads `findings-dismissed` and the command exits `0`, as today. The
   operator cannot dispose of a finding that validation already dismissed.
6. **A Reviewer Lineage.** Reviews in one checkout, with the same provider and
   merge base, each at a head that descends from the last reviewed head, are
   one lineage. Round 1 reviews the full candidate diff. Round 2 reviews only
   the diff from the round-1 head, and its prompt carries every round-1
   finding with its validation and disposition. A rebase, a changed merge
   base, a changed provider, a head that does not descend, or a blocked
   previous review starts a new lineage at round 1.
7. **A continued reviewer session.** When round 1 leaves a finding standing,
   its Agent Session stays open, and round 2 continues it through acpx on the
   same selection. The record names the session and the ACP session id of
   each round. It reports `continued` only when the two ids match. A runtime
   that cannot resume gets a new session, and round 2 still carries the
   recorded conversation. The session is closed after round 2, after a round
   that leaves nothing standing, when a reuse resolves to `findings-dismissed`
   and when a new lineage replaces it.
8. **The two-round ceiling.** A third review in a lineage never calls the
   reviewer. It records `ceiling-closed` and exits `0` when every standing
   round-2 finding has an operator disposition: fixed by a commit the head
   contains, or dismissed with evidence. Otherwise it prints a blocked result
   that names each finding to dispose and exits `2`. It keeps the round-2
   record in place, so `roundfix review dispose` still reads it. The Delivery
   Queue advances on `ceiling-closed` as it does on `findings-dismissed`.
9. **Documentation.** The review reference of the Roundfix skill, its mirror
   and the review command guide describe the conventions, the grammar,
   validation, the lineage, the continued session and the ceiling. The
   skill's version rises with its text.

## Declared breaks

- Findings that name a line outside the candidate diff no longer park an item,
  and existing tests whose fake reviewer answers cite such lines are updated to
  cite a line of their fixture diff, keeping their names.
- An answer whose findings carry no readable anchor was `findings` and is now
  `blocked`.
- The record gains `anchor` and `validation` on each finding, and `validation`
  and `lineage` at the top level. Its outcome set gains `ceiling-closed`. A
  Roundfix older than this Spec refuses to read a `ceiling-closed` record,
  which makes it review again.
- A round-2 review sends the delta diff instead of the full candidate diff.
- A review that leaves findings standing no longer closes its Agent Session.

## Non-Goals / Out of Scope

- Changing the verdict grammar's first line, `No findings.` or `Findings:`, or
  the rule that the final message is the answer.
- Changing the Pre-PR Review Policy, the providers, the review profile or its
  fallback order, or adding CodeRabbit support.
- Raising or configuring the two-round ceiling.
- A command that reinstates a finding validation dismissed. The finding's text
  stays in the record for anyone who disagrees.
- Using the TypeSafe judgment model as the validator. The maintainer's data
  grant covers Spec artifacts, not source code or diffs, and that model never
  decides a gate.
- Speaking ACP directly instead of through acpx.
- Closing reviewer sessions of checkouts that were removed while a lineage was
  still open.
- Changing PR-feedback Review Runs, `watch`, or their Rounds.

## Success Metrics

1. Success Metric: an answer whose only finding restates `C2`, anchored on a
   Task's `status` line, records `findings-dismissed` and exits `0` when the
   validator dismisses it. A finding with a failure clause, anchored on a
   Task's Requirements, stands and exits `1`, whatever the validator says.
2. Success Metric: an answer with one finding anchored outside the candidate
   diff and one anchored inside records the first as `dismissed-by-validation`
   (`unanchored`) and the second as `stands`, and exits `1`. An answer whose
   findings carry no readable anchor records `blocked` and exits `2`.
3. Success Metric: with a validator that fails, times out or uses a tool,
   every finding it was asked about stands and the record states the reason.
4. Success Metric: a second review at a descendant head sends a prompt that
   holds the delta diff, not the round-1 diff, and lists every round-1 finding
   with its validation and disposition. The record reports round 2.
5. Success Metric: after round 1 leaves a finding standing, no session is
   closed. Round 2 prepares the same session name and closes it. `continued`
   is true only when the runner reports the same ACP session id in both
   rounds.
6. Success Metric: a third review never calls the runner. It records
   `ceiling-closed` and exits `0` when every standing round-2 finding has a
   disposition, and prints a blocked result that names the missing ones, with
   exit `2`, when one does not. The round-2 record stays in place.
7. Success Metric: replaying the 34 findings of this repository's disposition
   ledger against their candidates' merge-base diffs shows how many anchor in
   the diff. It shows which of them the convention regions admit: the five
   convention dismissals, and none of the 23 `fixed` findings.

## Recorded limits

- Anchor validation alone would have dismissed none of the 34 recorded
  findings, because all of them cite a line the diff shows. Its value is the
  guard it gives. Validation's gain comes from the conventions.
- A convention dismissal is still an Agent's judgment. It is bounded by the
  convention's region, and it never applies outside it.
- Whether a runtime resumes a session after acpx's idle queue owner exits is
  measured on each real round 2 and recorded, not proven by this Spec's tests.
  codex-acp 2.0.1, claude-agent-acp 0.84.0 and opencode advertise both
  `loadSession` and session resume. Cursor's ACP mode is reported to fail
  `session/load`, and on it round 2 runs with the recorded conversation only.
- A lineage left at round 1 by a checkout that is then removed keeps its acpx
  session record open. The adapter process exits after acpx's idle limit.
- The validator's sealed session carries the sealed-prompt name prefix that
  Baseline analysis already uses.
- A reviewer that re-raises a finding round 1 dismissed is validated again,
  not suppressed.

## Decisions

- **Validate, then park.** A finding that parks must be anchored in the diff,
  and it must not merely restate a Delivery Convention. See ADR-0196.
- **Bound the Agent's judgment by a mechanical region.** The validator can
  dismiss only inside a convention's region, or only a finding with no failure
  clause. Anything else stands, so a validator mistake cannot hide a finding
  on authored code.
- **Fail closed.** A validator that cannot answer leaves findings standing.
  An answer that ignores the grammar entirely blocks.
- **Recorded conversation first, resume second.** Round 2's prompt always
  carries round 1's findings and dispositions. A continued session adds the
  reviewer's own memory and is measured, never assumed. See ADR-0197.
- **Enforce the ceiling without a third round.** The change after round 2 is
  the final correction. It closes only on recorded operator dispositions.
  See ADR-0197.
- **No reinstate command.** Anyone who disagrees with a dismissal can read it
  in the record and raise it again. A command would add surface without a
  measured need.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight:

- a convention dismissal outside the convention's region;
- a no-failure dismissal of a finding that has a failure clause;
- a failed validator that drops a finding;
- a round-2 prompt that still holds the full diff;
- a third round that calls the reviewer;
- a ceiling that closes with a finding that has no disposition.

The outside-evidence rows rest on sources this Spec did not produce:

- The disposition ledger at
  `~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review-dispositions.jsonl`,
  written by other sessions on 2026-09-29 and 2026-09-30. It holds 34
  dispositions, 23 `fixed` and 11 `dismissed`. The candidate heads it names
  are still in this repository's object store. Replayed during authoring
  against each head's merge base with `origin/main`, every one of the 34
  anchors cites a line of its diff. The five convention dismissals fall in the
  regions of `C1` (Specs 0182 and 0199 QA Reports), `C2` (0199's `status`
  line) and `C4` (0188's all-pending Task files). None of the 23 `fixed`
  findings does.
- The Agent Client Protocol's session setup page
  (<https://agentclientprotocol.com/protocol/v1/session-setup>). It says that
  `session/resume` restores a session "without replaying the conversation
  history", and that a client must check `sessionCapabilities.resume` before
  calling it.
- The installed adapters, read without running them. codex-acp 2.0.1 and
  claude-agent-acp 0.84.0 return `loadSession: true` and
  `sessionCapabilities.resume`. acpx 0.19.3 reattaches an open named session
  and, when it must respawn the adapter, tries resume, then load, then a new
  session.
- GitHub's REST reference for pull request review comments
  (<https://docs.github.com/en/rest/pulls/comments>). It describes review
  comments as "comments made on a portion of the unified diff", which is the
  anchor rule this Spec adopts.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "pre-PR review false positive findings validation persistent
reviewer session" --all --files --min-score 0.3`. It returned this
repository's mirrors of Specs 0126, 0143, 0179 and 0185 and ADR-0153, which
this Spec builds on. A second query, in Portuguese, on validating agent review
findings and a reviewer session that continues across rounds, returned the
concept page `wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md`.
That page argues that a checker must fail a known negative case and that a
deterministic component should admit a model's claim only against observed
state. The research capture of 2026-09-30 on Spec orchestrators lists "a
reviewer session that continues across rounds, with a findings file validated
by a script" as a mechanism worth studying. It also records a verifier pattern
in which a model proposes a verdict with a literal citation and code proves
the citation exists in the source. Exa found published practice for the same
shape:

- GitHub's review comment API, above, rejects a comment whose line is not part
  of the diff.
- A pull request to an open-source coding agent adds a false-positive
  exclusion list and a per-finding verification step before findings are
  reported (<https://github.com/QwenLM/qwen-code/pull/2687>).
- A 2026 paper finds that LLM reviewers systematically reject correct code
  when asked for explanations. It reports that grounding verdicts in
  executable evidence reduces those false rejections
  (<https://arxiv.org/abs/2603.00539>).

No code, name or text from those projects is used.

## Technical candidate

The [_techspec.md](_techspec.md) records the record shape, the validator and
lineage algorithms, the command transcripts, coverage and build order.
