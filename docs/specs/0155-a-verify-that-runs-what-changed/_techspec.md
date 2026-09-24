---
spec: 0155-a-verify-that-runs-what-changed
status: active
created: 2026-09-24
surfaces: [backend, docs]
---

# A verify that runs what changed

## Executive Summary

A small Go selector decides which of two test sets a change can affect. A new
`make verify-changed` target always formats, vets and builds the whole tree and
runs only the selected sets. `make verify` stays complete. Runs and pull request
CI switch to the selective target; `main` and releases keep the complete one.

## Project Constraints

- Identifier strategy: not applicable — no persisted identity is created or
  renamed. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and Go tooling only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection with
  `internal/speccheck/governed.go` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The two sets

The **Baseline set** is the packages `./internal/baseline/...`,
`./internal/baselineacp/...` and `./skills/...`, the test functions declared in
`internal/cli/baseline_*_test.go`, and the skill checks `skills-sync-check` and
`skills-check`.

The **core set** is every other package, and every other test function in
`internal/cli`.

`internal/cli` is the one package split between them. Its Baseline tests are
selected by name, and the name list is derived at run time from the
`func Test...` declarations in `internal/cli/baseline_*_test.go` — never kept by
hand, so a new Baseline test joins its set by being written in the right file.

## The selector

`internal/verifyselect` holds the logic; `cmd/verify-select` exposes it. Given a
base ref it lists changed paths — committed since the merge base, plus unstaged
and untracked — and classifies each:

- a Baseline path (`internal/baseline/`, `internal/baselineacp/`, `skills/`,
  `.agents/skills/`, `internal/cli/baseline_*`) selects the Baseline set;
- `go.mod`, `go.sum` and `Makefile` select both;
- any other Go source or test selects the core set;
- documentation outside the embedded directories selects neither.

It prints the selected sets and, on request, the package list and test-name
pattern of a set, so the Makefile and the contract read one definition. When the
base cannot be resolved or Git fails, it selects both.

## The targets

- `make verify-changed` runs `fmt-check`, `vet` and `build` over the whole tree,
  then the core set and the Baseline set as the selector decides.
- `make verify` is unchanged: the complete gate.

## Where each is used

- `.roundfixrc.yml` `defaults.verification` becomes `make verify-changed`, which
  is what the Daemon runs as repository Verification in every Run and review
  Batch.
- The pull request job in `.github/workflows/ci-verify.yml` runs
  `make verify-changed` against the pull request's base.
- Pushes to `main` and `.github/workflows/release.yml` keep `make verify`.

## API Contracts

1. The two sets partition the suite: every package and every `internal/cli` test
   is selected by exactly one.
2. A core-only change selects no Baseline-set test; a Baseline-only change
   selects no core-set test; `go.mod`, `go.sum` or `Makefile` select both.
3. An unresolvable base selects both sets.

## Coverage Map

- Goal 1 → The selector; API Contract 2.
- Goal 2 → The two sets; API Contract 1.
- Goal 3 → The targets; Where each is used.
- Core Feature 1 → The two sets; API Contract 1.
- Core Feature 2 → The selector; The targets; API Contract 2.
- Core Feature 3 → The selector; API Contract 3.
- Core Feature 4 → The targets.
- Core Feature 5 → Where each is used.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → The two sets, The selector.

## Integration Points

- **`.roundfixrc.yml`.** One value changes; its comment is updated to say the
  complete gate is `make verify`.
- **CI.** The pull request job gains a base for the selector; the full-history
  checkout it needs is already the pattern `release.yml` uses.
- **Generated guides.** Untouched: they describe `make verify`, which keeps its
  meaning.

## Testing Approach

1. **Selection.** Fixture change sets — core only, Baseline only, both, module
   files, documentation only — each select the expected sets.
2. **Partition contract.** Every package from `go list ./...` and every test
   function in `internal/cli` is selected by exactly one set; negative controls
   show the contract failing when an item is dropped from both or added to both.
3. **Fail safe.** An unresolvable base and a failing Git invocation each select
   both sets.
4. **Wiring.** The Makefile target exists and invokes the selector;
   `.roundfixrc.yml` and the pull request job name `make verify-changed`; `main`
   and release still name `make verify`.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The selector and its classification tests (depends on: none).
2. The partition contract (depends on: 1).
3. The Makefile target and the Daemon configuration (depends on: 1, 2).
4. The pull request CI job (depends on: 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A blind gate.** A selector that drops a test from both sets makes
  verification faster by making it blind. The partition contract is the
  control, with negative cases.
- **Coupling the sets miss.** A core change can break Baseline behavior without
  touching a Baseline path. `go build ./...` and `go vet ./...` catch interface
  breaks; behavioral coupling is caught by the complete gate on `main` and at
  release, which is the trade selective execution makes.
- **Base drift in a Run worktree.** The selector's base defaults to the
  repository's main branch; an unresolvable base fails safe to both sets.
