# Tooling authorization — the Secondbrain is a consultation source (2026-08-08)

On 2026-08-08 the maintainer specified what the Secondbrain is for on the
reading side:

> Quero também reforçar nas rules o uso do secondbrain como fonte para validação
> de estratégia, obtenção de dados de outros projetos, conhecimento geral,
> conhecimento técnico, livros, papers e tudo mais que é necessário para criar
> specs mais robustas e tomar melhores decisões considerando o ecossistema de
> projetos.

## What this covers

The module now contracts capture well and consumption narrowly. Its consumption
clauses say **how** to read — index first, then `qmd query`, then the pointed-to
files — and **when not to** — the guidance is bounded to "when repository code
and documentation do not fully answer the task". Nothing says the brain is an
input to a design decision.

That bound is the defect. It frames the brain as a fallback for questions the
repository cannot answer, when its highest-value content is precisely the
content no repository holds: books, papers, and the decisions other projects in
the same ecosystem already made and paid for. A Spec authored without it can be
internally consistent and still repeat a mistake another project recorded in an
ADR two weeks earlier.

The gap is measured on this Spec. Consulting the brain while diagnosing the
managed-refresh defect returned `verificacao-adversarial-e-oraculos-de-agentes`,
which states the failure mode the defect is an instance of — a gate trusted
because it exits zero rather than because it discriminates — and
`convergencia-observavel-em-sistemas-operacionais`, whose convergence rule
supplied the acceptance criterion the Spec now uses. Neither was reachable from
this repository. The consultation happened because the maintainer asked for it,
not because any rule required it.

## Authorized paths

- `internal/baseline/assets/modules/secondbrain.json`, limited to adding one
  consultation clause and bumping the version of the rule that receives it.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md`
  and the same Source Baseline's `manifest.json`, limited to seating that
  clause.

The Source Baseline paths are included from the start: the catalog refuses a
mandatory clause that no Source Baseline row carries, and the regenerator
maintains rows but never creates them.

## Bounded by purpose

The clause obliges consultation at the points where a decision is being formed —
authoring a Spec, choosing an approach, validating a strategy — and names what
the brain is expected to supply there: prior decisions from sibling projects,
literature, and general technical knowledge. It must not restate the query
order, which already has an owner, and it must not weaken the read-only
boundary or the citation obligation. It must not make consultation a gate that
blocks work when the brain is unavailable; an unreachable brain is a reported
condition, never a reason to stop.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edits, per ADR-0081. The maintained Source Baseline
entry count moves with the new corpus entry and its fixture expectation moves
with it.

## Consuming Spec

This authorization is consumed by Spec `0084-an-update-that-can-run`.

## Commit choreography

This record lands as its own commit, before the commit that changes the module.
