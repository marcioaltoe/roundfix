---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A cancelled Spec leaves an Archive Record and returns its sources

On 2026-10-10 an adopter dropped eight Specs that already had Task Graphs, one
of them with 5 Tasks completed. No command could record it. The Supersede
Command needs a delivering Spec, the QA Archive Override needs every non-QA
Task completed, and Task status belongs to the Daemon. So the adopter deleted
the folders, and what the dropped Spec delivered has no history
(`docs/references/2026-10-12-a-decomposed-spec-dropped-without-a-record.md`).

On 2026-10-11 the maintainer chose "Comando com registro":
`roundfix cancel <slug> --reason <text>` writes an Archive Record with the
disposition `cancelled`. The record keeps the reason, the date, and the Tasks
completed and pending at that moment. The command returns the adopted sources
under the adoption contract and refuses while a Run or a Delivery Queue item
is active for the Spec.

- **The record.** A cancelled record keeps the schema
  `roundfix/archive-record/v1`. It adds four fields that only `cancelled`
  carries: `cancel_reason`, `completed_tasks`, `pending_tasks` and
  `returned_sources`. `archived` is the cancellation date. The QA fields
  record what was observed, and the outcome paragraph is the PRD's opening
  paragraph, because a cancelled Spec has no outcome of its own.
- **Source Return.** Each adopted source indexed as `backlog` goes back to
  `docs/backlog/` with `status: open` and `spec: null`. Each source indexed as
  `finding` goes back to `docs/findings/` with `status: pending`, a new
  `updated_at` and a dated addendum, because a Finding is evidence and never
  becomes an open Backlog Entry. Every other byte of a returned source is
  kept. Other indexed sources, such as an Inbox Entry, stay in Git with the
  folder.
- **What a cancelled record means to readers.** Only the Tasks it lists as
  completed count as completed. It never meets a Prerequisite Spec. It never
  delivers another Spec's content through the Supersede Command. A Delivery
  Retry refuses it. A review, the Run causes report, the Spec check and the
  History Sanitize Command read it as they read any record. Sanitize converts
  only Legacy Archive Folders, so it never rewrites the record.
- **Supersession stays separate.** A Spec that carries a Supersession is
  archived, not cancelled. A Spec replaced by another one is cancelled with a
  reason that names the replacement, and the replacement adopts the returned
  sources.
- **Refusals.** It refuses while an Active Run holds the Spec, while a
  Delivery Queue item for the Spec is unfinished (any stage but `merged` and
  `parked`), while an active Spec names it in `requires`, and while a
  non-Markdown file names its active directory.

This extends ADR-0247 and ADR-0248, which list the dispositions, and ADR-0193:
"A prerequisite is met when its archived Spec folder is on the refreshed
default branch". A cancelled record is on that branch but meets no
prerequisite.

## Considered Options

- Deleting the folder by hand, which is what the adopter did, leaves no
  record.
- Archiving with a widened QA Archive Override would claim archive
  eligibility for Tasks that never ran. ADR-0154 reserves the override for a
  QA gate whose non-QA Tasks completed.
- Letting the Supersede Command accept a decomposed Spec would invent a
  delivering Spec where none exists.
- Returning a Finding to `docs/backlog/` would turn evidence into intent
  without a human choice, which the docs layout forbids.

## Consequences

A binary older than this decision refuses a `cancelled` record with "unknown
archive disposition", so every reader of a repository with a cancelled record
needs this version. A parked Delivery Queue item does not block a
cancellation, so its branch and any open pull request stay until the operator
closes them; a Delivery Retry of that item is refused. The docs layout guide
of the Baseline describes the command, so adopters learn it through the
Managed Refresh.
