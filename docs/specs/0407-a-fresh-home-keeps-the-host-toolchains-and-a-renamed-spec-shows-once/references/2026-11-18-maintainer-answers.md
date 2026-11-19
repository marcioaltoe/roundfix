# The maintainer's answers that concern Spec 0407

Recorded by the operator from the maintainer's answers of 2026-11-18 and
copied here by the author. Only the answers that concern this Spec are kept,
verbatim.

## In the adopted Backlog Entry

The entry's own `## Decisions`, written when the maintainer reported the two
defects:

> 2026-11-18: "Spec após a release"; authorization "Mesmas do BE (Recommended)" (to be recorded in the Spec's `_authorization.md`).

The release it waited on was v0.148.0, published on 2026-11-18.

## `/roundfix --follow-up`, 2026-11-18 about 19:55 BRT (AskUserQuestion)

> 5. Backlog 2026-11-16 Delivery Queue com branch estrita (Quaestor): "Spec fix após a corretiva (Recommended)" — "Mesmas do BE".
> 6. Ordem: "Corretiva → regras → Indexus (Recommended)"; o fix da fila entra logo após a corretiva.

"A corretiva" is this Spec: it delivers first, then Spec 0408 (the Delivery
Queue fix), then the rules Spec and the Indexus Spec.

## Resumed planning, 2026-11-18 about 22:35 BRT (AskUserQuestion)

The author measured that the two defects share no component and the
operator asked the maintainer whether to split the Spec. The label, verbatim:

> Split da 0407 (os dois defeitos não compartilham componente, medido pelo autor): "Uma Spec só (Recommended)" — mantém as duas correções numa Spec, um release patch.

Both corrections stay in this Spec and ship in one patch release.

## What "Mesmas do BE" grants in cycle BG

The operator's brief for cycle BG states the grant for all four Specs:
implement, commit, push, pull request, merge and release, with a `patch`
Release Ceiling for a fix Spec; up to five corrections per Spec; the
Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter
configuration only if a correction requires one; owned skills, the guides
under `docs/agents/` and Baseline source with every derived file where the
design needs them. The operator sends every Spec Event Announcement and every
message.
