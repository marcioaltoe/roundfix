---
spec: 0245-generated-records-that-hold-across-specs-and-platforms
prd: _prd.md
created: 2026-10-07
---

# Generated records that hold across Specs and platforms — Technical Spec

## Executive Summary

Both records become outputs a command computes from the tree, never values an
author types. A new test in `internal/baseline` keeps the Module Version
Record, `internal/baseline/module-versions.json`: it refuses module content
that is not recorded under its version, and with `-record-module-versions` it
chooses the version, the way Spec 0228's record chooses an owned skill's. The
Baseline catalog's derived-path declaration runs that step before `make
baseline-digests` and covers each module's top-level `version` line, so a
merge of two items that changed one module ends on the next free version. The
Coverage Record stops running `go test -list` for its content: `go list`
names each package's test files under every release `GOOS`, and the Go parser
finds the tests, so any host writes the same bytes; a host-only anchor test
keeps the parser honest against the toolchain. The trade-off accepted is a
second test enumeration that must agree with Go's own (the anchor test pins
it), in exchange for a record that no longer depends on the machine that
wrote it (ADR-0250).

## Project Constraints

- Identifier strategy: applicable — new names: the record file
  `internal/baseline/module-versions.json` with schema
  `roundfix/baseline-module-versions/v1`, the test
  `TestEveryBaselineModuleVersionIsRecorded` and its flag
  `-record-module-versions`, and the Coverage Record fields `platforms` and
  `platformTests`. Module identifiers and integer versions keep their scheme.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and the local `go`
  binary only (`go list`, `go test -list`, `go test`); no test or Verification
  command opens a network connection, and fixture modules have no
  dependencies. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0250 (this Spec) decides the record
  step and the platform-neutral collection, ADR-0250: "the record step chooses
  it" and ADR-0250: "The coverage record lists tests for every release
  platform". ADR-0189 is the rule applied to modules, ADR-0189: "A version
  names one content". ADR-0233 supplies line-scoped derived paths, ADR-0233:
  "may then change a line-scoped path only on matching lines", and this Spec
  declares the module version lines that way. A conflict confined to declared
  derived paths is still resolved by regeneration under ADR-0192, ADR-0192:
  "runs each matched declaration's regeneration command", now also for the new
  record and the Coverage Record. ADR-0149 resolves a command's outputs from
  `_ownership.yml`, ADR-0149: "the grant names the command"; the new record
  sits outside every scanned digest root, so no ownership record changes.
  ADR-0130 and ADR-0179 bound the Governed Paths, ADR-0182 settles each Task
  on its Verification, ADR-0184 binds Surface Transcripts and this Spec
  declares none, and the gate is bound by ADR-0080, ADR-0091, ADR-0096,
  ADR-0097, ADR-0104 and ADR-0167. ADR-0166 records a path a Task changed
  without declaring it, ADR-0240 decides when a QA partial qualifies, and
  ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this
  Spec by citation, stage and receipt. ADR-0194, ADR-0195 and ADR-0210 cite
  ADR-0097 but decide what a QA row records, when it is observed again and its
  evidence snapshot, ADR-0229 cites ADR-0167 but decides operator archives,
  and ADR-0237 cites ADR-0229 but decides a Delivery Retry after an outside
  merge; this Spec changes none of them, so none applies. ADR-0247 and
  ADR-0248 govern archives and the history sanitize; this Spec changes
  neither, so neither applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — 11 declared files are Governed Paths,
  measured with `GovernedPath` through a `go test -overlay` probe that wrote
  nothing. Express maintainer authorization: "Autorizar os dois"
  (2026-09-30), "Concedo" and "Pode seguir nessa ordem" (2026-10-07). No
  Makefile, lint, CI or `go.mod` change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0245-generated-records-that-hold-across-specs-and-platforms/_authorization.md`;
  bounded files: `.roundfixrc.yml`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, `internal/baseline/assets/modules/backend.json`, `internal/baseline/assets/modules/cli-surface.json`, `internal/baseline/assets/modules/external-triage.json`, `internal/baseline/assets/modules/monorepo.json`, `internal/baseline/assets/modules/rust.json`, `internal/baseline/assets/modules/tui-surface.json`, `internal/spec/coverage_test.go`.

