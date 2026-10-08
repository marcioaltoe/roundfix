---
task: task_01
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: pending
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
