---
status: accepted
knowledge: true
amended_by: [ADR-0434]
created_at: 2026-11-16T14:05:00Z
updated_at: 2026-11-19T02:15:09Z
deprecated_at: null
superseded_by: null
---

# The status view lists Specs being authored and open Backlog Entries

The status view of the Follow-up Command shows one table of the repository's
Spec queue. After the Delivery Queue items and the planned Specs it lists the
Specs being authored, read from the Phase Record and the Spec folders, and the
open Backlog Entries. It amends ADR-0353, under which the status view read no
Backlog Entry, and ADR-0407, under which Roundfix reported nothing from the
Phase Record.

## Context

The maintainer plans by a queue that is longer than the one the status view
reads. On 2026-11-15 his next two pieces of work were a Spec planned only in a
Session Handoff and a feature recorded as an open Backlog Entry, and the
status showed neither. On 2026-11-16 he decided which rows the table gains:
"Entram Specs em Ideia, PRD ou TechSpec (pelo registro de fase da 0397) e
backlog aberto, marcados 'ainda não é Spec', com a espera declarada e sem
previsão de fim."

## Decision

**Specs being authored.** A Spec is being authored when it has no PRD in the
checkout's Spec folders and either a Spec folder there or a phase entry of an
Authoring Phase in the last seven days of the Phase Record. The entry counts
when its repository is the checkout or a working tree with the same `origin`
remote. The row names the latest Authoring Phase recorded and when it began.

**Amended by [ADR-0434](0434-a-renamed-spec-keeps-its-spec-number-and-is-one-row-of-the-spec-queue.md) (2026-11-18).** A phase entry for a slug whose folder is not under the Spec Root no longer counts when another slug of the same Spec Number has a folder there, an Archive Record, a Delivery Queue item or a later counted entry: that slug is a renamed Spec's earlier slug. The seven days, the repository match and the row a Spec folder makes still hold.

**Open Backlog Entries.** Every open Backlog Entry is a row marked as not yet
a Spec, a waiting one included, with its Waiting Declaration as its note and
no forecast.

**What stays.** The status view writes nothing, accepts no Session Handoff,
consumes no message, asks no judge and refuses `--all` and `--messages`. It
reads no Finding and no inbox note, and prints no waiting count. From the
Phase Record it reads phase entries only, and it reports no usage.

## Considered Options

- **Rows from the Session Handoff and from ai-memory messages.** Refused by
  the maintainer: the table shows what the repository and the Phase Record
  hold.
- **Only the Delivery Queue and the Specs with a PRD**, as before. Refused by
  the maintainer.
- **Matching a phase entry by its repository path alone.** Refused: an author
  often works in another clone of the repository, and its entries would never
  match.
- **Matching a phase entry by its Spec slug alone.** Refused: the Phase Record
  is one for the machine, and a slug does not say which repository it belongs
  to.
- **Every phase entry the record holds.** Refused: a Spec dropped during
  authoring would stay a row until Run Retention removed its day file.

## Consequences

- A Spec whose authoring recorded no phase, and whose folder is in another
  clone, is not a row until its plan reaches the checkout.
- A working tree that no longer exists is not asked for its remote, so its
  entries stop counting.
- The follow-up list is unchanged: it still holds Waiting Items back unless
  `--all` is given.
