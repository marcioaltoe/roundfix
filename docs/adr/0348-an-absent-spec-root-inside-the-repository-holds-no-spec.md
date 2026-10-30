---
status: accepted
created_at: 2026-10-30T15:00:00Z
updated_at: 2026-10-30T15:00:00Z
deprecated_at: null
superseded_by: null
---

# An absent Spec Root inside the repository holds no Spec

ADR-0247 made the archive remove the Spec folder in the archive commit, and
Git keeps no empty directory. A repository whose last Spec was archived
therefore has no Spec Root in a fresh clone, a Run or item worktree, or CI,
and every command that resolved the root refused with `specs.root resolved to
"<path>", which does not exist`. On 2026-10-30 Quaestor reported its
`roundfix follow-up` `specs` source `unavailable` in that healthy state; the
same refusal stopped `roundfix spec check`, `roundfix deliver plan`,
`roundfix deliver status`, `roundfix history sanitize` and the history plan
of `roundfix baseline update`.

## Decision

A Spec Root is absent when its resolved path has no filesystem entry and
the nearest existing ancestor of that path, after symbolic links are
evaluated, lies inside the repository. An absent Spec Root holds no Spec.

- A command that enumerates Specs reads an absent Spec Root exactly as an
  empty one: the follow-up list and status views, `spec check` and
  `deliver plan` without a slug, the `deliver` subcommands that do not name a
  Spec, `history sanitize` and the history plan of `baseline update`.
- A command given a named Spec, and every other reader, keeps refusing with
  the same text, because the Spec it names cannot be there.
- A Spec Root that exists but cannot be read, a path that is a file, a
  symbolic link whose target is missing, and a configured root outside the
  repository that does not exist are still refused: those are configuration
  faults, and the external root of ADR-0035 is another repository that must
  be present.
- `roundfix archive` and `roundfix cancel` keep removing exactly the Spec
  folder and write no placeholder file.

## Considered Options

- The archive leaves a `.gitkeep`. The docs layout defines no such file,
  and ADR-0247's exact-archive proof accepts a commit that "removes exactly
  the Spec folder" and adds the record; a placeholder would widen that proof
  and every reader of it, and a repository that archived before the change
  would still lack the root.
- Every caller of the resolver accepts an absent root. The named-Spec
  commands would trade a precise refusal for a later "Spec not found" in more
  than twenty call sites, with no reported friction to justify it.
- Any missing root reads as empty. A broken knowledge-workspace link or
  an external root that was never cloned would then report zero Specs and
  hide the fault.
