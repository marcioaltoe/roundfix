---
status: accepted
created_at: 2026-10-21T04:00:00Z
updated_at: 2026-10-21T04:00:00Z
deprecated_at: null
superseded_by: null
---

# A queue archive subject and pull request title keep the slug while it fits 100 characters

The Delivery Queue's archive commit `docs: archive <slug>` and its pull request
title `feat: deliver <slug>` grow with the slug, and a squash merge appends
` (#<number>)` to the title, so a long slug passes the 100-character
`header-max-length` of the conventional commitlint configuration. This
repository's own squash commit `9e92ac39` has a 117-character subject. The
maintainer chose "Slug se couber (Recomendado)": both messages keep the slug
while it fits, and name the Spec by its Spec Number otherwise, as ADR-0302 did
for the Pull Request Record Commit.

## Decision

**Length is measured as commitlint measures it.** A header's length is its
count of UTF-16 code units, the JavaScript `length` commitlint's
`header-max-length` reads. The limit is 100.

**The archive subject.** `docs: archive <slug>` when it is at most 100 long,
byte for byte as before. Otherwise `docs: archive the <spec number> Spec`, or
`docs: archive the Spec` for a slug without a Spec Number. The Spec Number is
ADR-0302's: one to ten ASCII digits that start the slug and are followed by
`-`, read by the same helper as the record subject.

**The pull request title.** A squash merge's subject is the title followed by
` (#<number>)`, and the number is unknown when the title is written, so the
title reserves the width of ` (#99999)`, 9 units, and is at most 91 long.
`feat: deliver <slug>` is used when it is at most 91 long, byte for byte as
before. Otherwise the title is `feat: deliver <spec number>: <Spec title>`, or
`feat: deliver <Spec title>` without a Spec Number. The Spec title is the
Archive Record's `title`, which the archive took from the PRD's first heading.
Its whitespace runs collapse to one space; when the title would pass 91 it is
cut at the last space that keeps it within 91, or inside its first word when
that word alone is too long; trailing `.`, `,`, `;`, `:` and spaces are then
removed. An empty result gives `feat: deliver <spec number>`, or
`feat: deliver the Spec`. The Pull Request body still names the full slug.

**Readers keep working.** Source Revision Recovery and the legacy delivery
reader read only the trailing `(#<number>)`; the release plan reads only the
type before the first colon; the queue adopts an existing archive commit by
proving its change, never by reading its message, and finds an existing pull
request by its head branch, never by its title.

## Considered Options

- Always naming the Spec by number was rejected by the maintainer's decision:
  short slugs, almost every slug, keep today's bytes.
- Cutting the slug was rejected, as in ADR-0302: a cut slug names no Spec
  reliably.
- Reserving the ten digits Source Revision Recovery accepts was rejected: it
  shortens every long title by five more characters for repositories past
  99,999 pull requests, which no adopter approaches.
- Renaming the pull request after it opens, with the real number, was
  rejected: it is an extra forge write on every delivery and races with the
  merge.

## Consequences

- A squash subject stays within 100 characters for every slug while the pull
  request number is below 100,000; past that, each extra digit adds one.
- A pull request an older binary opened keeps its title; the queue finds it by
  branch and never renames it.
- A slug of 78 to 86 units, whose title fits 100 but whose squash subject did
  not, now gets the number form.
