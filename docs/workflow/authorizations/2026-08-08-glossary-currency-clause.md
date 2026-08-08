# Tooling authorization — a glossary-currency clause in the context-workflow module (2026-08-08)

On 2026-08-08 the maintainer specified a rule for keeping the project glossary
current:

> Mais uma regra normativa: O uso das skills grilling e domain-modeling para
> atualização dos arquivos CONTEXT.md para manter o glossario atualizado. Não é
> obrigatório a interação humana mas é importante que seja observada a
> necessidade de atualização a cada SPEC ou feat ou refatoração ou fix para
> evitar incongruencias.

## What this covers

`rule.context.domain-docs` carries exactly one clause, and it governs only
reading: `clause.context.read-domain-contract` requires reading the domain
context before naming concepts. Nothing obliges a session to check, when the
work closes, whether that work introduced, changed, or retired a term the
glossary should carry.

The `domain-modeling` skill is already dispatched, but its trigger is "Defining
or changing domain vocabulary, ownership, or bounded-context relationships" —
which fires only once a session already knows it is changing vocabulary. The
incongruence the maintainer describes comes from the opposite case: a term
enters the code without anyone noticing it is a term.

`grilling` is dispatched for stress-testing ideas, not for vocabulary. Its role
here is sharpening a term that is ambiguous enough to need it, which is
discretionary rather than a step.

## Authorized paths

- `internal/baseline/assets/modules/context-workflow.json`, limited to adding
  one clause to `rule.context.domain-docs` and bumping that rule's version.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/domain.md`
  and the same Source Baseline's `manifest.json`, limited to seating that
  clause.

This is a distinct purpose from the same day's grant on this module, which was
bounded to removing a routing sentence from an inbox clause.

## Bounded by purpose

The clause makes the *check* mandatory at the close of a Spec, feature,
refactor, or fix, and the *update* conditional on that check finding something.
It must not make human interaction mandatory: the maintainer said explicitly
that it is not required, so the clause obliges observation and update, not
approval. It must not restate the existing read obligation, which already has an
owner in the same rule.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edits, per ADR-0081.

## Commit choreography

This record lands as its own commit, before the commit that changes the module.
