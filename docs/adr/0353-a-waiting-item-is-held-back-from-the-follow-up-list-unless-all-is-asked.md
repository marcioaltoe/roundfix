---
status: accepted
created_at: 2026-11-03T18:00:00Z
updated_at: 2026-11-03T18:00:00Z
deprecated_at: null
superseded_by: null
---

# A waiting item is held back from the follow-up list unless --all is asked

The Follow-up Command holds Waiting Items back by default and counts them on
their source lines. `--all` lists them after ranked work, with their reason
and date, so an outside dependency is visible without competing for action.

## Context

[ADR-0319](0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md)
orders Follow-up Items in four classes: "The order is computed in code, never
asked." [ADR-0338](0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md)
requires the follow-up mode to claim the handoff and account for every
handed-in message: "No message the mode handed in stays pending when the turn
ends". Neither records an item that waits on something outside the session.

On 2026-11-02 Fluxus listed three Findings and a Backlog Entry that each
waited on someone else, with only free text explaining why. On 2026-11-03 the
maintainer accepted `waiting_on` and `since` and answered:

> Os tres e por padrão sem waiting quando usado só --follow-up, mas preciso
> de uma flag ou parametro --all (Algo assim) para incluir tudo, incluindo os
> waiting

## Decision

An open Backlog Entry or a pending or partial Finding may carry a Waiting
Declaration: `waiting_on`, one non-empty line naming who or what it waits
on, and `since`, the `YYYY-MM-DD` date the wait began. A valid declaration
makes it a Waiting Item in a fifth class, `waiting`, outside the four ranked
classes of ADR-0319.

The default list holds Waiting Items back and counts them on each source's
line; its total count still includes every item read. `--all` lists Waiting
Items after every ranked item, oldest `since` first, then reference in byte
order. A Waiting Item is never sent to the judge, with or without --all.
The status view refuses `--all`, reads no Backlog Entry or Finding, and
prints no waiting count.

The command that reads a declaration validates it. A malformed declaration
leaves its item in its ordinary class, with the validation reason on the
item, and never makes a source unavailable. An item declaring neither field
keeps its previous report bytes.

## Considered Options

- List Waiting Items by default: rejected because the maintainer asked for
  the default list to leave them out and a flag to include them.
- Drop them without a count: rejected because each item read must remain
  visible as listed work or in its source's held-back count.
- Judge Waiting Items under `--all`: rejected because listing an outside
  dependency does not make it ranked work for the session.
- Add a Spec Check finding for a malformed declaration: rejected in this
  change; validation belongs to the command that reads it.
- Change the Baseline contracts in the same change: deferred to the open
  Backlog Entry
  `docs/backlog/2026-11-03-the-backlog-and-findings-contracts-name-the-waiting-declaration.md`.
- Add a waiting count to status: rejected because status reads neither
  Backlog Entries nor Findings.

## Consequences

The list gains `--all`, the optional source JSON key `waiting` and the
optional item keys `waiting_on`, `since` and `waiting_error`. A Waiting Item
has priority `5`, name and basis `waiting`, and null answers, blocks and
reason. Judge totals count ranked items only.

The follow-up mode writes the declaration on an item its plan leaves waiting
(and updates a Finding's `updated_at`), and its report names each listed
Finding's and Backlog Entry's destination: a new Spec, an existing Spec it
was grouped into, or its wait and date. Malformed declarations are discovered
when the command runs, rather than by Spec Check.
