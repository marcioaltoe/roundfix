---
status: accepted
created_at: 2026-10-20T16:30:00Z
updated_at: 2026-10-20T16:30:00Z
deprecated_at: null
superseded_by: null
---

# GC sanitation classifies the Artifact Root each Run wrote to

On 2026-10-20 `roundfix gc sanitize` on the maintainer's Roundfix Home listed
one Artifact Root, the empty string, classified `unsafe`, because every
Implement Run records no Artifact Directory while it writes under
`~/.roundfix/artifacts/<repository id>`. About 605 MiB sat under those default
roots, and sanitation never reached them. Run Retention already resolved an
empty recorded root with its own code, and `roundfix storage report` gained a
second copy the same day. The maintainer chose "Resolver + eu aplico
(Recomendado)": sanitation resolves the empty root the way Run Retention does
and keeps its existing proofs, and the maintainer runs the dry run on the real
home after the release and applies it only if it comes out clean.

## Decision

**One resolver.** The Run Database package owns the only function that turns
a Run's recorded Artifact Directory into the Artifact Root the Run wrote to. A
blank recorded root resolves to the repository key's default root, or to the
recorded checkout's default root when only that one holds `runs/<run-id>` or
the key's default cannot be derived; a relative root resolves under the Run's
Git root; an absolute root stays as recorded. It reads the Run Database and
Git metadata, checks whether `runs/<run-id>` exists, and never creates,
writes or removes anything. The storage report, GC sanitation and Run
Retention all call it, so the three cannot disagree about where a Run's
artifacts live.

**Sanitation classifies the resolved root with the proofs it already had.**
Runs are grouped by their resolved root. Each group then passes the same
proofs as before, in the same order. The root must be a clean absolute path,
inside Roundfix Home and not equal to it, a physical directory with a physical
`runs` directory. It must be the key's or the checkout's default root for
every Run that resolves there. Only an `active` or `orphaned` root yields
candidates, and only `runs/<run-id>` directories: those of retention-eligible
terminal Runs recorded there and those of Run IDs absent from the Run
Database. A Run whose root cannot be resolved makes its group `unsafe`, so the
whole group is kept. Nothing outside `runs/` is touched.

**Absence is proven again before removal.** A directory proven only because
its Run ID is absent from the Run Database is checked against the Run Database
once more immediately before it is removed, and is kept when the Run now
exists. A Run's row is written before its `runs/<run-id>` directory, so a Run
that started after discovery is found by that check. The dry run names each
resolved root, each recorded value that resolved to it and what the root would
give back.

## Considered Options

- Recording the resolved Artifact Directory on new Implement Runs was rejected:
  it leaves every Run already recorded with an empty root unreachable.
- A third copy of the resolution inside sanitation was rejected: three copies
  had already started to drift on how they probe a directory.
- Trusting an empty root's default without the default-root proof was rejected:
  a Run that recorded an empty root may have written under a root configured at
  the time, and the proof is what keeps a root another repository shares.

## Consequences

This extends what ADR-0255 calls "proven the way `roundfix gc sanitize`
proves one": Run Retention now finds a relative recorded root under the Run's
Git root rather than under the process's working directory. The first
`roundfix gc sanitize --apply` after this decision can remove directories that
earlier versions always kept, so the operator reads its dry run first.
