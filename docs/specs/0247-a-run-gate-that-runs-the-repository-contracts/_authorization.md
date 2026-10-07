---
status: approved
granted: 2026-10-07
action: let make verify-changed run the Repository Contract Tests a change makes relevant, declared by Contract Relevance directives that cmd/verify-select reads, and describe the rule in the glossary and the repository rules
consuming: 0247-a-run-gate-that-runs-the-repository-contracts
paths:
  - Makefile
  - docs/agents/specific-repository.md
  - internal/baseline/derived_regeneration_repocontract_test.go
  - internal/docscontract/publicdocs_test.go
  - internal/speccheck/governed_repocontract_test.go
  - skills/owned_skill_edit_repocontract_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0247

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. The standing answer for the Governed Paths each Spec declares
is "Concedo".

On 2026-10-07 Spec 0245's Run and QA passed. GitHub CI then failed on
repository contract tests that the Run's Verification, `make verify-changed`,
never compiled. The maintainer decided:

- On the gate: "Contratos sempre no gate (Recommended)". `verify-changed`
  runs the repository contract tests of the affected packages, plus the cheap
  global contracts such as the suiteguard installation contract. The gate
  stays selective and fast and must not grow by more than about 20 s on a
  typical change.
- Grants: the Governed Paths this Spec declares, including the Makefile, and
  the skills and Baseline guides if needed. The CI workflows under
  `.github/workflows` were explicitly authorized too. This Spec does not need
  them, so it does not bound them.
- Platform-only failures stay out of scope. They are documented as a known
  limit, because CI is the Linux gate.
- A `qa_override` for an environment-only partial is standing.

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use temporary repositories, temporary homes and fake
runners.

The governed set was measured with `GovernedPath` on the authoring branch at
`59548f13`. The probe was a `go test -overlay` test that wrote nothing to the
repository, run against every file the Tasks declare.

## Why each governed path is unavoidable

- `Makefile` holds `verify-changed`. The new `verify-changed-contracts` target
  and its invocation are the change itself.
- `docs/agents/specific-repository.md` holds the hard rule that repository
  contracts validate at the pull request boundary. That rule must say that
  the selective gate now runs the relevant ones first.
- `internal/baseline/derived_regeneration_repocontract_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/speccheck/governed_repocontract_test.go` and
  `skills/owned_skill_edit_repocontract_test.go` are Repository Contract Tests
  that the Governed Path set already protects. Each gains its Contract
  Relevance directive in its file header. Their test bodies do not change.

## What is not governed

These paths are ordinary:

- `internal/verifyselect/contracts.go`,
  `internal/verifyselect/contracts_test.go`,
  `internal/verifyselect/contracts_repository_test.go` and
  `internal/verifyselect/verifyselect.go`.
- The other contract test headers: `internal/suiteguard/installation_repocontract_test.go`,
  `internal/baseline/repository_gate_repocontract_test.go`,
  `internal/config/regeneration_declared_test.go`, and the seven other
  `docscontract` files in `internal/docscontract`.
- `CONTEXT.md` and
  `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`.

## Limits

- No new dependency in `go.mod`. No change to `.roundfixrc.yml`, the CI
  workflows, the lint or formatter configuration, any skill or any Baseline
  module.
- No contract test changes its body or assertions. A file gains only a
  directive line and, where its header comment says the test runs only at the
  pull request boundary, the corrected sentence.
- `make verify-docs` and `repo-test` keep running every Repository Contract
  Test.
- No test, Verification command or QA row reaches a provider, starts a real
  Agent Session, reads a credential, or reads or writes the real
  `~/.roundfix`.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
