---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A retired skill leaves the Baseline and the adopter removes its copy

On 2026-10-06 the maintainer retired three skills: "remover: the-fool,
autoresearch, council". `the-fool` and `council` were required by the
`context-workflow` module, dispatched by its guide, and listed by every Setup
Snapshot. `council` was also a Roundfix-owned skill shipped in the binary.
`autoresearch` was installed only in this repository and no module named it.
The upstream skills catalog still lists `the-fool` and `council` in its setup
lists, so ADR-0206's rule (drop a skill when upstream removes it) does not
apply, and ADR-0191's sync would put both names back into every snapshot.

The Baseline now has Retired Skills:

- A Retired Skill is a skill the Baseline stopped requiring while the
  upstream catalog may still list it. Roundfix keeps the set in code, beside
  the repository-owned set the asset sync already holds. It starts as
  `council` and `the-fool`.
- No module requires or dispatches a Retired Skill, no activation bundle
  names one, and the asset sync drops one from every Setup Snapshot it
  writes, even when the upstream list names it and even when the recorded
  snapshot holds it as repository-owned. Every other skill still follows its
  upstream list by name (ADR-0191).
- A Retired Skill that was Roundfix-owned leaves the binary's bundle, its
  canonical copy, its mirror and the owned-version record. Owned skills that
  pointed at it stop doing so in the same change.
- Roundfix never deletes an installed skill tree (ADR-0191). `roundfix
  baseline update` reports each Retired Skill the repository still holds, by
  its `.agents/skills/<name>` or `.claude/skills/<name>` entry or its
  `skills-lock.json` entry, under `Skills retired`, with the paths to delete.
  The report is read-only and never changes the update's state, category or
  exit code, and `--no-skills` skips it with the rest of the skills stage.
  Doctor ignores unrequired installed skills, as before.
- The adopter removes the copy: the user guide and the release note give the
  commands. A lock entry is deleted by hand, because the skills CLI's
  project-scope `remove` has left the lock entry behind.
- `autoresearch` never reached an adopter through the Baseline, so it is not
  a Retired Skill. This repository drops its lock entry, recommended-list
  line and tree.

Two alternatives were rejected. Leaving the two names in the snapshots as
unrequired entries, as `firecrawl` is, would keep `council` recorded as a
repository-owned skill the binary no longer ships, and the asset sync keeps a
repository-owned entry by name. Deleting the adopter's installed copies, or
editing its lock, during `baseline update` would break the rule that Roundfix
never deletes an installed skill tree, and an adopter may still use the skill
on its own.

## Consequences

A future retirement adds a name to the set, and the catalog tests then
refuse any module, trigger, bundle or snapshot that names it. If the
upstream catalog later removes a Retired Skill, the entry stays harmless until
a later change drops it. An adopter that keeps a retired copy sees
`Skills retired` on every update until it deletes the copy.
