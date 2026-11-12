---
status: accepted
created_at: 2026-11-09T12:00:00Z
updated_at: 2026-11-11T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Only an ADR a Spec applies or changes draws in the ADRs that cite it

`SC-ADR-RELATED` reports an accepted ADR that cites an ADR a Spec lists under
Active ADR obligations when the Spec does not list the citing ADR itself
(ADR-0093), and ADR-0168 limits it to ADRs that predate the Spec. A Spec
answers the gap by listing the citing ADR, and until now a listing that only
records why an ADR does not apply counted the same as one the Spec applies.
Each such listing drew in the ADRs that cite the listed one, and the next
check drew in the ADRs that cite those.

On 2026-11-07 Fiscus reported the effect. After its Spec 0069 merged two ADRs
that cite a third, its Spec 0081 failed the strict check before
`roundfix deliver start`, and answering the gap took seven rounds and 72 ADRs,
every one of them recorded as not applying. In this repository Spec 0358
listed 34 ADRs as "unchanged; accounted for in Project Constraints" for the
same reason. On 2026-11-09 the maintainer chose to amend the rule: "Só ADRs
aplicadas ou alteradas pela Spec entram no fecho; uma ADR “não se aplica”
deixa de propagar."

## Decision

**Only an ADR the Spec applies or changes seeds the related-ADR candidates.**

- An ADR the PRD lists in its Active ADR obligations row is applied unless
  the TechSpec's `## ADR Alternatives` line for it, `- ADR-NNNN: <relation>`,
  begins with `not applicable`, matched without regard to case.
- An ADR marked `not applicable` stays listed. `SC-ADR-RELATED` never reports
  it, it satisfies `SC-ADR-UNLISTED`, and `SC-ADR-ALTERNATIVES-UNRELATED`
  still requires its line.
- `SC-ADR-RELATED` reports an unlisted accepted ADR that cites an applied ADR.
  An ADR that cites only ADRs marked `not applicable` is not reported.
- A Spec without a TechSpec has no ADR Alternatives lines, so every listed ADR
  seeds the candidates, as before.
- The horizon of ADR-0168, the depth of one citation edge and the gap
  severity are unchanged.

The check reads a fixed leading token in a line the Spec already writes for
every listed ADR, so it still compares written declarations and never judges
whether an ADR is relevant (ADR-0093).

**Amended by [ADR-0389](0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md) (2026-11-11).** The leading token is read from a closed list: `not applicable`, `não se aplica` with a precomposed or decomposed `ã`, and `nao se aplica`, matched without regard to case and ending before a character that is not an ASCII letter, digit or underscore. Everything else in this decision still holds.

## Considered Options

- A marker phrase inside the PRD's Active ADR obligations row, such as a
  "Not applicable:" sentence: rejected, because the row is free prose with no
  boundary between ADRs, and a phrase written in one Spec's language binds no
  other. ADR-0389 keeps this option rejected: it reads the Portuguese token
  only in the TechSpec's ADR Alternatives line, from a closed list every
  checkout shares.
- A fifth Project Constraints row for the ADRs that do not apply: rejected,
  because it changes the four-row contract that every adopter's templates and
  the `SC-CONSTRAINT-*` checks hold.
- A limit on the rounds or the number of ADRs the candidates may reach:
  rejected, because the check already stops at one citation edge; the growth
  came from what seeded it.
- Removing `SC-ADR-RELATED`: rejected, because the omission ADR-0093 measured
  in Spec 0056, ADR-0055, is still caught when the cited ADR is one the Spec
  applies.
