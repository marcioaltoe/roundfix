---
type: fix
status: promoted
created: 2026-10-13
spec: 0272-the-history-holds-only-archive-records
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

Revised on 2026-10-16 by Triage of the Secondbrain Inbox Entry
`inbox/roundfix/2026-10-15-history-sanitize-rewrites-techspec-links-to-the-record-itself.md`
(origin conexus). `roundfix baseline update` 0.82.0, converting conexus's
Legacy Archive Folders, rewrote a body link to the Spec's `_techspec.md` so it
points at the Archive Record itself: four of 33 records
(`0001-correcoes-sync-e-contrato-canonico`, `0002-acesso-do-dashboard-aos-dados`,
`0004-fixtures-de-sessao-para-qa`, `0015-sessao-autenticada-para-o-gate-de-qa`)
read, for example, ``O documento técnico é [`_techspec.md`](0015-sessao-autenticada-para-o-gate-de-qa.md)``.
The TechSpecs exist at conexus `85ec98030c416af8072b13b630c7b3e96f67446b`
(tag `history-full`); conexus fixed them by hand in its PRs #61 and #62 with a
"full text in Git at <revision>: <path>" pointer. So a destination inside the
removed folder that is not promoted must name the file in Git at the record's
`source_revision`, never the record, and no record may link to itself.
Detection over a converted tree:
`for f in docs/history/specs/*.md; do b=$(basename "$f"); grep -n -E "\]\(\.?/?$b\)" "$f"; done`.