## Measured starting point

Measured on 2026-10-07 at `66f2d85b`:

- Nothing reads a module's top-level `version` beyond the catalog check that
  it is an integer of at least 1 (`internal/baseline/catalog_load.go`
  `validateVersions`). The number is a review convention: 0241 was caught by
  review, and 0242 and 0243 changed `context-workflow.json` without raising
  it.
- Ten modules hold `  "version": <n>,` on its own line. Six (`backend`,
  `cli-surface`, `external-triage`, `monorepo`, `rust`, `tui-surface`) put
  `schemaVersion`, `id`, `version`, `kind` and `title` on one line. In a
  scratch clone, splitting those six lines, `make baseline-digests` and the
  Managed Refresh changed exactly `catalog.digest`,
  `catalog.normalized.json`, four plan-characterization goldens and
  `docs/agents/setup-context.json`; a second refresh changed nothing, and
  `go test ./...` passed (8,511 tests).
- Test files with an operating-system constraint: `internal/cli`
  (`detach_fixture_group_darwin_test.go`, `_linux`, `_other`,
  `orphan_unix_test.go`) and `internal/store` (`process_info_unix`,
  `reclaim_unix`, `process_unix`, `process_windows`, `process_linux`). The
  other constrained files carry `repocontract` or `docscontract` tags.
- A static collection (`go list` per `GOOS` with `CGO_ENABLED=0`, then the Go
  parser) found 4,401 tests on `darwin`, 4,403 on `linux` and 4,362 on
  `windows`, 48 of them platform-limited: one macOS-only, three Linux-only,
  three Windows-only and 41 Unix-only. On macOS it matched `go test -list`
  with zero differences. `GOOS=linux go test -list` on macOS fails with
  `exec format error`.
- `dist/npm/platforms.json` ships `darwin`, `linux` and `windows`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Module Version Record | new `internal/baseline/module_versions_test.go`, new `internal/baseline/module-versions.json` | check, record step, seeded record |
| Module layout | six module JSON files | top-level header split, one key per line |
| Coverage Record | `internal/spec/coverage_test.go`, new `internal/spec/coverage_platform_test.go`, `docs/references/coverage-record.json` | static collection, platform fields, anchor test |
| Derived declarations | `.roundfixrc.yml`, `internal/config/verification_tools_test.go` | record step before digests, module version lines, coverage declaration |
| Guidance | `docs/agents/specific-repository.md`, `CONTEXT.md` | the repository rule and two terms |

The record and the collection are test code, as the owned-skill record is: no
production package, flag or exported symbol is added.

## Implementation Design

### Interfaces

```go
// internal/baseline/module_versions_test.go
var recordModuleVersions = flag.Bool("record-module-versions", false, "...")
const moduleVersionsSchema = "roundfix/baseline-module-versions/v1"
const moduleVersionsPath = "module-versions.json" // relative to internal/baseline
const recordModuleVersionsCommand = "go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1"
type moduleVersionEntry struct {
	Version int64  `json:"version"`
	Digest  string `json:"digest"` // "sha256:<64 hex>"
}
type moduleVersionRecord struct {
	SchemaVersion string                          `json:"schemaVersion"`
	Modules       map[string][]moduleVersionEntry `json:"modules"`
}
func moduleContent(data []byte) (moduleVersionEntry, error)
func checkModuleVersionRecord(record moduleVersionRecord, current map[string]moduleVersionEntry, recording bool) (moduleVersionRecord, error)
func recordModuleVersionFiles(record moduleVersionRecord, modulesDir string) (moduleVersionRecord, error)
```

