---
status: accepted
created_at: 2026-10-01T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Baseline takes upstream setup names and drops a removed skill

On 2026-10-01 the upstream skills catalog renamed its setup lists `go-tui` to
`go`, `rust-cli` to `rust` and `typescript-bun` to `typescript`, merged
`go-cli` into `go`, and removed the `review` and `triage` skills. Every
built-in Baseline Profile named a setup that no longer existed, the `core`
module required `review`, and the `typescript` module required `triage`, so the
asset sync could no longer refresh any snapshot.

The Baseline now follows the reshaped catalog:

- A setup snapshot takes the upstream list's new name, as ADR-0191 requires.
  Built-in profile identifiers do not change: `go-cli-tui` takes `go`,
  `rust-cli` takes `rust`, and `standard-typescript-monorepo` takes
  `typescript`. The `go-cli` snapshot retires, because no built-in profile
  selected it and its upstream list is gone.
- A skill the catalog removes is dropped from every module, trigger and
  snapshot. It is never followed to a skill with a similar purpose under
  another name, so `review` is not replaced by `code-review`.
- When a removed skill carried behavior the Baseline still needs, that
  behavior becomes Normative Clauses of the module that owns the domain. The
  `external-triage` module therefore states, with force, how an item is
  classified, how it moves through triage states, when to ask and stop, how a
  forge comment discloses that an AI agent wrote it, how a `wontfix` decision
  is recorded, how an external pull request is triaged, and where accepted
  work goes. The `triage.external` decision stays a boolean. The repository
  records its state-to-label mapping in its own external-triage guide, outside
  the setup markers. A `wontfix` decision is recorded as a declined Backlog
  Entry, so no new directory is introduced. The old rule's Source Baseline
  entry is declared `replaced` by the first new clause.
- A skill that every built-in setup lists and every stack can use is required
  by `core`. `typesafe-ai` and `crafting-effective-readmes` are required there,
  and `crafting-effective-readmes` leaves the `typescript` module.

Two alternatives were rejected. Adding `crafting-effective-readmes` to the
`go` and `rust` modules would give one skill three triggers that say the same
thing, and a composed profile that selects two stacks would render two of
them. Keeping `triage` as a vendored copy would keep a skill that depends on
another removed skill, routes work outside the Spec workflow, and cannot be
invoked by the model.

## Consequences

A future built-in setup that omits `typesafe-ai` or
`crafting-effective-readmes` is refused by the check that a profile's setup
lists every skill its modules dispatch. An adopter keeps its installed
`review` or `triage` tree until it deletes the tree, because Roundfix never
deletes an installed skill tree (ADR-0191). `baseline skills reconcile` removes
the obsolete lock entry.
