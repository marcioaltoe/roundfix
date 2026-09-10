---
type: fix # feat | fix | perf | refactor
status: promoted
created: 2026-08-26
spec: 0119-spec-contained-authorization
reason: null
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

---

Addendum 2026-08-26, same sweep, same module: two more lifecycle rules the
maintainer stated that no Baseline guide carries.

- **Findings lifecycle.** `docs/findings/` holds pending work only. A finding
  that is done or deferred leaves the active directory for
  `docs/history/findings/` — and the existing archive-license contract
  (`SC-ARCHIVE-LICENSE`) requires every archived finding to carry `absorbed_by`
  resolving to an active Rollup or a Spec slug, so deferral routes through
  absorption into a thematic Rollup. Measured live: moving 24 deferred findings
  raised 261 then 54 license errors until the 6 Rollups were restored as active
  anchors and the 18 unabsorbed findings were absorbed member-by-member.
- **Handoffs lifecycle.** Handoffs are archived to `docs/history/handoffs/`
  only on the user's confirmation, and when confirmed, all of them archive
  together — the active directory holds only the capture door.
- **Backlog lifecycle.** `docs/backlog/` holds `open` or `promoted` entries
  only; any other status archives to `docs/history/backlog/`. Unlike findings,
  archived backlog entries carry no absorption license — intent that was set
  aside needs no absorber, only the history root.

---

Delivered 2026-08-26: the three lifecycle rules above (findings, backlog,
handoffs) were encoded into the Baseline the same day — clause extensions in
`internal/baseline/assets/modules/context-workflow.json`, digests regenerated,
the rendered guide carries them, and the managed refresh converges. What
remains open in this entry is the original gap alone: no Baseline guide names
a home for authorization records.

## Addendum — 2026-09-08 — Current triage

The current guide still omits an authorization-record home while the audit
reads legacy `docs/workflow/authorizations/` records. The user now explicitly
requests authorizations inside their Specs and asks for that contract in the
canonical rules. That instruction changes the proposed destination in the
original Shape; it does not erase the original observation.

Keep this entry open for provisional P1, authorization and knowledge lifecycle.
The new design must resolve Spec-contained grants while preserving valid legacy
references and the boundary between an authorization record and its governed
change. No Spec number or final grant schema is assigned by this triage.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0119-spec-contained-authorization](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
