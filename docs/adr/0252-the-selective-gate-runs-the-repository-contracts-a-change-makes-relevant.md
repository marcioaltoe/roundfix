---
status: accepted
created_at: 2026-10-07T00:00:00Z
updated_at: 2026-10-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The selective gate runs the repository contracts a change makes relevant

`make verify-changed` is the Verification of Runs and of the QA gate, and it
compiled no test file under the `docscontract` or `repocontract` build tags.
Those tags hold the Repository Contract Tests, which check the repository
itself: its documents, derived artifacts and test wiring. Only
`make verify-docs` ran them, and in practice that meant GitHub CI. On
2026-10-07 a Spec's Run and QA passed, and then CI failed on
`TestEverySpawningPackageInstallsTheSuiteGuard`: a new `internal/config` test
spawned a process without `suiteguard.Main`. Each such miss cost a correction
round after the Pull Request opened. The maintainer chose "Contratos sempre no
gate (Recommended)": the selective gate runs the contracts of the affected
packages plus the cheap global ones, and stays selective and fast.

**Contract Relevance.** Each Repository Contract Test declares when the
selective gate runs it. It uses a `//verify:` line comment, either in the
file's header before the `package` clause, which covers every test in the
file, or in one test's doc comment, which overrides the file's:

- `//verify:always` runs on every change. It is for contracts whose inputs
  are spread across the repository and that cost about a second.
- No directive means the package class. The test runs when a changed path is
  a file directly in its package directory.
- `//verify:relevant <pattern>...` adds declared inputs to the package
  directory. A pattern ending in `/` matches everything under that directory.
  Any other pattern is matched by `path.Match` against the whole path.
- `//verify:boundary <reason>` keeps the test out of the selective gate. It
  runs only where it ran before.

A change to `go.mod`, `go.sum` or `Makefile` selects every contract that is
not `boundary`. So does a change list that Git cannot produce. A malformed
directive fails the gate rather than skipping a test. `cmd/verify-select`
reads the directives from the source with `go/parser`, so the declaration sits
next to the test and has no central list to drift from. A later change can
make an expensive contract relevant by giving it narrow inputs.

**What is always and what is relevant.** The following was measured on
2026-10-07 on the maintainer's 10-core machine, with Go's build cache warm.
The second figure in each pair is the extra compile cost when only the
untagged packages were already built.

| Contract | Class | Cost (warm / extra compile) |
| --- | --- | --- |
| Suiteguard installation audit in `internal/suiteguard` | always | 0.7 s / 1.0 s |
| Every `docscontract` test in `internal/docscontract` | always | 0.8 s / 7.4 s |
| The four governed-set contracts in `internal/speccheck` | always | 1.1 s / 2.4 s |
| `TestRepositoryGateRunsTheAnalyzer` | relevant: `internal/baseline/analyzer/` | 0.6 s / 2.1 s |
| The two derived-regeneration contracts in `internal/baseline` | relevant: Baseline assets, Baseline testdata, owned skills | 43 s wall |
| `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical` | relevant: the same inputs | 15 s |
| Each package's own suiteguard check | package | about 0.4 s per package |
| `TestRegenerationIsDeclared` | boundary | 43 s |

Before this change, `make verify-changed` ran a typical core change, one
comment line in `internal/config`, with `GOFLAGS=-count=1` and a warm cache.
The always set adds about 3 s to that, and at most about 11 s when its tagged
binaries must compile. That is inside the budget of about 20 s the
maintainer set for a typical change. A change to Baseline assets or an owned
skill also selects the regeneration contracts, about 45 s of wall time. This
cost is deliberate: those are the contracts whose verdict such a change can
move, and CI already pays it. `TestRegenerationIsDeclared` stays `boundary`
here, and selecting it by its inputs is left to a later change.

## Consequences

The selective gate still runs on the host. A failure that only another
platform produces is a known limit: CI is the Linux gate. One example is
`go list` printing `go: downloading` on stderr for modules the host never
fetched. `make verify-docs` keeps running every Repository Contract Test at
the Pull Request boundary, so CI runs the selected contracts twice. That costs
a few seconds on a typical change. We rejected two alternatives. Running every
contract on every change would cost about 45 s even for a typical change.
Keeping a central list of tests and paths in `cmd/verify-select` would drift
from the tests it describes.
