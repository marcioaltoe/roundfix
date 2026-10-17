---
status: accepted
created_at: 2026-10-16T00:00:00Z
updated_at: 2026-10-17T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A retirement deletes, and the History Root holds only Archive Records

On 2026-10-14 `docs/history` held 383 files: 253 Archive Records, 98 Reduced
History Entries of Findings and 32 of Backlog Entries. A Jev pass that day
judged none of them reusable knowledge, and a second pass against the full ADR
and record text found 7 whose lesson no ADR states; those were captured in the
Secondbrain on 2026-10-15. No code reads the body of a Reduced History Entry,
while the Delivery Queue, Merge Evidence and the suite guard read Archive
Records. The maintainer decided on 2026-10-14 that retired material stops
being a reference ("remover definitivamente quando passarem a ser inuteis")
and on 2026-10-15 that the history keeps only Archive Records.

**A retirement deletes.** A Finding, Rollup or Backlog Entry that reaches any
terminal status, `deferred` included, a handoff and a Review Artifact leave the
tree by deletion in the change that records the status; a rejected or
deprecated ADR leaves the same way, and a superseded ADR once an accepted
successor names it with its date. Git keeps every text. No Reduced History
Entry is written, and nothing is moved into `docs/history/findings/`,
`docs/history/backlog/` or `docs/history/adr/`. Once a repository's Baseline
carries this rule, a file found there, a terminal Finding left in
`docs/findings/` and a relative link from an Archive Record or a Live Document
to a missing file fail the Spec Consistency Check, whose diagnostic prints the
Git Locator to use; before that update the binary does not refuse what the
repository's guide still asks for. A Spec that absorbs a source adopts it, so
its Archive Record names it; any other retirement names its reason in the
commit that deletes the file.

**A link to a deleted file becomes a Git Locator.** A Live Document or an
Archive Record that cites a file the change deletes is rewritten in the same
change: the link becomes its text followed by the Git Locator
`` `<path>@<commit>` `` written as code, so no relative link points at a
missing file. The commit is the provenance commit the file already names when
the default branch or a tag reaches it, and otherwise the newest commit the
default branch reaches that holds the file. A file no reachable commit holds,
such as one that lived only on a squash-merged delivery branch, is named by
its path and the record's pull request instead. This also applies when the
Archive Command removes a Spec folder: a destination inside the folder that is
not promoted, in a link or in a prose path, becomes a Git Locator, never a link
to the record itself.

**Existing history goes the same way.** In this repository the reduced entries
and the ten `superseded` Archive Records of Specs 0120 to 0129, never delivered,
cited by nothing and holding no regeneration grant, are deleted by the History
Sanitize Command. The same change gives the 40 broken relative links
already in Archive Records their Git Locators. In an adopter, `roundfix
baseline update` plans the same removals, including `superseded` records that nothing cites and that hold no
regeneration grant, rewrites every citation the removal breaks, and runs only
after the History Full Tag exists. A removed path must be held by the tag or by
a commit the default branch reaches, so files written after the tag still
leave. Retired files the Baseline would have relocated into the old families
are removed in the same plan instead.

**Lessons leave before the text.** A Finding that closes with a lesson its ADR,
Archive Record or reference does not already state is captured as one pending
Secondbrain Inbox Entry, its source pinned as `<repository>@<commit>:<path>`.

Industry ADR practice keeps superseded ADRs in the log and OpenSpec keeps each
completed change folder. The decision keeps what that practice protects, the
lineage, through the successor's consolidation line, the Archive Record's
sources and the Git Locator, and drops the copies an Agent would otherwise read
as current. This decision narrows ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0254,
ADR-0258, ADR-0260, ADR-0261 and ADR-0277, each of which carries a dated note.

**Extended by [ADR-0282](0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md) (2026-10-17).** When a superseded ADR's accepted successor does not name it, the history plan appends the dated `Supersedes` line to the successor and removes the superseded file; an ADR with no accepted successor is refused by its own path while the rest of the unit applies.
