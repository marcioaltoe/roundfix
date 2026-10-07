---
status: accepted
created_at: 2026-10-07T00:00:00Z
updated_at: 2026-10-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Existing history is sanitized in batches after a history-full tag

ADR-0247 made an archive leave an **Archive Record** and remove the Spec
folder, but it left alone everything already archived. On 2026-10-07
`docs/history` still held 223 Spec folders from before that decision (5,790
files, 58.2 MB), 50 Review Artifact folders (501 files, 1.6 MB), 11 handoffs,
31 terminal Backlog Entries and 98 terminal Findings. On 2026-10-06 the
maintainer decided that the clean-up "deve atuar também no que já existe no
diretório", and chose "Aplicar em lotes via PR (Recommended)": after a tag
marks the full history, the existing history is applied in batches of Specs,
one Pull Request per batch, each with `make verify` green and each revertible.

`roundfix history sanitize` does that migration, and only when asked:

- **Dry run by default.** Without `--apply` it lists each pending unit and
  writes nothing in the repository. A unit is one legacy Spec folder, or one
  of the other history kinds as a whole. For a folder it names the files,
  bytes, the record it would write and the candidate files to promote. With
  `--advise`, Jev classifies the candidates of the next batch through the
  Spec judge's key, ceiling and log. The advice never gates, as in ADR-0247.
- **Batches through Pull Requests.** `--apply --batch <n>` converts the next
  `n` units in a fixed order and stops. It never runs implicitly, and it
  refuses a dirty working tree. The operator commits the batch on its own
  branch and opens one Pull Request for it.
- **A tag before the first batch.** The operator creates the annotated tag
  `history-full` on the last commit before the first applied batch and pushes
  it. `--apply` refuses unless that tag exists, is annotated, is an ancestor
  of `HEAD` and holds every path the batch removes or rewrites. Every removed
  byte therefore stays reachable from a tag, whatever later merges do.
- **One record per legacy folder.** Each folder becomes `<slug>.md` through
  `BuildArchiveRecord`, the builder the archive itself uses. Its `source` is
  the folder's own path under the archive root and its `source_revision` is
  `HEAD` at apply, which holds the folder. The delivery commit and Pull
  Request are read from Git when a first-parent commit took the Spec out of
  the Spec Root or first added it to the archive root. A commit that only
  relocated an archived folder is not one. Thirteen folders predate the
  archive stamp, and two of them never had a QA Report. A folder with neither
  a QA Report, an override nor a supersession gets the new disposition
  `no-qa` instead of an invented `pass`.
- **The other kinds, by what still reads them.** A retired Finding or Backlog
  Entry becomes a **Reduced History Entry**: its front matter unchanged, its
  title, its first paragraph and the revision that holds its full text. The
  Spec check still validates every `absorbed_by` license and Rollup member
  against those names, and each `absorbed_by` names a Spec slug that its
  Archive Record keeps resolving. Retired Review Artifacts and handoffs have
  no reader, so they are removed. Retired ADRs stay whole, because they are
  decision rationale.

Readers of a legacy folder stay. `ReadArchivedSpec`'s folder branch, the
exact-move proof, ADR-0230's link pass and the pinned-history reads of
ADR-0247 still serve adopters whose history was never sanitized. Only this
repository's corpus holds no legacy folder once the last batch merges. The
Secondbrain export keeps its exclusions while any legacy folder remains, and
mirrors `docs/` whole once none does. A documentation contract holds both
directions.

Two alternatives were rejected. One commit converting all 223 folders would
make a 58 MB diff no person can review and a single revert the only way back.
Leaving the existing history alone would keep 99% of `docs/history` that the
maintainer called irrelevant, and the Secondbrain mirror would keep its
exclusions forever.

## Consequences

This narrows ADR-0215's "Nothing archived is removed" for the existing
history. The maintainer's approval is the decision quoted above, and each
batch's Pull Request is the "later change" ADR-0215 asks for. It also narrows
the Findings contract's "immutable history" for retired Findings only: a
reduced entry keeps every front-matter field and its first paragraph, and the
original text stays in Git at the named revision and at the `history-full`
tag. The Archive Record schema stays `roundfix/archive-record/v1` and gains
the disposition `no-qa`. Git history is not rewritten and the pack does not
shrink ("Não agora", 2026-10-06).
