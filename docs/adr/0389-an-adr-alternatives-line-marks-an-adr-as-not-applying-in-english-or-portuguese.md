---
status: accepted
created_at: 2026-11-11T00:00:00Z
updated_at: 2026-11-11T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An ADR Alternatives line marks an ADR as not applying in English or Portuguese

Amends ADR-0377 (2026-11-11): the leading token that marks a listed ADR as not
applying is read from a closed list of English and Portuguese forms, not from
the English token alone.

ADR-0377 made only an ADR a Spec applies or changes seed the related-ADR
candidates of `SC-ADR-RELATED`, and read the exception from one fixed leading
token, `not applicable`, in the TechSpec's `## ADR Alternatives` line. On
2026-11-10 Fiscus, which writes its Specs in Portuguese, reported that the 76
ADR lines behind the cascade ADR-0377 removed begin with "não se aplica". The
check reads them as applied ADRs, so a Portuguese Spec keeps the large closure
until its author rewrites every line with the English token, and nothing tells
the author: a line read as applied raises no finding, it only draws in more
ADRs. The maintainer chose the change the same day ("Aceitar o marcador em
português nas ADR Alternatives").

## Decision

**A closed list of leading tokens marks an ADR as not applying.** An ADR
Alternatives line `- ADR-NNNN: <relation>` marks its ADR as not applying when
the relation begins with one of these tokens, matched without regard to case:

- `not applicable`;
- `não se aplica`, with `não` written either with the precomposed `ã` or with
  `a` followed by the combining tilde, the two canonically equivalent forms;
- `nao se aplica`, the same token without the accent.

The token ends at the end of the line or before a character that is not an
ASCII letter, digit or underscore, so `não se aplicam` and `not applicable`
followed by a letter are not tokens. Every other part of ADR-0377 holds: the
token is read only in the TechSpec's `## ADR Alternatives` line, an ADR it
marks stays listed and satisfies `SC-ADR-UNLISTED`, `SC-ADR-ALTERNATIVES-UNRELATED`
still requires its line, and a Spec without a TechSpec seeds with every listed
ADR. The `SC-ADR-RELATED` fix text names both the English and the accented
Portuguese token, and the authoring skills prescribe the English token first
and accept the Portuguese one.

The list is part of the check, not of a repository's configuration, so every
checkout judges a Spec the same way; adding a language is an amendment of
this decision.

## Considered Options

- A Project Config key naming each adopter's own token: rejected, because the
  same Spec would pass or fail by the checkout's configuration, and the two
  languages the fleet writes Specs in are known.
- No code change and a skill rule that prescribes the English token whatever
  the Spec's language: rejected, because a line written otherwise is read as
  applied without any finding, so the rule fails silently, and Fiscus already
  holds 76 such lines.
- Normalizing the line's Unicode before matching: rejected, because the only
  character with two forms in the list is the `ã` of `não`, which the list
  names in both forms, and normalization would need a module this repository
  does not require.
- Accepting synonyms such as `não aplicável` or `n/a`: rejected, because an
  open list of phrasings turns a declaration into a guess; the list stays
  closed and every member is a fixed token.
- Reading the marker in the PRD's Active ADR obligations row: still rejected
  for ADR-0377's reason, that the row is free prose with no boundary between
  ADRs.

## Consequences

A Portuguese Spec whose lines already begin with "não se aplica" draws in no
ADR that cites only those ADRs, with no rewrite. A Spec in a third language
still writes the English token until this decision is amended. A line that
begins with a token in the list but means the ADR applies, an unlikely use,
now stops seeding; the author fixes it by rewording the relation.
