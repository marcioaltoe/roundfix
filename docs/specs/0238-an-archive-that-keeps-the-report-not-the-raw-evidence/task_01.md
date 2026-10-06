---
task: task_01
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
status: pending
type: test
complexity: low
---

# Task 01: No test or record reads archived QA evidence

## Overview

Spec 0214's ablation removed archived QA evidence in a disposable clone and
broke both repository gates. `make verify` failed because the coverage record
lists three Go packages that live inside archived evidence, and
`make verify-docs` failed because the repository-copy helper stops on a
tracked file missing from the working tree. This Task removes both
dependencies, so a later cut of archived evidence cannot break either gate.
It changes no product behavior.

## Requirements

1. MUST make the coverage record count no Go package whose import path is
   under `roundfix/docs/`: `collectCoverageRecord` in
   `internal/spec/coverage_test.go` skips them, so packages that live inside
   any Spec's evidence never enter the record or its comparison (TechSpec
   Invariant 11).
2. MUST re-record `docs/references/coverage-record.json` with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record`.
   The diff removes the three entries under `roundfix/docs/history/specs/`
   and may add test names the suite gained since the last record; it MUST
   NOT remove any other package or test name. The record is never edited by
   hand.
3. MUST make `copyTrackedRepository` in
   `skills/baseline_skill_contract_test.go` skip a path that `git ls-files`
   lists but that does not exist in the working tree. Every other read,
   stat or write error still fails the calling test.
4. MUST add `TestCoverageRecordCountsNoPackageUnderDocs` to
   `internal/spec/coverage_test.go`, a unit test of the package filter with
   a path under `roundfix/docs/history/specs/`, one under
   `roundfix/docs/specs/`, and ordinary packages that stay counted.
5. MUST add `skills/copy_tracked_repository_test.go` with
   `TestCopyTrackedRepositorySkipsATrackedFileMissingFromTheWorkingTree`.
   It builds a temporary Git repository with isolated configuration, commits
   two files, deletes one from the working tree without committing, copies
   the repository with `copyTrackedRepository`, and asserts that the kept
   file is copied byte-for-byte and the deleted one is absent.
6. MUST leave `internal/cli/review_scope_test.go` unchanged. It writes its
   evidence paths into a disposable repository and reads no archived file.
7. MUST NOT delete, move or rewrite any file under `docs/history/`.

## Subtasks

- [ ] Skip packages under `docs/` in the coverage collector.
- [ ] Re-record the coverage record with its flag.
- [ ] Let the repository copy skip a tracked file missing from the working tree.
- [ ] Add the two tests.

## Acceptance Criteria

- [ ] The coverage record lists no package under `roundfix/docs/`, and
      `TestCoverageEquivalence` passes.
- [ ] A repository copy with a tracked file deleted from the working tree
      succeeds and leaves that file out.
- [ ] No file under `docs/history/` changed.

## Context

- instruction: `docs/references/archived-evidence-measurement.md`
- interface: `internal/spec/coverage_test.go`
- interface: `docs/references/coverage-record.json`
- interface: `skills/baseline_skill_contract_test.go`
- creates: `skills/copy_tracked_repository_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCoverageEquivalence|TestCoverageRecordCountsNoPackageUnderDocs)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCoverageEquivalence TestCoverageRecordCountsNoPackageUnderDocs; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestCopyTrackedRepositorySkipsATrackedFileMissingFromTheWorkingTree$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestCopyTrackedRepositorySkipsATrackedFileMissingFromTheWorkingTree " || { printf 'missing pass: TestCopyTrackedRepositorySkipsATrackedFileMissingFromTheWorkingTree\n' >&2; exit 1; }; ! grep -q 'roundfix/docs/' docs/references/coverage-record.json || { printf 'coverage record still lists a package under docs\n' >&2; exit 1; }` — expected: exit 0; before this Task the two new tests do not exist and the record lists three packages under `roundfix/docs/`, so the command fails; after it both new tests and the coverage equivalence pass and the record names no package under `docs/`.

## References

- `_prd.md` → Goals; User Story 5; Core Feature 7; Success Metric 4; Acceptance evidence
- `_techspec.md` → Invariant 11; System Architecture; Testing Approach; Build Order 1
- ADR-0243; ADR-0215
