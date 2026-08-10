---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The preflight reads the catalog fixtures as repository carriers

## Opportunity

Baseline preflight inventories every instruction carrier under the repository
and warns when one nests managed markers inside unmanaged bytes. It applies that
rule to the Baseline's own embedded assets. Measured on 2026-08-10, both
`roundfix baseline plan` and `roundfix baseline apply` emitted:

```
baseline.inventory.nested-carrier-conflict:
  internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md
baseline.inventory.nested-carrier-conflict:
  internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/AGENTS.md
```

Those two paths are a formatter golden and a Source Baseline corpus entry. They
carry markers because they are fixtures *of* marked output, not because this
repository's agents read them.

## Value

The warnings are unactionable by construction: the resolution a
`nested-carrier-conflict` asks for — review the unmanaged bytes and decide — has
no meaning for a fixture whose whole purpose is to hold those exact bytes. Every
plan and apply in this repository emits both, so the operator learns to read
past a warning class that is also the one real signal for a genuine nested
carrier. That is the failure mode: a check that always fires teaches its own
reader to ignore it.

The cost is specific to Roundfix, which is the only repository that embeds the
Baseline catalog it also adopts. It grows with every fixture added under
`internal/baseline/assets/`.

## Shape

The inventory could exclude the Baseline's own embedded assets from carrier
discovery, the way it already scopes what counts as a recognized repository-rule
carrier. A narrower alternative is to keep discovering them but classify a
fixture-owned carrier as evidence rather than as a conflict, so a genuine nested
carrier elsewhere still warns. Worth settling in the same work: whether the
exclusion should be path-based or declared by the assets' own ownership records,
since `_ownership.yml` already states which trees are fixtures that nothing
regenerates. This shape is non-binding.
