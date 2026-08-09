# Tooling authorization — an outlive-the-turn clause in the autonomous-work module (2026-08-08)

On 2026-08-08 the maintainer asked that the lesson from the previous night
become Baseline guidance, and chose how to split it:

> Então o ideal é ter em docs/agents/autonomous-work.md uma instrução para
> montar o comando de construção do prompt para o /loop antes de iniciar o
> trabalho […] Isso deve ser, também, uma parte do canônico do context-driven.
> Princípio no catálogo, receita aqui.

## What this covers

On 2026-08-07 a supervising session ended a turn with "boa noite — vou
trabalhar" and no tool call, nothing detached, and no scheduled wake-up. The
maintainer slept believing four Specs were being implemented. Nothing ran: the
queue, the pull request, and both pending Specs were byte-identical the next
morning.

The cause is structural, not incidental. A turn that ends is indistinguishable,
from the waiting side, from a turn still working. A session that intends to
continue unattended must arm something that outlives the turn — a detached Run,
a scheduled wake-up, a daemon — or say plainly that it is stopping.

Nothing in the catalog states that today. `rule.autonomous.loop` carries six
clauses about how a Spec is implemented and none about how a session keeps
itself alive between Specs.

## Authorized paths

- `internal/baseline/assets/modules/autonomous-work.json`, limited to adding one
  clause to the existing `rule.autonomous.loop` rule and bumping that rule's
  version.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/autonomous-work.md`
  and the same Source Baseline's `manifest.json`, limited to seating that
  clause.

The Source Baseline paths are included from the start because 2026-08-07 proved
they are required: the catalog refuses a mandatory clause that no Source
Baseline row carries, and the regenerator maintains rows but never creates them.
Learning that mid-edit once is enough.

## Bounded by purpose

The clause states the obligation in runtime-agnostic terms. It must not name
`/loop`, which is a Claude Code mechanism: the module's existing clauses are
deliberately neutral about which runtime supervises, and an adopting repository
whose Supervisor is Codex would receive an instruction for a command it does not
have.

The concrete recipe — how to build the dynamic `/loop`, arm a Monitor on the
Run's terminal state, and template the queue prompt — belongs to this
repository's `docs/agents/autonomous-work.md`, which is repository-authored and
needs no grant.

This authorization does not permit changing any other clause, rule, decision,
capability, or template selection.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edits, per ADR-0081. The maintained Source Baseline
entry count moves with the new corpus entry and its fixture expectation moves
with it.

## Commit choreography

This record lands as its own commit, before the commit that changes the module.
