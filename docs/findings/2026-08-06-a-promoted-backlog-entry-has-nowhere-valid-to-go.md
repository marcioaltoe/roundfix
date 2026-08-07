---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
---

# A promoted backlog entry has nowhere valid to go

**Date:** 2026-08-06
**Found by:** the maintainer asking whether Spec 0079 keeps the adoption
mechanism that moves a Spec's source documents into
`docs/specs/<slug>/references/`.

## Two shipped contracts disagree

The Backlog Operational Contract, adopted into `docs/agents/docs-layout.md`,
says what happens when a Spec claims an entry:

> When a Spec adopts an entry, set `status: promoted` and `spec` to that Spec's
> slug, then move the entry to `docs/specs/<slug>/references/`.

The adoption index that governs that same directory accepts two source types
and no more. `archive-spec`'s self-containment precondition validates every row
of `references/_index.md` and rejects anything else:

```awk
if (type != "inbox" && type != "finding") {
  reject("type must be `inbox` or `finding` at line " NR ": " type)
}
```

So a promoted backlog entry has three possible destinations and all three are
wrong: indexed as `backlog` it fails the archive; indexed as `finding` the
index lies about what the document is; unindexed inside `references/` it
becomes a file the self-containment check cannot account for.

## How it surfaced

Three entries were promoted on 2026-08-06 —
`2026-08-03-verification-performance-contract.md` and
`2026-08-06-two-stage-qa-gate-economics.md` to Spec 0080, and
`2026-08-06-event-journal-payload-economics.md` to Spec 0081. Their frontmatter
was set correctly and the move was not performed, which means the promoting
session violated the clause without any check noticing. Nothing in
`make verify`, `roundfix spec check`, or the QA gate observes it, because the
rule lives in prose and has no detector.

The entries are currently in the honest half-state: `status: promoted` with the
owning Spec named, still in `docs/backlog/`.

## Why the seam exists

The two contracts were authored for different problems and met for the first
time in practice. The adoption index predates the typed backlog: it was built
when the only things a Spec adopted were raw inbox notes and findings. Spec
0075 then added the backlog and reused the same destination sentence without
the index learning the new type.

The distinction the index encodes is still the right one and worth preserving:
adoption transfers *ownership* of source material into the Spec that commits to
implementing it, and the document travels with the Spec through archive. A
finding adopted this way leaves `docs/findings/` entirely — which is exactly
why the 2026-08-06 legacy sweep never touched Spec 0079's adopted finding, and
why rollup archival with `absorbed_by` is a different act from adoption, not a
competing one.

## What a fix has to settle

1. Whether `backlog` joins the index type set, which is a one-word change in
   the validator and a contract change to the index's meaning.
2. Whether the clause should require the move at all, or whether a promoted
   entry is better left in `docs/backlog/` with its `spec` pointer as the
   forward link — the pointer already answers "what is this Spec built from"
   without moving bytes.
3. Whichever wins, the rule needs a detector. A clause no check enforces is
   how three entries ended up promoted-but-unmoved on the day the clause was
   being cited.

The archive-spec skill is Roundfix-owned, so changing its validator needs an
express tooling authorization with bounded files.

---

## Addendum — 2026-08-07, resolved

The seam is closed. The adoption index learned the type the typed backlog
introduced: `backlog` joins `inbox` and `finding` in the `archive-spec`
validator and in the `write-prd` adoption steps, with the source-shape check
requiring a `docs/backlog/` path for it.

The three entries promoted on 2026-08-06 moved into their Specs' `references/`
with an `_index.md` each, so `docs/backlog/` again holds only live intent.

The third question this finding raised — that the rule needs a detector,
because its absence is why three entries sat promoted-but-unmoved on the day
the clause was being cited — is answered by `SC-BACKLOG-UNMOVED`. It reads only
declared values (the entry's own `status` and `spec`), reports the destination
in its fix line, and skips when no backlog directory exists. It fails an entry
that is promoted without naming a Spec, one that names a Spec that does not
resolve, and one that names a real Spec but never left `docs/backlog/`.

