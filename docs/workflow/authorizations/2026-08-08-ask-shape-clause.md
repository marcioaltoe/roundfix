# Tooling authorization — sharpen the ask clause in the core module (2026-08-08)

On 2026-08-08 the maintainer specified how an agent must ask, and where it
belongs:

> Mais uma regra que deve estar no docs/agents/agent-instructions.md e baseline:
> O uso do AskUserQuestion tool para perguntas ao usuário. Se estiver usando o
> codex que não possui essa ferramenta, deve fazer uma pergunta por vez,
> apresentando as opções e dando opções de escolha com numeros ou letras
> (simulando o askuserquestion tool).

## What this covers

`clause.core.ask-user-answerable-decisions` already exists and already says to
use the available user-interaction tool. What it does not say is the shape the
question takes when no such tool exists: "ask plainly and stop" leaves a session
free to ask three questions in one message, in prose, with no enumerated
options — which is what makes a question expensive to answer and easy to answer
partially.

That gap matters because the fleet's runtimes differ. A Claude Code supervisor
has `AskUserQuestion`; a Codex supervisor does not, and it is the Codex sessions
that most need the shape spelled out.

The clause is sharpened rather than duplicated: adding a second clause about
asking would leave two clauses governing one act, and the existing one already
carries the `stop-and-ask` enforcement and the coverage mapping.

## Authorized paths

- `internal/baseline/assets/modules/core.json`, limited to rewriting the
  guidance of `clause.core.ask-user-answerable-decisions` and bumping the
  version of the rule that carries it.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`
  and the same Source Baseline's `manifest.json`, limited to keeping that
  clause's corpus entry in agreement.

This is a distinct purpose from the 2026-08-07 grant on the same module, which
was bounded to adding a review-request clause.

## Bounded by purpose

The rewrite states: ask through the user-interaction tool when one exists; when
none does, ask one question at a time with enumerated options the user can
answer by number or letter. It does not change the clause's enforcement, its
identity, or any other clause.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edits, per ADR-0081.

## Commit choreography

This record lands as its own commit, before the commit that changes the module.
