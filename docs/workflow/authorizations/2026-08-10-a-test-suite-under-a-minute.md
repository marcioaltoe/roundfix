---
granted: 2026-08-10
action: speed-up-verification-suite
paths:
  - internal/baseline/derived_ownership_test.go
  - skills/baseline_skill_contract_integration_test.go
  - .github/workflows/ci-verify.yml
  - Makefile
  - internal/baseline/repository_test.go
  - internal/baseline/plan_test.go
  - internal/cli/baseline_plan_test.go
  - internal/cli/baseline_human_test.go
  - internal/cli/baseline_release_gate_test.go
consuming: direct
---

# Tooling authorization — a test suite under a minute (2026-08-10)

On 2026-08-10 the maintainer set hard performance targets for verification:

> Quero algo mais rápdido e eficiente para os testes do verify e CI. 4 minutos
> em um projeto com GO é inaceitável e extremamente lento. […] Temos que ter
> uma meta de no maximo 1 minuto de make verify e no maximo 2 minutos para o
> CI.

Shown the measured diagnosis and the three-phase proposal with these exact
bounded files, the maintainer granted the vehicle:

> Pode criar uma nova branch para fazer as mudanças das 3 fases

## What this covers

Measured on 2026-08-10: the fresh local suite takes 148s of wall clock at 339%
CPU on 12 cores, and `internal/baseline` alone accounts for 146s. Three tests
spawn nested `go test`/`make` invocations against an empty `GOCACHE` inside a
copied repository — 52.4s, 46.4s, and 24.2s in-suite — and their cold compiles
saturate every core, inflating neighbours up to 15× (`TestAssetsSyncCheck…`:
1.8s isolated, 28.1s in-suite). The same fork-storm is the measured cause of
the CI flakes: five different spawn-adjacent tests blew 90-second budgets
across five runs, one of them on a one-line JSON change.

The three phases:

1. **Nested invocations reuse the ambient build cache.** The regeneration
   tests' contract is byte-identical output from the declared command, not a
   cold compile; the explicit empty-`GOCACHE` overrides are removed so nested
   builds hit the same warm cache the suite already uses, while
   `GOFLAGS=-count=1` on those invocations keeps every nested test executing
   rather than replaying a cached result.
2. **CI restores what it already computed.** `actions/setup-go` gains
   `cache: true`, the job aligns `GOCACHE` with the restored path, and the
   fixed `-parallel 16` becomes hardware-derived. CI keeps `-count=1`: the
   verdict on the whole tree stays a full execution; only compilation and
   module downloads stop being repaid per run.
3. **Git fixture templates and parallel gaps**, only as far as the
   re-measured critical path demands: the repository-creation helpers named
   above may build one golden template per shape and copy it instead of
   spawning `git init`/`add`/`commit` per test.

## Bounded by purpose

The purpose is "melhoria na performance garantindo a confiabilidade": every
test keeps executing, no assertion is weakened, and a rule may move or become
cheaper but never disappear. This grant does not authorize deleting tests,
skipping tests, weakening any assertion, or changing what `make verify`
verifies.

## Consuming Spec

Applied directly on branch `ma/test-suite-under-a-minute` at the maintainer's
direction, superseding the vehicle question for the open test-performance
campaign.

## Commit choreography

This record lands as its own commit, before any authorized change.
