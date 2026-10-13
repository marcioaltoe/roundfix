---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# An Archive Record's outcome keeps links relative to the removed Spec folder

## Symptom

`roundfix archive` copies the outcome paragraph of the PRD or QA Report into
the Archive Record byte for byte. A relative link in that paragraph, such as
`[the Backlog Entry](references/<file>.md)`, was written relative to the Spec
folder. The record lives one level up under the archive root and the folder
is removed in the same change, so the link resolves to nothing. On main at
`9963ffd3`, `docs/history/specs/0262-a-run-worktree-under-a-symlinked-root-is-cleaned.md`
line 25 links `references/2026-10-12-a-run-worktree-under-a-symlinked-root-is-never-cleaned.md`,
which does not exist.

## Where

`BuildArchiveRecord` as the Archive Command calls it (not the Legacy Archive
Folder path, which Spec 0271 fixes), and the Delivery Queue's exact-archive
proof, which accepts a record that agrees with the removed `_prd.md`.

## Expected

The Archive Command applies the same Legacy Link Rebase rule to the outcome it
writes: a destination inside the Spec folder points at its promoted copy or at
the record, and any other relative destination keeps naming the same
repository path. The exact-archive proof accepts that rewrite. Records already
written keep their bytes.

## Evidence

Found while authoring Spec 0271 on 2026-10-13: 29 records under
`docs/history/specs/` carry a relative link to `references/` or `../` in their
outcome, 28 of them written by the History Sanitize Command and one, 0262, by
the Archive Command.
