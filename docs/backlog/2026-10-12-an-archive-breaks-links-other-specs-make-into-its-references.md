---
type: fix
status: open
created: 2026-10-12
spec: null
---

# An archive breaks the links other Specs make into its `references/`

## Symptom

`write-prd` and `write-techspec` tell a secondary Spec to link the owner
Spec's post-adoption copy of a shared source and adopt nothing. Since v0.53,
`roundfix archive` replaces the owner's folder with an Archive Record and
removes `references/`. The archive rewrites only the links that leave the
archived Spec (`prepareArchiveLinks` in `internal/spec/archive.go`). It does
not rewrite the links that arrive from other active Specs or live documents,
and `spec check` on the secondary Spec does not flag the dependency.

In the vortex, on 2026-10-10, open Specs 0079 and 0080 linked
`../0078-…/references/…`. Specs 0096 and 0097 each had four links to
`../0095-…/references/…`. Every one of those links breaks as soon as the owner
is archived, which is the normal order in a queue.

## Where

`roundfix archive` (`internal/spec/archive.go`), `roundfix spec check`, and the `write-prd`, `write-techspec` and `write-tasks` skills.

## Expected

Maintainer decision, 2026-10-11 ("Archive promove e reescreve"): when the
owner is archived, a `references/` file that a live document still cites
(an active Spec, an open Backlog Entry or Finding, `docs/references/`) moves to
`docs/references/`, and the inbound links are rewritten to the new path. The
source keeps a single owner. Use the same citation detector as the
`baseline update` fix.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-10-10-spec-secundaria-linka-a-references-de-outra-spec-e-quebra-no-archive.md`
  (vortex, 2026-10-10).
