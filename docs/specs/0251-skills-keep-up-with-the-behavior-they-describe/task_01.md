---
task: task_01
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: completed
type: backend
complexity: medium
---

# Task 01: The skill coverage package parses the map and the record and finds Lagging Surfaces

## Overview

Adds `internal/skillcoverage`, the pure package every later Task reads: it
parses and validates the Skill Coverage Map and the Behavior Surface Record,
computes fingerprints, matches a path to the surfaces whose sources it
touches, and compares a base and a target revision into changed surfaces and
their outcomes. It runs no process and no Git, so it is verifiable on its own
through table tests.

## Requirements

1. MUST create the package with the constants, types and functions of
   `_techspec.md` → Interfaces, with exactly those exported names and
   signatures. It MUST import no package of this repository and MUST NOT run
   a process or read the file system.
2. MUST implement `ParseMap` and `ParseRecord` with the validation of
   `_techspec.md` → Data Models 1 and 2: exact schema value, unknown fields
   refused, ids non-empty, unique and sorted, exactly one of `skills` or
   `uncovered` per entry, clean relative slash paths without `..` in `skills`
   and `sources`, valid `path.Match` patterns, and fingerprints of the form
   `sha256:` plus 64 lowercase hex digits. Each refusal MUST name the
   offending id or field.
3. MUST implement `EncodeRecord` as Data Model 2 states, so that parsing and
   re-encoding a record yields identical bytes, and `Fingerprint` as
   `sha256:` plus the lowercase hex SHA-256 of the bytes.
4. MUST implement `Map.SurfacesForPath` with the source grammar of
   Data Model 1: a source ending in `/` matches every path under it, any
   other source is matched by `path.Match` against the whole path.
5. MUST implement `Compare` exactly as `_techspec.md` → API Contract 1 states,
   including the `added`, `changed` and `removed` kinds, the judging entry,
   the outcome order `uncovered`, `described`, `reviewed`, `lagging`, the rule
   that a Coverage Review counts only when its text differs from the base
   entry's, and the ordering by id.
6. MUST add these tests to `internal/skillcoverage/skillcoverage_test.go`:
   - `TestParseMapRefusesEveryMalformedEntry`: one case per refusal of
     Requirement 2, each asserting the error names the id or field, and a
     valid map that parses.
   - `TestRecordRoundTripsByteForByte`: a record with unsorted input keys
     encodes, parses and re-encodes to the same bytes, ending in one newline.
   - `TestSurfacesForPathFollowTheSourceGrammar`: a directory source, a
     wildcard source and an exact source each match and miss as stated.
   - `TestCompareJudgesEveryChangedSurface`: an unchanged surface is absent;
     added, changed and removed surfaces each produce one Change; a changed
     covering skill gives `described`; an uncovered entry gives `uncovered`;
     a new or changed review gives `reviewed`; an unchanged review, a removed
     surface with an unchanged skill and an id with no entry give `lagging`.

## Subtasks

- [ ] Create the package with its constants, types and functions.
- [ ] Implement parsing, validation, encoding and fingerprints.
- [ ] Implement source matching and the comparison.
- [ ] Write the four table tests.

## Acceptance Criteria

- [ ] A malformed map or record is refused with an error that names the entry.
- [ ] A record re-encodes to identical bytes.
- [ ] `Compare` returns one Change per changed surface with the outcome API Contract 1 states, and a Coverage Review counts only in the range in which it changed.

## Context

- creates: `internal/skillcoverage/skillcoverage.go`
- creates: `internal/skillcoverage/skillcoverage_test.go`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestParseMapRefusesEveryMalformedEntry|TestRecordRoundTripsByteForByte|TestSurfacesForPathFollowTheSourceGrammar|TestCompareJudgesEveryChangedSurface)$' ./internal/skillcoverage 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestParseMapRefusesEveryMalformedEntry TestRecordRoundTripsByteForByte TestSurfacesForPathFollowTheSourceGrammar TestCompareJudgesEveryChangedSurface; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/skillcoverage` — expected: exit 0. Before this Task the package does not exist, so `go test` fails. After it, the four tests pass and the package builds for Windows.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 4; User Story 3
- `_techspec.md` → Interfaces; Data Model 1; Data Model 2; API Contract 1; Invariant 6; Build Order 1
- ADR-0256; ADR-0252

## Result

Implemented the pure `internal/skillcoverage` package with the TechSpec's
exported constants, types and signatures. Its imports are standard-library
only; it reads supplied bytes and runs no process or filesystem operation.
The pre-change inspection found that the package directory did not exist.
The only pre-existing modified path was this Task file, whose daemon-owned
`status: in_progress` is preserved.

Acceptance evidence:

- Malformed map or record refusal: `TestParseMapRefusesEveryMalformedEntry`
  exercises schemas, unknown fields, blank and duplicate ids, map ordering,
  exclusive non-empty coverage, invalid relative paths and source patterns,
  malformed fingerprints, missing/null surfaces and trailing JSON. Every
  negative case asserts that the error names its id or offending field;
  valid map and record cases parse.
- Byte-identical record encoding: `TestRecordRoundTripsByteForByte` starts
  with unsorted map keys and asserts exact two-space-indented output, sorted
  keys, one final newline and identical parse/re-encode bytes. It also
  checks SHA-256 against the known digest of `abc` and round-trips an empty
  surfaces object.
- Comparison and range-limited reviews:
  `TestCompareJudgesEveryChangedSurface` checks added, changed and removed
  records; target entries for additions/changes and base entries for
  removals; outcome precedence; changed, new and unchanged reviews; missing
  entries; sorted output; nil snapshots; and omission of unchanged surfaces.
  `TestSurfacesForPathFollowTheSourceGrammar` checks recursive directory
  prefixes, whole-path wildcards and exact sources, including misses and
  one result when two sources match the same surface.

Focused checks:

- `GOCACHE=/private/tmp/roundfix-task01-skillcoverage-cache rtk proxy go test -count=1 -cover ./internal/skillcoverage`
  exited 0; all four tests passed; statement coverage was 95.6%.
- `GOCACHE=/private/tmp/roundfix-task01-skillcoverage-cache rtk proxy go test -count=1 -v ./internal/skillcoverage`
  exited 0 and reported each of the four required tests passing.
- `rtk proxy gofmt -l internal/skillcoverage/skillcoverage.go internal/skillcoverage/skillcoverage_test.go`
  emitted no paths; `rtk proxy git -c core.fsmonitor=false diff --check`
  exited 0.
- Initial focused runs exposed a wrong test import prefix (corrected to the
  module's `roundfix` prefix) and sandbox denial of the default Go build
  cache. The task-specific writable cache resolved the latter.
- `rtk make verify-incremental` initially exited 2: the sandbox denied
  process-table access in two CLI process-owner tests, and suiteguard
  detected implementation/Result edits made while that run was active.
  The rerun with host process-table access and the worktree held unchanged
  exited 0. Go analysis, the repository tests, skill synchronization and
  skill checks, and the CLI build passed. The final package implementation
  was unchanged throughout that successful run.

The authored Verification command, including its Windows build, remains
for the Daemon. No Task status, Task Graph, other Task file, commit, push or
Pull Request was changed by this implementation turn.
