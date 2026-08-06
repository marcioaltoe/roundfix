---
task: task_03
spec: 0059-run-storage-compaction-and-global-sanitation
status: pending
type: backend
complexity: high
---

# Task 03: Discover and classify every Artifact Root

## Overview

`PruneTerminalRuns` works machine-wide over the Run Database, but
`internal/cli/gc.go` resolves artifact cleanup against the current
repository's Artifact Root. One side sees every repository, the other sees one,
so roots belonging to other repositories accumulate forever and are never even
listed.

This slice closes the scope mismatch, dry-run first, and preserves anything it
cannot prove.

## Requirements

1. MUST discover every Roundfix-owned Artifact Root from durable Run metadata.
2. MUST NOT scan the filesystem for roots outside recorded metadata. Guessing
   would delete directories Roundfix never created.
3. MUST classify each root as active, orphaned, missing, overridden, outside
   Roundfix Home, or unsafe, and report the classification with the evidence
   that produced it.
4. MUST remove only directories proven to belong to eligible or absent Runs.
5. MUST preserve Review Artifacts and every ambiguous path, with the reason.
6. MUST be dry-run by default; mutation requires an explicit apply.
7. MUST leave per-repository GC available and unchanged.
8. MUST assert idempotence — a second run reclaims zero — rather than asserting
   any recorded byte count.

## Subtasks

- [ ] Add discovery from durable Run metadata.
- [ ] Add the six classifications with their evidence.
- [ ] Assert idempotence and one preservation case per classification.

## Acceptance Criteria

- [ ] Discovery returns roots from other repositories, not only the current one.
- [ ] Each of the six classifications is produced for a fixture case and
      reported with its evidence.
- [ ] Only proven-eligible or absent directories are removed.
- [ ] Review Artifacts are never removed, asserted directly.
- [ ] An ambiguous root is preserved with a stated reason.
- [ ] A second sanitation run reclaims zero, asserted as a relation.
- [ ] Dry-run removes nothing, asserted on the filesystem.
- [ ] Per-repository GC behaves exactly as it does today.

## Context

- interface: `internal/cli/gc.go`
- interface: `internal/store/journal.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/cli ./internal/store -count=1 -run 'Sanitat|ArtifactRoot|Discover|Classif' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the sanitation tests ran and passed.
- `go test ./internal/cli ./internal/store -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Feature 2; User Stories 3 and 5; Success Metric 2.
- `_techspec.md` → System Architecture; Build Order 3; Risks & Considerations.
- ADR-0053.
