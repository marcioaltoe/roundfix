---
status: approved
granted: 2026-10-08
action: make the tests of seven packages run in parallel under a repository rule, convert the sequential residue of internal/cli to per-test values, narrow two pinned-history fixtures, and run both Repository Contract Test tags in one go test
consuming: 0254-a-faster-make-test
paths:
  - docs/agents/specific-repository.md
  - internal/baseline/derived_ownership_test.go
  - internal/cli/cli_test.go
  - internal/spec/coverage_test.go
  - internal/speccheck/mechanical_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0254

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. The standing answer for the Governed Paths each Spec declares
is "Concedo".

On 2026-10-08 the maintainer placed this Spec last in the cycle: "A spec de
make test é importante, mas acho que pode ser feita, se complexa, após o final
das demais specs e releases". It is approved for delivery without further
consultation, within the bound of four implementation Tasks, with grants for
the governed paths this Spec declares, including the Makefile and the CI
workflows had the design needed them, and a `qa_override` is approved for a
QA partial caused only by the environment.

The design needs neither the Makefile nor a CI workflow. The governed set was
measured on the authoring branch at `2ce5abe8` with `GovernedPath`, through a
`go test -overlay` probe in `internal/speccheck` that wrote nothing to the
repository. It ran against every path the Tasks declare and every test file
the authoring prototype changed, and five of them are governed: four
historically bounded test files that a conversion changes, and the
repository-owned guide.

No live provider call is authorized by this record. Tests, Verification and QA
use temporary repositories, temporary clones of this repository and the
installed Go toolchain. No adopter repository is read or written.

## Why each governed path is unavoidable

- `docs/agents/specific-repository.md` gains one repository rule, outside
  every setup marker, so whoever writes the next test in a Parallel Test
  Package meets the rule before the rule test fails. It also documents the
  `internal/cli` test surface task_01 changes.
- `internal/cli/cli_test.go` holds sequential top-level tests that task_01
  converts or marks, and task_03 may replace one of their process-wide values
  with a per-test value.
- `internal/baseline/derived_ownership_test.go`,
  `internal/spec/coverage_test.go` and
  `internal/speccheck/mechanical_test.go` hold sequential top-level tests that
  task_02 converts.

## What is not governed

These paths are ordinary: `CONTEXT.md`,
`docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`,
`internal/testfixture/parallel_tests_test.go`, `internal/verifyselect/**`,
and every other test file the conversions change.

## Sanctioned regeneration

None. No Task runs a generator, and no derived file changes.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`,
  `.roundfixrc.yml`, the lint, formatter or test-runner configuration, or any
  production Go file other than `internal/verifyselect/contracts.go` and, for
  corrective task_06 only, `internal/cli` files that route a process-global read
  through `commandDependencies` with today's read as the default (operator
  decision 2026-10-09 on QA finding F2; no behavior change).
- No assertion, expected value, golden or fixture meaning changes.
- No change to archived Specs, existing Archive Records, existing history
  entries or `CHANGELOG.md`.
- No test, Verification command or QA row opens a network connection other
  than the QA gate's optional read of public CI logs, reads or writes the real
  `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
