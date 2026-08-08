# Tooling authorization — capture clauses in the secondbrain module (2026-08-08)

On 2026-08-08 the maintainer supplied transcripts from four fleet sessions
(fiscus, fluxus, vortex, oraculum) that captured feedback correctly, and
directed that their behaviour become the default:

> Analise o comportamento provocado em cada projeto e faça com que isso seja o
> comportamento padrão para feedbacks de bugs/erros, ajustes/refatorações para
> melhorias e novas features.

and chose the scope:

> As duas acima + a regra de baseline-owned.

## What this covers

The module carries twelve clauses. Eleven contract consumption or permission;
exactly one contracts production, and it covers only external research
(`clause.secondbrain.research-capture`). The single general write clause is
permissive — "Sessions **MAY** create files under `inbox/**`" — so nothing
obliges a session that observes a defect, an improvement, or a feature idea to
capture it at all.

The gap is measured, not theorised. A vortex session named itself as the
evidence: it wrote feedback as a finding into a project checkout with an Active
Run — the precise collision ADR-0095 exists to prevent — then deleted it and
rerouted. A roundfix session on 2026-08-07 wrote five findings straight into
`docs/findings/` and captured one entry to the inbox, and only the captured one
was prompted by anything.

The behaviour to be made default, drawn from the four sessions:

- capture is triggered by observing a defect, an improvement, or a feature idea,
  not merely permitted;
- the entry is routed to the **owner's** namespace, which is frequently not the
  project the session is running in — `origin` and `destination` differ;
- existing entries are read first, and a strong verified match extends rather
  than duplicates;
- the entry is self-contained, because the triaging session lacks the author's
  context;
- an entry that belongs to another session is never committed by this one;
- guidance delivered inside `setup-context-driven` markers is baseline-owned, so
  changing it is an inbox entry to the baseline's owner, never a local edit the
  next update overwrites.

## Authorized paths

- `internal/baseline/assets/modules/secondbrain.json`, limited to adding three
  clauses and bumping the version of the rule that receives them.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md`
  and the same Source Baseline's `manifest.json`, limited to seating those
  clauses.

The Source Baseline paths are included from the start: the catalog refuses a
mandatory clause that no Source Baseline row carries, and the regenerator
maintains rows but never creates them.

## Bounded by purpose

The three clauses are the capture trigger with owner routing, the entry's
self-containment and ownership boundary, and the baseline-owned guidance rule.
This authorization does not permit changing any consumption clause, the query
order, the permission boundary, the secret prohibition, or any decision,
capability, or template selection.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edits, per ADR-0081. The maintained Source Baseline
entry count moves with the new corpus entries and its fixture expectation moves
with it.

## Commit choreography

This record lands as its own commit, before the commit that changes the module.
