---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-26
---

# The Baseline does not name a home for authorization records

`docs/workflow/authorizations/` is machine-read contract data: the QA gate's
changed-path audit resolves tooling grants from the records PRDs cite there,
and `internal/speccheck/governed.go` names one record by path. Yet
`grep -rn "docs/workflow" docs/agents/*.md` returns nothing — no Baseline guide
provides for the directory, so an adopting repository learns the location only
by reading Go source or an existing PRD.

Measured 2026-08-26 while answering where the Baseline provides for
`docs/workflow`: it does not. The same sweep found the directory's other
former contents were guidance duplicates and were removed; the authorization
records are the one load-bearing resident.

Fix: one mandatory clause in the docs-layout module naming
`docs/workflow/authorizations/` as the home for tooling authorization records
(shape: dated file, `consuming:` Specs, bounded `paths:`), rendered into
`docs/agents/docs-layout.md`. Ride it with Spec 0116's decomposition, which
must write an authorization record there anyway for its in-session grant —
one module edit, one digest regeneration, one Task.
