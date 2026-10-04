---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An archived Spec keeps its relative links

`roundfix archive` moves a Spec one directory deeper, to
`docs/history/specs/<slug>/` for the built-in Spec Root or
`<spec-root>/_archived/<slug>/` for any other. A relative Markdown link that
leaves the Spec then resolves one level short. fluxus measured it on
2026-10-02: Spec 0096 reached `main` with links that can no longer be fixed in
place, and Specs 0097 and 0098 needed three links fixed by hand in their
archive commits.

Before it writes anything, the archive now resolves every relative link
destination in the Spec's Markdown files against the file that holds it. A
destination inside the Spec keeps its bytes. A destination that leaves the
Spec and reaches an existing path is rewritten to the relative path that
reaches the same path from the archived location. A destination whose target
does not exist from the active location, but which already resolves unchanged
from the archived location, keeps its bytes, because its target was archived
into the mirrored history layout first. When any other destination that leaves
the Spec resolves from neither location, the archive refuses before it stamps
or moves anything and names each such link with its file and line. The
Delivery Queue accepts an archive commit whose only content changes besides the
archive stamp are these rewrites.

Two alternatives were rejected. Leaving links to a later checker would let
history accumulate broken links that the archived-Spec rule forbids fixing in
place. Rewriting only the links that resolve, without refusing, would hide a
broken link inside a moved Spec.

## Consequences

The refusal includes links that were already broken before the archive, so a
Spec carrying one cannot archive until its author fixes or removes the link,
and the Delivery Queue parks such an item with the archive's message. Evidence
scripts that climb a fixed number of directories, HTML anchors, wiki links and
absolute destinations are not rewritten. A rewritten link that reaches an
active document breaks again if that document is archived later; the archive
of a target does not rewrite links in other archived Specs. Specs archived
before this decision keep their bytes.
