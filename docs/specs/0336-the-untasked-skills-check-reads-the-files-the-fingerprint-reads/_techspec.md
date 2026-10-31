---
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
prd: _prd.md
created: 2026-10-31
---

# The untasked-skills check reads the files the fingerprint reads — Technical Spec

## Executive Summary

The untasked-skills rule maps each path a Task declares to Behavior Surfaces
with `skillcoverage.Map.SurfacesForPath`, which matches the map's source
patterns literally. The fingerprint test filters Go test files out of the
`config keys` read with its own `strings.HasSuffix(entry.Name(), "_test.go")`.
The two file sets differ by exactly that filter. The fix adds one exported
predicate, `skillcoverage.FingerprintReads`, makes `SurfacesForPath` return
no surface for a file the predicate refuses, and makes the fingerprint test
call the same predicate in place of its own filter. `internal/speccheck`
changes no production line; its tests gain the regression cases.

The trade-off accepted: the lookup itself refuses Go test files, rather than
the Spec Check filtering before it calls the lookup, so every reader of the
lookup gets the fingerprint's file set. The Spec Check is its only reader
today (`_prd.md` → Decisions).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0256 defines the sources of a
  Behavior Surface and the authoring rule; the design follows it, as
  `_prd.md` → Project Constraints states and ADR Alternatives below relates.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AS batch 2 answers of 2026-10-31 quoted in the authorization record; the
  Governed Paths are the Roundfix skill's `spec` reference and its `SKILL.md`
  pair; no Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-31 at `0c27cdc9`; `_prd.md` → Validation records the runs.

1. `detectSkillsDeclaration` (`internal/speccheck/skills_declaration.go`)
   calls `coverage.SurfacesForPath(ref.Path)` for every path a non-QA Task
   declares under `interface:`, `creates:` or `deletes:`, and reports
   `SC-SKILLS-UNTASKED` for each matched covered surface without a covering
   Task.
2. `SurfacesForPath` (`internal/skillcoverage/skillcoverage.go`) matches a
   source ending in `/` as a prefix and any other source with `path.Match`.
   It never looks at the file name, so `internal/config/*.go` matches
   `internal/config/config_test.go`.
3. `computeBehaviorSurfaces` (`internal/docscontract/skill_coverage_test.go`)
   fingerprints the surfaces:
   - each command surface from the built CLI's help, which the Go toolchain
     compiles without any Go test file;
   - `config keys` from the `yaml` tags of the top-level `.go` files of
     `internal/config`, skipping directories and `_test.go` files;
   - `exit codes` from `internal/cli/cli.go`;
   - each guide surface from the bytes of one Markdown file under
     `docs/user-guide/`, nested directories included.
4. Other exclusions do not diverge. The `*` of `internal/config/*.go` never
   crosses a `/` and requires `.go`, which matches the `config keys` read's
   directory and suffix filters. Every other source in the map is an exact
   non-test file. No fingerprint excludes `testdata` or generated files; a
   `testdata` exclusion would split the guide fingerprints, which read every
   Markdown file under `docs/user-guide/`.
5. `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` writes an
   unparsable `internal/config/config_test.go` and requires the `config keys`
   fingerprint to ignore it.

## System Architecture

No package, command, flag or file format is added. Changes by component:

- `internal/skillcoverage/skillcoverage.go`: the predicate and the lookup's
  guard.
- `internal/docscontract/skill_coverage_test.go`: the `config keys` read
  calls the predicate.
- `internal/skillcoverage/skillcoverage_test.go`,
  `internal/speccheck/skills_declaration_test.go`: the regression tests.
- The spec command guide and the Roundfix skill's spec reference.

## Implementation Design

### Interfaces

```go
// internal/skillcoverage/skillcoverage.go

// FingerprintReads reports whether a Behavior Surface fingerprint can depend
// on the file at the slash-separated filePath.
func FingerprintReads(filePath string) bool

// SurfacesForPath: unchanged signature; nil when !FingerprintReads(filePath).
func (m Map) SurfacesForPath(filePath string) []Surface
```

### Data Models

None. No record or file format changes.

### API Contracts

1. API Contract: the predicate. `FingerprintReads(filePath)` returns `false`
   exactly when the last element of `filePath` ends in `_test.go`, and `true`
   for every other path, `testdata` paths and names that merely contain
   `_test` or `test` included.
2. API Contract: the lookup. `SurfacesForPath` returns nil for a path
   `FingerprintReads` refuses and otherwise matches as today.
   `detectSkillsDeclaration` keeps its code, summary, location and fix.
3. API Contract: the fingerprint. The `config keys` read in
   `computeBehaviorSurfaces` skips a file when `FingerprintReads` refuses its
   name, in place of its own `_test.go` suffix test; the computed
   fingerprints, and so `docs/references/behavior-surfaces.json`, keep their
   values for unchanged sources.
4. API Contract: the rule in writing. `docs/user-guide/commands/spec.md`, in
   the paragraph that introduces the skills rule, and
   `.agents/skills/roundfix/references/spec.md`, in the paragraph on the
   Skills Declaration, add one sentence containing `is never a Behavior
   Surface source`: a Go test file, whose name ends in `_test.go`, is never
   a Behavior Surface source, because no fingerprint reads it, so declaring
   one needs no skills Task even when a source pattern matches it.

### Surface Transcripts

None. No command output changes; the check only stops emitting a finding.

### Invariants

