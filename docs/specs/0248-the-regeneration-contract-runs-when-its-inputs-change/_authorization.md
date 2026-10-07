---
status: approved
granted: 2026-10-07
action: select TestRegenerationIsDeclared by the inputs its regeneration reads, add the Full Contract Run that CI runs on every push to main and before every release, name the selective gate's exclusions, and describe the rule in the glossary and the repository rules
consuming: 0248-the-regeneration-contract-runs-when-its-inputs-change
paths:
  - .github/workflows/ci-verify.yml
  - .github/workflows/release.yml
  - Makefile
  - docs/agents/specific-repository.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0248

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. The standing answer for the Governed Paths each Spec declares
is "Concedo".

On 2026-10-07, after Spec 0247 merged, the maintainer decided how
`TestRegenerationIsDeclared` runs: "Rodar só quando relevante (Recommended)".
The option read: "Roda quando o diff toca módulos do Baseline, perfis,
fixtures ou o .roundfixrc.yml, e sempre no CI do main e das releases. Nos
demais casos fica de fora do verify-changed." The same day the maintainer
granted:

- the Governed Paths this Spec declares, including the `Makefile`;
- changes to the CI workflows under `.github/workflows`, which were
  explicitly authorized;
- the skills and the Baseline guides if needed. This Spec changes neither a
  skill nor a Baseline module, so it bounds none of them;
- a standing `qa_override` for an environment-only `partial`.

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use temporary repositories, temporary homes, stub
commands and fake runners.

The governed set was measured with `GovernedPath` on the authoring branch at
`ae56aba0`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository, run against every file the Tasks
declare.

## Why each governed path is unavoidable

- `Makefile` gains the `verify-contracts` target, the Full Contract Run, next
  to `verify-changed-contracts`.
- `.github/workflows/ci-verify.yml` runs the Full Contract Run on every push
  to main. Today that workflow runs no step that reaches
  `TestRegenerationIsDeclared`.
- `.github/workflows/release.yml` runs the Full Contract Run after its
  `make verify` gate and before any build or publication. Today the release
  workflow runs no Repository Contract Test at all.
- `docs/agents/specific-repository.md` holds the hard rule that repository
  contracts validate at the pull request boundary. It says that
  `make verify-docs` runs all Repository Contract Tests, which was never true
  of `TestRegenerationIsDeclared`, and it must name the Full Contract Run.

## What is not governed

These paths are ordinary:

- `internal/verifyselect/verifyselect.go`, `internal/verifyselect/contracts.go`,
  `internal/verifyselect/contracts_test.go`,
  `internal/verifyselect/contracts_repository_test.go` and the new
  `internal/verifyselect/full_contract_run_test.go`;
- `internal/config/regeneration_declared_test.go`, whose header directive
  changes;
- `CONTEXT.md` and
  `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`.

## Limits

- No new dependency in `go.mod`. No change to `.roundfixrc.yml`, the lint or
  formatter configuration, any skill or any Baseline module.
- `TestRegenerationIsDeclared` changes only its header directive. No contract
  test changes its body or its assertions.
- The `verify`, `verify-changed`, `verify-changed-contracts`, `verify-docs`,
  `docs-test`, `repo-test` and `spec-budget` recipes and `REPO_CONTRACT_TESTS`
  stay byte-identical. The pull request CI path is unchanged: it runs
  `make verify-changed` and `make verify-docs` as before.
- No test, Verification command or QA row reaches a provider, starts a real
  Agent Session, reads a credential, or reads or writes the real
  `~/.roundfix`. No row triggers a GitHub workflow.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
