---
status: accepted
created_at: 2026-11-04T00:00:00Z
updated_at: 2026-11-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An installed owned skill keeps an unpinned lock entry, and Doctor names a stale one

ADR-0189 made the bundle install keep `skills-lock.json` true by rewriting the
`computedHash` of each existing entry that names an owned skill, and kept
`source`, `ref` and every other field. The `ref` names a commit of the skills
catalog that the upstream tool installed, not the copy Roundfix installed, so
the entry said "this hash, at that commit". The published skills CLI (1.7.6)
clones the catalog at the entry's `ref`, finds a different hash and
reinstalls the copy at that commit: on 2026-11-04 Quaestor's `skills update`
undid a `roundfix baseline update` and `roundfix doctor` failed again. The
maintainer chose "Lock coerente + doctor avisa (Recommended)".

## Decision

**The bundle install removes `ref`.** After installing the bundle,
`roundfix baseline update` and `roundfix skills install --target project`
write, in each existing lock entry that names an installed owned skill, the
installed tree's `computedHash` and no `ref` member. `source`, `sourceType`,
`skillPath`, every other member, their order and every other entry keep their
bytes. No entry or lock is created, as before.

An entry without `ref` is the skills CLI's own form for a default-branch
install: its `skills update` then compares the catalog's default branch with
the installed copy, keeps it when the hashes match, and replaces it otherwise.
The catalog syncs the owned skills from Roundfix's `main` once a day, so the
unpinned entry follows the copy the catalog publishes and pins no older one.

**Doctor names a Stale Owned Lock Entry.** An existing entry naming an
installed owned skill that carries a `ref`, or whose `computedHash` is not the
installed tree's, is stale: `skills update` would replace the installed copy
with another one. When the Setup Manifest is present, Doctor appends
`DR-SKILL-LOCK-STALE` with the sorted skill names to the `skills:` line, turns
an `ok` line into `warn`, keeps a `failed` or `unversioned` status, and adds
`roundfix skills install --target project`. A warning never changes Doctor's
exit code.

## Rejected alternatives

- Write a catalog commit as `ref`. The binary carries the owned skills, not a
  catalog commit: the Setup Snapshot marks owned skills `type: repo` with no
  `ref`, and the catalog copies them after Roundfix releases. Finding the
  commit needs the network during `baseline update` and still fails before the
  catalog's next sync. Roundfix never invents one.
- Change `sourceType` to `local` or drop `skillPath` so the CLI skips the
  entry. Both change the provider the format records and break
  `experimental_install`, which ADR-0189 refused.
- Fail Doctor on a stale entry. The installed copy is current and usable until
  someone runs `skills update`; a failure would block a repository that is
  ready.

## Consequences

The first bundle install after the upgrade removes `ref` from the owned
entries of every adopter whose lock has them; external entries keep their
bytes. Until then Doctor warns. Inside the catalog's daily sync window right
after a Roundfix release, an unpinned entry can still be replaced by the
catalog's previous copy; Doctor's version check then fails as before and one
bundle install restores it. This narrows the field rule of ADR-0189's
addendum of 2026-10-27 and keeps the rest of it.