```text
1. A non-test file matched by a covered surface's source still gives SC-SKILLS-UNTASKED.
2. A Go test file never gives SC-SKILLS-UNTASKED, whatever source pattern matches it.
3. The config keys fingerprint and the lookup refuse the same files, through one function.
4. No recorded fingerprint changes for unchanged sources.
```

### Declared breaks

- `SC-SKILLS-UNTASKED` no longer reports a Task whose only surface-matched
  declarations are Go test files. No existing test asserts that report:
  the prototype ran `./internal/skillcoverage`, `./internal/speccheck` and
  the `internal/docscontract` skill coverage tests unchanged and green.

## Coverage Map

- Goal 1, Core Features 1 and 2, Success Metrics 1 and 3 →
  `FingerprintReads` and `SurfacesForPath` (API Contracts 1
  and 2), Invariants 2.
- Goal 2, Success Metric 2 → `SurfacesForPath` (API Contract
  2), Invariants 1.
- Goal 3, Core Feature 3, Success Metric 4 → `computeBehaviorSurfaces`
  (API Contract 3), Invariants 3 and 4.
- Goal 4, Core Feature 4, Success Metric 5 → the guide and the reference
  (API Contract 4).

Deviations from the PRD: None.

## Absence Evidence

1. A shared predicate for the fingerprint's file set — command: `grep -rn '_test.go' internal/skillcoverage internal/speccheck/skills_declaration.go`; result: no match; the only `_test.go` filter is the inline `strings.HasSuffix` in `internal/docscontract/skill_coverage_test.go`.
2. A file-name guard in the lookup — command: `grep -n 'func (m Map) SurfacesForPath' -A20 internal/skillcoverage/skillcoverage.go`; result: the body only matches sources by prefix or `path.Match`.
3. Another reader of the lookup — command: `grep -rn 'SurfacesForPath' --include='*.go' .`; result: only `internal/speccheck/skills_declaration.go` and the lookup's own test.

## ADR Alternatives

- ADR-0256: followed; it defines a surface's sources as the paths whose change can change it, which a Go test file cannot, and its rejected alternatives (a Git-only source comparison, fingerprints inside the map, a Baseline clause, a model judging the diff) stay rejected, since the rule still reads the authored map and needs no model.

## Integration Points

None external. The check and the fingerprint test read local files.

## Testing Approach

All tests are hermetic: temporary repositories through the existing
`skillsRepo` fixture, in-memory maps, and no network.

- `internal/skillcoverage/skillcoverage_test.go`:
  `TestAGoTestFileBelongsToNoBehaviorSurface`, over a map with a directory
  source `internal/cli/`, a wildcard `internal/config/*.go` and an exact
  `internal/cli/run_test.go`: `internal/config/config_test.go`,
  `internal/cli/deep/run_test.go` and `internal/cli/run_test.go` are refused
  with no surface; `internal/config/config.go`,
  `internal/config/test_helpers.go`, `internal/config/config_testing.go` and
  `internal/cli/deep/run.go` are read with their surface;
  `docs/user-guide/commands/spec.md` is read with none.
- `internal/speccheck/skills_declaration_test.go`:
  `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask`, on `skillsRepo`
  (whose map holds `internal/cli/example_*.go` and `internal/cli/sub/`), for
  each of `interface`, `creates`, `deletes`: `internal/cli/example_test.go`
  and `internal/cli/sub/other_test.go` give zero `SC-SKILLS-UNTASKED`;
  `internal/cli/example_testing.go` and `internal/cli/sub/other.go` give one
  naming the path; then a test-only `task_01` beside a `task_02` declaring
  `internal/cli/example.go` gives one naming `task_02`.
- The existing `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent`
  and `TestTheSkillCoverageMapIsCurrent` cover the fingerprint side.
- Sabotages, each run on the prototype: `FingerprintReads` always `true`
  (the speccheck test fails on `internal/cli/example_test.go`, the
  skillcoverage test on `internal/config/config_test.go`, and the bytes test
  on parsing the fixture `config_test.go`); always `false` (the speccheck
  test fails on `internal/cli/example_testing.go`, the bytes test on the
  `config keys` fingerprint); `!strings.Contains(filePath, "_test")` (both
  new tests fail on `example_testing.go` and `config_testing.go`); the guard
  removed from `SurfacesForPath` (both new tests fail on the test-file
  cases).

## Build Order

1. The predicate, the lookup guard, the fingerprint call, the tests, the
   guide and the reference, with the skill version and the fingerprints
   re-recorded last (task_01).
2. QA gate (task_02) (depends on: 1).

Specs 0332 to 0335 merge before this Spec runs. Specs 0333, 0334 and 0335
raise the Roundfix skill version and re-record the Behavior Surface Record;
task_01 raises the version above whatever they recorded through the record
command and re-records the fingerprints after its last guide edit, never by
hand. No step depends on their code.

## Risks & Considerations

- A future surface whose fingerprint reads Go test files would need the
  predicate to change; the bytes test and `TestTheSkillCoverageMapIsCurrent`
  would then fail on the stale fingerprint rather than pass silently.
- A future directory source could cover a `testdata` file that the help
  fingerprint never reads; the check then asks for a skills Task, which is
  safe, and the Non-Goals record why `testdata` stays in.

## Glossary

None.

## Decisions

- One exported predicate in `internal/skillcoverage`, read by the lookup and
  the fingerprint; no second copy in the Spec Check.
- Go test files only; `testdata` and generated files stay in.
- One implementation Task, because the code and its documentation share the
  Roundfix skill version record.