```go
// internal/spec/coverage_test.go
type CoverageRecord struct {
	Platforms     []string                       `json:"platforms,omitempty"`
	Packages      map[string][]string            `json:"packages"`
	PlatformTests map[string]map[string][]string `json:"platformTests,omitempty"`
}
func coveragePlatforms(repoRoot string) ([]string, error)
func collectCoverageRecordFor(repoRoot string, platforms []string) (CoverageRecord, error)
func collectCoverageRecord(repoRoot string) (CoverageRecord, error) // every release platform
func collectToolchainCoverage(repoRoot string) (CoverageRecord, error) // today's go test -list parse
```

### Invariants

1. `moduleContent` decodes a module with `UseNumber` into a JSON object whose
   top-level `version` is an integer of at least 1, deletes that key,
   re-encodes with `encoding/json` and digests the bytes with SHA-256.
   Whitespace, key order and the version number do not change the digest;
   any other byte of content does.
2. Every module file holds exactly one line matching `^  "version": [0-9]+,$`,
   and its number is the decoded version. Otherwise the check and the record
   step refuse, naming the module and that pattern.
3. Without the flag, the check reads every `assets/modules/*.json` and the
   record and refuses, each message ending with `run <recordModuleVersionsCommand>`:
   a missing record; a schema other than `moduleVersionsSchema`; a recorded
   module with entries not strictly ascending; a recorded module no longer in
   the catalog; content whose digest differs from the one recorded under its
   version (`<id>: content changed under version <n>`); and a version that is
   not recorded (`<id>: version <n> is not recorded`).
4. With the flag: a missing record starts empty, and a recorded module no
   longer in the catalog is dropped. For each module in name order, content
   recorded under its version is left alone; a version above every recorded
   one (or a module with no entries) is appended as it is; otherwise the
   version becomes the highest recorded plus one, the module's version line is
   rewritten to it with every other byte kept, and that version and digest are
   appended. A recorded entry is never changed or removed while its module is
   in the catalog.
5. The record step validates every module and the whole history before its
   first write, so a refusal writes nothing. It writes the record with
   two-space indentation and a final newline, only when it changed. A second
   run changes nothing.
6. `coveragePlatforms` returns the sorted, unique `goos` values of
   `dist/npm/platforms.json` and refuses an unreadable or empty matrix.
7. `collectCoverageRecordFor` runs, per platform, `go list -tags docscontract
   -f '{{.ImportPath}}\t{{.Dir}}\t{{join .TestGoFiles ","}}\t{{join .XTestGoFiles ","}}' ./...`
   in `repoRoot` with `GOOS=<platform>`, `CGO_ENABLED=0` and `GOWORK=off`. It
   skips packages under `roundfix/docs/`, as today, and parses each listed
   file once.
8. A test is a top-level function without receiver, type parameters or
   results, whose name passes `isTestFunctionName`, and whose single
   parameter has type `*<name>.T`, where `<name>` is the file's local name for
   the `"testing"` import. `TestMain` is therefore excluded.
9. The record lists every package any platform lists. `Packages[p]` holds the
   tests built on every platform; `PlatformTests[p][t]` holds the sorted
   platforms of a test built on only some. A package without such a test has
   no `PlatformTests` key, and the record is the same bytes on any host.
10. The comparison works per test on platform sets. A record with no
    `platforms` counts as one implicit platform, so today's comparison tests
    pass unedited. A platform the record holds and the collection lacks is a
    regression; a new one is an addition. When the lost or gained platforms
    are all of the record's, the message is today's text; otherwise it ends
    with ` on <p1>, <p2>`. Package-level messages are unchanged.
11. `validateCoverageRecord` also refuses: platforms that are not sorted,
    unique, lower-case tokens; a `PlatformTests` package missing from
    `Packages`; a platform list that is empty, unsorted, not a proper subset
    of `Platforms`, or present without `Platforms`; and a test listed in both
    places.
12. `TestCoverageEquivalence` compares the record with
    `collectCoverageRecord`; `-update-coverage-record` writes that collection.
    `TestCoverageCollectionMatchesTheToolchainOnThisPlatform` compares
    `collectCoverageRecordFor(repo, []string{runtime.GOOS})` with
    `collectToolchainCoverage(repo)` and names every difference.

### Data Models

