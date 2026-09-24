---
spec: 0155-a-verify-that-runs-what-changed
status: active
created: 2026-09-24
surfaces: [backend, docs]
---

# A verify that runs what changed

Every repository Verification runs every test. A change to the Daemon runs the
Baseline's 81 seconds, the Baseline command tests inside `internal/cli` (35
seconds) and the skill tests; a change to a Baseline profile runs the whole
core. Measured on 2026-09-24, the Baseline and skills layer is about 26 percent
of test time (roughly 126 of 478 seconds), 25 percent of the Go code and 38
percent of the governed paths, and every one of those seconds is paid by every
change, including the ones that cannot affect it.

The cost compounds inside a Run. `.roundfixrc.yml` sets
`defaults.verification: make verify`, and that command is what the terminal QA
gate runs as the repository Verification of every Run. A correction to one core
function pays for the whole Baseline suite on each attempt.

The maintainer approved on 2026-09-24 a structure scored by Jev at 0.80 against
a separate Go module (0.18) and a separate repository (0.02): one repository,
with verification that selects the tests a change can affect. External research
found the same pattern — most monorepo CI speedup comes from selective
execution, not from splitting repositories.

## Project Constraints

- Identifier strategy: not applicable — no persisted identity is created or
  renamed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and Go tooling only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — a selective gate is build-tool and
  Verification configuration. Express maintainer authorization: the decision of
  2026-09-24 for selective gates in one repository, and the standing grant of
  2026-09-21 for governed paths a slice needs, recorded in
  [_authorization.md](_authorization.md); bounded files: `Makefile`,
  `.roundfixrc.yml`, `.github/workflows/ci-verify.yml`. Measured with
  `GovernedPath` over every changed path; the selector package is ordinary
  source. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- A change pays only for the tests it can affect.
- No test escapes verification: every test belongs to exactly one set.
- The complete gate stays available and stays the gate for `main` and release.

## Core Features

1. **Two test sets that cover everything once.** The core set and the Baseline
   set partition the suite: every package, and every test function in
   `internal/cli`, is selected by exactly one of them. A contract fails when a
   package or test is in neither or both.
2. **Selection by what changed.** `make verify-changed` compares the candidate
   with its base and runs the core set when core paths changed, the Baseline set
   when Baseline or skill paths changed, and both when both did. Formatting,
   `go vet ./...` and `go build ./...` always cover the whole tree, so a broken
   interface surfaces even when its tests are not selected.
3. **Fail safe.** When the base cannot be resolved or the change cannot be
   listed, `make verify-changed` runs everything.
4. **The complete gate stays complete.** `make verify` keeps running the whole
   suite, and remains the gate for pushes to `main` and for releases.
5. **The Daemon and pull requests use the selective gate.** Repository
   Verification in Runs and the pull request CI job run `make verify-changed`.

## Non-Goals / Out of Scope

- Splitting the Baseline and skills into another module or repository.
- Changing what `make verify` means, or the generated guides that describe it.
- Selecting tests at a finer grain than the two sets.
- Changing `make verify-docs`.

## Success Metrics

1. A change touching only core paths runs no Baseline-set test, measured by the
   selector's output for a fixture change set.
2. The partition contract fails when a package or an `internal/cli` test is
   removed from both sets or added to both.
3. An unresolvable base selects both sets.

## Recorded limits

The corrective ceiling, already extended once by the maintainer for this Spec,
was spent on Tasks 06 to 09. The second pre-PR review of 2026-09-24 found one
more gap, carried to Spec 0166:

- A change only under `docs/**`, or to Markdown at the root, selects no test
  set, yet tests read those files: `TestDurableTableLifecyclePolicyCoversEveryTable`
  reads `docs/user-guide/run-database-lifecycle.md`, and
  `TestBaselineExamplesParse` reads `README.md`. Such a change passes
  `make verify-changed` and fails only in the complete `make verify` that pushes
  to `main` still run. Reproduction: remove the
  `<!-- durable-table-lifecycle:begin -->` marker in a commit and run
  `go run ./cmd/verify-select -base HEAD~1`; it prints no set.

## Decisions

- **Keep `make verify` complete and add a selective target.** Changing the
  meaning of `make verify` would falsify the generated guides that describe it
  as the complete gate, which is governed Baseline output.
- **Select in a small Go program, not in shell.** The selection is logic that
  must be tested; a Makefile conditional cannot carry a contract.
- **Compile everything, test selectively.** `go build ./...` and `go vet ./...`
  are cheap and catch interface breaks across the boundary; the tests are where
  the time goes.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a selector that quietly drops a
test from both sets would make verification faster by making it blind.

## Research basis

Measured on this repository on 2026-09-24: per-package test time from a full
verify totalled 478 seconds, of which `internal/baseline` was 81, the Baseline
command tests in `internal/cli` 35 and `skills` 9. The QA gate's repository
Verification is `defaults.verification` in `.roundfixrc.yml`, read in
`internal/cli/implement.go` and passed to the Daemon as
`RepositoryVerification`. The decision and its scoring are recorded in the
secondbrain at `inbox/roundfix/2026-09-24-reestruturacao-da-fila-com-jev.md`.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
