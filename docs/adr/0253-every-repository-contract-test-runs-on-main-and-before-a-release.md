---
status: accepted
created_at: 2026-10-07T00:00:00Z
updated_at: 2026-10-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Every Repository Contract Test runs on main and before a release

ADR-0252 gave each Repository Contract Test a Contract Relevance and left
`TestRegenerationIsDeclared` as `boundary`. That test proves that every path
the declared regeneration commands rewrite is covered by a `derived_paths`
declaration in `.roundfixrc.yml`. It copies the repository and runs those
commands, which took 52.8 s on 2026-10-07. No Makefile target and no workflow
ran it: `docs-test` compiles only `internal/docscontract`, and `repo-test`
runs a hand-kept list of `repocontract` tests that never named it. The
maintainer chose "Rodar só quando relevante (Recommended)": run it when a
change touches the Baseline modules, profiles, fixtures or `.roundfixrc.yml`,
always in the CI of main and of releases, and leave it out of the selective
gate otherwise.

**Relevance derived from the declarations.** The test becomes
`//verify:relevant` over every input its regeneration reads. Those are the
declarations file, every path a declaration names as an output or as a
line-scoped input, and the code the declared commands run as their
generators: `internal/baseline/`, `skills/`, `.agents/skills/`,
`docs/agents/`, `docs/references/coverage-record.json`, `internal/spec/`,
`internal/cli/baseline_*`, `cmd/roundfix/`, `internal/suiteguard/` and
`internal/suiteguardcontract/`. A repository test asserts that a change to
any declared output or line path selects the test, so a new declaration
cannot fall outside it unnoticed. Of the last 100 commits on main, 29 touched
one of these inputs; 13 touched the narrower set the maintainer named, and 26
already selected ADR-0252's regeneration contracts.

**The Full Contract Run.** `cmd/verify-select -contracts -all` prints one
`go test` invocation per build tag for every contract that discovery finds,
whatever its class, and `make verify-contracts` runs them. Discovery replaces
a hand-kept list as the source of "every contract", which is how
`TestRegenerationIsDeclared` fell out of every gate. CI runs the Full Contract
Run on every push to main and in the release workflow before any build or
publication. Published practice does the same: CircleCI's test impact
analysis runs all tests on the default branch and only impacted tests on
feature branches, and the usual safety net under affected-test selection is a
full run on every push to the default branch.

**Exclusions are named.** In its selective mode the selector's summary line
keeps its counts and adds the name and package of each `relevant` contract it
did not select and of each `boundary` contract. A skipped expensive contract
is therefore visible in every Run gate log, not only as a count.

## Consequences

A pull request whose change touches none of the inputs pays nothing for the
test, and `make verify-docs` keeps its scope and its measured 45.5 s warm
wall time. A pull request that touches them pays about 52 s more in
`make verify-changed`, on top of ADR-0252's 43 s regeneration contracts when
those are selected too, because the selector runs one tag after the other. A
miss of the selection rule now surfaces on the push to main, after the merge,
rather than before it; the Full Contract Run on main is that safety net, and
it costs about 100 s of CI time per push to main and per release. We rejected
three alternatives. Adding the test to `make verify-docs` would make that
target about 98 s on every pull request. A hand-kept list of every contract
for main would drift as `REPO_CONTRACT_TESTS` did. An environment switch that
widens `make verify-docs` on main would hide in the workflow what the step
runs.