The Module Version Record, seeded with each module's current version:

```json
{
  "schemaVersion": "roundfix/baseline-module-versions/v1",
  "modules": {
    "autonomous-work": [
      {
        "version": 12,
        "digest": "sha256:<64 hex>"
      }
    ]
  }
}
```

The Coverage Record gains two fields:

```json
{
  "platforms": ["darwin", "linux", "windows"],
  "packages": {"roundfix/internal/store": ["TestAgentSelection..."]},
  "platformTests": {
    "roundfix/internal/store": {"TestParseProcStatStartTimeRejectsMissingField": ["linux"]}
  }
}
```

The derived declarations in `.roundfixrc.yml`, in this order:

```yaml
delivery:
  derived_paths:
    - paths:
        - internal/baseline/module-versions.json
        - internal/baseline/testdata/catalog.digest
        - internal/baseline/testdata/catalog.normalized.json
        - internal/baseline/testdata/catalog.diagnostics.golden.json
        - internal/baseline/testdata/plan-characterization/*.golden.json
      lines:
        paths: [internal/baseline/assets/modules/*.json]
        match: '^  "version": [0-9]+,$'
      regenerate: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 && make baseline-digests
    # setup-context.json and owned-skill-versions.json declarations unchanged
    - paths:
        - docs/references/coverage-record.json
      regenerate: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
```

Two items that change one module both touch `catalog.digest`, so the first
declaration always matches their conflict; its command records the merged
content before the digests are rebuilt, and the Setup Manifest declaration
runs after it.

### API Contracts

1. API Contract: record step. Input: `go test ./internal/baseline -run
   '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions
   -count=1` on a tree where a module's content is not recorded under its
   version. Output: exit 0; the module's version line holds the chosen
   version and the record gains that version and digest. Failure: exit 1 and
   no file written for a module layout or history refusal.
2. API Contract: module check. Input: the same test without the flag.
   Output: exit 0 on a recorded tree. Failure: exit 1 with Invariant 3's
   message naming the module and the record command.
3. API Contract: Coverage Record. Input: `go test ./internal/spec -run
   '^TestCoverageEquivalence$' [-update-coverage-record]`. Output: exit 0
   and, with the flag, the same record bytes on any host. Failure: exit 1
   naming each regression and, where a test is platform-limited, its
   platforms.

### Surface Transcripts

None. No `roundfix` command, flag, output or exit code changes. The two
contracts above are `go test` invocations of repository tests, asserted by
the tests the Tasks add.

## Coverage Map

- Goal 1 → Invariants 1-5; API Contracts 1 and 2; task_01.
- Goal 2 → Data Models (the first declaration); Invariant 4; task_03; QA
  merge ablation.
- Goal 3 → Invariants 6-9; API Contract 3; task_02.
- Goal 4 → Invariants 10-12; task_02.
- Core Feature 1 → Module Version Record; task_01.
- Core Feature 2 → Coverage Record collection; task_02.
- Core Feature 3 → Derived declarations; task_03.
- Core Feature 4 → Guidance; task_03 and ADR-0250.
- Success Metric 1 → Testing Approach 1; task_01.
- Success Metric 2 → `TestEveryBaselineModuleVersionIsRecorded`; Invariant 5; task_01.
- Success Metric 3 → Data Models (the first declaration); QA merge ablation in task_04.
- Success Metric 4 → Invariant 9; task_02.
- Success Metric 5 → Testing Approach 2; Invariant 12; task_02.

## Integration Points

- The local `go` binary: `go list` under several `GOOS` values, `go test
  -list` on the host, and `go test` for the record step.
- The derived merge of ADR-0192 and ADR-0233, unchanged in code; only this
  repository's declarations change.
- `make baseline-digests` and the Managed Refresh, run unchanged after the
  header split and after a record step.

## Testing Approach

