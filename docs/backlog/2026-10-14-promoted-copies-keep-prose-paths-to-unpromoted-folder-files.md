---
type: feat
status: open
created: 2026-10-14
spec: null
reason: null
---

# Promoted copies keep prose paths to unpromoted folder files

## Symptom

After a Legacy Archive Folder conversion, a promoted copy that names an
unpromoted, non-core file of the removed folder in prose (for example
`docs/history/specs/0004-links/notes.md`) keeps that path, while a Markdown
link to the same file is rewritten to the Archive Record. The path names a
file that no longer exists on the branch; it still exists at the record's
`source_revision`.

## Where

`internal/spec` Legacy Link Rebase and the Archive Relink's path text pass
(`ArchiveCitationMoves`), which maps only the folder, core artifacts and
promotions.

## Expected

Decide one rule for prose paths to removed, unpromoted files in both the
Archive Command and legacy conversion: keep them as historical paths, or
point them at the Archive Record. ADR-0277 currently keeps them ("No other
byte changes").

## Evidence

Spec 0271's QA (2026-10-14, finding QA-02 of the first QA report) reproduced
it with the eight-folder fixture through `history sanitize --apply` and
`baseline update --yes`. The behavior follows ADR-0277 and the Archive
Relink, so the Spec's QA was amended to record it as observed.
