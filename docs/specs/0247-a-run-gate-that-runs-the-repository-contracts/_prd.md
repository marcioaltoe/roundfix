---
spec: 0247-a-run-gate-that-runs-the-repository-contracts
status: active
created: 2026-10-07
surfaces: [backend, cli, docs]
---

# A Run gate that runs the repository contracts

A Run's Verification and its QA gate use `make verify-changed`, the
`defaults.verification` of `.roundfixrc.yml`. That target compiles no test
file under the `docscontract` or `repocontract` build tags. Only
`make verify-docs` and `repo-test` run those tests, and in practice that means
GitHub CI. On 2026-10-07 Spec 0245's Run and QA passed, and then CI failed
twice on checks the Run never ran:

- The suiteguard installation contract
  `TestEverySpawningPackageInstallsTheSuiteGuard` refused `internal/config`.
  A new test there spawned a process without `suiteguard.Main`, and the
  package was not in the guarded list.
- A Linux-only behavior: `go list` printed `go: downloading` on stderr for
  modules the host had not fetched.

Each failure cost a manual correction round after the merge was attempted.
The operator log records both, in entries 237 and 238 of
`~/.roundfix-operator/queue-interventions.md`. The Archive Record is
`docs/history/specs/0245-generated-records-that-hold-across-specs-and-platforms.md`.

This is a bug fix in the repository's own Verification. Roundfix's product
behavior does not change. The maintainer decided on 2026-10-07: "Contratos
sempre no gate (Recommended)". This minimal PRD exists for the downstream
artifact contract. The design lives in [_techspec.md](_techspec.md) and
ADR-0252.

## Prerequisites

None. No other Spec is active. Spec 0248, which will select
`TestRegenerationIsDeclared` by its inputs, builds on this Spec's directive
and is authored after this one merges.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The
  directive names `always`, `relevant` and `boundary` are lowercase words in a
  Go line comment, and the selector prints the existing build-tag names.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The selector reads local Go sources and Git, and every test
  uses a temporary repository and fake commands. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0252 (this Spec) decides Contract
  Relevance, the always set, the time budget and the known platform limit.
  ADR-0182 runs "the configured repository Verification, as the last
  Verification command" of every Settlement Check, so the contracts selected
  here become facts each Task settles on. ADR-0184: "A TechSpec states a
  command surface as a transcript", answered by the TechSpec's Surface
  Transcripts. `docs/agents/specific-repository.md`: "repository contracts
  validate at the pull request boundary". `make verify-docs` keeps running
  every contract, and task_03 adds that the selective gate runs the relevant
  ones first. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable. task_02 changes the `Makefile` and four
  Repository Contract Test files the Governed Path set protects. task_03
  changes `docs/agents/specific-repository.md`. Express maintainer
  authorization covers these: "Contratos sempre no gate (Recommended)" and the
  grant of the Governed Paths this Spec declares, including the Makefile,
  both on 2026-10-07, and the standing "Concedo". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`. The
  Spec-contained authorization record is
  `docs/specs/0247-a-run-gate-that-runs-the-repository-contracts/_authorization.md`.
  Bounded files: `Makefile`, `docs/agents/specific-repository.md`,
  `internal/baseline/derived_regeneration_repocontract_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/speccheck/governed_repocontract_test.go`,
  `skills/owned_skill_edit_repocontract_test.go`.

## Goals

- `make verify-changed` runs, on every change, the cheap global contracts.
  These are the suiteguard installation audit, every `docscontract` test in
  `internal/docscontract` and the governed-set contracts.
- `make verify-changed` runs each other Repository Contract Test when a
  changed path is a file in its package directory or matches its declared
  inputs.
- A contract can declare the paths that make it relevant, next to its test
  code, so a later change can opt an expensive contract in on narrow inputs.
- A typical core change grows the gate by no more than about 20 s on the
  maintainer's machine.
- A malformed declaration, or a change list Git cannot produce, never narrows
  the gate.

## Core Features

1. **Contract discovery and selection.** `cmd/verify-select -contracts`
   finds each Repository Contract Test and reads its directive. It selects
   tests from the changed paths and prints one `go test` invocation per build
   tag (ADR-0252).
2. **The gate runs the selection.** `make verify-changed` ends with
   `verify-changed-contracts`, which runs every printed invocation and fails
   on the first failure. The contract test files declare their relevance
   (ADR-0252).
3. **Docs and glossary.** `CONTEXT.md` gains **Repository Contract Test** and
   **Contract Relevance**. The repository rules say that the selective gate
   runs the relevant contracts and that `make verify-docs` still runs all of
   them.

## Non-Goals / Out of Scope

- Platform-only failures. CI is the Linux gate, and the selective gate does
  not emulate Linux. This is documented as a known limit (ADR-0252).
- Selecting `TestRegenerationIsDeclared` by its inputs. It stays `boundary`,
  as unrun by the selective gate as it is today, and Spec 0248 owns that
  change.
- Changing `make verify`, `make verify-docs`, `repo-test`, `docs-test`, the CI
  workflows or `.roundfixrc.yml`.
- Changing any contract test's body or assertion.

## Success Metrics

1. Success Metric: in a temporary clone, take a branch where `internal/config`
   lost its suiteguard wiring, which is the shape of the CI failure of
   2026-10-07. `make verify-changed` passed on that branch before this Spec,
   and it fails after it, naming `internal/config`.
2. Success Metric: on a typical core change (one comment line in
   `internal/config/config.go`, `GOFLAGS=-count=1`, a warm build cache),
   `make verify-changed` takes no more than 20 s longer than before.
3. Success Metric: a change only to `docs/user-guide` runs the always set and
   no `relevant` regeneration contract. A change to
   `.agents/skills/roundfix/SKILL.md` also runs the three regeneration
   contracts.
4. Success Metric: a malformed directive makes `make verify-changed` exit
   non-zero with the file and test named. An unresolvable base ref selects
   every contract that is not `boundary`.
5. Success Metric: every existing `./internal/verifyselect` test passes, and
   `make verify-docs` passes.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The operator's queue log, entries 237 and 238 of
  `~/.roundfix-operator/queue-interventions.md`. It records the two CI
  failures of 2026-10-07 and the correction each needed.
- Microsoft's description of Test Impact Analysis, "Accelerated Continuous
  Testing with Test Impact Analysis – Part 3"
  (devblogs.microsoft.com/devops, 2017). When a commit holds a file type the
  mapping does not know, it falls back to running all tests.
- Datadog's Test Impact Analysis documentation
  (docs.datadoghq.com/tests/test_impact_analysis). Its "tracked files", such
  as Makefiles and dependency manifests, run every test when they change.
- CircleCI's test impact analysis guide
  (circleci.com/docs/guides/test/set-up-test-impact-analysis). Its
  `test-selection-rules` extend selection to non-source files, or always run
  specific tests.

## Glossary

- adds: **Repository Contract Test**
- adds: **Contract Relevance**

## Decisions

- Each contract declares its relevance in a `//verify:` directive that
  `cmd/verify-select` parses; see ADR-0252.
- The always set is the suiteguard audit, the `internal/docscontract` tests
  and the governed-set contracts; see ADR-0252.
- The regeneration contracts run when Baseline assets, Baseline testdata or
  owned skills change; see ADR-0252.
