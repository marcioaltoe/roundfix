---
type: perf
status: promoted
created: 2026-10-22
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
reason: null
---

# The documentation contract suite costs about a minute per active Spec and runs twice in one CI job

## Slow

`internal/docscontract` runs per-Spec contract work for every active Spec under `docs/specs/`. Roundfix's own `ci-verify` job runs the package twice, through the Make targets `verify-changed-contracts` (step "Verify changed") and `docs-test` (step "Verify docs"). The maintainer and the Delivery Queue wait on that job for every pull request, so its time grows with the number of active Specs rather than with the size of the change.

## Measured

On 2026-10-21, with five active Specs, the package took about 764 seconds in "Verify changed" and about 541 seconds in "Verify docs" of the CI runs of pull request #619 (Spec 0292). The job took 30m21s, past its former 30-minute limit, and 32m20s after PR #620 raised the limit to 60 minutes. The Delivery Queue's 30-minute check wait parked Spec 0292's item `checks-timeout` twice while CI was still running. With zero to two active Specs the same job stays under 30 minutes.

## Target

The job's documentation contract time stays within a few minutes of its zero-active-Spec time with five active Specs, or the per-Spec work runs once per CI job. A change to the Makefile or `.github/workflows` needs the maintainer's specific authorization. Kept as a backlog entry by the maintainer's cycle AO decision ("Backlog por enquanto (Recomendado)", 2026-10-22).
