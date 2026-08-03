---
task: task_03
spec: 0074-git-spawn-economy
status: completed
type: backend
complexity: medium
---

# Task 03: Batch object reads in assets-sync provenance

## Overview

The second per-file loop: `internal/baseline/assets_sync_git.go` reads
source files with one `git cat-file blob` spawn per tree entry while
computing the provenance digest. This Task applies the batch-reader shape
proven in task_02 to this loop, keeping the reader package-private and the
assets-sync error vocabulary intact.

## Requirements

1. MUST route the provenance loop through a batch reader with the same
   framing-exact semantics task_02 established.
2. MUST preserve the assets-sync error vocabulary: unsafe path, non-blob
   entries, and read failures report exactly as before.
3. MUST keep the provenance digest byte-identical for identical inputs —
   the existing sync and compatibility tests are the characterization and
   pass unmodified.
4. MUST cover the same framing edge cases where the shape differs from
   task_02's reader; shared cases need no duplication if the reader is
   shared within the package.

## Subtasks

- [ ] Route the loop through the batch reader.
- [ ] Preserve the error vocabulary.
- [ ] Confirm digests are byte-identical via the existing tests.

## Acceptance Criteria

- [ ] Provenance over N files spawns one `cat-file` process, not N.
- [ ] `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` and the
      whole assets-sync family pass unmodified.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`
      and this task file.

## Verification

- `go test ./internal/baseline -count=1 -run 'AssetsSync' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `grep -n "cat-file" internal/baseline/assets_sync_git.go | grep -qv -- "--batch" && exit 1 || exit 0`
  — expected: exit 0; no per-file spawn remains.

## References

- `_prd.md` → Core Feature 1.
- `_techspec.md` → Build Order 3 (reuses the reader shape proven in
  task_02).
- ADR-0090.

## Result

Assets-sync provenance now validates the committed tree before opening one
shared `git cat-file --batch` reader, reads every blob in tree order, closes
the reader on success and failure, and computes the same portable digest from
the returned bytes. The Task 02 runner seam was renamed from restore-specific
to package-generic names so both provenance loops use the same framing-exact
reader without duplicating its protocol implementation or unit cases.

The pre-change focused signal was
`GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -count=1 -run '^TestAssetsSyncCommittedTreeDigestReadsManyFilesThroughOneBatchProcess$'`;
it failed with `batch process starts = 0, want 1`, proving that the provenance
loop still used the per-file path before the implementation change.

Focused implementation evidence by acceptance criterion:

1. One batch process for N files:
   `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -count=1 -run '^TestAssetsSyncCommittedTreeDigest(ReadsManyFilesThroughOneBatchProcess|KeepsTreeAndReadErrors)$'`
   passed 5 cases. The three-file case observed one batch open, one `ls-tree`
   call, one close, and the three object requests in tree order. The companion
   cases preserved the exact unsafe-path and non-blob diagnostics and the
   injected read error chain.
2. Existing digest and compatibility characterization:
   `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -count=1 -run '^TestAssetsSyncCompatibilityMatchesMaintainedPythonContract$'`
   passed the unmodified maintained-Python compatibility case. The unmodified
   refresh and provenance cases passed 7 cases through
   `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -count=1 -run '^Test(BaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent|AssetsSyncProvenanceAndPreMutationRefusals)$'`.
   The shared framing command
   `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -count=1 -run '^TestBatchObjectReader(ReturnsMultipleObjectsInRequestOrder|ReportsMissingObject|ReturnsZeroByteBlob|PreservesFramingDelimitersInContent|ReportsProcessDeathMidStream)$'`
   passed all 5 Task 02 protocol cases without duplicating them for assets
   sync.
3. Changed-path scope: `git -c core.fsmonitor=false status --porcelain`
   listed only this task file and files under `internal/baseline/`;
   `git diff --check` exited 0.

The Daemon-owned commands under `## Verification`, including the whole
assets-sync family and whole `internal/baseline` package, were not run in this
Agent turn.
