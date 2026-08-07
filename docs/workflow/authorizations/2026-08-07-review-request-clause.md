# Tooling authorization — a review-request clause in the core module (2026-08-07)

On 2026-08-07 the maintainer directed the hand-opened-pull-request review rule
to become Baseline guidance rather than a rule of this repository alone:

> Essa rule DEVE fazer parte do context-driven para todos os usuários e não
> somente nesse repo.

and chose its shape:

> Condicional, sem decisão nova.

## What this covers

`.coderabbit.yaml` in this repository disables automatic review. Spec 0078
taught Roundfix to request the review after its own Final Push, which covers the
watch and resolve loops and nothing else. A pull request opened directly gets no
review, and its check reports `Review skipped: automatic reviews are disabled` —
a green result that reads like approval. Three pull requests merged unreviewed
here on 2026-08-07 before anyone noticed.

No module in the catalog carries any review-request guidance today; the concern
is absent, not merely weak.

## Authorized paths

- `internal/baseline/assets/modules/core.json`, limited to adding one clause to
  the existing `rule.core.git-delivery` rule and bumping that rule's version.

Extended on the same day, after the first attempt revealed that the clause
cannot be added without them. The catalog validator refuses a required clause
that no Source Baseline row carries, and its own message states that the
regenerator maintains manifest rows but never creates them:

- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`

Both limited to seating `clause.core.request-review-explicitly`. This is the
same sequence pull request #35 used when it last added a clause: corpus text
first, then the manifest row, then regeneration to compute offsets and digests.
No other clause, row, or corpus passage may change.

That rule already owns delivery authority — it carries
`clause.core.preserve-git-scope`, `clause.core.ask-before-destructive-git`, and
`clause.core.ask-before-delivery` — so the concern belongs to it rather than to
a new rule.

## Bounded by purpose

The clause is conditional by explicit maintainer choice: it applies to every
adopting repository but only bites where review automation does not run on its
own, and it forbids reading a skipped or absent review check as approval. It
introduces **no new Baseline decision**, so no adoption or update gains a
prompt — the maintainer rejected that option specifically because the day's work
was spent removing exactly that friction.

This grant does not authorize changing any other clause, decision, capability,
required skill, or template selection in `core.json` or in any other module.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic consequences
of the authorized asset edit, per ADR-0081. Hand-edited pin values remain
unauthorized mutations.

## Relationship to the repository's own rule

`docs/agents/specific-repository.md` carries a repository-specific version of
this rule, written before the maintainer widened the scope. It names this
repository's concrete mechanism — `.coderabbit.yaml`, `request_review`, the
Final Push — which portable Baseline guidance cannot assume. Both may stand: the
catalog clause states the obligation, the repository rule states how it is met
here.

## Commit choreography

This record lands as its own commit, before the commit that changes the module.
