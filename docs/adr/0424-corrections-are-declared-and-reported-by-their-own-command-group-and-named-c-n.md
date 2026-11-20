---
status: accepted
knowledge: true
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T20:30:00Z
deprecated_at: null
superseded_by: null
---

# Corrections are declared and reported by their own command group and named `C<n>`

A Correction is declared with `roundfix corrections declare` and read with
`roundfix corrections report`, a command group of its own. The command that
records a fixed finding takes the same declaration as flags. A Correction is
named `C<n>`, counted from 1 within its Spec, and a finding is named with the
head its review examined. The names are public CLI contract, so they are
fixed here once.

## Context

The PRD of Spec 0404 shows the content of each command and leaves the names
to the TechSpec. `roundfix runs causes` reads Runs and assigns a Cause Class
by signature, and the PRD keeps every byte of it. `roundfix review dispose`
resolves no Spec: it acts on the review record of the checkout it runs in,
and that record lists the Specs the reviewed candidate changes. A review
names its findings `F1`, `F2` and so on, anew at each reviewed head.

## Decision

**The group.** `roundfix corrections declare` writes one declaration and
`roundfix corrections report` is the Corrections Report. Both take `--spec`;
the report also takes `--since`, `--until` and `--format`, with the meanings
they have in `roundfix runs causes`.

**The declaration.** `--class` is one of `plan`, `implementation`,
`environment`, `operator` and `unclear`, or the answer `not-a-correction`
(ADR-0417); `--artifact` is one of `prd`,
`techspec` and `tasks`; `--reason` is the one line; `--finding`, `--commit`
and `--task` each name a Correction Sign and may repeat; `--correction <id>`
declares an existing Correction again, `--withdraw` withdraws it, and `--by`
gives the name of the person on whose word it is declared.

**With the disposition.** `roundfix review dispose <finding-id> --fixed-by
<commit>` takes `--class`, `--reason` and `--artifact`, or `--correction <id>`
in their place, and `--spec` when the reviewed candidate does not name exactly
one Spec. It checks every one of them before it writes, records the
disposition as before, and then the Correction, which names the finding and
its fixing commit.

**The identifier.** `C<n>`, where `n` is the next free whole number of the
Spec's corrections record, from 1, never reused.

**A finding's name.** `F<n>@<reviewed head>`. A command accepts `F<n>` alone
when exactly one fixed finding of that name was reviewed at an ancestor of
`HEAD`, and otherwise lists the candidates.

**Exit codes.** A refusal exits 2 and writes nothing; a record that cannot be
read or written exits 1; the report exits 0 with a source it could not read
named as unavailable.

## Considered Options

- **The report as `roundfix runs corrections`, or as a flag of `roundfix runs
  causes`.** Rejected: a Correction belongs to a Spec's delivery, not to a
  Run, and `roundfix runs causes` keeps its bytes.
- **The declaration only as flags of `roundfix review dispose`.** Rejected: a
  repaired QA finding, a commit by hand and a Corrective Task have no
  disposition.
- **`roundfix deliver declare`.** Rejected: a declaration is accepted for a
  Spec that no queue holds.
- **An identifier derived from the content or the time.** Rejected: it is
  typed by a person, and `C2` is typed without a mistake.
- **A finding named by its identifier alone.** Rejected: two reviews of one
  delivery each hold an `F1`.

## Consequences

- The Roundfix skill gains one reference for the group, and the user guide one
  command page.
- A Correction's identifier means nothing outside its Spec.
- An older Roundfix refuses the new flags of `roundfix review dispose` as
  unknown flags, with its present exit code.
