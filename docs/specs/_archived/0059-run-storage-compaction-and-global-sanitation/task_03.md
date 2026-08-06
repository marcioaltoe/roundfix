---
task: task_03
spec: 0059-run-storage-compaction-and-global-sanitation
status: completed
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

## Result

Implemented metadata-only Artifact Root discovery and the global
`roundfix gc sanitize [--apply]` surface. The default invocation reports a
dry-run; apply removes only retention-eligible terminal Run directories or
directories whose Run ID is absent from the complete durable Run index.
Classification happens before scanning, and unsafe, missing, overridden, and
outside-Home roots stay preserved with their evidence.

Focused checks:

- `rtk go test ./internal/store -count=1 -run '^TestDiscoverArtifactRootsReturnsEveryRecordedRepository$' -v`
  passed. The first sandboxed attempt could not access the system Go build
  cache; the same focused command passed after granting that filesystem
  boundary.
- `rtk go test ./internal/cli -count=1 -run '^TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories$' -v`
  passed after the final test edit.
- `rtk go test ./internal/cli -count=1 -run '^TestRunGC(DryRunListsEligibleRunsAndChangesNothing|PrunesEligibleJournalsArtifactsAndOrphans|SkipsWhenJournalRetentionIsZero|Help)$' -v`
  passed all four existing per-repository GC cases.
- The commands under `## Verification` were not run; the Daemon owns them.

Acceptance evidence:

- Cross-repository discovery: `TestDiscoverArtifactRootsReturnsEveryRecordedRepository`
  creates Runs for two repositories and observes both recorded roots and their
  Run evidence without creating either root on disk.
- Six classifications and evidence: the sanitation CLI fixture observes
  `active`, `orphaned`, `missing`, `overridden`, `outside Roundfix Home`, and
  `unsafe`, plus a stated evidence line for each.
- Proven deletion boundary: apply removes an old terminal Run directory and
  an absent-Run directory, while active, recent, overridden, outside-Home, and
  unsafe paths remain.
- Review Artifact preservation: the fixture directly asserts the file under
  `reviews/pr-123/round-001/issue.md` exists after both apply invocations.
- Ambiguity preservation: overridden, outside-Home, and unsafe fixture paths
  remain after apply, and the report states why each root was preserved.
- Idempotence: the fixture derives the first reclaimed directory and byte
  counts, asserts both are non-zero, then asserts the second apply reclaims
  zero directories and zero bytes.
- Dry-run safety: every eligible and preserved fixture path is asserted present
  after the default sanitation invocation.
- Per-repository compatibility: all four pre-existing `TestRunGC...` cases pass
  unchanged for dry-run, apply, zero retention, and help behavior.

Follow-up boundary: Task 05 owns the authorized Roundfix Skill synchronization
for the new CLI contract; no Skill or generated digest file changed in this
slice.