1. `internal/baseline/module_versions_test.go`, over temporary module
   directories and records: `TestModuleContentDigestIgnoresVersionAndLayout`,
   `TestAModuleChangedUnderARecordedVersionIsRefused`,
   `TestAnUnrecordedModuleVersionIsRefused`,
   `TestRecordingRaisesAModuleAboveTheHighestRecordedVersion` (content
   changed under a recorded version, and a lower unrecorded version),
   `TestRecordingKeepsAHigherModuleVersion`,
   `TestRecordingNeverRewritesARecordedModuleEntry`,
   `TestAModuleVersionNotOnItsOwnLineIsRefused`,
   `TestRecordingRewritesOnlyTheTopLevelVersionLine` (nested `version` lines
   keep their bytes) and `TestRecordingWritesNothingWhenAModuleIsRefused`;
   plus `TestEveryBaselineModuleVersionIsRecorded` over this repository.
2. `internal/spec/coverage_platform_test.go`:
   `TestCoverageCollectionListsEachPlatformOnlyTest` builds a temporary
   module with its own `go.mod`, `dist/npm/platforms.json` and test files
   constrained by `_darwin`, `_linux`, `//go:build windows` and `unix`, and
   asserts the record on any host; `TestCoverageCollectionMatchesTheToolchainOnThisPlatform`;
   `TestCompareCoverageRecordsReportsAPlatformRegression`;
   `TestCoverageRecordRefusesAMalformedPlatformEntry`; and
   `TestCoveragePlatformsComeFromTheReleaseMatrix`. The existing tests of
   `internal/spec/coverage_test.go` pass unedited.
3. `internal/config/verification_tools_test.go`:
   `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` expects the new
   declarations, the one declared break there.

## Build Order

1. The Module Version Record, its tests, the six header splits, the seeded
   record, `make baseline-digests` and the Managed Refresh, task_01 (depends
   on: none).
2. The platform-neutral Coverage Record, its tests and the re-recorded
   record, task_02 (depends on: none; it shares no file with step 1).
3. The derived declarations, their config test, the repository rule and the
   glossary terms, task_03 (depends on: 1, 2, whose commands it declares).
4. Terminal QA, task_04 (depends on: 3).

## Risks & Considerations

- This Spec's own merge: The queue resolves a conflict with the
  declarations of the default branch. Until this Spec merges, its own item
  meets the old declaration: if Spec 0244 changed a module first, the operator
  runs the record command once by hand on the item. Every later item gets
  the new declaration.
- A skipped version: A Task that records, edits again and records again
  takes the next number, and the first stays recorded, as ADR-0233 accepts.
- Parser drift: A Go release that changes what `go test` treats as a test
  shows up first as an anchor-test failure on the host, not as a silent
  record change.
- Collection time: Three `go list` runs and one parse of each test file
  replace one `go test -list ./...` in `TestCoverageEquivalence`; the anchor
  test keeps that one toolchain run.

## Glossary

- adds: **Module Version Record**
- adds: **Coverage Record**

## Research basis

- Secondbrain: `wiki/index.md` was read, then
  `qmd query "parallel specs generated record version collision coverage record platform build constraints" --all --files --min-score 0.3`
  returned the adopted Backlog Entry and ADR-0233 (score 0.62) and nothing
  newer; a second query on Go build constraints returned `docs/agents/go.md`,
  whose rule "build the affected non-test packages for every operating system
  its constraints name" this design extends to the record.
- Exa: the Go command documentation (<https://pkg.go.dev/cmd/go>) and
  `go/build` (<https://pkg.go.dev/go/build>) state that a `_GOOS` file name is
  an implicit build constraint, that `go list` reports `TestGoFiles` and
  `XTestGoFiles`, and that `go test` executes a test binary per package, which
  is why the collection uses `go list` and the parser.
- Spec 0228's record at `a30dc847` and ADR-0233 supplied the record-step
  pattern this Spec applies to modules.

## Decisions

- The Module Version Record chooses a module's version and is regenerated
  before the digests at merge. See ADR-0250.
- The Coverage Record is collected statically for every release `GOOS` and
  anchored to `go test -list` on the host. See ADR-0250.
- The six one-line module headers are split so every module's version is a
  line-scoped path. See ADR-0250.
