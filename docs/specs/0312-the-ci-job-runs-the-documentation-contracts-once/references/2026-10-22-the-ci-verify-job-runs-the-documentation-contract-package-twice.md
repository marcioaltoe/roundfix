---
type: perf
status: promoted
created: 2026-10-22
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
reason: null
---

# The ci-verify job runs the documentation contract package twice

## Slow

Roundfix's own `ci-verify` job runs `internal/docscontract` twice for every
pull request: once through the Make target `verify-changed-contracts`, which
selects it because its tests are `//verify:always` (its Contract Relevance), in the step
"Verify changed", and once through `docs-test` in the step "Verify docs". The
second run repeats every corpus Check of the first on the same tree, and the
maintainer and the Delivery Queue wait for both.

## Measured

On 2026-10-21, with five active Specs, the package took about 764 s in
"Verify changed" and about 541 s in "Verify docs" of the CI runs of pull
request #619. Spec 0305
(`0305-the-documentation-contracts-cost-seconds-per-active-spec`) makes each
run cheaper inside Go code; it leaves the second run in place because a
Makefile or workflow change was not authorized in cycle AP.

## Target

The package runs once per `ci-verify` job, for example by excluding it from
one of the two steps or by giving `verify-docs` a mode that skips contracts the
selective step already ran, while `make verify-docs` alone still runs it
before a pull request. The change needs the maintainer's specific
authorization for the Makefile or `.github/workflows`.
