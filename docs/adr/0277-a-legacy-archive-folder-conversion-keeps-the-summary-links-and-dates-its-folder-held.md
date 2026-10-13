---
status: accepted
created_at: 2026-10-13T00:00:00Z
updated_at: 2026-10-13T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Legacy Archive Folder conversion keeps the summary, links and dates its folder held

On 2026-10-12 an adopter ran `roundfix history sanitize --apply` in three
batches, and the pre-PR review flagged records the generator had written: a
multi-line blockquote joined with `>` left mid-text, a record with a title and
no summary because the legacy `_prd.md` went straight to its first section,
empty `qa_override_approval` and `qa_override_reason` read as an unauthorized
override, summary links to `references/<file>.md` that no longer resolve
although the same command promoted the file, `sources: []` beside a promoted
adopted source, and `created: ""` although the folder's `_techspec.md` holds a
date. The maintainer decided on 2026-10-13 that only generation from Legacy
Archive Folders changes and that records already written stay byte-identical.

Lenient Legacy Reading therefore reads five more things from a Legacy Archive
Folder, wherever it becomes a record (the History Sanitize Command, the history
section of `roundfix baseline update` and the folder branch of
`ReadArchivedSpec`):

- **Summary.** Each line of a block quote loses its block quote markers (up to
  three spaces, `>` and one optional space, repeated for nesting) before the
  paragraph is joined. When the title has no paragraph of its own, the summary
  is the first paragraph after the first second-level heading.
- **Override fields.** A legacy `qa_override: true` without an approval or a
  reason records `legacy: not recorded` in the missing field. The other
  override fields keep the legacy value, so `ArchivedTaskCompleted` reads
  them as before.
- **Sources.** When `references/_index.md` lists no path, `sources` names each
  Markdown file directly under `references/` whose front matter has the shape
  of a Backlog Entry (`type`, `status`, `created`) or of a Finding (`status`,
  `created_at`).
- **Created.** An empty PRD `created` falls back to the TechSpec's `created`
  date, then, in the two commands that read Git, to the author date of the
  oldest commit that added the folder's `_prd.md` under the Spec Root or the
  archive root.
- **Legacy Link Rebase.** In the record's outcome and in every promoted copy,
  each relative link destination is resolved against the file it was written
  in. A destination inside the folder points at its promoted copy, or at the
  Archive Record when the file was not promoted, which names `source` and
  `source_revision` where its bytes remain. Any other destination is rewritten
  to name the same repository path from the new location. Path text naming the
  folder, a core artifact or a promoted file changes as the Archive Relink's
  path pass changes it. No other byte changes, so a file without such links
  is copied byte-identically.

The Archive Command, active Specs and every record already written are
unchanged.

## Considered Options

- Rewriting the records an earlier sanitize wrote was rejected by the
  maintainer: their folders are gone, and a record is never hand-edited.
- Changing the Archive Command too would change what the Delivery Queue's
  exact-archive proof accepts ("adds a record that agrees with the removed
  `_prd.md`"); it is recorded as a separate Backlog Entry.
- Teaching the review convention that empty override fields are legal for
  legacy records would hide the gap from every other reader; an explicit
  value says what the record knows.
- Reporting broken links as History Citations instead of rewriting them would
  leave every promoted copy broken in the adopter's live `docs/references/`,
  which the adopter fixed by hand.

## Consequences

This narrows ADR-0260's "a promoted file keeps its bytes and its own outward
relative links are not rebased" for Legacy Archive Folder conversions only: a
promoted copy written by the History Sanitize Command or by `baseline update
--promote` now keeps its bytes except for rebased relative destinations and
the Archive Relink's path text. ADR-0247's "Promoted copies must equal removed
blobs" binds the Delivery Queue's proof of an Archive Command, which never
sees a conversion. The Plan Digest of `baseline update` already binds each
record's bytes, so a changed record changes the digest of a plan that selects
it. The schema stays `roundfix/archive-record/v1`.
