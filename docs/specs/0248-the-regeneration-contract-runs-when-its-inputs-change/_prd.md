---
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: active
created: 2026-10-07
surfaces: [backend, cli, docs]
---

# The regeneration contract runs when its inputs change

`TestRegenerationIsDeclared` proves that every path the declared regeneration
commands rewrite is covered by a `derived_paths` declaration in
`.roundfixrc.yml`. Spec 0245 added it after a module edit regenerated paths
no declaration covered, which entry 233 of the operator's queue log records.
Spec 0247 gave it the Contract Relevance `boundary`. The Archive Record is
`docs/history/specs/0247-a-run-gate-that-runs-the-repository-contracts.md`.

`boundary` was meant to keep the test where it already ran, but it ran
nowhere. The test carries the `docscontract` build tag in `internal/config`.
`make verify-docs` compiles `docscontract` files only in
`internal/docscontract`. Its `repo-test` step runs a hand-kept list of
`repocontract` tests that does not name it. Neither CI workflow runs any
other target that reaches it, and the release workflow runs no Repository
Contract Test at all. So the guard Spec 0245 paid for has run only by hand.

This is a bug fix in the repository's own Verification. Roundfix's product
behavior does not change. On 2026-10-07 the maintainer decided "Rodar só
quando relevante (Recommended)": run the test when a change touches the
Baseline modules, profiles, fixtures or `.roundfixrc.yml`, always in the CI
of main and of releases, and leave it out of `make verify-changed`
otherwise. This minimal PRD exists for the downstream artifact contract. The
design lives in [_techspec.md](_techspec.md) and ADR-0253.

## Prerequisites

None. No other Spec is active. Spec 0247 has merged, and its Contract
Relevance directives and `verify-changed-contracts` target are on main.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The new
  `-all` flag and the `verify-contracts` target are lowercase command words,
  and the selector prints existing build-tag names and Go package paths.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The selector reads local Go sources and Git, and every test
  uses temporary repositories and stub commands. No Task or QA row triggers a
  GitHub workflow. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0253 (this Spec) decides the
  derived relevance, the Full Contract Run and the named exclusions.
  ADR-0252: "A later change can make an expensive contract relevant by giving
  it narrow inputs"; this Spec is that change. ADR-0182 runs "the configured repository Verification, as the last
  Verification command" of every Settlement Check, so the selection becomes a
  fact each Task settles on. ADR-0184: "A TechSpec states a command surface
  as a transcript", answered by the TechSpec's Surface Transcripts.
  `docs/agents/specific-repository.md`: "repository contracts validate at the
  pull request boundary". task_04 corrects that rule's claim about
  `make verify-docs` and names the Full Contract Run. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable. task_03 changes the `Makefile`,
  `.github/workflows/ci-verify.yml` and `.github/workflows/release.yml`, and
  task_04 changes `docs/agents/specific-repository.md`. Express maintainer
  authorization covers them: "Rodar só quando relevante (Recommended)", the
  grant of the Governed Paths this Spec declares including the `Makefile`,
  the explicit authorization of CI workflow changes, all on 2026-10-07, and
  the standing "Concedo". Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. The Spec-contained authorization record is
  `docs/specs/0248-the-regeneration-contract-runs-when-its-inputs-change/_authorization.md`.
  Bounded files: `.github/workflows/ci-verify.yml`,
  `.github/workflows/release.yml`, `Makefile`,
  `docs/agents/specific-repository.md`.

## Goals

- `make verify-changed` runs `TestRegenerationIsDeclared` when a change
  touches an input its regeneration reads, and not otherwise.
- The inputs are derived from the `derived_paths` declarations, and a
  repository test fails when a declared path would not select the test.
- Every push to main and every release runs every Repository Contract Test,
  found by discovery rather than by a hand-kept list.
- The selective gate names each expensive contract it leaves out, so no
  exclusion is silent.
- `make verify-docs` keeps its scope and stays under 60 s of wall time on the
  maintainer's machine.

## Core Features

1. **Full Contract Run.** `cmd/verify-select -contracts -all` prints the
   invocations of every discovered contract, and `make verify-contracts` runs
   them. CI runs it on every push to main and in the release workflow before
   any build or publication (ADR-0253).
2. **Named exclusions.** The selective summary line names each `relevant`
   contract it did not select and each `boundary` contract (ADR-0253).
3. **Derived relevance.** `TestRegenerationIsDeclared` is `relevant` to the
   declarations file, every declared output and line path, and the code its
   declared commands run, and a repository test keeps that true (ADR-0253).
4. **Docs and glossary.** `CONTEXT.md` gains **Full Contract Run** and revises
   **Repository Contract Test** and **Contract Relevance**. The repository
   rules stop claiming that `make verify-docs` runs every contract.

## Non-Goals / Out of Scope

- Making `TestRegenerationIsDeclared` itself faster. It copies the repository
  and runs `make baseline-digests`; that cost is listed as a next contributor.
- Running the two build tags of the selective gate in one `go test`, which
  would overlap the regeneration contracts. Listed as a next contributor.
- Changing `make verify`, `make verify-changed`, `make verify-docs`,
  `docs-test`, `repo-test`, `REPO_CONTRACT_TESTS` or `.roundfixrc.yml`.
- Changing the pull request CI path.
- Changing any contract test's body or assertion.

## Success Metrics

1. Success Metric: in a temporary clone, a committed change that removes the
   formatter fixture declaration from `.roundfixrc.yml` makes
   `make verify-changed-contracts` run `TestRegenerationIsDeclared` and exit
   non-zero. Before this Spec the same change exits 0, because the test is
   never selected.
2. Success Metric: a change only to `docs/user-guide` or to
   `internal/daemon` does not select `TestRegenerationIsDeclared`, and the
   summary line names it as a relevant contract not selected. A change to
   `.roundfixrc.yml`, a Baseline module, a profile or a formatter fixture
   selects it.
3. Success Metric: `make verify-contracts` runs every contract that
   discovery finds, `TestRegenerationIsDeclared` included, and exits 0 at the
   audited head. Its wall time is recorded.
4. Success Metric: the push-to-main job of `ci-verify.yml` and the release
   workflow each run `make verify-contracts`; the release runs it after
   `make verify` and before the publication preflight. The pull request path
   of `ci-verify.yml` is unchanged.
5. Success Metric: `make verify-docs` exits 0 and takes no more than 60 s of
   wall time, warm, on the maintainer's machine (45.5 s before this Spec).
6. Success Metric: every existing `./internal/verifyselect` test passes.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The operator's queue log, entries 233 and 238 of
  `~/.roundfix-operator/queue-interventions.md`. Entry 233 records the Spec
  0245 QA failure in which a module edit regenerated paths no declaration
  covered, the class `TestRegenerationIsDeclared` guards. Entry 238 records
  its cost: 46 s with a shared build cache.
- CircleCI's test impact analysis guide
  (circleci.com/docs/guides/test/set-up-test-impact-analysis). By default the
  default branch runs all tests and feature branches run only impacted ones.
- WarpBuild's guide to running only the tests a change affects
  (warpbuild.com/guides/test-impact-analysis-github-actions). It names the
  safety net as the full suite on every push to the default branch.

## Glossary

- adds: **Full Contract Run**
- changes: **Repository Contract Test**
- changes: **Contract Relevance**

## Decisions

- `TestRegenerationIsDeclared` is `relevant` to inputs derived from the
  declarations; see ADR-0253.
- CI runs the Full Contract Run on every push to main and before every
  release; see ADR-0253.
- The selective summary line names the `relevant` contracts it leaves out and
  every `boundary` contract; see ADR-0253.
