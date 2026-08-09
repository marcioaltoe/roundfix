# Tooling authorization — restore the structural Baseline clauses (2026-08-07)

On 2026-08-07 the maintainer stated the goal that motivates this work:

> O que eu quero é que ao atualizar o setup do projeto, a estrutura do
> context-driven e autonomous work esteja totalmente respeitada. Que só fique
> aberto personalizações como banco de dados escolhido, HTTP Contract, etc.

and authorized the bounded tooling for it:

> Conceder para módulos + retention.

## What this covers

Spec 0082 proved the update mechanism works: refreshing a real
`standard-typescript-monorepo` copy accounts 89 of 103 prior managed clauses
mechanically, with no human input and no analyzer. The remaining 14 are
unaccounted because the current catalog no longer emits them:

```
clause.backend.boundary-contracts          clause.domain.canonical-language
clause.backend.http-independent-use-cases  clause.domain.layout-decision
clause.backend.layered-architecture        clause.frontend.organize-by-system
clause.backend.persistence-owner           clause.frontend.public-system-boundary
clause.backend.prohibit-generic-layers     clause.spec.local-task-tracker-only
clause.backend.thin-http-handlers          clause.spec.status-only-in-task
rule.backend.boundary-contracts            rule.monorepo.context-boundaries
```

All fourteen are structural Context-Driven guidance, not per-project
customization. `internal/baseline/assets/modules/backend.json` currently
declares no `clause.backend.*` at all. The fail-closed retention gate is
therefore behaving exactly as ADR-0058 requires; the loss is upstream, in the
catalog.

## Authorized paths

- `internal/baseline/assets/modules/*.json` — to restore the structural
  Normative Clauses named above to their owning modules.
- `internal/baseline/assets/retention/**` — to record any accounting a restored
  clause needs when its identity or owner legitimately moved.

Read-only reference: `internal/baseline/assets/source-baselines/**` is the
record of what the prior baseline emitted and is not modified by this grant.

## Bounded by purpose

This grant restores structural clauses and records their accounting. It does
not authorize changing any module's decisions, capabilities, required skills,
skill dispatch, or template selection for an unrelated reason, and it does not
authorize weakening any clause that is currently emitted.

Per-project customization stays a decision, never a clause: identifier
strategy, HTTP contract, authentication provider, and database choice remain
answered by the maintainer and are out of scope here.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic
consequences of the authorized asset edits, per ADR-0081. Hand-edited pin
values remain unauthorized mutations.

## Not authorized

- `internal/baseline/assets/profiles/**`.
- `internal/baseline/assets/decisions.json`.
- Mass reasoned-rejection of the fourteen clauses. The maintainer explicitly
  chose restoration over rejection; declaring them discontinued would need a
  new decision and a new record.

## Consuming Spec

This authorization is consumed by Spec `0084-an-update-that-can-run`.

## Commit choreography

This record lands as its own commit, before the commit of any Task it
authorizes.
