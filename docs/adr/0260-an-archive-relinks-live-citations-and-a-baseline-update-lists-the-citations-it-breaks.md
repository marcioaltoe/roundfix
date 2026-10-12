---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An archive relinks live citations, and a Baseline update lists the citations it breaks

Two retirements cut files that live documents still cite. On 2026-10-09 an
adopter's `roundfix baseline update` sanitized 79 Pending History units, and
46 citations in 23 tracked Markdown files outside the History Root stopped
resolving; the plan named none of them. On 2026-10-10 the same adopter had
open Specs linking into the `references/` of an owner Spec, as ADR-0083 asks
("secondary Specs link the owner's copy"), and each link breaks when the owner
archives, because ADR-0247 removes the folder ("removes the Spec folder in the
same change"). The maintainer decided on 2026-10-11 "Archive promove e
reescreve" and on 2026-10-12 that the update lists every citation and promotes
only when the operator passes `--promote`.

A **Live Document** is a tracked Markdown file inside an active Spec folder,
an open Backlog Entry, an unresolved Finding, or a file under
`docs/references/`. A citation is a relative Markdown link destination or a
repository-relative path that resolves to a file the retirement removes, found
by the same detector the History Sanitize Command already uses.

**The archive relinks.** Before it writes anything, `roundfix archive` finds
every citation of the Spec folder from tracked Markdown outside that folder
and outside the History Root. A file outside the core artifacts that a Live
Document cites is promoted to `docs/references/<basename>`, exactly as
`--promote` copies it. A citation of the folder itself or of a core artifact
points at the Archive Record instead. Each Live Document citation of a
promoted file or of the folder is rewritten in place: only the link
destination or the path text changes, and every other byte stays. The record
lists the rewritten files under an optional `relinked` key, absent when empty,
so every earlier record keeps its bytes. A citation in any other file is
reported as a warning and left alone. The archive refuses, writing nothing,
when a promotion destination already exists, when two promotions share a
basename, or when a Live Document it must rewrite has changes not committed at
`HEAD`. The Delivery Queue's archive proof accepts a relinked file only when
its committed bytes equal the archive's rewrite of the parent's bytes, as
ADR-0230 already accepts outward link rewrites ("whose only content changes
besides the archive stamp are these rewrites"). The source keeps one owner:
secondary Specs still link the owner's copy, and the archive moves the copy
upstream when the owner leaves.

**The update lists.** The history section of `roundfix baseline update` lists
every citation that a selected unit's removal breaks, in text and JSON, and
the Plan Digest binds them, as ADR-0173 already does for History Relocations
("The warnings are part of the plan, so the Plan Digest binds them"). A citation from a Live Document of a promotable
file carries the offer `--promote <path>`. The update copies a file to
`docs/references/` only when the operator passes that flag, as the History
Sanitize Command does, and it never rewrites a citing file. When the History
Full Tag exists, the update also reports each citation that already fails to
resolve while the tag still holds its target, so an adopter finds what an
earlier update or batch broke; that report is advice outside the digest.

## Consequences

This narrows ADR-0254's "The history section changes only files under the
History Root" when `--promote` is passed: the promoted copy lands under
`docs/references/`. ADR-0247's "Promoted copies must equal removed blobs"
still holds, so a promoted file keeps its bytes and its own outward relative
links are not rebased. An archive now changes files outside the Spec folder,
which the Delivery Queue stages from the record.

Four alternatives were rejected. A Spec check refusing links into another
Spec's `references/`, or a copy of the source in each secondary Spec, would
give one source two homes, which ADR-0083 rejected. An archive that refuses
while a live document cites the folder would stall the normal queue order,
owner first. Automatic promotion by the update was declined by the maintainer
on 2026-10-12. Rewriting citing files from the update would change live
documents in a plan whose subject is the History Root.
