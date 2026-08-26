---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# An excluded artifact outlives the decision that excluded it

## Opportunity

Answering a Baseline decision so that it excludes an artifact removes that
artifact from the Setup Manifest but leaves its bytes on disk. Measured on
2026-08-10: setting `triage.external` to `false` and running
`roundfix baseline update --yes` dropped the `external-triage` module,
`root.external-triage`, and `guide.external-triage` from
`docs/agents/setup-context.json`, and reported one file change. It did not
delete `docs/agents/external-triage.md`, and it did not remove the
`root.external-triage` region from `AGENTS.md`. Both survived with their
`setup-context-driven` markers intact.

Only the greenfield pass removed them, and only because greenfield rebuilds
every managed carrier from nothing.

## Value

Between the two operations the repository states a rule the Baseline no longer
recognizes, in a file whose markers claim setup owns it. An Agent reading
`docs/agents/` cannot tell that the guidance was retired: the marker is the
signal that setup manages those bytes, and here it outlived the management.
`AGENTS.md` kept pointing at the retired guide, so the root instructions routed
readers to it.

The same gap applies to any decision flip that excludes an artifact, not only
this one. `spec.scaffold`, `secondbrain.enabled`, and
`repository.extension.enabled` all carry `excludeArtifacts` effects.

## Shape

Exclusion could plan a deletion the way readoption already plans one: the plan
already computes `DeletePaths` for repository-rule carriers, so the same
mechanism could cover a managed artifact whose module or artifact left the
manifest, and could strip an orphaned managed region from a shared carrier such
as `AGENTS.md`. Preservation would then converge with greenfield instead of
diverging from it. Worth settling in the same work: whether a removed artifact
should be deleted outright or reported as an orphan for the maintainer to
confirm, since deletion of repository bytes on a decision flip is the kind of
change the Baseline otherwise asks about. This shape is non-binding.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
