---
status: accepted
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T19:30:00Z
deprecated_at: null
superseded_by: null
---

# The corrections record is one append-only file per Spec under Roundfix Home

The declarations of a Spec's Corrections and the Correction Signs of its
delivery are kept in one JSON Lines file per Spec under Roundfix Home, at
`corrections/<repository key>/<slug>.jsonl`. A line is appended and never
rewritten, the standing state of a Correction is its last declaration line,
and the Run Retention Sweep removes a Spec's file whole. The record lives
there because nothing else outlives a delivery: the squash merge keeps one
commit, the Item Branch is released, and a new Delivery Queue replaces the
finished one.

## Context

ADR-0417 leaves the place of the record to the TechSpec and asks that it
outlive the Delivery Queue item and the Item Branch. Read at `a8236a87`: the
review disposition ledger is keyed by the reviewed head and names no Spec; a
queue item stores its candidate heads and is deleted with its queue; and the
Phase Record already keeps one JSON object per line under Roundfix Home, which
the Run Retention Sweep trims (ADR-0407).

Published practice treats a correction as a new record, not as an edit. The
Secondbrain's copy of Designing Data-Intensive Applications, chapter 11, says
of accountants: "If a mistake is made, accountants don't erase or change the
incorrect transaction in the ledger — instead, they add another transaction
that compensates for the mistake". POSIX says of a file opened with
`O_APPEND` that "the file offset shall be set to the end of the file prior to
each write and no intervening file modification operation shall occur between
changing the file offset and the write operation".

## Decision

**One file per Spec.** The file is
`<Roundfix Home>/corrections/<repository key>/<slug>.jsonl`, where the
repository key is the first 16 hexadecimal characters of the Repository
Identity (ADR-0421). Every line is one JSON object under the schema
`roundfix/corrections/v1`, of kind `declaration` or `signs`.

**A declaration line is the whole state of its Correction.** It holds the
identifier, the class, the artifact, the reason, every Correction Sign the
Correction names, whether it is withdrawn, and who declared it and when. A
later declaration of the same identifier adds the signs it names to the ones
already named and replaces the rest, so a reader takes the last line of each
identifier and joins nothing. A withdrawal is one more line, and it names no
sign.

**A signs line is the delivery as the Queue Owner read it.** Each time an item
reaches its merge the Queue Owner appends the signs of that head, with the
number of dismissed findings. The last signs line counts. A Spec with no signs
line was not delivered by this machine's queue.

**Minting an identifier takes the lock; nothing else does.** `C<n>` is the
next free number of the file, read and appended under the file lock the
disposition ledger already uses. A signs line is one `O_APPEND` write.

**A reader skips a line it cannot decode and says how many it skipped.** It
never rewrites the file.

**Retention removes a file whole.** The Run Retention Sweep removes a Spec's
file when its newest line is older than the cutoff, in the same step that
removes Phase Record day files, with or without a Run Database. A delivery's
record is therefore whole or absent, and an identifier is never minted twice
for a file that still exists.

## Considered Options

- **A trailer on the correction commit as the record.** Rejected: the squash
  merge keeps one commit, as ADR-0417 says.
- **Git notes on the Item Branch's commits.** Rejected: a note is attached to
  an object, and the squash commit is another object; the branch's commits are
  unreachable once the branch is released.
- **Tables in the Run Database.** Rejected: a new table needs a schema
  migration, the Delivery Queue's own rows are deleted with their queue, and
  ADR-0255 keeps the queue tables out of retention.
- **A file beside the disposition ledger, in the Artifact Directory.**
  Rejected: that directory is keyed by a checkout path, an `artifact_dir`
  override moves it out of Roundfix Home, the storage report measures it, and
  nothing reclaims the ledger.
- **Day files, as the Phase Record.** Rejected: the sweep would remove the
  signs of a delivery and keep a declaration made a day later, and an
  identifier could be minted again once its first day was swept.
- **Fields of the Archive Record.** Not taken: ADR-0417 leaves the Archive
  Record as it is.

## Consequences

- The record belongs to one machine. A Spec delivered elsewhere shows what was
  declared for it here and no sign.
- Neither the storage report nor GC sanitation reads the directory: it lies
  outside every Artifact Root and outside the Run Database.
- A declaration costs one lock, one read of a small file and one append.
- A torn or damaged line costs that line, and the report names the count.
