# Tooling authorization — five rules the Baseline should have carried (2026-08-09)

On 2026-08-09, reading `docs/agents/specific-repository.md` after the greenfield
adoption, the maintainer worked through it section by section and ruled on where
each rule belongs. On the branch prefix:

> prefixo ma/ de branch deve ser configuravel pelo usuário e deve ser regra
> normativa do baseline.

On the sections that duplicated canonical guidance:

> Agent selection é parte de docs/agents/autonomous-work.md
> [## Profile-led Task routing] Acho que é dispensavel.
> [## Frontend Tasks: Claude Opus 5 xhigh] Informação repetida.

Then, on two more rules still sitting in the repository file:

> Regras que devem fazer parte da baseline:
> - Commit and PR titles are unscoped Conventional Commits subjects here;
>   `cog.toml` sets `scopes = []`.
> - **HARD RULE — release planning**: release work starts with the read-only
>   `roundfix release plan` before changelog, version, tag, push, package,
>   asset, or GitHub Release mutation. […]

And on the tooling choreography rule:

> [## HARD RULE — protected-tooling commit choreography] Talvez tenha que fazer
> parte da baseline também.

Finally, on external trackers:

> deve ser removido do base line qualquer sugestão ou orientacao de uso do
> linear ou github project, não vamos ter suporte a tools externas como o
> linear ou github issues nesse momento.

## What this covers

**One new configurable decision.** `branch.prefix` is a string decision with the
default `ma/`, selected by every Profile and rendered into the
`guide.agent-instructions` template beside `verification.gate`. Setup suggests
the default and the maintainer answers it, so the prefix is configurable per
repository rather than hard-coded for one.

**Five clauses promoted from one repository to the whole fleet.** Two under
`rule.core.verification-integrity` — never let a pipe hide a gate's exit status,
and an assertion reads the constant it means. Two under `rule.core.git-delivery`
— Conventional Commits subjects for commits and pull request titles, and the
read-only release plan before any release mutation. One under
`rule.core.tooling-authority` — the commit choreography that keeps an
authorization record, a prerequisite fix, and a consequent fix in separate
commits around the authorized change. Each was written portably: the roundfix
incident dates, the `cog.toml` scope value, and the `roundfix release plan`
command name stay in the repository file, which now only binds the portable rule
to this repository's concrete command and runbook.

**No external-tracker guidance to remove.** Measured on 2026-08-09, the catalog
contains zero references to Linear or GitHub Projects, and
`clause.spec.local-task-tracker-only` already forbids external triage labels and
external issue status as Task state. The one external surface was the
`triage.external` decision, which this repository answered `true` while shipping
no forge-label behavior in production code — only test strings mention labels.
That answer is now `false`, which drops `root.external-triage` and
`guide.external-triage` from the manifest.

## Authorized paths

- `internal/baseline/assets/decisions.json`, limited to adding `branch.prefix`.
- `internal/baseline/assets/templates/index.json` and
  `internal/baseline/assets/templates/guides/agent-instructions.md`, limited to
  rendering the `branch.prefix` token.
- `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/rust-cli.json`, and
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, limited
  to selecting `branch.prefix`.
- `internal/baseline/assets/modules/core.json`, limited to the five clauses
  named above and to requiring `branch.prefix`.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`
  and its sibling `manifest.json`, limited to the entries and rows those five
  clauses require.
- `docs/agents/setup-context.json`, limited to answering `branch.prefix` and
  changing `triage.external`.

Every derived pin rewritten by `make baseline-digests` — catalog digests,
plan-characterization goldens, formatter fixtures, and the source-baseline
index — is sanctioned fallout of these source edits under ADR-0081.

## Bounded by purpose

This grant covers moving rules that already bind this repository into the
portable catalog, and making the branch prefix a decision instead of a constant.
It does not authorize changing what any rule requires, adding rules the
maintainer did not name, or altering any other module.

## Consuming Spec

Applied directly rather than through a Spec: the maintainer named each rule and
its destination in the same session, and the edits are additive clauses plus one
decision.

## Commit choreography

This record lands as its own commit, before the commit that changes the catalog.
