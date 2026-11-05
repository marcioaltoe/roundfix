---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A required external skill is checked only for presence, and only an owned skill is held to the binary's copy

Roundfix checks the content of the skills it writes itself, the Roundfix-owned
skills, against the copy inside the running binary. A required external skill
is checked only for presence: Roundfix reads neither its content nor its lock
hash, never compares it with its Setup Snapshot, and never replaces a present
copy unasked.

## Context

ADR-0221 held each required external skill to the `computedHash` of its
`skills-lock.json` entry and to the `treeDigest` of its Setup Snapshot.
ADR-0264 turned a snapshot difference into an informational
`DR-SKILL-DIFFERS-FROM-SNAPSHOT` that the managed refresh lists under
`skills.differing`, and kept restoring a copy that no longer matched its lock.
On 2026-11-05, with Roundfix 0.132.0, Fluxus reported both for
`systematic-debugging` and `ui-ux-pro-max`, whose copies come from
`marcioaltoe/skills`, the source the maintainer chose; a restore would have
replaced them with the snapshot copy. Vortex showed the same two skills the
same morning. A copy that differs from its lock is no better evidence: the
published skills CLI hashes the source folder, not the installed one, so its
`computedHash` cannot always be recomputed from the installed copy
(vercel-labs/skills issue 806). The maintainer asked that Roundfix check
version and digest only of the skills it authors and only presence of the
others, and answered the question of which copy the owned skills are compared
with:

> Cópia do binário, sem rede (Recommended) — O catálogo marcioaltoe/skills
> copia as skills do Roundfix do main todo dia, então a cópia dentro do
> binário é a mesma ou mais nova. O doctor aceita digest igual ou versão igual
> ou maior que a do binário, sem acessar a rede.

## Decision

**A required external skill is checked only for presence.** It is present
when `skills-lock.json` has an entry that names it and `.agents/skills/<name>`
is a directory. Roundfix reads neither the folder's files nor the entry's
`computedHash`. A present copy is ready whatever its content; a skill without
its folder or without its entry is a Missing Required Skill, and Doctor,
`roundfix skills check` and the `roundfix baseline update` preview report it
as before. An absent, unreadable or malformed lock, and a skill root that is a
symbolic link, keep their current refusals.

**Nothing compares an external skill with its Setup Snapshot.** Doctor no
longer appends `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or `snapshot comparison
unavailable`. The managed refresh no longer prints `Skills differing:` or the
JSON field `skills.differing`, and restores from the Setup Snapshot only a
Missing Required Skill, never a present copy.

**A deliberate restore keeps working.** `roundfix baseline skills restore
--profile <id> --skill <name>` still replaces a present external copy with its
Setup Snapshot copy, with its own preview and Plan Digest confirmation. No
command suggests it for a present copy.

**An owned skill keeps the version rule of ADR-0189.** An installed owned copy
passes when the version it declares is equal to or higher than the version the
running binary carries, read offline from the embedded `SKILL.md`. A copy whose
digest equals the embedded copy passes by the same rule, because the
owned-skill version record ties every shipped version to one digest. Roundfix reads
no catalog over the network.

This repository still holds its own Vendored Skills to their Setup Snapshot
and to their lock entries in its contract tests, because the Catalog Skill
Bump writes the three together (ADR-0224, ADR-0360). That rule is about this
repository's bump, not about an adopter's copies.

## Considered Options

- Drop only the snapshot comparison and keep the lock hash: rejected because
  the maintainer asked for no digest check of an external skill, the skills
  CLI's hash does not always describe the installed copy, and `baseline update
  --yes` would still restore a lock-mismatched copy from the snapshot.
- Count a folder alone as presence and stop reading the lock: rejected because
  the lock entry is the skills CLI's record that it installed the skill and
  the source a reinstall uses, and dropping it would turn every repository
  without a lock from a lock refusal into a list of every required skill.
- Compare an owned skill with the `marcioaltoe/skills` catalog over the
  network: rejected by the maintainer's answer above.
- Keep `skills.differing` as an always-empty field: rejected because a field
  that can never hold a value describes a comparison that no longer runs; its
  removal is declared in the release notes.

## Consequences

An adopter that installs its external skills from `marcioaltoe/skills` with
its own tool sees `skills: ok` while every required skill is installed, and
`baseline update` never offers to replace one. A hand-edited external skill is
no longer reported; the divergence order of ADR-0264 still decides how a
repository corrects one. Doctor's `outdated:` detail and the reason "no longer
matches its skills-lock.json entry" no longer occur for an external skill. The
owned-skill check, the Stale Owned Lock Entry warning of ADR-0359 and the
deletion of Retired Skills of ADR-0350, which reads the lock hash to prove a
copy unedited, are unchanged.
