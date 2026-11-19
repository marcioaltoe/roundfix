---
status: accepted
created_at: 2026-11-19T11:00:00Z
updated_at: 2026-11-19T11:00:00Z
deprecated_at: null
superseded_by: null
---

# A rule the user keeps carries a Kept Line, and is asked about again only with a new fact

When the user answers a Rule Challenge with keep, the rule gains one dated
line that says what was asked and why the rule stays. A later author reads
the line and asks about that rule again only with a fact the line does not
hold. The line is added to an accepted ADR in place, without a new status and
without a new ADR.

## Context

[ADR-0436](0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md)
makes an author ask the user whether a rule that blocks the goal is kept,
changed or retired. A keep answer that lives only in the Spec that asked is
lost: an archived Spec leaves the tree, and the next author meets the same
rule with no trace of the question. On 2026-11-18 the maintainer chose where
the answer lives, "Na própria regra (Recommended)": the kept rule or ADR
gains a dated line with what was asked and why it stayed, and it is asked
about again only with a new fact. With the approval of the Idea he also
placed the line of a kept Baseline rule in an adopter: it goes in the
repository's own `docs/agents/specific-repository.md`.

Published practice is divided. AWS Prescriptive Guidance
(https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html)
says "When the team accepts an ADR, it becomes immutable", and Fowler
(https://martinfowler.com/bliki/ArchitectureDecisionRecord.html) says an
accepted ADR "should never be reopened or changed - instead it should be
superseded", while naming the record as the place for "any changes in the
product context that should trigger the team to reevaluate the decision". A
published ADR of the Shadewright project, "Amending accepted ADRs"
(https://shadewright.io/adr/0021-amending-accepted-adrs/), relaxes the strict
rule for a change under which "the decision itself stands" and gives every
amendment "a dated entry", because under the strict rule "small corrections
tend not to be made at all". This repository already annotates an accepted
ADR in place: an amended ADR carries a paragraph that starts `**Amended by`
and names the amending ADR and its date. A Kept Line changes no decision; it
records that the decision was questioned and confirmed.

## Decision

**The form.** A Kept Line is one line, `Kept (YYYY-MM-DD): <what was asked>;
<why it stays>.`, written after the rule's text. The date is the date of the
user's answer, and the reason is the user's. It names no Spec, because a
guide never references one and an archived Spec can be deleted.

**Lines stack.** A rule keeps every Kept Line, the newest last.

**An ADR takes it in place.** An accepted ADR gains the line after its
Decision text, its `updated_at` rises and its status stays `accepted`. No
status, no frontmatter field and no new ADR is added for a keep answer.

**A rule the Baseline owns.** A clause inside a Managed Region cannot hold
the line, in an adopter or in the repository that owns the Baseline, because
a Baseline guide says only what holds for every repository that reads it. Its
Kept Line goes in `docs/agents/specific-repository.md`, under a heading of
its own, naming the rule.

**Who writes it.** The Spec that asked: its PRD requires the line and one of
its Tasks adds it. When no Spec follows the question, the session that asked
adds the line.

**Asking again needs a new fact.** An author asks no Rule Challenge about a
rule that carries a Kept Line unless the question names a fact the earlier
one did not have. The record of that question holds the fact under `new
fact`.

## Considered Options

- **Keep the answer only in the Spec's record.** Rejected: the record leaves
  the tree with the Spec, and the rule is asked about again.
- **A new ADR for each keep answer.** Rejected: it scatters one decision over
  several files to say that nothing changed.
- **A review date in the ADR frontmatter.** Rejected: a date says when and
  not why, and the ADR layout belongs to the documentation rules, which this
  decision leaves as they are.
- **Strict immutability of an accepted ADR.** Rejected: the line reverses
  nothing, and this repository already writes amendment notices in place.
- **One place for every Kept Line, apart from the rules.** Rejected by the
  maintainer's choice: an author reads the rule, so the line is where the
  rule is, except for a rule the repository cannot edit.

## Consequences

- A rule that is asked about and kept is the first to say why it exists.
- Kept Lines accumulate on a rule that is asked about often; several lines on
  one rule are a sign that the rule deserves a change answer.
- What counts as a new fact is the author's judgment, and the record shows
  it.
- A repository's own rules guide gains a heading for the Kept Lines of
  Baseline rules.
