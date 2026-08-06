---
status: accepted
created_at: 2026-08-06T00:00:00Z
updated_at: 2026-08-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The secondbrain inbox is the single capture door

Every fleeting observation a fleet session produces — for its own project or
for another — is born as a file in the secondbrain's `inbox/<destination>/`,
committed there at the moment of capture. Project repositories hold only
triaged artifacts: a Finding when the entry is evidence, a Backlog Entry when
it is intent, minted by a session of the destination project and committed by
that project alone. Origin is frontmatter, never commit authorship.

The decision was made on 2026-08-06, at the close of the v0.4.0 session, and
the trade-off was measured before it was accepted. Capture inside project
repositories collides with Active Runs by construction — review Batch commits
stage every path changed since their snapshot, and a dirty findings file
refused a `resolve` Preflight that same night. A gitignored project inbox
removes the collision but leaves the only copy of an observation on one
machine until triage. The brain inbox removes both: no file touches a project
checkout until the destination deliberately triages, and durability is a git
commit away from the moment of capture, with the existing push giving remote
backup and the existing `qmd` index making the pending queue searchable.

Two boundaries keep this from repeating a failure this fleet already paid
for. The repositories are siblings, never nested: nothing in the brain is
symlinked or embedded into a project tree, which is the structure that
produced the conexus git shims and the symlink pathspec failures the Daemon's
commit boundary now refuses. And only the capture layer lives in the brain:
findings and backlog stay in the destination repository, where the spec
pipeline, the consistency rules, the QA gates, and the pull requests that act
on them operate. The brain receives the triaged result anyway, curated and
versioned, through the mirror sync that already runs on every merge — the
loop closes with zero new pipeline.

Consequences accepted with the decision: the brain becomes load-bearing for
all capture, so the contract for repositories without a brain — the
`docs/_inbox/` fallback — is defined once, in the Baseline module, rather
than improvised per repository; the brain's own contract gains a narrow
carve-out, consumers may create files under `inbox/**` and touch nothing
else; and the door must stay cheaper than the bypass it replaces, which is
why its mechanics live in a skill and a one-line helper rather than in
discipline. Spec 0079 carries the implementation and the pilot that must
prove the save–retrieve–triage cycle before the rules bind the fleet.
