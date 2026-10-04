---
type: fix
status: promoted
created: 2026-10-03
spec: 0226-an-archived-spec-keeps-its-links
reason: null
---

# Archiving a Spec breaks its relative links that leave the Spec

## Symptom

`roundfix archive` (0.26.0) moves `docs/specs/<slug>/` to `docs/history/specs/<slug>/`, one level deeper. Relative Markdown links that leave the Spec, such as `../../history/findings/<file>.md` in a PRD or a Task, then resolve to `docs/history/history/...` and break. Evidence scripts that climb a fixed number of directories break the same way. Measured on fluxus on 2026-10-02:

- Spec 0096, already on main: broken links that can no longer be fixed in place.
- Specs 0097 and 0098: three broken links, fixed in the archive commit.
- Spec 0096's `journey.py`: climbed six directories to find `.env`; the pre-PR review caught it.

## Expected

The archive rewrites relative links that leave the moved directory, as source adoption already does in `write-prd` step 8. It refuses the move when any such link would still be broken.

## Authority

The change lives in `internal/spec/archive.go`, a Governed Path. It needs a named maintainer grant before a Spec can implement it. Spec 0219 relates (archive and requeue traps) but did not adopt this for that reason.

## Source

Secondbrain `inbox/roundfix/2026-10-02-o-archive-nao-reescreve-links-relativos.md` (fluxus).
