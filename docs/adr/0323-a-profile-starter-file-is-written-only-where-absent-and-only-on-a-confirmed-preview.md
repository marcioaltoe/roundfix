---
status: accepted
created_at: 2026-10-23T12:00:00Z
updated_at: 2026-10-23T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A profile Starter File is written only where absent and only on a confirmed preview

The maintainer asked, on 2026-10-22, for a full Python Baseline profile that
also generates the repository's Makefile, Project Config and CI workflow in the
shape Quaestor wrote by hand. The Baseline writes only guidance today: managed
guides, root blocks and the Setup Manifest, refreshed on every update. Its own
tooling clause forbids an Agent to create or edit repository tooling
configuration without express maintainer authorization and says that a profile
does not grant it. A Makefile or workflow refreshed like a guide would
overwrite the repository's own decisions, and one written as part of every
plan would let a profile grant what the clause forbids.

## Decision

**A separate, confirmation-gated command writes a profile's Starter Files.** A
module may declare Starter Files: a repository path and a template of kind
`starter-file`. `roundfix baseline starters --profile <id>` lists, for the
selected profile, each Starter File it would create where nothing exists at the
path and each one it keeps because something does, with a Plan Digest over the
profile, the paths, the actions and the content identities. It writes only with
`--confirm-plan` equal to that digest, creates each absent file exclusively,
writes nothing when the digest differs (a path that appeared since the preview
changes the digest), and removes what it created if a later write fails. The
maintainer who confirms the exact preview gives the express authorization the
tooling clause requires for exactly those paths; an Agent runs the
confirmation only with that authorization, as for any other tooling change.

**Once written, the repository owns the file.** A Starter File is never part
of the Change Plan, the Setup Manifest or a managed refresh, carries no
markers, and is never rewritten, compared or removed by Roundfix.

## Considered Options

- Writing the files in `roundfix baseline apply` was rejected: the Change Plan
  is approved as guidance, and folding tooling into it would turn a guidance
  refresh into a tooling change.
- Managing the files like guides, with markers, was rejected: a Makefile and a
  workflow are repository tooling the repository changes freely, and a refresh
  must never overwrite them.
- Documentation-only examples were rejected: the maintainer asked for
  generation, and a copied example drifts from the profile's guide.

## Consequences

- An improved Starter File reaches only repositories that have not written
  that path yet; the profile's guide, which is managed, carries the rules that
  must reach every adopter.
- Profile alignment's next action for an undeclared Verification command names
  the command when the profile has Starter Files, so a new repository finds
  the way to a working gate.
