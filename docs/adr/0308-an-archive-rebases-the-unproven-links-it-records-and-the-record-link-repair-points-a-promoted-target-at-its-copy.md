---
status: accepted
created_at: 2026-10-21T04:30:00Z
updated_at: 2026-10-21T04:30:00Z
deprecated_at: null
superseded_by: null
---

# An archive rebases the unproven links it records, and the Record Link Repair points a promoted target at its copy

An Archive Record lists, under `unproven`, the `satisfied-by` text of each
`## Unreachable Acceptance` declaration when the QA verdict is a qualifying
`partial` (ADR-0247). That text is written in the Spec's `_prd.md`, so a
relative link in it, such as `[references/x.md](references/x.md)`, resolves
against the Spec folder. Until this decision the archive copied the text
unchanged into `docs/history/specs/<slug>.md`, where the same link names
`docs/history/specs/references/x.md`, a file that never exists, while the
Legacy Link Rebase already rewrote the outcome paragraph (ADR-0277, ADR-0278).
On 2026-10-20 an adopter archived a Spec with 12 such declarations and
`roundfix spec check --strict` failed with 15 `SC-HISTORY-RETAINED` findings.
`roundfix history sanitize` planned Git Locators for the removed files but
pointed the file the archive had promoted to `docs/references/` at
`docs/history/specs/references/<file>` with the record's pull request, which
also never exists. The maintainer chose "Archive + alvo promovido
(Recomendado)" (option text: "O archive reescreve os `satisfied-by` (Git
Locator para arquivo removido, `docs/references/` para arquivo promovido), e o
sanitize passa a apontar arquivos promovidos para `docs/references/`. A
exigência da tag fica como está; ela é a âncora do Git Locator.").

**The archive rebases every `unproven` entry as it rebases the outcome.**
Wherever an Archive Record is built with a Legacy Link Rebase, by the Archive
Command and inside a Run or a Delivery Queue item, and by the conversion of a
Legacy Archive Folder, each `unproven` entry is rewritten as if it were
written in the folder's `_prd.md`. A destination inside the removed folder
points at its promoted copy under `docs/references/` when the archive promotes
it, and otherwise becomes the folder file's Git Locator, with the record's
pull request when no reachable commit holds it. Every other destination keeps
naming the same repository path. An entry without a relative link keeps its
bytes, and the entries keep their order and count. It is the same rewrite the
outcome receives, one code path and no copy.

**The Record Link Repair points a promoted target at its copy.** The
`record-links` unit of `roundfix history sanitize` and of the history section
of `roundfix baseline update` repairs an Archive Record's broken relative
links. For a broken link whose destination, resolved against the record's
`source` folder, names a file under that folder whose `docs/references/<file
name>` copy the record lists under `promoted` and the working tree holds as a
regular file, it now writes the relative link from the record to that copy,
and the plan names `docs/references/<file name>` as the rewrite. Every other
broken link keeps today's resolution: the record-relative target, then the
source target, then the record-relative path with the pull request.
`SC-HISTORY-RETAINED` prints the same promoted copy in its fix line. The
annotated `history-full` tag that `--apply` requires is unchanged (ADR-0248):
it anchors every Git Locator the repair and the conversion write.

Records already written are not rewritten by any command except the Record
Link Repair when an operator applies it.

## Considered Options

- Leaving the record text unchanged and asking `roundfix history sanitize`
  to repair every new record was rejected: a repository without the
  `history-full` tag cannot apply it, and the strict check would fail after
  every qualifying `partial` archive.
- Quoting each `satisfied-by` text as code, so no link survives, would hide
  the evidence the declaration names and still leave the promoted copy
  unlinked.
- A separate rewrite for `unproven` beside the outcome's would let the two
  drift; Spec 0288 had to make the outcome's rewrite resolve paths from the
  checkout root, and the entries inherit that fix only by sharing the code.
- Pointing the repair at the source path with the pull request for every
  promoted file would name a path that a later squash merge leaves only under
  `refs/pull/<N>/head`, although a surviving copy is in the tree.

## Consequences

An `unproven` entry with a link to a removed file now reads
`` `references/x.md` (`docs/specs/<slug>/references/x.md@<commit>`) ``, the
label followed by the locator, exactly as the outcome reads. The Delivery
Queue's exact-archive and Pull Request Record Commit proofs do not compare
`unproven`, so they accept the new records unchanged. The repair matches a
promoted copy by file name, because the archive names each copy
`docs/references/<file name>` and the record keeps no source path. A link to
an unpromoted file of the same Spec that shares a promoted file's name in
another directory is therefore pointed at that copy; the archive already
refuses two promotions that share a file name, so no link can match two
copies.
