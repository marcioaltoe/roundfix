---
granted: 2026-08-10
action: split-docs-validation-from-verify
paths:
  - Makefile
  - .github/workflows/ci-verify.yml
  - internal/docscontract
  - internal/cli/cli_test.go
  - internal/cli/baseline_documentation_contract_test.go
  - internal/cli/releaseplan_documentation_contract_test.go
  - internal/cli/baseline_release_gate_test.go
  - internal/speccheck/constraints_characterization_test.go
  - internal/speccheck/testdata/corpus-golden.json
  - internal/spec/archive_test.go
  - docs/agents/specific-repository.md
  - docs/references/coverage-record.json
  - internal/spec/coverage_test.go
consuming: direct
---

# Tooling authorization — markdown contracts move to the pull request boundary (2026-08-10)

On 2026-08-10, shown that the five most expensive test packages all read the
repository's documentation trees — so that virtually every Spec-loop commit
invalidates their cached results and `make verify` re-executes minutes of
unrelated tests — the maintainer directed the split:

> Vamos separar a validação dos markdowns da validação feita no make verify.
> Essa validação dos markdowns só devem ser feitas antes de abrir o pull
> request. Alem disso, nunca deve ser validado o que estiver em
> **/_archived/**/*

## What this covers

**One new invalidation domain.** The tests that validate repository markdown —
the public-documentation contracts from `internal/cli`, the published
Decision-Document example, the release-plan documentation contract, and the
Spec-corpus characterization from `internal/speccheck` — move into a dedicated
`internal/docscontract` package. `make verify` excludes that package, and a new
`make verify-docs` target runs it together with `roundfix spec check` and the
corpus time budget. No validation disappears; each moves to the boundary where
its subject changes matter: prose is validated when a pull request is about to
publish it, and CI runs `verify-docs` on every pull request so the boundary
stays fail-closed rather than advisory.

**Nothing under `_archived` is asserted.** The corpus characterization stops
materializing and sweeping `docs/specs/_archived`; the golden pins active
counts only. The Spec-0058 replay stops asserting anything about the archived
tree — the whole-tree before/after snapshot and the byte-comparison against
the archived report both leave. Reading one frozen archived file as fixture
input is not validating the archive and stays; asserting the archive is, and
goes.

**The pre-pull-request obligation is written where agents read.**
`docs/agents/specific-repository.md` gains the hard rule that `make
verify-docs` passes before any pull request opens.

Moving or renaming Go tests re-records `docs/references/coverage-record.json`
with the sanctioned update command in the same commit, per the standing rule.

## Addendum — 2026-08-10 — coverage enumeration sees the tagged domain

The moved tests carry the `docscontract` build tag so `go test ./...` excludes
them. `TestCoverageEquivalence` enumerates with `go test -list ./...`, which
would silently drop the moved tests from the record — the exact loss that test
exists to prevent. Its enumeration gains `-tags docscontract`, so the record
keeps tracking every moved test.

## Bounded by purpose

The purpose is confining cache invalidation, not weakening validation: every
moved assertion still runs, on every pull request, and fails closed. This
grant does not authorize deleting an assertion, weakening a contract string,
or exempting any active artifact from validation.

## Consuming Spec

Applied directly on branch `ma/docs-validation-at-the-pr-boundary` at the
maintainer's direction, as the next phase of the verification-speed campaign.

## Commit choreography

This record lands as its own commit, before any authorized change.
